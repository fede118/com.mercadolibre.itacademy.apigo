package domains

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"../utils"
	"sync"
	"time"
)

type Country struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Locale             string `json:"locale"`
	CurrencyID         string `json:"currency_id"`
	DecimalSeparator   string `json:"decimal_separator"`
	ThousandsSeparator string `json:"thousands_separator"`
	TimeZone           string `json:"time_zone"`
	GeoInformation struct {
		Location struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location"`
	} `json:"geo_information"`
	States []struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"states"`
}

func (country *Country) Get() *utils.ApiError {
	if country.ID == "" {
		return &utils.ApiError{
			Message: "Site ID is Empty.",
			Status:  http.StatusBadRequest,
		}
	}

	url := fmt.Sprintf("%s%s", utils.UrlCountry, country.ID)

	response, err := http.Get(url)
	if err != nil {
		return &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	if err := json.Unmarshal(data, &country); err != nil {
		return &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}


	return nil
}

func (country *Country) GetWithWaitGroup (waitGroup *sync.WaitGroup, apiError *utils.ApiError) {
	if country.ID == "" {
		apiError.Message = "Site ID is Empty."
		apiError.Status = http.StatusBadRequest
		return
	}

	url := fmt.Sprintf("%s%s", utils.UrlCountry, country.ID)

	response, err := http.Get(url)
	if err != nil {
		apiError.Message = err.Error()
		apiError.Status = http.StatusInternalServerError
		return
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		apiError.Message = err.Error()
		apiError.Status = http.StatusInternalServerError
		return
	}

	if err := json.Unmarshal(data, &country); err != nil {
		apiError.Message = err.Error()
		apiError.Status = http.StatusInternalServerError
		return
	}

	waitGroup.Done()
	return
}

func (country *Country) GetWithChannel(channel chan Result) {
	if country.ID == "" {
		channel <- Result {
			ApiError: &utils.ApiError{
				Message: "Site ID is Empty.",
				Status: http.StatusBadRequest,
			},
		}

		return
	}

	url := fmt.Sprintf("%s%s", utils.UrlMockCountries, country.ID)

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	response, err := client.Get(url)

	if err != nil {
		utils.CircuitBreakerInstance.PlusError()
		channel <- Result {
			ApiError: &utils.ApiError{
				Message: err.Error(),
				Status: http.StatusInternalServerError,
			},
		}

		return
	}

	if response.StatusCode == 500 {
		utils.CircuitBreakerInstance.PlusError()
		channel <- Result {
			ApiError: &utils.ApiError{
				Message: "500 back from server",
				Status: http.StatusInternalServerError,
			},
		}
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		channel <- Result {
			ApiError: &utils.ApiError{
				Message: err.Error(),
				Status: http.StatusInternalServerError,
			},
		}

		return
	}

	if err := json.Unmarshal(data, &country); err != nil {
		channel <- Result {
			ApiError: &utils.ApiError{
				Message: err.Error(),
				Status: http.StatusInternalServerError,
			},
		}

		return
	}

	result := Result{
		Country: country,
	}

	channel <- result

	return
}

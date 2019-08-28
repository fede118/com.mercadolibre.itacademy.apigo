package domains

import(
	"../utils"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"sync"
	"time"
)

type Site struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	CountryID          string   `json:"country_id"`
	SaleFeesMode       string   `json:"sale_fees_mode"`
	MercadopagoVersion int      `json:"mercadopago_version"`
	DefaultCurrencyID  string   `json:"default_currency_id"`
	ImmediatePayment   string   `json:"immediate_payment"`
	PaymentMethodIds   []string `json:"payment_method_ids"`
	Settings           struct {
		IdentificationTypes      []string `json:"identification_types"`
		TaxpayerTypes            []string `json:"taxpayer_types"`
		IdentificationTypesRules []struct {
			IdentificationType string `json:"identification_type"`
			Rules              []struct {
				EnabledTaxpayerTypes []string `json:"enabled_taxpayer_types"`
				BeginsWith           string   `json:"begins_with"`
				Type                 string   `json:"type"`
				MinLength            int      `json:"min_length"`
				MaxLength            int      `json:"max_length"`
			} `json:"rules"`
		} `json:"identification_types_rules"`
	} `json:"settings"`
	Currencies []struct {
		ID     string `json:"id"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	Categories []struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"categories"`
}

func (site *Site) Get() *utils.ApiError {
	if site.ID == "" {
		return &utils.ApiError{
			Message: "Site ID is Empty.",
			Status: http.StatusBadRequest,
		}
	}

	url := fmt.Sprintf("%s%s", utils.UrlSite, site.ID)

	response, err := http.Get(url)
	if err != nil {
		return &utils.ApiError{
			Message: err.Error(),
			Status: http.StatusInternalServerError,
		}
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &utils.ApiError{
			Message: err.Error(),
			Status: http.StatusInternalServerError,
		}
	}

	if err := json.Unmarshal(data, &site); err != nil {
		return &utils.ApiError{
			Message: err.Error(),
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

func (site *Site) GetWithWaitGroup (waitGroup *sync.WaitGroup, apiError *utils.ApiError) {
	if site.ID == "" {
		 apiError.Message = "Site ID is Empty."
		 apiError.Status = http.StatusBadRequest
		 return
	}

	url := fmt.Sprintf("%s%s", utils.UrlSite, site.ID)

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

	if err := json.Unmarshal(data, &site); err != nil {
		apiError.Message = err.Error()
		apiError.Status = http.StatusInternalServerError
		return
	}

	waitGroup.Done()

	return
}

func (site *Site) GetWithChannel(channel chan Result) {
	if site.ID == "" {
		channel <- Result {
			ApiError: &utils.ApiError{
				Message: "Site ID is Empty.",
				Status: http.StatusBadRequest,
			},
		}

		return
	}

	url := fmt.Sprintf("%s%s", utils.UrlMockSites, site.ID)

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

	if err := json.Unmarshal(data, &site); err != nil {
		channel <- Result {
			ApiError: &utils.ApiError{
				Message: err.Error(),
				Status: http.StatusInternalServerError,
			},
		}

		return
	}

	result := Result{
		Site: site,
	}

	channel <- result

	return
}
package services

import(
	"../domains"
	"../utils"
	"sync"
)

const (
	numberOfRoutines = 2
)

func GetResult (userId int) (*domains.Result, *utils.ApiError) {
	user := domains.User{
		ID: userId,
	}
	err := user.Get()
	if err != nil {
		return nil, err
	}

	country := domains.Country{
		ID: user.CountryID,
	}

	site := domains.Site{
		ID: user.SiteID,
	}

	err = site.Get()
	if err != nil {
		return nil, err
	}
	err = country.Get()
	if err != nil {
		return nil, err
	}

	resp := domains.Result{
		User: &user,
		Country: &country,
		Site: &site,
	}

	return &resp, nil
}

func GetResultWithWaitGroup (userId int) (*domains.Result, *utils.ApiError) {
	var waitGroup sync.WaitGroup
	user := domains.User{
		ID: userId,
	}
	err := user.Get()
	if err != nil {
		return nil, err
	}

	country := domains.Country{
		ID: user.CountryID,
	}

	site := domains.Site{
		ID: user.SiteID,
	}

	apiError := utils.ApiError{}
	waitGroup.Add(numberOfRoutines)
	go site.GetWithWaitGroup(&waitGroup, &apiError)
	go country.GetWithWaitGroup(&waitGroup, &apiError)

	waitGroup.Wait()

	resp := domains.Result{
		User: &user,
		Country: &country,
		Site: &site,
	}

	return &resp, nil
}

func GetResultWithChannel (userId int) (*domains.Result, *utils.ApiError) {
	user := domains.User{
		ID: userId,
	}
	err := user.Get()
	if err != nil {
		return nil, err
	}

	country := domains.Country{
		ID: user.CountryID,
	}

	site := domains.Site{
		ID: user.SiteID,
	}


	channel := make(chan domains.Result, numberOfRoutines)

	go site.GetWithChannel(channel)
	go country.GetWithChannel(channel)

	finalResult := domains.Result{}
	for i := 0; i < numberOfRoutines; i++ {

		result := <-channel
		if result.ApiError != nil {
			return nil, result.ApiError
		}

		if result.Site != nil {
			finalResult.Site = result.Site
		}

		if result.Country != nil {
			finalResult.Country = result.Country
		}

	}

	finalResult.User = &user

	return &finalResult, nil
}


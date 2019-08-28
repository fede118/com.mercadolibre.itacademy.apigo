package utils

import (
	"time"
	"net/http"
)

const (
	CLOSED = 0
	HALFOPEN = 1
	OPEN = 2
	ERRORTRESHOLD = 3
	TIMEOUT = time.Duration(10)
)

var (
	CircuitBreakerInstance = CircuitBreaker{}
)

type CircuitBreaker struct {
	State 			int
	ErrorTreshold 	int
	ErrorCounter	int
	TimeOut			time.Duration
}

func NewState() CircuitBreaker {
	return CircuitBreaker{
		State: CLOSED,
		ErrorTreshold: ERRORTRESHOLD,
		ErrorCounter: 0,
		TimeOut: 10,
	}
}


func main() {

}


func (circuitBreaker *CircuitBreaker) SetState(newState int) {
	//pasa de estado
	if newState != CLOSED && newState != OPEN  && newState != HALFOPEN {
		// error
		return
	}

	circuitBreaker.State = newState

}

func (circuitBreaker *CircuitBreaker) Reset() {
	circuitBreaker.ErrorCounter = 0
}

func (circuitBreaker *CircuitBreaker) PlusError() {
	println("Added 1 Error to CircuitBreaker")
	circuitBreaker.ErrorCounter++
	if circuitBreaker.ErrorCounter >= 3 {
		go circuitBreaker.CloseConnectionWithTimeOut()
	}
}

func (circuitBreaker *CircuitBreaker) CloseConnectionWithTimeOut() {
	for {
		println("Setting state to OPEN")
		circuitBreaker.SetState(OPEN)
		time.Sleep(time.Second * circuitBreaker.TimeOut)
		circuitBreaker.Reset()

		println("Setting state to HALFOPEN")
		circuitBreaker.SetState(HALFOPEN)


		println("Pinging servers")

		timeout := time.Duration(5 * time.Second)
		client := http.Client{
			Timeout: timeout,
		}

		pongUsers, err := client.Get(UrlMockUsersPing)
		//fmt.Println("Users Pong: ", pongUsers.StatusCode)
		if err != nil || pongUsers.StatusCode != 200 {
			continue
		}

		pongSites, err := client.Get(UrlMockSitesPing)
		//fmt.Println("Sites Pong: ", pongSites.StatusCode)
		if err != nil || pongSites.StatusCode != 200 {
			continue
		}


		pongCountries, err := client.Get(UrlMockCountriesPing)
		//fmt.Println("Countries Pong: ", pongCountries.StatusCode)
		if err != nil || pongCountries.StatusCode != 200 {
			continue
		}


		if pongUsers.StatusCode == 200 && pongCountries.StatusCode == 200 && pongSites.StatusCode == 200 {
			println("setting state to Closed")
			circuitBreaker.SetState(CLOSED)
			circuitBreaker.Reset()
			return
		}
	}
}

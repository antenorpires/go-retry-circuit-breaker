package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sony/gobreaker"
)

type Result struct {
	Status string `json:"status"`
}

var retryClient *retryablehttp.Client
var cb *gobreaker.CircuitBreaker

func main() {

	retryClient = retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = 500 * time.Millisecond
	retryClient.RetryWaitMax = 2 * time.Second
	retryClient.Logger = nil

	cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "ServiceA-CB",
		MaxRequests: 5,
		Timeout:     30 * time.Second,

		ReadyToTrip: func(counts gobreaker.Counts) bool {
			if counts.Requests < 10 {
				return false
			}
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= 0.5
		},

		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("[CB] %s: %s → %s\n", name, from, to)
		},
	})

	http.HandleFunc("/", Home)
	http.HandleFunc("/process", Process)

	log.Println("Service A running on http://localhost:9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func Home(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/home.html"))
	_ = t.Execute(w, Result{})
}

func Process(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	log.Println("Processing ID:", id)

	result, err := callWithCircuitBreaker("http://localhost:9091", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	t := template.Must(template.ParseFiles("templates/home.html"))
	_ = t.Execute(w, result)
}

func callWithCircuitBreaker(url string, id string) (Result, error) {

	res, err := cb.Execute(func() (interface{}, error) {
		return callWithRetry(url, id)
	})

	if err != nil {
		return Result{}, err
	}

	return res.(Result), nil
}

func callWithRetry(urlMicroService string, id string) (Result, error) {

	values := url.Values{}
	values.Add("id", id)

	req, err := retryablehttp.NewRequest(
		http.MethodPost,
		urlMicroService,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return Result{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := retryClient.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return Result{}, errors.New("provider error 5xx")
	}

	var result Result
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Result{}, err
	}

	return result, nil
}

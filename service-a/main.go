package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type Result struct {
	Status string `json:"status"`
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/process", Process)

	log.Println("Request service running on :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func Home(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/home.html"))
	t.Execute(w, Result{})
}

func Process(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	log.Println("Processing ID:", id)

	result, err := makeHttpCall("http://localhost:9091", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t := template.Must(template.ParseFiles("templates/home.html"))
	t.Execute(w, result)
}

func makeHttpCall(urlMicroService string, id string) (Result, error) {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5                         // Máximo de 5 tentativas de retry
	retryClient.HTTPClient.Timeout = 2 * time.Second // Timeout de 2 segundos para cada requisição

	values := url.Values{}
	values.Add("id", id)

	reqBody := strings.NewReader(values.Encode())

	req, err := retryablehttp.NewRequest("POST", urlMicroService, reqBody)
	if err != nil {
		return Result{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := retryClient.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer res.Body.Close()

	var result Result
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return Result{}, err
	}

	return result, nil
}

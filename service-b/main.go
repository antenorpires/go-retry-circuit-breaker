package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Result struct {
	Status string `json:"status"`
}

func main() {
	http.HandleFunc("/", Home)

	log.Println("Service B running on http://localhost:9091")
	log.Fatal(http.ListenAndServe(":9091", nil))
}

func Home(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	log.Println("Received ID:", id)

	time.Sleep(2 * time.Second)

	if id == "fail" {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	result := Result{
		Status: "failed",
	}

	if id == "123" {
		result.Status = "success"
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Result struct {
	Status string `json:"status"`
}

func main() {
	http.HandleFunc("/", Home)
	log.Println("Provider running on :9091")
	log.Fatal(http.ListenAndServe(":9091", nil))
}

func Home(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")

	result := Result{
		Status: "failed",
	}

	if id == "123" {
		result.Status = "success"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp`
}

type Location struct {
	City string `json:"city"`
	ID   int    `json:"id"`
}

var objects = []Location{
	{ID: 1, City: "default"},
	{ID: 2, City: "default"},
	{ID: 3, City: "default"},
}

var currentLocation string

func handleObjects(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getObjects(w, r)
	case "PUT":
		setObjectLocation(w, r)
	case "PATCH":
		resetObjectLocation(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
	}
}

func getObjects(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func setObjectLocation(w http.ResponseWriter, r *http.Request) {

}

func resetObjectLocation(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing city parameter"})
		return
	}

	var updatedObjects []Location
	for i, obj := range objects {
		if obj.City == city {
			objects[i].City = "default"
			updatedObjects = append(updatedObjects, objects[i])
		}
	}

	if len(updatedObjects) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "No objects found in the given city"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedObjects)
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction Transaction

	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if transaction.Amount <= 0 {
		http.Error(w, "amount can be in negative. Bad Request", http.StatusBadRequest)
		return
	}

	if transaction.Timestamp.After(time.Now()) {
		http.Error(w, "can not process the request entity", http.StatusUnprocessableEntity)
	}
	timeDifference := time.Since(transaction.Timestamp)

	if timeDifference.Seconds() > 60 {
		http.Error(w, "No content", http.StatusNoContent)
		return
	}

}

func handleTransactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		deleteTransaction(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
	}
}

func deleteTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing ID parameter"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID parameter"})
		return
	}

	var deletedTx Transaction
	found := false
	var transactions []int
	for i, val := range transactions {
		if val.ID == id {
			transactions = append(transactions[:i], transactions[i+1:]...)
			deletedTx = val
			found = true
			break
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Transaction not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deletedTx)
}

func statisticsHandler(w http.ResponseWriter, r *http.Request) {
	// check if location is set and match with current location
	if currentLocation != "" {
		// get client's IP address and check if it's in the same location as the currentLocation
		ip := getClientIP(r)
		if !isLocationAllowed(ip, currentLocation) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

}

func setLocationHandler(w http.ResponseWriter, r *http.Request) {
	// parse the request body
	var data struct {
		City string `json:"city"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// set the current location

	currentLocation = data.City

	// send a response with status code 201 (created)
	w.WriteHeader(http.StatusCreated)
}

func isLocationAllowed(ip, location string) bool {
	// use a third-party API to get the location of the IP address
	resp, err := http.Get("https://ipapi.co/" + ip + "/json/")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// parse the JSON response
	var data struct {
		City string `json:"city"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return false
	}

	// check if the city matches the current location
	return data.City == location
}

func getClientIP(r *http.Request) string {
	// get the IP address from the X-Forwarded-For header
	ips := r.Header.Get("X-Forwarded-For")
	if ips != "" {
		return strings.Split(ips, ",")[0]
	}

	// if X-Forwarded-For header is not present, use the remote address
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	return ip
}

func main() {
	transactions := []Transaction{
		{Amount: 10.0},
		{Amount: 20.0},
		{Amount: 30.0},
		{Amount: 40.0},
		{Amount: 50.0},
	}

	var sum, min, max, count float64
	for i, tx := range transactions {
		sum += tx.Amount
		if i == 0 || tx.Amount < min {
			min = tx.Amount
		}
		if i == 0 || tx.Amount > max {
			max = tx.Amount
		}
		count++
	}
	average := sum / count

	currentLocation = ""

	http.HandleFunc("/statistics", statisticsHandler)
	http.HandleFunc("/setLocation", setLocationHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))

	fmt.Println("Sum:", sum)
	fmt.Println("Average:", average)
	fmt.Println("Min:", min)
	fmt.Println("Max:", max)
	fmt.Println("Count:", count)

	http.HandleFunc("/transactions", createTransaction)
	http.HandleFunc("/transactions", handleTransactions)
	http.HandleFunc("/transactions", deleteTransaction)
}

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetLocationAPI(t *testing.T) {
	// create a test server with the handler function for the API
	ts := httptest.NewServer(http.HandlerFunc(handleObjects))

	// make a sample request body
	data := map[string]string{
		"city": "Bangalore",
	}
	body, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// make a request to the API with the sample request body
	res, err := http.Post(ts.URL+"/setLocation", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// check that the response status code is 201 (created)
	if res.StatusCode != http.StatusCreated {
		t.Errorf("unexpected status code: got %v, want %v", res.StatusCode, http.StatusCreated)
	}

	// check that the location was set correctly
	var currentLocation string
	if currentLocation != "Bangalore" {
		t.Errorf("unexpected location: got %v, want %v", currentLocation, "Bangalore")
	}
}

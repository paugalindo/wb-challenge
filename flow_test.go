package wbchallenge_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"
	"time"
	"wb-challenge/cmd/api/bootstrap"
)

const (
	srvPort = "90"
	baseURL = "http://localhost:" + srvPort
)

func TestFlow(t *testing.T) {
	cfg := bootstrap.GetConfigFromEnv()
	cfg.ServicePort = srvPort
	go bootstrap.Run(context.Background(), cfg, log.Default())

	time.Sleep(1 * time.Second)

	s, _, err := doReq(http.MethodGet, baseURL+"/status", "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if s != http.StatusOK {
		t.Errorf("Expected status 200, got %d", s)
	}

	// Add some vehicles
	jsonData := `
[
  {
    "id": 1,
    "seats": 4
  },
  {
    "id": 2,
    "seats": 6
  }
]`
	s, _, err = doReq(http.MethodPut, baseURL+"/evs", jsonData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if s != http.StatusOK {
		t.Errorf("Expected status 200, got %d", s)
	}

	// Create a journey for a group
	jsonData = `
{
  "id": 1,
  "people": 6
}`
	s, _, err = doReq(http.MethodPost, baseURL+"/journey", jsonData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if s != http.StatusAccepted {
		t.Errorf("Expected status 202, got %d", s)
	}

	time.Sleep(time.Second)

	// Locate first group
	jsonData = `
{
  "id": 1
}`
	s, b, err := doReq(http.MethodPost, baseURL+"/locate", jsonData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if s != http.StatusOK {
		t.Errorf("Expected status 202, got %d", s)
	}

	response := struct {
		VehicleID int `json:"vehicle_id"`
	}{}
	err = json.Unmarshal([]byte(b), &response)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if response.VehicleID != 2 {
		t.Errorf("Expected group 1 locate in vehicle 2, got %d", response.VehicleID)
	}

	// Create a journey for another group
	jsonData = `
{
  "id": 2,
  "people": 4
}`
	s, _, err = doReq(http.MethodPost, baseURL+"/journey", jsonData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if s != http.StatusAccepted {
		t.Errorf("Expected status 202, got %d", s)
	}

	time.Sleep(time.Second)

	// Locate second group
	jsonData = `
{
  "id": 2
}`
	s, b, err = doReq(http.MethodPost, baseURL+"/locate", jsonData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if s != http.StatusOK {
		t.Errorf("Expected status 202, got %d", s)
	}

	response = struct {
		VehicleID int `json:"vehicle_id"`
	}{}
	err = json.Unmarshal([]byte(b), &response)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if response.VehicleID != 1 {
		t.Errorf("Expected group 1 locate in vehicle 1, got %d", response.VehicleID)
	}

	// Create a journey for another group
	jsonData = `
{
  "id": 3,
  "people": 5
}`
	s, _, err = doReq(http.MethodPost, baseURL+"/journey", jsonData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if s != http.StatusAccepted {
		t.Errorf("Expected status 202, got %d", s)
	}

	time.Sleep(time.Second)

	// Locate third group
	jsonData = `
{
  "id": 3
}`
	s, _, err = doReq(http.MethodPost, baseURL+"/locate", jsonData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if s != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", s)
	}

	// Dropoff first group
	jsonData = `
{
  "id": 1
}`
	s, _, err = doReq(http.MethodPost, baseURL+"/dropoff", jsonData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if s != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", s)
	}

	time.Sleep(time.Second)

	// Locate third group
	jsonData = `
{
  "id": 3
}`
	s, b, err = doReq(http.MethodPost, baseURL+"/locate", jsonData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if s != http.StatusOK {
		t.Errorf("Expected status 202, got %d", s)
	}

	response = struct {
		VehicleID int `json:"vehicle_id"`
	}{}
	err = json.Unmarshal([]byte(b), &response)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if response.VehicleID != 2 {
		t.Errorf("Expected group 3 locate in vehicle 2, got %d", response.VehicleID)
	}
}

func doReq(method, url, jsonData string) (int, string, error) {
	client := &http.Client{}

	var req *http.Request
	var err error
	if jsonData != "" {
		req, err = http.NewRequest(method, url, bytes.NewBuffer([]byte(jsonData)))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return 0, "", err
	}

	if jsonData != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}

	return resp.StatusCode, string(body), nil
}

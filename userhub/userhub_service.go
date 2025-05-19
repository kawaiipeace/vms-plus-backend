package userhub

import (
	"bytes"
	"encoding/json"
	"net/http"
	"vms_plus_be/config"
)

type Request_1000101001 struct {
	Phone string `json:"phone" example:"0818088770"`
}
type Response_1000101001 struct {
	IsValid bool   `json:"is_valid" example:"true"`
	Message string `json:"message" example:"success"`
}

func CheckPhoneNumber(phoneNumber string) (bool, error) {
	if phoneNumber != "" {
		return true, nil
	}
	client := &http.Client{}

	// Create request body
	reqBody := Request_1000101001{
		Phone: phoneNumber,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return false, err
	}

	// Create request
	req, err := http.NewRequest("POST", config.AppConfig.UserHubEndPoint+"/service/1000101001", bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ServiceKey", config.AppConfig.UserHubSecretKey)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Parse response
	var response Response_1000101001
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return false, err
	}

	return response.IsValid, nil
}

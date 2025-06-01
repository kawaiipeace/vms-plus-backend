package userhub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"vms_plus_be/config"
	"vms_plus_be/models"
)

func CheckPhoneNumber(phoneNumber string) (bool, error) {
	client := &http.Client{}

	// Create request body
	reqBody := ServiceCheckPhoneNumberRequest{
		ServiceCode: "vms",
		Phone:       phoneNumber,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return false, err
	}

	// Create request
	req, err := http.NewRequest("POST", config.AppConfig.UserHubEndPoint+"/service/check-phone-number", bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ServiceKey", config.AppConfig.UserHubServiceKey)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Parse response
	var response ServiceCheckPhoneNumberResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return false, err
	}

	return response.IsValid, nil
}

func LoginUser(loginBy string, empID string, identityNo string, phone string, ipAddress string) (ServiceUserInfoResponse, error) {
	client := &http.Client{}

	// Create request body
	reqBody := ServiceLoginUserRequest{
		ServiceCode: "vms",
		LoginBy:     loginBy,
		EmpID:       empID,
		IdentityNo:  identityNo,
		Phone:       phone,
		IpAddress:   ipAddress,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return ServiceUserInfoResponse{}, err
	}

	// Create request
	req, err := http.NewRequest("POST", config.AppConfig.UserHubEndPoint+"/service/login-user", bytes.NewBuffer(jsonBody))
	if err != nil {
		return ServiceUserInfoResponse{}, err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ServiceKey", config.AppConfig.UserHubServiceKey)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return ServiceUserInfoResponse{}, err
	}
	defer resp.Body.Close()

	// Parse response
	var response ServiceUserInfoResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return ServiceUserInfoResponse{}, err
	}

	return response, nil
}

func GetUserInfo(empID string) (models.AuthenUserEmp, error) {
	client := &http.Client{}

	// Create request body
	reqBody := ServiceUserInfoRequest{
		ServiceCode: "vms",
		EmpID:       empID,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return models.AuthenUserEmp{}, err
	}

	// Create request
	req, err := http.NewRequest("POST", config.AppConfig.UserHubEndPoint+"/service/get-user-info", bytes.NewBuffer(jsonBody))
	if err != nil {
		return models.AuthenUserEmp{}, err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ServiceKey", config.AppConfig.UserHubServiceKey)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return models.AuthenUserEmp{}, err
	}
	defer resp.Body.Close()
	fmt.Println(resp.Body)
	// Parse response
	var response models.AuthenUserEmp
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return models.AuthenUserEmp{}, err
	}
	response.IsEmployee = true
	return response, nil
}

func GetUserList(request ServiceListUserRequest) ([]models.MasUserEmp, error) {
	client := &http.Client{}

	// Create request body
	reqBody := ServiceListUserRequest{
		ServiceCode:   "vms",
		Search:        request.Search,
		UpperDeptSap:  request.UpperDeptSap,
		BureauDeptSap: request.BureauDeptSap,
		BusinessArea:  request.BusinessArea,
		LevelCodes:    request.LevelCodes,
		Role:          request.Role,
		Limit:         request.Limit,
	}
	fmt.Println("userhub list user", reqBody)
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return []models.MasUserEmp{}, err
	}

	// Create request
	req, err := http.NewRequest("POST", config.AppConfig.UserHubEndPoint+"/service/list-user", bytes.NewBuffer(jsonBody))
	if err != nil {
		return []models.MasUserEmp{}, err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ServiceKey", config.AppConfig.UserHubServiceKey)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return []models.MasUserEmp{}, err
	}
	defer resp.Body.Close()

	// Parse response
	var response []models.MasUserEmp
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return []models.MasUserEmp{}, err
	}
	for i := range response {
		response[i].IsEmployee = true
	}
	return response, nil
}

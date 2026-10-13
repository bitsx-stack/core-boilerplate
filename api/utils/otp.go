package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	BaseURL = "https://portal.textify.africa/api"
	APIKey  = "7be5005f080c22955a50335dd0931d6152e33cd8669fc5c06eaa24aae02d358fa1c3febd39094655a4c0771444f3e5a7"
)

type SendOTPRequest struct {
	PhoneNumber string `json:"phone_number"`
	BrandName   string `json:"brand_name"`
	SenderName  string `json:"sender_name"`
}

type VerifyOTPRequest struct {
	PhoneNumber string `json:"phone_number"`
	OTP         string `json:"otp"`
}

// =======================
// SEND OTP
// =======================
func SendOTP(phone string) error {
	url := BaseURL + "/otp/create"

	payload := SendOTPRequest{
		PhoneNumber: phone,
		SenderName:  "BITS TEXTS",
		BrandName:   "BITS INFORMATICS LTD",
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to send otp: status %d", resp.StatusCode)
	}

	return nil
}

// =======================
// VERIFY OTP
// =======================
func VerifyOTP(phone string, otp string) error {
	url := BaseURL + "/otp/verify"

	payload := VerifyOTPRequest{
		PhoneNumber: phone,
		OTP:         otp,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle responses properly based on docs
	switch resp.StatusCode {
	case 200:
		return nil
	case 400:
		return fmt.Errorf("invalid otp")
	case 404:
		return fmt.Errorf("no pending otp")
	case 401:
		return fmt.Errorf("otp expired")
	case 429:
		return fmt.Errorf("max attempts reached")
	default:
		return fmt.Errorf("unexpected error: %d", resp.StatusCode)
	}
}

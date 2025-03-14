package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	loginPath        = "/auth/login"
	availabilityPath = "/me/appointment-availabilities"
)

type loginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type HLAUser struct {
	Token      string `json:"token"`
	CustomerID string `json:"customer_id"`
}

type HLA struct {
	baseURL        string
	provinceID     int
	healthCentreID int
	agreementID    int
	formatID       int
	// api returns the first 3 dates from this date
	initialDate string
}

func NewHLA() *HLA {
	baseURL := os.Getenv("HLA_BASE_URL")
	if baseURL == "" {
		log.Fatalf("HLA_BASE_URL is required")
	}
	provinceID := mustLoadIntEnvVar("HLA_PROVINCE_ID")
	healthCentreID := mustLoadIntEnvVar("HLA_HEALTH_CENTRE_ID")
	agreementID := mustLoadIntEnvVar("HLA_AGREEMENT_ID")
	formatID := mustLoadIntEnvVar("HLA_FORMAT_ID")

	hla := &HLA{
		baseURL:        baseURL,
		provinceID:     provinceID,
		healthCentreID: healthCentreID,
		agreementID:    agreementID,
		formatID:       formatID,
	}

	// initialDate is optional, defaults to today
	if initialDate := os.Getenv("HLA_INITIAL_DATE"); initialDate != "" {
		hla.initialDate = initialDate
	} else {
		year, month, day := time.Now().Date()
		hla.initialDate = fmt.Sprintf("%d/%02d/%02d", year, month, day)
	}

	return hla
}

func mustLoadCredentials() (string, string) {
	username := os.Getenv("HLA_USERNAME")
	if username == "" {
		log.Fatalf("HLA_USERNAME is required")
	}
	password := os.Getenv("HLA_PASSWORD")
	if password == "" {
		log.Fatalf("HLA_PASSWORD is required")
	}
	return username, password
}

// Login logs in to the HLA service and returns a token and customer ID
func (h *HLA) Login() (*HLAUser, error) {
	loginURL := h.baseURL + loginPath

	username, password := mustLoadCredentials()
	loginCredentials := loginPayload{
		Username: username,
		Password: password,
	}

	payload, err := json.Marshal(loginCredentials)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling credentials: %v", err)
	}

	req, err := http.NewRequest("POST", loginURL, io.NopCloser(strings.NewReader(string(payload))))
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Login failed: %v: %s", resp.Status, string(body))
	}

	var user *HLAUser
	if err = json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("Error unmarshalling login response: %v", err)
	}

	return user, nil
}

type availabilityParams struct {
	FormatIDs      int    `json:"format_ids"`
	SpecialtyID    int    `json:"specialty_id"`
	InitialDate    string `json:"initial_date"`
	InitialTime    string `json:"initial_time"`
	EndTime        string `json:"end_time"`
	AgreementID    int    `json:"agreement_id"`
	ProvinceID     int    `json:"province_id"`
	HealthCentreID int    `json:"health_centre_id"`
}

func (p availabilityParams) toQueryParams() string {
	qparams := map[string]string{
		"format_ids":       fmt.Sprintf("%d", p.FormatIDs),
		"specialty_id":     fmt.Sprintf("%d", p.SpecialtyID),
		"initial_date":     p.InitialDate,
		"initial_time":     p.InitialTime,
		"end_time":         p.EndTime,
		"agreement_id":     fmt.Sprintf("%d", p.AgreementID),
		"province_id":      fmt.Sprintf("%d", p.ProvinceID),
		"health_centre_id": fmt.Sprintf("%d", p.HealthCentreID),
	}
	qparamsStr := ""
	for k, v := range qparams {
		qparamsStr += fmt.Sprintf("%s=%s&", k, v)
	}

	return qparamsStr
}

func (h *HLA) defaultParams() availabilityParams {
	return availabilityParams{
		FormatIDs:      h.formatID,
		AgreementID:    h.agreementID,
		InitialDate:    h.initialDate,
		InitialTime:    "07:00",
		EndTime:        "21:00",
		ProvinceID:     h.provinceID,
		HealthCentreID: h.healthCentreID,
	}
}

type Availability struct {
	AvailabilityID   string `json:"availability_id"`
	DateTime         string `json:"date_time"`
	FormatName       string `json:"format_name"`
	DoctorName       string `json:"doctor_full_name"`
	LocationName     string `json:"location_name"`
	ConsultationName string `json:"consultation_name"`
}

// AvailabilityCheck checks if the HLA service is available
func (h *HLA) AvailabilityCheck(token string, specialtyID int) ([]Availability, error) {
	params := h.defaultParams()
	params.SpecialtyID = specialtyID

	availabilityURL := h.baseURL + availabilityPath + "?" + params.toQueryParams()

	req, err := http.NewRequest("GET", availabilityURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("language", "es")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Availability check failed: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response: %v", err)
	}

	var availability []Availability
	if err = json.Unmarshal(body, &availability); err != nil {
		return nil, fmt.Errorf("Error unmarshalling availability response: %v", err)
	}

	return availability, nil
}

func mustLoadIntEnvVar(name string) int {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("%s is required", name)
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Error parsing %s: %v", name, err)
	}
	return intValue
}

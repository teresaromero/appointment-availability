package hlaservice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	availabilityPath = "/me/appointment-availabilities"
)

// availabilityCheckRequest checks if the HLA service is available
func (h *HLA) availabilityCheckRequest(token string, healthcenterID, specialtyID int) ([]Availability, error) {
	params := h.defaultParams()
	params.SpecialtyID = specialtyID
	params.HealthCentreID = healthcenterID

	availabilityURL := h.baseURL + availabilityPath + "?" + params.toQueryParams()

	req, err := http.NewRequest("GET", availabilityURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("language", "es")

	resp, err := h.client.Do(req)
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

type Availability struct {
	AvailabilityID   string `json:"availability_id"`
	DateTime         string `json:"date_time"`
	FormatName       string `json:"format_name"`
	DoctorName       string `json:"doctor_full_name"`
	LocationName     string `json:"location_name"`
	ConsultationName string `json:"consultation_name"`
}

type availabilityParams struct {
	FormatIDs      int    `json:"format_ids"`
	SpecialtyID    int    `json:"specialty_id"`
	InitialDate    string `json:"initial_date"`
	InitialTime    string `json:"initial_time"`
	EndTime        string `json:"end_time"`
	AgreementID    int    `json:"agreement_id"`
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
		"health_centre_id": fmt.Sprintf("%d", p.HealthCentreID),
	}
	qparamsStr := ""
	for k, v := range qparams {
		qparamsStr += fmt.Sprintf("%s=%s&", k, v)
	}

	return qparamsStr
}

func (h *HLA) defaultParams() availabilityParams {
	year, month, day := time.Now().Date()
	return availabilityParams{
		FormatIDs:   h.formatID,
		AgreementID: h.agreementID,
		InitialDate: fmt.Sprintf("%d/%02d/%02d", year, month, day),
		InitialTime: "07:00",
		EndTime:     "21:00",
	}
}

package hlaservice

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type sendNotificationFn func(ctx context.Context, msg string)

type HLA struct {
	client      *http.Client
	username    string
	password    string
	baseURL     string
	agreementID int
	formatID    int
	notifierFn  sendNotificationFn
}

func New(client *http.Client, baseURL, username, password string, notifierFn sendNotificationFn) *HLA {
	return &HLA{
		client:   client,
		baseURL:  baseURL,
		username: username,
		password: password,
		// agreementID represents the type of insurance.
		agreementID: 90002,
		// HLAFormatID represents the type of appointment. 1-Presential
		formatID:   1,
		notifierFn: notifierFn,
	}
}

func (h *HLA) Run(ctx context.Context, healthcenterIDList, specialtyIDList []int) error {
	if h.baseURL == "" {
		log.Default().Println("HLA: baseURL is empty")
		return nil
	}
	user, err := h.loginRequest(h.username, h.password)
	if err != nil {
		return fmt.Errorf("Error logging in: %v", err)
	}

	for _, healthcenter := range healthcenterIDList {
		for _, specialtyID := range specialtyIDList {
			if err := h.runJob(ctx, user.Token, healthcenter, specialtyID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *HLA) runJob(ctx context.Context, token string, healthcenter, specialtyID int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	data, err := h.availabilityCheckRequest(token, healthcenter, specialtyID)
	if err != nil {
		return fmt.Errorf("Error checking availability: %v", err)
	}

	msgHLA := "HLA: No appointment available for specialty ID: " + fmt.Sprintf("%d", specialtyID)
	if len(data) > 0 {
		msgHLA = "HLA: Appointment available for specialty ID: " + fmt.Sprintf("%d", specialtyID)
		for _, a := range data {
			msgHLA += "\n" + a.DateTime + " " + a.FormatName + " " + a.DoctorName + " " + a.LocationName + " " + a.ConsultationName
		}
	}
	h.notifierFn(ctx, msgHLA)

	return nil
}

package main

import (
	"appointment-availability/internal/bot"
	"appointment-availability/internal/services"
	hla "appointment-availability/internal/services"
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tgBot := bot.New()
	defer tgBot.Close(ctx)

	// HLA services
	hlaService := hla.NewHLA()
	user, err := hlaService.Login()
	if err != nil {
		log.Fatalf("Error logging in: %v", err)
	}

	specialtyIDList := os.Getenv("HLA_SPECIALTY_ID_LIST")
	// Check availability for each specialty
	for s := range strings.SplitSeq(specialtyIDList, ",") {
		// wait for each request to avoid rate limiting
		time.Sleep(5 * time.Second)

		specialtyId, err := strconv.Atoi(s)
		if err != nil {
			log.Fatalf("Error converting specialty ID: %v", err)
		}

		avail, err := hlaService.AvailabilityCheck(user.Token, specialtyId)
		if err != nil {
			log.Fatalf("Error checking availability: %v", err)
		}

		msgHLA := "HLA: No appointment available for specialty ID: " + s
		if len(avail) > 0 {
			msgHLA = "HLA: Appointment available for specialty ID: " + s
			for _, a := range avail {
				msgHLA += "\n" + a.DateTime + " " + a.FormatName + " " + a.DoctorName + " " + a.LocationName + " " + a.ConsultationName
			}
		}

		// Send message to telegram
		tgBot.SendNotification(ctx, msgHLA)
	}

	hcURL := os.Getenv("HC_URL")
	if hcURL != "" {
		// HC services
		msgHC, err := services.GetHCAvailablePeople(ctx, hcURL)
		if err != nil {
			log.Fatalf("Error getting content: %v", err)
		}
		tgBot.SendNotification(ctx, msgHC)
	}
}

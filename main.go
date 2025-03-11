package main

import (
	"appointment-availability/bot"
	hla "appointment-availability/services"
	"context"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tgBot := bot.New()

	hlaService := hla.NewHLA()
	user, err := hlaService.Login()
	if err != nil {
		log.Fatalf("Error logging in: %v", err)
	}

	avail, err := hlaService.AvailabilityCheck(user.Token)
	if err != nil {
		log.Fatalf("Error checking availability: %v", err)
	}

	msg := "No appointment available"
	if avail != nil {
		msg = "Appointment available"
		for _, a := range avail {
			msg += "\n" + a.DateTime + " " + a.FormatName + " " + a.DoctorName + " " + a.LocationName + " " + a.ConsultationName
		}
	}

	// Send message to telegram
	tgBot.SendNotification(ctx, msg)
}

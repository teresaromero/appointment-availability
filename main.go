package main

import (
	hla "appointment-availability/services"
	"log"
)

func main() {
	hlaService := hla.NewHLA()
	user, err := hlaService.Login()
	if err != nil {
		log.Fatalf("Error logging in: %v", err)
	}

	ok, err := hlaService.AvailabilityCheck(user.Token)
	if err != nil {
		log.Fatalf("Error checking availability: %v", err)
	}

	if ok {
		log.Printf("Appointment available")
	} else {
		log.Printf("No appointment available")
	}
}

package services

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

func GetHCAvailablePeople(ctx context.Context, url string) (string, error) {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	// Variables para almacenar el valor inicial y el valor seleccionado
	var people []string

	if err := chromedp.Run(ctx,
		// Navegar a la URL
		chromedp.Navigate(url),
		// Esperar que el elemento select esté visible
		chromedp.WaitVisible(`#centro`, chromedp.ByID),
		chromedp.SetValue(`#centro`, "2", chromedp.ByID),
		chromedp.Sleep(3*time.Second),

		chromedp.WaitVisible(`#aseguradora`, chromedp.ByID),
		chromedp.SetValue(`#aseguradora`, "3", chromedp.ByID),
		chromedp.Sleep(3*time.Second),

		chromedp.WaitVisible(`#colectivo`, chromedp.ByID),
		chromedp.SetValue(`#colectivo`, "4", chromedp.ByID),
		chromedp.Sleep(3*time.Second),

		chromedp.WaitVisible(`#especialidad`, chromedp.ByID),
		chromedp.SetValue(`#especialidad`, "16", chromedp.ByID),
		chromedp.Sleep(3*time.Second),

		// concepto 61 consulta
		chromedp.WaitVisible(`#concepto`, chromedp.ByID),
		chromedp.SetValue(`#concepto`, "61", chromedp.ByID),
		chromedp.Sleep(3*time.Second),

		chromedp.WaitVisible(`#profesional`, chromedp.ByID),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('#profesional option')).map(option => option.innerText)`, &people),
	); err != nil {
		return "", err
	}

	msg := "HC: No appointment available for specialty ID: 16"
	if len(people) > 1 {
		msg = "HC: Appointment available for specialty ID: 16"
		for _, p := range people {
			msg += "\n" + p
		}
	}

	return msg, nil
}

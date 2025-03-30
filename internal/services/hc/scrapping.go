package hcservice

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// scrap scrapes the HC website for available doctors based on the provided parameters.
// It navigates to the URL, selects the appropriate options for health center,
// insurance, group, and specialty, and retrieves the list of available doctors.
func scrap(ctx context.Context, url, centerID, insuranceID, groupID, specialtyID string) (string, error) {
	ctx, cancelWithTimeout := context.WithTimeout(ctx, 60*time.Second)
	defer cancelWithTimeout()

	ctx, cancelChromedp := chromedp.NewContext(ctx)
	defer cancelChromedp()

	if err := selectOptions(ctx, url, centerID, insuranceID, groupID, specialtyID); err != nil {
		return "", fmt.Errorf("HC: error selecting form options: %w", err)
	}

	list, err := availableDoctors(ctx)
	if err != nil {
		return "", fmt.Errorf("HC: error getting available doctors: %w", err)
	}

	if len(list) == 0 {
		return fmt.Sprintf("🔴 HC: No doctors available for center %s - specialty %s", centerID, specialtyID), nil
	}

	if ok, err := isDatePickerEnabled(ctx); err != nil {
		return "", fmt.Errorf("HC: error checking date picker: %w", err)
	} else if !ok {
		return fmt.Sprintf("🔴 HC: no dates available for center %s - specialty %s", centerID, specialtyID), nil
	}

	return fmt.Sprintf("🎉 HC: doctors available for center %s - specialty %s: %s", centerID, specialtyID, strings.Join(list, ", ")), nil
}

// selectOptions get a list of doctors available for a given health center and specialty
// if only one doctor, fallback to the input field
func availableDoctors(ctx context.Context) ([]string, error) {
	var doctors []string
	list, err := getProfessionalsFromSelect(ctx)
	if err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		single, err := getSingleProfessionalFromInput(ctx)
		if err != nil {
			return nil, err
		}
		if single == "" {
			return nil, nil
		}
		doctors = append(doctors, single)
	}
	// The list can contain the doctor's name and surname, separated by a comma.
	// Or a "not available" message.
	if len(list) > 0 {
		for _, d := range list {
			if strings.Contains(d, ",") {
				doctors = append(doctors, d)
			}
		}
	}
	return doctors, nil
}

// selectFormOptions selects the form options for the health center, insurance, group, and specialty
// by setting the values of the corresponding input elements.
// It waits for the elements to be visible before setting the values.
func selectOptions(ctx context.Context, url, centerID, insuranceID, groupID, specialtyID string) error {
	return chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#centro`, chromedp.ByID),
		chromedp.SetValue(`#centro`, centerID, chromedp.ByID),
		chromedp.Sleep(5*time.Second),
		chromedp.WaitVisible(`#aseguradora`, chromedp.ByID),
		chromedp.SetValue(`#aseguradora`, insuranceID, chromedp.ByID),
		chromedp.Sleep(5*time.Second),
		chromedp.WaitVisible(`#colectivo`, chromedp.ByID),
		chromedp.SetValue(`#colectivo`, groupID, chromedp.ByID),
		chromedp.Sleep(5*time.Second),
		chromedp.WaitVisible(`#especialidad`, chromedp.ByID),
		chromedp.SetValue(`#especialidad`, specialtyID, chromedp.ByID),
		chromedp.Sleep(5*time.Second),
		chromedp.WaitVisible(`#concepto`, chromedp.ByID),
		chromedp.SetValue(`#concepto`, "61", chromedp.ByID),
		chromedp.Sleep(5*time.Second),
	)
}

// isDatePickerEnabled checks if the date picker is enabled or disabled
// by checking the readOnly attribute of the input element with id "dia".
// If the attribute is set to true, the date picker is disabled.
func isDatePickerEnabled(ctx context.Context) (bool, error) {
	var datePickerReadonly bool
	if err := chromedp.Run(ctx,
		chromedp.Evaluate(`document.getElementById("dia").readOnly`, &datePickerReadonly),
	); err != nil {
		return false, err
	}
	return !datePickerReadonly, nil
}

// getSingleProfessionalFromSelect retrieves the professional from the input element
// with id "profesionaloTX" and returns it as a string.
// If the element is not found or the value is empty, it returns an error.
func getSingleProfessionalFromInput(ctx context.Context) (string, error) {
	var professionalInput string
	if err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			lctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			// Run chromedp.WaitVisible with a child timeout context.
			// If it returns an non-nil error, the tasks behind won't run.
			return chromedp.WaitVisible(`#profesionaloTX`, chromedp.ByID).Do(lctx)
		}),
		chromedp.Evaluate(`document.getElementById('profesionaloTX').value`, &professionalInput),
	); err != nil {
		return "", err
	}
	return professionalInput, nil
}

// getProfessionalsFromSelect retrieves the professionals from the select element
// with id "profesional" and returns them as a slice of strings.
// If the element is not found or the value is empty, it returns an error.
func getProfessionalsFromSelect(ctx context.Context) ([]string, error) {
	var professionals []string
	if err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			lctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			// Run chromedp.WaitVisible with a child timeout context.
			// If it returns an non-nil error, the tasks behind won't run.
			return chromedp.WaitVisible(`#profesional`, chromedp.ByID).Do(lctx)
		}),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('#profesional option')).map(option => option.innerText)`, &professionals),
	); err != nil {
		return nil, err
	}
	return professionals, nil
}

package hcservice

import (
	"context"
	"log"

	"golang.org/x/sync/errgroup"
)

type sendNotificationFn func(ctx context.Context, msg string)

type HC struct {
	notifierFn  sendNotificationFn
	baseURL     string
	insuranceID string
	groupID     string
}

func New(baseURL string, notifierFn sendNotificationFn) *HC {
	return &HC{
		baseURL:     baseURL,
		notifierFn:  notifierFn,
		insuranceID: "3",
		groupID:     "4",
	}
}

func (h *HC) Run(ctx context.Context, healthcenterList, specialtyList []string) error {
	if h.baseURL == "" {
		log.Default().Println("HC: baseURL is empty")
		return nil
	}

	errgroup, ctx := errgroup.WithContext(ctx)
	for _, centerID := range healthcenterList {
		for _, specialtyID := range specialtyList {
			errgroup.Go(func() error {
				msg, err := scrap(ctx, h.baseURL, centerID, h.insuranceID, h.groupID, specialtyID)
				if err != nil {
					return err
				}
				h.notifierFn(ctx, msg)
				return nil
			})
		}
	}
	return errgroup.Wait()
}

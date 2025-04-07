package service_test

import (
	"sweng-task/internal/model"
	"sweng-task/internal/service"
	"testing"

	"go.uber.org/zap"
)

func TestAdService_GetWinningAds(t *testing.T) {
	logger, _ := zap.NewProduction()
	log := logger.Sugar()

	lineSvc := service.NewLineItemService(log)
	adSvc := service.NewAdService(lineSvc, log)

	testCases := []struct {
		name            string
		lineItems       []model.LineItemCreate
		placement       string
		category        string
		keyword         string
		limit           int
		expectedAdCount int
	}{
		{
			name: "Exact match on category and keyword, limit=1",
			lineItems: []model.LineItemCreate{
				{
					Name:         "Test Ad 1",
					AdvertiserID: "adv1",
					Bid:          2.0,
					Budget:       100,
					Placement:    "homepage_top",
					Categories:   []string{"electronics"},
					Keywords:     []string{"discount"},
				},
				{
					Name:         "Test Ad 2",
					AdvertiserID: "adv2",
					Bid:          1.5,
					Budget:       50,
					Placement:    "homepage_top",
					Categories:   []string{"sale"},
					Keywords:     []string{"summer"},
				},
			},
			placement:       "homepage_top",
			category:        "electronics",
			keyword:         "discount",
			limit:           1,
			expectedAdCount: 1,
		},
		{
			name: "No match on category/keyword, limit=2",
			lineItems: []model.LineItemCreate{
				{
					Name:         "Ad A",
					AdvertiserID: "advA",
					Bid:          3.0,
					Budget:       100,
					Placement:    "homepage_top",
					Categories:   []string{"sports"},
					Keywords:     []string{"shoes"},
				},
				{
					Name:         "Ad B",
					AdvertiserID: "advB",
					Bid:          2.5,
					Budget:       80,
					Placement:    "homepage_top",
					Categories:   []string{"fashion"},
					Keywords:     []string{"jeans"},
				},
			},
			placement:       "homepage_top",
			category:        "electronics",
			keyword:         "discount",
			limit:           2,
			expectedAdCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lineSvc = service.NewLineItemService(log)
			adSvc = service.NewAdService(lineSvc, log)

			for _, li := range tc.lineItems {
				_, err := lineSvc.Create(li)
				if err != nil {
					t.Fatalf("failed to create line item: %v", err)
				}
			}

			ads, err := adSvc.GetWinningAds(tc.placement, tc.category, tc.keyword, tc.limit)
			if err != nil {
				t.Fatalf("GetWinningAds returned an error: %v", err)
			}
			if len(ads) != tc.expectedAdCount {
				t.Fatalf("expected %d ads, got %d", tc.expectedAdCount, len(ads))
			}
		})
	}
}

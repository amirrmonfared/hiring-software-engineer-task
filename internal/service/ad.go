package service

import (
	"fmt"
	"sort"
	"sweng-task/internal/model"

	"go.uber.org/zap"
)

type AdService struct {
	lineItemService *LineItemService
	log             *zap.SugaredLogger
}

func NewAdService(lineItemService *LineItemService, log *zap.SugaredLogger) *AdService {
	return &AdService{
		lineItemService: lineItemService,
		log:             log,
	}
}

type scoredItem struct {
	item  *model.LineItem
	score float64
}

func (s *AdService) GetWinningAds(placement, category, keyword string, limit int) ([]model.Ad, error) {
	items, err := s.lineItemService.FindMatchingLineItems(placement, category, keyword)
	if err != nil {
		return nil, err
	}

	scored := make([]scoredItem, 0, len(items))
	for _, li := range items {
		score := s.calculateScore(li, category, keyword)
		scored = append(scored, scoredItem{
			item:  li,
			score: score,
		})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	if limit > len(scored) {
		limit = len(scored)
	}
	winners := make([]model.Ad, 0, limit)
	for i := 0; i < limit; i++ {
		li := scored[i].item
		winners = append(winners, model.Ad{
			ID:           li.ID,
			Name:         li.Name,
			AdvertiserID: li.AdvertiserID,
			Bid:          li.Bid,
			Placement:    li.Placement,
			ServeURL:     fmt.Sprintf("/ad/serve/%s", li.ID),
		})
	}

	return winners, nil
}

func (s *AdService) calculateScore(li *model.LineItem, category, keyword string) float64 {
	score := li.Bid

	if category != "" {
		for _, c := range li.Categories {
			if c == category {
				score += 1.0
				break
			}
		}
	}
	if keyword != "" {
		for _, k := range li.Keywords {
			if k == keyword {
				score += 1.0
				break
			}
		}
	}
	return score
}

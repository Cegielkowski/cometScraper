package entity

import (
	"cometScraper/tools/scraper/pkg/applicant"
	"time"
)

type CometScraper struct {
	Uuid      string              `json:"uuid"`
	Status    string              `json:"status"`
	Applicant applicant.Candidate `json:"applicant"`
	TimeTaken string              `json:"time_taken"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

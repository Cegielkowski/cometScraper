package crawler

import (
	"cometScraper/tools/scraper/pkg/applicant"
)

type Response struct {
	Uuid      string              `json:"uuid"`
	Status    string              `json:"status"`
	Applicant applicant.Candidate `json:"applicant"`
	TimeTaken string              `json:"time_taken"`
}

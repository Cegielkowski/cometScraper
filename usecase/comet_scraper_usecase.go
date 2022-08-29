package usecase

import (
	"cometScraper/tools/scraper/pkg/crawler"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"cometScraper/entity"
	"cometScraper/repository/pgsql"
	"cometScraper/repository/redis"
	"cometScraper/transport/request"
	"cometScraper/utils"
)

// CometScraperUsecase represent the craper's usecase contract
type CometScraperUsecase interface {
	StartProcess(ctx context.Context, request *request.CreateCometScraperReq) (string, error)
	GetByID(ctx context.Context, id string) (entity.CometScraper, error)
	Fetch(ctx context.Context) ([]entity.CometScraper, error)
	Update(ctx context.Context, cometScraper *entity.CometScraper) error
	UpsertStatus(id string, status string) error
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, cometScraper entity.CometScraper) error
}

type cometScraperUsecase struct {
	cometScraperRepo pgsql.CometScraperRepository
	redisRepo        redis.RedisRepository
	cometCrawler     crawler.CometScraper
}

// NewCometScraperUsecase will create new an cometScraperUsecase object representation of CometScraperUsecase interface
func NewCometScraperUsecase(cometScraperRepo pgsql.CometScraperRepository, redisRepo redis.RedisRepository, cometCrawler crawler.CometScraper) CometScraperUsecase {
	return &cometScraperUsecase{
		cometScraperRepo: cometScraperRepo,
		redisRepo:        redisRepo,
		cometCrawler:     cometCrawler,
	}
}

func (c *cometScraperUsecase) StartProcess(ctx context.Context, request *request.CreateCometScraperReq) (string, error) {
	credentials := crawler.Credentials{
		Email: request.Email,
		Pass:  request.Password,
	}

	cr := make(chan crawler.Response)
	done := make(chan struct{})
	processUuid := c.cometCrawler.GetUuid()
	err := c.Create(ctx, entity.CometScraper{Uuid: processUuid, Status: entity.Start, TimeTaken: "0"})
	if err != nil {
		return "", utils.NewInternalServerError("Some internal error happened, please contact support")
	}

	go c.cometCrawler.StartCrawling(processUuid, credentials, cr, done)
	go c.HandleAsync(processUuid, cr, done)

	return processUuid, nil
}

func (c *cometScraperUsecase) HandleAsync(processUuid string, cr chan crawler.Response, done chan struct{}) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		select {
		case <-time.After(80 * time.Second):
			log.Println("Time Out")
			_ = c.UpsertStatus(processUuid, entity.TimeOut)
			return
		case <-done:
			log.Println("Finished")
			return
		case response := <-cr:
			err := c.Update(ctx, &entity.CometScraper{
				Uuid:      response.Uuid,
				Status:    response.Status,
				Applicant: response.Applicant,
				TimeTaken: response.TimeTaken,
			})
			if err != nil {
				log.Println(err)
				return
			}
		case <-ctx.Done():
			log.Println("Done")
			_ = c.UpsertStatus(processUuid, entity.Fail)
			return
		}
	}
}

func (c *cometScraperUsecase) Update(ctx context.Context, cometScraper *entity.CometScraper) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	comet, err := c.cometScraperRepo.GetByID(ctx, cometScraper.Uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.NewNotFoundError("scraper not found")
			return
		}
		return
	}

	if comet.Applicant.Name != "" {
		cometScraper.Applicant.Name = utils.Encrypt(cometScraper.Uuid, cometScraper.Applicant.Name)
	}

	comet.Status = cometScraper.Status
	comet.TimeTaken = cometScraper.TimeTaken
	comet.Applicant = cometScraper.Applicant
	cometScraper.UpdatedAt = time.Now()

	err = c.cometScraperRepo.Update(ctx, &comet)
	_ = c.redisRepo.Delete("cometScrapers")
	return
}

func (c *cometScraperUsecase) UpsertStatus(id string, status string) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	comet, err := c.cometScraperRepo.GetByID(ctx, id)
	comet.Status = status
	if err != nil {
		if err == sql.ErrNoRows {
			comet.Uuid = id
			err = c.Create(ctx, comet)
			return
		}
		return
	}

	comet.UpdatedAt = time.Now()

	err = c.cometScraperRepo.UpdateStatus(ctx, &comet)
	_ = c.redisRepo.Delete("cometScrapers")

	return
}

func (c *cometScraperUsecase) Create(ctx context.Context, cometScraper entity.CometScraper) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err = c.cometScraperRepo.Create(ctx, &entity.CometScraper{
		Uuid:      cometScraper.Uuid,
		Status:    cometScraper.Status,
		Applicant: cometScraper.Applicant,
		TimeTaken: cometScraper.TimeTaken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	_ = c.redisRepo.Delete("cometScrapers")

	return
}

func (c *cometScraperUsecase) GetByID(ctx context.Context, id string) (cometScraper entity.CometScraper, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cometScraper, err = c.cometScraperRepo.GetByID(ctx, id)
	if err != nil && err == sql.ErrNoRows {
		err = utils.NewNotFoundError("process not found")
		return
	}

	if cometScraper.Applicant.Name != "" {
		cometScraper.Applicant.Name = utils.Decrypt(cometScraper.Uuid, cometScraper.Applicant.Name)
	}
	return
}

func (c *cometScraperUsecase) Fetch(ctx context.Context) (cometScrapers []entity.CometScraper, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cometScrapersCached, _ := c.redisRepo.Get("cometScrapers")
	if err = json.Unmarshal([]byte(cometScrapersCached), &cometScrapers); err == nil {
		return
	}

	cometScrapers, err = c.cometScraperRepo.Fetch(ctx)
	if err != nil {
		return
	}

	cometScrapersString, _ := json.Marshal(&cometScrapers)
	_ = c.redisRepo.Set("cometScrapers", cometScrapersString, 30*time.Second)

	return
}

func (c *cometScraperUsecase) Delete(ctx context.Context, id string) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	_, err = c.cometScraperRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.NewNotFoundError("Process not found")
			return
		}
		return
	}

	_ = c.redisRepo.Delete("cometScrapers")

	return
}

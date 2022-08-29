package crawler

import (
	"cometScraper/entity"
	"cometScraper/tools/scraper/pkg/applicant"
	"cometScraper/tools/scraper/pkg/element"
	"context"
	"errors"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	uuid "github.com/satori/go.uuid"

	"log"
	"strconv"
	"strings"
	"time"
)

func generateActions(child int, elementAndValue []element.ElementAndValue) []chromedp.Action {
	childString := strconv.Itoa(child)
	var actions []chromedp.Action
	for _, value := range elementAndValue {
		actions = append(actions, chromedp.Text(strings.ReplaceAll(value.Element, "HIREME", childString), value.FutureValue(child-1), chromedp.ByQuery))
	}

	return actions
}

func GetActionsBaseInfo(resumeUrl string, elements element.Elements, ap *applicant.Candidate, ok *bool, nodesSkill, nodesExperience *[]*cdp.Node) []chromedp.Action {
	resumeSection := elements.GetResumeSection()
	return []chromedp.Action{
		chromedp.Navigate(resumeUrl),
		chromedp.Sleep(4 * time.Second),
		chromedp.Text(resumeSection.Name, &ap.Name, chromedp.ByQuery),
		chromedp.Text(resumeSection.Description, &ap.Description, chromedp.ByQuery),
		chromedp.Text(resumeSection.Role, &ap.Role, chromedp.ByQuery),
		chromedp.Text(resumeSection.TimeOfExperience, &ap.TimeOfExperience, chromedp.ByQuery),
		chromedp.AttributeValue(resumeSection.Image, `src`, &ap.ImageUrl, ok, chromedp.ByQuery),
		chromedp.Nodes(resumeSection.Skills, nodesSkill, chromedp.ByQueryAll),
		chromedp.Nodes(resumeSection.Experiences, nodesExperience, chromedp.ByQueryAll),
	}
}

func rangeAndGetActions(rangeSize int, eAndVal []element.ElementAndValue) []chromedp.Action {
	var actions []chromedp.Action

	for i := 1; i <= rangeSize; i++ {
		actions = append(actions, generateActions(i, eAndVal)...)
	}

	return actions
}

func GetActionsLogin(elements element.Elements, credentials Credentials, currentUrl *string) []chromedp.Action {
	return []chromedp.Action{
		chromedp.Navigate(elements.GetUrls().StartPage),
		chromedp.WaitVisible(elements.GetButtons().AcceptCookie, chromedp.ByQuery),
		chromedp.Sleep(3 * time.Second),
		chromedp.Click(elements.GetButtons().AcceptCookie, chromedp.ByQuery),
		chromedp.Sleep(2 * time.Second),
		chromedp.WaitNotPresent(elements.GetButtons().AcceptCookie, chromedp.ByQuery),
		chromedp.SendKeys(elements.GetInputs().Email, credentials.Email, chromedp.ByQuery),
		chromedp.SendKeys(elements.GetInputs().Password, credentials.Pass),
		chromedp.Sleep(2 * time.Second),
		chromedp.Click(elements.GetButtons().Login, chromedp.ByQuery),
		chromedp.Sleep(5 * time.Second),
		chromedp.Location(currentUrl),
	}
}

func GetActionsResume(elements element.Elements, resumeUrl *string, ok *bool) []chromedp.Action {
	return []chromedp.Action{
		chromedp.Navigate(elements.GetUrls().FreelanceProfile),
		chromedp.WaitVisible(elements.GetButtons().Resume, chromedp.ByQuery),
		chromedp.Sleep(2 * time.Second),
		chromedp.AttributeValue(elements.GetButtons().Resume, `href`, resumeUrl, ok, chromedp.ByQuery),
	}
}

func GetActionsToGetSkillAndExp(lenSkills, lenExperiences int, eAndValSkills, eAndValExperience []element.ElementAndValue, resumeUrl string) chromedp.Tasks {
	var actions chromedp.Tasks

	actions = append(actions, chromedp.Navigate(resumeUrl))
	actions = append(actions, chromedp.Sleep(4*time.Second))
	if lenSkills > 0 {
		actions = append(actions, rangeAndGetActions(lenSkills, eAndValSkills)...)
	}
	if lenExperiences > 0 {
		actions = append(actions, rangeAndGetActions(lenExperiences, eAndValExperience)...)
	}
	return actions
}

type Credentials struct {
	Email string
	Pass  string
}

type cometScraper struct {
	elements  element.Elements
	applicant applicant.Applicant
}

type CometScraper interface {
	StartCrawling(id string, credentials Credentials, cr chan Response, done chan struct{})
	GetUuid() string
}

func NewCometCrawler(elements element.Elements, applicant applicant.Applicant) CometScraper {
	return &cometScraper{
		elements:  elements,
		applicant: applicant,
	}
}

func (c *cometScraper) login(ctx context.Context, credentials Credentials) error {
	var currentUrl string

	err := chromedp.Run(ctx, GetActionsLogin(c.elements, credentials, &currentUrl)...)
	if err != nil {
		return err
	}

	if currentUrl != c.elements.GetUrls().FreelancerDashboard {
		return errors.New(entity.FailedCredentials)
	}

	return nil
}

func (c *cometScraper) getBaseInfo(ctx context.Context) (int, int, string, error) {
	var resumeUrl string
	var nodesSkill, nodesExperience []*cdp.Node
	var ok bool
	err := chromedp.Run(ctx, GetActionsResume(c.elements, &resumeUrl, &ok)...)
	if err != nil {
		log.Println(err)
		return 0, 0, "", err
	}

	err = chromedp.Run(ctx, GetActionsBaseInfo(resumeUrl, c.elements, c.applicant.Get(), &ok, &nodesSkill, &nodesExperience)...)
	if err != nil {
		log.Println(err)
		return 0, 0, "", err
	}

	lenSkills := len(nodesSkill)
	lenExperiences := len(nodesExperience)
	return lenSkills, lenExperiences, resumeUrl, nil
}

func (c *cometScraper) getSkillsAndExp(ctx context.Context, lenSkills, lenExperiences int, resumeUrl string) error {
	c.applicant.InitializeSkillAndExperience(lenSkills, lenExperiences)
	eAndValExperience := c.applicant.GenerateExperienceElementsAndValue(c.elements.GetExperienceElements())
	eAndValSkills := c.applicant.GenerateSkillsElementsAndValue(c.elements.GetSkillsElements())
	err := chromedp.Run(ctx, GetActionsToGetSkillAndExp(lenSkills, lenExperiences, eAndValSkills, eAndValExperience, resumeUrl))
	if err != nil {
		log.Println(err)
		return err
	}
	c.applicant.Clear()
	return nil
}

func (c *cometScraper) GetUuid() string {
	return uuid.NewV4().String()
}

func (c *cometScraper) StartCrawling(id string, credentials Credentials, cr chan Response, done chan struct{}) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)

	defer cancel()

	response := Response{
		Uuid:      id,
		Status:    entity.Start,
		Applicant: *c.applicant.Get(),
	}

	c.crawl(ctx, credentials, response, cr, done)
}

func (c *cometScraper) crawl(ctx context.Context, credentials Credentials, res Response, cr chan Response, done chan struct{}) {
	start := time.Now()
	err := c.login(ctx, credentials)
	if err != nil {
		res.Status = entity.FailedCredentials
		res.TimeTaken = time.Since(start).String()
		if err.Error() != entity.FailedCredentials {
			res.Status = entity.Fail
		}
		cr <- res
		close(done)
		return
	}

	res.Status = entity.Logged
	res.TimeTaken = time.Since(start).String()
	cr <- res

	lenSkills, lenExperiences, resumeUrl, err := c.getBaseInfo(ctx)
	if err != nil {
		res.Status = entity.Fail
		res.TimeTaken = time.Since(start).String()
		cr <- res
		close(done)
		return
	}

	res.Status = entity.Basic
	res.TimeTaken = time.Since(start).String()
	res.Applicant = *c.applicant.Get()
	cr <- res

	if lenSkills+lenExperiences > 0 {
		err = c.getSkillsAndExp(ctx, lenSkills, lenExperiences, resumeUrl)
		if err != nil {
			res.Status = entity.Fail
			res.TimeTaken = time.Since(start).String()
			cr <- res
			close(done)
			return
		}
	}

	res.Status = entity.Success
	res.Applicant = *c.applicant.Get()
	res.TimeTaken = time.Since(start).String()
	cr <- res
	close(done)
}

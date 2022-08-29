package applicant

import (
	"cometScraper/tools/scraper/pkg/element"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
)

func clearString(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "·", "")
	s = strings.TrimSpace(s)
	return s
}

type Skill struct {
	Name string `json:"name"`
	Time string `json:"time"`
}

type Job struct {
	Title       string `json:"title"`
	Skill       string `json:"skill"`
	Desc        string `json:"desc"`
	Period      string `json:"period"`
	PeriodCount string `json:"period_count"`
}

type Applicant interface {
	Get() *Candidate
	GetImageUrl() *string
	GetName() *string
	GetRole() *string
	GetTimeOfExperience() *string
	GetDescription() *string
	GetJobTitle(key int) *string
	GetJobSkill(key int) *string
	GetJobDesc(key int) *string
	GetJobPeriod(key int) *string
	GetJobPeriodCount(key int) *string
	GetSkillName(key int) *string
	GetSkillTime(key int) *string
	Clear()
	GenerateExperienceElementsAndValue(e element.ExperienceElements) []element.ElementAndValue
	GenerateSkillsElementsAndValue(s element.SkillsElements) []element.ElementAndValue
	InitializeSkillAndExperience(lenSkills, lenExperiences int)
}

type Candidate struct {
	ImageUrl         string  `json:"image_url"`
	Name             string  `json:"name"`
	Role             string  `json:"role"`
	Experience       []Job   `json:"experience"`
	Description      string  `json:"description"`
	Skill            []Skill `json:"skill"`
	TimeOfExperience string  `json:"time_of_experience"`
}

// Value Make the Candidate struct implement the driver.Valuer interface. This method simply returns the JSON-encoded representation of the struct.
func (c Candidate) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan Make the Candidate struct implement the sql.Scanner interface. This method simply decodes a JSON-encoded value into the struct fields.
func (c *Candidate) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &c)
}

func (c *Candidate) Get() *Candidate {
	return c
}

func NewApplicant() Applicant {
	return &Candidate{}
}

func (c *Candidate) InitializeSkillAndExperience(lenSkills, lenExperiences int) {
	c.Skill = make([]Skill, lenSkills)
	c.Experience = make([]Job, lenExperiences)
}

func (c *Candidate) GetImageUrl() *string {
	return &c.ImageUrl
}

func (c *Candidate) GetName() *string {
	return &c.Name
}

func (c *Candidate) GetRole() *string {
	return &c.Role
}

func (c *Candidate) GetDescription() *string {
	return &c.Description
}

func (c *Candidate) GetTimeOfExperience() *string {
	return &c.TimeOfExperience
}

func (c *Candidate) GetJobTitle(key int) *string {
	return &c.Experience[key].Title
}

func (c *Candidate) GetJobSkill(key int) *string {
	return &c.Experience[key].Skill
}

func (c *Candidate) GetJobDesc(key int) *string {
	return &c.Experience[key].Desc
}

func (c *Candidate) GetJobPeriod(key int) *string {
	return &c.Experience[key].Period
}

func (c *Candidate) GetJobPeriodCount(key int) *string {
	return &c.Experience[key].PeriodCount
}

func (c *Candidate) GetSkillName(key int) *string {
	return &c.Skill[key].Name
}
func (c *Candidate) GetSkillTime(key int) *string {
	return &c.Skill[key].Time
}

func (c *Candidate) Clear() {
	for key, value := range c.Skill {
		c.Skill[key].Time = clearString(value.Time)
	}
	for key, value := range c.Experience {
		c.Experience[key].Desc = clearString(value.Desc)
	}
}

func (c *Candidate) GenerateExperienceElementsAndValue(e element.ExperienceElements) []element.ElementAndValue {
	return []element.ElementAndValue{
		{
			Element:     e.Title,
			FutureValue: c.GetJobTitle,
		},
		{
			Element:     e.Skill,
			FutureValue: c.GetJobSkill,
		},
		{
			Element:     e.Desc,
			FutureValue: c.GetJobDesc,
		},
		{
			Element:     e.Period,
			FutureValue: c.GetJobPeriod,
		},
		{
			Element:     e.PeriodCount,
			FutureValue: c.GetJobPeriodCount,
		},
	}
}

func (c *Candidate) GenerateSkillsElementsAndValue(s element.SkillsElements) []element.ElementAndValue {
	return []element.ElementAndValue{
		{
			Element:     s.Name,
			FutureValue: c.GetSkillName,
		},
		{
			Element:     s.Time,
			FutureValue: c.GetSkillTime,
		},
	}
}

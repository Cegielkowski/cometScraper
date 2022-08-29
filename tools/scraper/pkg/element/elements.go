package element

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

type ExperienceElements struct {
	Title       string `json:"title"`
	Skill       string `json:"skill"`
	Desc        string `json:"desc"`
	Period      string `json:"period"`
	PeriodCount string `json:"periodCount"`
}

type SkillsElements struct {
	Name string `json:"name"`
	Time string `json:"time"`
}

type Urls struct {
	StartPage           string `json:"startPage"`
	FreelancerDashboard string `json:"freelancerDashboard"`
	FreelanceProfile    string `json:"freelanceProfile"`
}

type Inputs struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Buttons struct {
	Resume       string `json:"resume"`
	AcceptCookie string `json:"acceptCookie"`
	Login        string `json:"login"`
}

type ResumeSection struct {
	Image            string `json:"image"`
	Name             string `json:"name"`
	Role             string `json:"role"`
	TimeOfExperience string `json:"timeOfExperience"`
	Description      string `json:"description"`
	Skills           string `json:"skills"`
	Experiences      string `json:"experiences"`
}

type elements struct {
	Urls               Urls               `json:"urls"`
	Inputs             Inputs             `json:"inputs"`
	Buttons            Buttons            `json:"buttons"`
	ResumeSection      ResumeSection      `json:"resumeSection"`
	ExperienceElements ExperienceElements `json:"experienceElements"`
	SkillsElements     SkillsElements     `json:"skillsElements"`
}

func (e *elements) GetUrls() Urls {
	return e.Urls
}

func (e *elements) GetInputs() Inputs {
	return e.Inputs
}

func (e *elements) GetButtons() Buttons {
	return e.Buttons
}

func (e *elements) GetResumeSection() ResumeSection {
	return e.ResumeSection
}

func (e *elements) GetSkillsElements() SkillsElements {
	return e.SkillsElements
}

func (e *elements) GetExperienceElements() ExperienceElements {
	return e.ExperienceElements
}

type Elements interface {
	GetUrls() Urls
	GetInputs() Inputs
	GetButtons() Buttons
	GetResumeSection() ResumeSection
	GetSkillsElements() SkillsElements
	GetExperienceElements() ExperienceElements
}

type ElementAndValue struct {
	Element     string
	FutureValue func(int) *string
}

func NewElement(jsonFile io.Reader) (Elements, error) {
	byteResult, _ := ioutil.ReadAll(jsonFile)

	var element elements

	err := json.Unmarshal(byteResult, &element)
	if err != nil {
		return nil, errors.New("something wrong with JSON")
	}

	return &element, nil
}

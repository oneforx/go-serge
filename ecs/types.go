package ecs

import "fmt"

type FeedBackType string

var (
	FB_ERROR   FeedBackType = "FB_ERROR"
	FB_SUCCESS FeedBackType = "FB_SUCCESS"
)

type FeedBack struct {
	Host    string // Le nom du bloc ou plus précisément le nom de la fonction qui host les jobs
	Job     string // Le nom de la fonction qui a provoqué le feedback
	Label   string
	Type    FeedBackType
	Comment string
	Data    interface{}
}

func (fb *FeedBack) String() string {
	return fmt.Sprintf("[%v][%v][%v][%v]: %v", fb.Type, fb.Host, fb.Job, fb.Label, fb.Comment)
}

type Composition []string

func (source Composition) Equals(target Composition) bool {
	if len(source) != len(target) {
		return false
	}
	seen := make(map[string]bool)
	for _, s := range source {
		seen[s] = true
	}
	for _, t := range target {
		if !seen[t] {
			return false
		}
	}
	return true
}

package helper

import (
	"encoding/json"
	"strings"

	"github.com/metallust/rms-be/internals/models"
	"github.com/rhnvrm/textsimilarity"
)

func Compare(profile models.Profile, job models.Job) float64 {
	profileString, _ := json.Marshal(profile)
	resume := string(profileString)
    desc := strings.Split(job.Description, ".")
    desc = append(desc, job.Title)
    corpus := append(desc, resume)
    ts := textsimilarity.New(corpus)
	result, _ := ts.Similarity(strings.Join(desc, " "), resume)
	return result
}

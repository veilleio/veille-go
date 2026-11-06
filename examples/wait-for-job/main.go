package main

import (
	"context"
	"fmt"
	"log"

	"github.com/veilleio/veille-go"
)

func main() {
	client, err := veille.NewClient("", nil) // reads VEILLE_API_KEY
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	submitted, err := client.IntelligenceEmail(ctx, &veille.IntelligenceEmailParams{Query: "jane.doe@example.com", IsAsync: veille.Bool(true)})
	if err != nil {
		log.Fatal(err)
	}
	var accepted struct {
		JobID string `json:"job_id"`
	}
	if err := submitted.Decode(&accepted); err != nil {
		log.Fatal(err)
	}
	job, err := client.WaitForJob(ctx, accepted.JobID, &veille.WaitForJobOptions{
		OnProgress: func(job map[string]any) { fmt.Println(job["status"], job["progress"]) },
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(job.Body))
}

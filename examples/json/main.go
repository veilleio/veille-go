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
	response, err := client.IntelligenceEmail(ctx, &veille.IntelligenceEmailParams{Query: "jane.doe@example.com"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(response.Body))
}

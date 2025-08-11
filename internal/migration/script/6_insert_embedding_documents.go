package script

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/aryadhira/otogenius-agent/internal/llm"
	"github.com/aryadhira/otogenius-agent/utils"
)

func InsertEmbeddingDocuments(db *sql.DB) error {
	content := getRAGDocuments()
	query := `
		INSERT INTO documents (content, embedding)
		VALUES ($1, $2)
	`
	embeddingUrl := os.Getenv("EMBEDDING_URL")
	llamacpp := llm.NewLlamaCpp(embeddingUrl, 0, 0, false)

	for _, each := range content {
		strContent := fmt.Sprintf("question: %s\nanswer: %s\nreason: %s", each.Question, each.Answer, each.Reason)

		embedding, err := llamacpp.GetEmbedding(strContent)
		if err != nil {
			return err
		}

		vectorStr := utils.ConvertEmbeddingToString(embedding)

		_, err = db.Exec(query, strContent, vectorStr)
		if err != nil {
			return err
		}

	}

	return nil
}

type Documents struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Reason   string `json:"reason"`
}

func getRAGDocuments() []Documents {
	return []Documents{
		{
			Question: "What is the best car category for a family trip?",
			Answer:   "MPV",
			Reason:   "MPVs have spacious interiors, multiple rows of seating, and large luggage capacity, making them ideal for long trips with many passengers.",
		},
		{
			Question: "Which car type is good for traveling with 6 or more people?",
			Answer:   "MPV",
			Reason:   "MPVs can carry more passengers comfortably, making them suitable for large groups.",
		},
		{
			Question: "I need a car for my small business deliveries and also for family vacations, what should I choose?",
			Answer:   "MPV",
			Reason:   "MPVs can transport goods during business operations and still offer comfort for family use.",
		},
		{
			Question: "Which car is suitable for parents, kids, and grandparents on road trips?",
			Answer:   "MPV",
			Reason:   "MPVs provide comfort, safety, and enough space for multiple generations in one vehicle.",
		},
		{
			Question: "What is the best car category for camping, hiking, and offroad adventures?",
			Answer:   "SUV",
			Reason:   "SUVs have higher ground clearance, all-terrain capability, and strong engines, making them ideal for rough roads.",
		},
		{
			Question: "Which car is suitable for mountain trips and uneven terrain?",
			Answer:   "SUV",
			Reason:   "SUVs are built to handle steep climbs and rocky paths while carrying passengers comfortably.",
		},
		{
			Question: "I want a vehicle that can go to remote campsites, which category is best?",
			Answer:   "SUV",
			Reason:   "SUVs have durable suspension and can navigate dirt or gravel roads effectively.",
		},
		{
			Question: "What is the best car category that is practical and easy to park?",
			Answer:   "Hatchback",
			Reason:   "Hatchbacks are compact, making them easy to maneuver in tight spaces.",
		},
		{
			Question: "Which car is best for daily school drop-offs?",
			Answer:   "Hatchback",
			Reason:   "Hatchbacks are fuel-efficient, easy to park near schools, and affordable to maintain.",
		},
		{
			Question: "I need a small car that saves fuel for daily commuting, what should I choose?",
			Answer:   "Hatchback",
			Reason:   "Hatchbacks have lightweight bodies and smaller engines, which consume less fuel.",
		},
		{
			Question: "Which car is easy to drive in crowded city traffic?",
			Answer:   "Hatchback",
			Reason:   "Hatchbacks handle well in urban areas and require less parking space.",
		},
		{
			Question: "What is the best car category that is sporty and stylish?",
			Answer:   "Sedan",
			Reason:   "Sedans often have sleek designs and offer a comfortable, smooth driving experience.",
		},
		{
			Question: "I want a stylish car for business meetings, which category should I choose?",
			Answer:   "Sedan",
			Reason:   "Sedans project a professional image and are comfortable for long drives.",
		},
		{
			Question: "Which car is good for someone who wants both performance and elegance?",
			Answer:   "Sedan",
			Reason:   "Sedans often have better handling and acceleration compared to larger vehicles.",
		},
		{
			Question: "What is the best transmission type for passing through traffic jams?",
			Answer:   "Automatic",
			Reason:   "Automatic transmissions reduce driver fatigue by eliminating the need to constantly shift gears in stop-and-go traffic.",
		},
		{
			Question: "Which transmission is easier for new drivers in the city?",
			Answer:   "Automatic",
			Reason:   "Automatic cars are easier to control in congested areas and allow the driver to focus on the road.",
		},
		{
			Question: "I want a car that is comfortable to drive during rush hours, which transmission should I pick?",
			Answer:   "Automatic",
			Reason:   "Automatic transmission makes driving smoother in slow-moving traffic.",
		},
	}
}

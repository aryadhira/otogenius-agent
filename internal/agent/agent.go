package agent

import "github.com/aryadhira/otogenius-agent/internal/models"

type Agent interface {
	Run(prompt string) (any, error)
	RunContinues(prompt string, messages []models.Message) (any, error)
}

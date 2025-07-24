package agent

type Agent interface {
	Run(prompt string) (any, error)
}

package contract

import (
	"oma/models"
)

type Opa interface {
	Eval(policy string, input string) (*models.EvalResult, error)
	Format(policy string) (string, error)
	Lint(policy string) (string, error)
}

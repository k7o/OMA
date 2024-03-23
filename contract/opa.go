package contract

import (
	"oma/models"
)

type Opa interface {
	Eval(bundle *models.Bundle, input string, options *models.EvalOptions) (*models.EvalResult, error)
	Format(policy string) (string, error)
	Lint(policy string) (string, []string, error)
}

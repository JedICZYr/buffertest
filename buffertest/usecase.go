package one

import (
	"context"

	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

var (
	validate *validator.Validate
)

type repository interface {
	CreateMetric(ctx context.Context, data string) error
	CloseRepository() error
}

// Usecase for interacting with pcmetricscpu
type Usecase struct {
	Repository repository
}

// CreateCPUMetrics creates a single metric
func (u *Usecase) CreateMetric(ctx context.Context, data string) error {
	validate = validator.New()
	validate.SetTagName("form")
	if err := validate.Struct(*data); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return validationErrors
	}

	if err := u.Repository.CreateMetric(ctx, data); err != nil {
		return errors.Wrap(err, "error creating new CPU Data")
	}

	return nil
}

func (u *Usecase) CloseRepository() error {

	if err := u.Repository.CloseRepository(); err != nil {
		return errors.Wrap(err, "error closing repository")
	}

	return nil
}

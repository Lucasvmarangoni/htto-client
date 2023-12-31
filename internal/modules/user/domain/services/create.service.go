package services

import (
	"context"

	"github.com/Lucasvmarangoni/financial-file-manager/internal/modules/user/domain/entities"
	"github.com/Lucasvmarangoni/financial-file-manager/internal/modules/user/infra/repositories"
)

type CreateService struct {
	Repository *repositories.UserRepositoryDb
}

func (c *CreateService) Create(name, lastName, cpf, email, password string, admin bool) error {
	newUser, err := entities.NewUser(name, lastName, cpf, email, password, admin)
	if err != nil {
		return err
	}
	newUser, err = c.Repository.Insert(newUser, context.Background())
	if err != nil {
		return err
	}
	return nil
}

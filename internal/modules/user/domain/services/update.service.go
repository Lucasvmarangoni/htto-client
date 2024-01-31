package services

import (
	"context"

	"github.com/Lucasvmarangoni/financial-file-manager/internal/modules/user/domain/entities"
	"github.com/Lucasvmarangoni/logella/err"
)

func (u *UserService) Update(id, name, lastName, email, password string) error {

	user, err := u.FindById(id, nil)
	if err != nil {
		return errors.ErrCtx(err, "u.FindById")
	}

	if name == "" {
		name = user.Name
	}
	if lastName == "" {
		lastName = user.LastName
	}
	if email == "" {
		email = user.Email
	}
	if password == "" {
		password = user.Password
	}

	newUser, err := entities.NewUser(name, lastName, user.CPF, email, password)
	if err != nil {
		return errors.ErrCtx(err, "entities.NewUser")
	}
	newUser.Update(user.ID, user.CreatedAt)

	err = u.Repository.Update(newUser, context.Background())
	if err != nil {
		return errors.ErrCtx(err, "Repository.Update")
	}
	return nil
}

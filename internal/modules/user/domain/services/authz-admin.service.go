package services

import (
	"context"

	"github.com/Lucasvmarangoni/logella/err"
)

func (u *UserService) AuthzAdmin(adminID, id string) error {
	admin, err := u.FindById(adminID, nil)
	if err != nil {
		return errors.ErrCtx(err, "u.FindById")
	}

	if admin.Admin == true {
		err := u.Repository.ToggleAdmin(id, context.Background())
		if err != nil {
			return errors.ErrCtx(err, "u.Repository.ToggleAdmin")
		}
	}
	return nil
}

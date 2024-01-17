package services_test

import (
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Lucasvmarangoni/financial-file-manager/config"
	"github.com/Lucasvmarangoni/financial-file-manager/internal/modules/user/domain/entities"
	pkg_entities "github.com/Lucasvmarangoni/financial-file-manager/pkg/entities"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Authn(t *testing.T) {
	userService, mockRepo, _ := prepare(t)
	id, err := pkg_entities.ParseID("52c599f3-af83-4fd9-bfd6-e532918f7b13")
	createdAt, _ := time.Parse(time.RFC3339Nano, "2024-01-17T01:04:23.883938Z")
	password := "hgGFHJ654*"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	emailToFind := "john.doe@example.com"


	mockRepo.EXPECT().
		FindByEmail(emailToFind, gomock.Any()).
		Return(&entities.User{
			ID:        id,
			Name:      "John",
			LastName:  "Doe",
			CPF:       "123.356.229-00",
			Email:     "john.doe@example.com",
			Password:  string(hashedPassword), 
			Admin:     false,
			CreatedAt: createdAt,
			UpdatedAt: nil,
		}, nil).Times(1)

		token, err := userService.Authn(emailToFind, password, config.GetTokenAuth(), 3600)

		assert.NoError(t, err)
		assert.NotEmpty(t, token, "The token should not be empty")
	
		parts := strings.Split(token, ".")
		assert.Equal(t, 3, len(parts), "The token should have three parts separated by dots")
}
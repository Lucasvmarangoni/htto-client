package handlers

import (
	"encoding/json"
	go_err "errors"
	"fmt"
	"github.com/Lucasvmarangoni/logella/err"
	"github.com/asaskevich/govalidator"
	"net/http"
	"sync"

	"github.com/Lucasvmarangoni/financial-file-manager/internal/modules/user/http/dto"
	"github.com/go-chi/jwtauth"
	"github.com/rs/zerolog/log"
)

// Create user godoc
// @Summary      Create user
// @Description  Create user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request     body      dto.UserInput  true  "user data"
// @Success      200
// @Failure      400
// @Router       /authn/create [post]
func (u *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user dto.UserInput
	var wg sync.WaitGroup

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error().Err(err).Msg("Error decode request")
		return
	}

	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "BadRequest",
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = u.userService.Create(user.Name, user.LastName, user.CPF, user.Email, user.Password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error().Stack().Err(err).Msg("Error create user ")
			return
		}
	}()
	wg.Wait()
	w.WriteHeader(http.StatusOK)
}

// Authentication godoc
// @Summary      Generate a user JWT
// @Description  Generate a user JWT. Requires either a CPF or an Email and Password.
// @Tags         Authn
// @Accept       json
// @Produce      json
// @Param        request     body      dto.AuthenticationInput  true  "Authentication input. Requires either a CPF or an Email and Password."
// @Success      200  {object}  dto.GetJWTOutput
// @Failure      400  {object}  string  "Both email and CPF are required for authentication."
// @Failure      401  {object}  string  "Unauthorized."
// @Router       /authn/jwt [post]
func (u *UserHandler) Authentication(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("JwtExpiresIn").(int)
	var user dto.AuthenticationInput

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error().Err(err).Msg("Error decode request")
		return
	}
	u.validatePassword(user.Password, w)

	err = u.validateUserUpdateInputForCPFAndEmail(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "BadRequest",
			"message": fmt.Sprintf("%v", err),
		})
	}

	unique := user.Email + user.CPF
	tokenString, err := u.userService.Authn(unique, user.Password, jwt, jwtExpiresIn)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Error().Stack().Err(err).Msg("Error authenticate user")
		return
	}

	accessToken := dto.GetJWTOutput{AccessToken: tokenString}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

func (u *UserHandler) GetSub(w http.ResponseWriter, r *http.Request) (string, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "", errors.ErrCtx(err, "Failed to get JWT claims")
	}
	id, ok := claims["sub"].(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return "", errors.ErrCtx(err, "sub claim is missing or not a string")
	}
	return id, nil
}

func (u *UserHandler) validateUserUpdateInputForCPFAndEmail(user *dto.AuthenticationInput) error {

	if user.Email == "" && user.CPF == "" {
		return go_err.New("An Email or a CPF is necessary") // Using the golang standard error type because it will be sent in the response
	}
	if user.Email != "" && user.CPF != "" {
		user.CPF = ""
	}
	if err := u.validateEmail(&user.Email); err != nil {
		return err // Using the golang standard error type because it will be sent in the response
	}

	if err := u.validateCPF(&user.CPF); err != nil {
		return err // Using the golang standard error type because it will be sent in the response
	}
	return nil
}

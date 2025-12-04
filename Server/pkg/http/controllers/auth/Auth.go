package auth

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/raiashpanda007/MailForge/pkg/types"
	"github.com/raiashpanda007/MailForge/pkg/utils"
)

type AuthController struct {
	service AuthService
}

func NewAuthController(s AuthService) *AuthController {
	return &AuthController{service: s}
}

func (auth *AuthController) Login(res http.ResponseWriter, req *http.Request) {
	slog.Info("LOGIN USER")
	var loginCreds types.LoginCredentials
	err := json.NewDecoder(req.Body).Decode(&loginCreds)
	if errors.Is(err, io.EOF) {
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "PLEASE PROVIDE CORRECT JSON", Data: utils.GeneralError(err, "PLEASE PROVIDE CORRECT JSON")})
	}

	if err != nil {
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "INVALID DATA INPUT", Data: err.Error()})
	}

	err = validator.New().Struct(loginCreds)
	if err != nil {
		validateErrs := err.(validator.ValidationErrors)
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "PLEASE PROVIDE AND RIGHT AND REQUIRED DATA", Data: utils.ValidationError(validateErrs)})
	}
	result, err := auth.service.Login(req.Context(), loginCreds.Email, loginCreds.Password)
	if err != nil {
		utils.WriteJson(res, http.StatusUnauthorized, utils.Data{Message: "Unable to login", Data: err})
	}
	http.SetCookie(res, &http.Cookie{
		Name:     "Access-Token",
		Value:    result.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   7 * 24 * 60 * 60,
	})
	utils.WriteJson(res, http.StatusOK, utils.Data{Message: "USER LOGGED IN SUCCESSFULLY", Data: nil})

}

func (auth *AuthController) SignUp(res http.ResponseWriter, req *http.Request) {
	slog.Info("SIGN UP USER")
	var signUpCreds types.SignUpCredentials
	err := json.NewDecoder(req.Body).Decode(&signUpCreds)

	if errors.Is(err, io.EOF) {
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "PLEASE PROVIDE CORRECT JSON", Data: utils.GeneralError(err, "PLEASE PROVIDE CORRECT JSON")})
	}

	if err != nil {
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "INVALID DATA INPUT", Data: err.Error()})
	}

	err = validator.New().Struct(signUpCreds)
	if err != nil {
		validateErrs := err.(validator.ValidationErrors)
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "PLEASE PROVIDE AND RIGHT AND REQUIRED DATA", Data: utils.ValidationError(validateErrs)})
	}
	result, err := auth.service.SignUp(req.Context(), signUpCreds.Email, signUpCreds.Name, signUpCreds.Password)

	if err != nil {
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "UNABLE TO SIGN UP", Data: utils.GeneralError(err, "UNABLE TO SIGN UP")})
	}

	//Now login

	loginResults, err := auth.service.Login(req.Context(), result.Email, signUpCreds.Password)
	if err != nil {
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "UNABLE TO SIGN UP IN LOGIN", Data: utils.GeneralError(err, "UNABLE TO SIGN UP IN LOGIN")})
	}

	http.SetCookie(res, &http.Cookie{
		Name:     "Access-Token",
		Value:    loginResults.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   7 * 24 * 60 * 60,
	})
	utils.WriteJson(res, http.StatusOK, utils.Data{Message: "USER CREATED AND LOGGED IN SUCCESSFULLY", Data: nil})
}

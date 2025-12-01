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

func (auth *AuthController) Login(res http.ResponseWriter, req *http.Request) {
	slog.Info("LOGIN USER")
	var loginCreds types.LoginCredentials
	err := json.NewDecoder(req.Body).Decode(&loginCreds)
	if errors.Is(err, io.EOF) {
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "PLEASE PROVIDE CORRECT JSON", Data: utils.GeneralError(err, "PLEASE PROVIDE CORRECT JSON")})
		return
	}

	if err != nil {
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "INVALID DATA INPUT", Data: err.Error()})
		return
	}

	err = validator.New().Struct(loginCreds)
	if err != nil {
		validateErrs := err.(validator.ValidationErrors)
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "PLEASE PROVIDE AND RIGHT AND REQUIRED DATA", Data: utils.ValidationError(validateErrs)})
	}

}

func (auth *AuthController) SignUp(res http.ResponseWriter, req *http.Request) {
	slog.Info("SIGN UP USER")
	var signUpCreds types.SignUpCredentials
	err := json.NewDecoder(req.Body).Decode(&signUpCreds)

	if errors.Is(err, io.EOF) {
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "PLEASE PROVIDE CORRECT JSON", Data: utils.GeneralError(err, "PLEASE PROVIDE CORRECT JSON")})
		return
	}

	if err != nil {
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "INVALID DATA INPUT", Data: err.Error()})
		return
	}

	err = validator.New().Struct(signUpCreds)
	if err != nil {
		validateErrs := err.(validator.ValidationErrors)
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "PLEASE PROVIDE AND RIGHT AND REQUIRED DATA", Data: utils.ValidationError(validateErrs)})
	}
}

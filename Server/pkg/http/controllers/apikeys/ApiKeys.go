package apikeys

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

type ApikeyController struct {
	service ApiKeyService
}

func NewApiKeysController(service ApiKeyService) *ApikeyController {
	return &ApikeyController{service: service}
}

func (r *ApikeyController) GenerateApiKeys(res http.ResponseWriter, req *http.Request) {
	slog.Info("GENERATING API KEY")
	var generateApiKeyCreds types.GenerateApiKeysCredentials
	err := json.NewDecoder(req.Body).Decode(&generateApiKeyCreds)
	if errors.Is(err, io.EOF) {
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "PLEASE PROVIDE CORRECT JSON", Data: utils.GeneralError(err, "PLEASE PROVIDE CORRECT JSON")})
		return
	}
	if err != nil {
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "INVALID DATA INPUT", Data: err.Error()})
		return
	}
	err = validator.New().Struct(generateApiKeyCreds)
	if err != nil {
		validateErrs := err.(validator.ValidationErrors)
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "PLEASE PROVIDE RIGHT AND REQUIRED DATA", Data: utils.ValidationError(validateErrs)})
		return
	}
	result, err := r.service.GenerateKey(req.Context(), generateApiKeyCreds.Organization, generateApiKeyCreds.EmailAppPassword)

	if err != nil {
		utils.WriteJson(res, http.StatusInternalServerError, utils.Data{Message: "UNABLE TO GENERATE API KEY", Data: utils.GeneralError(err, "UNABLE TO GENERATE AND SAVE API KEY")})
		return
	}

	utils.WriteJson(res, http.StatusCreated, utils.Data{Message: "YOUR API KEY IS GENERATED", Data: result})
}

func (r *ApikeyController) DeleteApiKeys(res http.ResponseWriter, req *http.Request) {
	slog.Info("DELETING API KEY")
	var deletingKeysCreds types.DeleteApiKeyCredentials
	err := json.NewDecoder(req.Body).Decode(&deletingKeysCreds)
	if errors.Is(err, io.EOF) {
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "PLEASE PROVIDE CORRECT JSON", Data: utils.GeneralError(err, "PLEASE PROVIDE CORRECT JSON")})
		return
	}
	if err != nil {
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "INVALID DATA INPUT", Data: err.Error()})
		return
	}
	err = validator.New().Struct(deletingKeysCreds)
	if err != nil {
		validateErrs := err.(validator.ValidationErrors)
		utils.WriteJson(res, http.StatusBadRequest, utils.Data{Message: "PLEASE PROVIDE RIGHT AND REQUIRED DATA", Data: utils.ValidationError(validateErrs)})
		return
	}

	result, err := r.service.DeleteKey(req.Context(), deletingKeysCreds.Id)

	if err != nil {
		utils.WriteJson(res, http.StatusInternalServerError, utils.Data{Message: "UNABLE TO DELETE API KEY ", Data: utils.GeneralError(err, "UNABLE TO DELETE API KEY ")})
		return
	}
	utils.WriteJson(res, http.StatusCreated, utils.Data{Message: "YOUR API KEY IS DELETED", Data: result})

}

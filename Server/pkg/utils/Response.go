package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Data struct {
	Message string
	Data    any
}

type Response struct {
	Status  string `json:"status"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

const (
	StatusError = "ERROR"
	StatusOK    = "OK"
)

func WriteJson(response http.ResponseWriter, status int, data Data) error {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)

	return json.NewEncoder(response).Encode(data)
}

func GeneralError(err error, message string) Response {
	return Response{
		Error:   err.Error(),
		Status:  StatusError,
		Message: message,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required field", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is  invalid", err.Field()))
		}
	}

	return Response{
		Status:  StatusError,
		Error:   strings.Join(errMsgs, ", "),
		Message: "Please provide all the required data",
	}

}

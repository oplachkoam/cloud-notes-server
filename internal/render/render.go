package render

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Detail string `json:"detail"`
}

func JSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(&v)
}

func Empty(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func EmptyJSON(w http.ResponseWriter, statusCode int) {
	JSON(w, statusCode, &struct{}{})
}

func Error(w http.ResponseWriter, statusCode int, err error) {
	detail := err.Error()
	JSON(w, statusCode, &ErrorResponse{
		Detail: detail,
	})
}

func ValidationError(w http.ResponseWriter, request any, err error) {
	var errs validator.ValidationErrors
	errors.As(err, &errs)

	fieldName := errs[0].Field()

	t := reflect.TypeOf(request)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	field, ok := t.FieldByName(errs[0].StructField())
	if ok {
		tag := field.Tag.Get("json")
		if tag != "" && tag != "-" {
			fieldName = strings.Split(tag, ",")[0]
		}
	}

	Error(w, http.StatusUnprocessableEntity,
		fmt.Errorf("invalid field: %s", fieldName))
}

func InvalidJSONError(w http.ResponseWriter) {
	Error(w, http.StatusBadRequest, fmt.Errorf("invalid json"))
}

func ServerError(w http.ResponseWriter, statusCode int) {
	Error(w, statusCode, errors.New(strings.ToLower(
		http.StatusText(statusCode))))
}

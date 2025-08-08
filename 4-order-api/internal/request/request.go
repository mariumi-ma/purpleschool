package request

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func HandleBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := decode[T](r.Body)
	if err != nil {
		return nil, err
	}

	err = validate(&body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}

func decode[T any](body io.ReadCloser) (T, error) {
	var payload T

	err := json.NewDecoder(body).Decode(&payload)
	if err != nil {
		return payload, err
	}

	return payload, nil
}

func validate[T any](body *T) error {
	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return err
	}

	return nil
}

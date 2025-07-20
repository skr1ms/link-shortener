package req

import (
	"net/http"

	"linkshortener/pkg/res"
)

func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	data, err := Decode[T](r.Body)
	if err != nil {
		res.Response(*w, 400, err.Error())
		return nil, err
	}

	if err := IsValid(data); err != nil {
		res.Response(*w, 400, err.Error())
		return nil, err
	}

	return &data, nil
}

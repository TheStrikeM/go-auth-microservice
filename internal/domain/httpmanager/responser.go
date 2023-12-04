package httpmanager

import (
	"encoding/json"
	"io"
	"net/http"
)

type ResponseForm[T interface{}] struct {
	Code   int `json:"code"`
	Result T   `json:"result"`
}

func Response[T interface{}](code int, result T) ([]byte, error) {
	return json.Marshal(
		ResponseForm[T]{
			Code:   code,
			Result: result,
		},
	)
}

func Request[T interface{}](req *http.Request, obj *T) error {
	requestBody, _ := io.ReadAll(req.Body)
	if err := json.Unmarshal(requestBody, &obj); err != nil {
		return err
	}
	return nil
}

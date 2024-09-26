package exp

import (
	"errors"
	"fmt"
	"net/http"
)

func IsHttpExp(err error) *HttpExp {
	if err == nil {
		return nil
	}
	var httpExp *HttpExp
	if errors.As(err, &httpExp) {
		return httpExp
	}
	return nil
}

type HttpExp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *HttpExp) Error() string {
	return fmt.Sprintf("Status: %d, Message: %s", e.Status, e.Message)
}

func NewHttpExp(status int, message string) *HttpExp {
	exp := &HttpExp{
		Status:  status,
		Message: message,
	}
	if exp.Message == "" {
		exp.Message = http.StatusText(exp.Status)
	}
	return exp
}

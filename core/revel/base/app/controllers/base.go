package base

import (
	"encoding/json"
	"errors"
	"github.com/revel/revel"
	"net/http"
)

var (
	ERROR_JSON_INPUT      = errors.New("Invalid JSON input.")
	ERROR_INTERNAL        = errors.New("Unexpected internal server error.")
	ERROR_DATABASE        = errors.New("Database error.")
	ERROR_NOT_FOUND       = errors.New("Not found.")
	ERROR_NO_ACTION       = errors.New("Action not found.")
	ERROR_NOT_IMPLEMENTED = errors.New("Action not implemented.")
)

type Base struct {
	*revel.Controller
}

func (b *Base) FromJson(obj interface{}, validator func(interface{}, *revel.Validation)) revel.Result {
	err := json.NewDecoder(b.Request.Body).Decode(&obj)
	if err != nil {
		return b.Error(ERROR_JSON_INPUT, http.StatusBadRequest)
	}

	if validator != nil {
		validator(obj, b.Validation)
		if b.Validation.HasErrors() {
			return b.RenderJson(b.Validation.Errors)
		}
	}

	return nil
}

func (b *Base) Error(err error, httpStatus int) revel.Result {
	b.Response.Status = httpStatus
	return b.RenderJson(err.Error())
}

func (b *Base) ErrorNotFound() revel.Result {
	return b.Error(ERROR_NOT_FOUND, http.StatusNotFound)
}

func (b *Base) ErrorNotFoundMsg(msg string) revel.Result {
	return b.Error(errors.New(msg), http.StatusNotFound)
}

func (b *Base) ErrorInternal() revel.Result {
	return b.Error(ERROR_INTERNAL, http.StatusInternalServerError)
}

func (b *Base) ErrorInternalMsg(msg string) revel.Result {
	return b.Error(errors.New(msg), http.StatusInternalServerError)
}

func (b *Base) ErrorNotImplemented() revel.Result {
	return b.Error(ERROR_NOT_IMPLEMENTED, http.StatusInternalServerError)
}

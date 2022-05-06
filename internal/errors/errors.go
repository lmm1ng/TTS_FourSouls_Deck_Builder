package errors

import (
	"log"
	"net/http"
)

var (
	InternalError  = NewError("internal error").HTTP(http.StatusInternalServerError)
	BarURL         = NewError("bad url")
	BadHTTPRequest = NewError("bad http request")

	BadName = NewError("bad name")

	GameExist          = NewError("game exist")
	GameNotExists      = NewError("game not exists").HTTP(http.StatusNoContent)
	GameInvalid        = NewError("game data invalid")
	GameInfoNotExists  = NewError("game info not exists").HTTP(http.StatusInternalServerError)
	GameImageNotExists = NewError("game image not exists").HTTP(http.StatusNoContent)

	CollectionExist          = NewError("collection exist")
	CollectionNotExists      = NewError("collection not exists").HTTP(http.StatusNoContent)
	CollectionInvalid        = NewError("collection data invalid")
	CollectionInfoNotExists  = NewError("collection info not exists").HTTP(http.StatusInternalServerError)
	CollectionImageNotExists = NewError("collection image not exists").HTTP(http.StatusNoContent)

	DeckExist     = NewError("deck exist")
	DeckNotExists = NewError("deck not exists").HTTP(http.StatusNoContent)

	UnknownImageType = NewError("unknown image type")
)

type Error struct {
	Message string `json:"message"`
	code    int
}

func NewError(message string) *Error {
	return &Error{
		Message: message,
		code:    http.StatusBadRequest,
	}
}

func (e Error) HTTP(code int) *Error {
	e.code = code
	return &e
}
func (e Error) AddMessage(message string) *Error {
	e.Message += ": " + message
	return &e
}

func (e *Error) GetCode() int       { return e.code }
func (e *Error) GetMessage() string { return e.Message }

func IfErrorLog(err error) {
	if err != nil {
		log.Output(2, err.Error())
	}
}

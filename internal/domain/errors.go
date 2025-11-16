package domain

import "errors"

var (
	ErrPokemonNotFound    = errors.New("pokemon not found")
	ErrInvalidPokemonData = errors.New("invalid pokemon data")
	ErrPokeAPIUnavailable = errors.New("pokeapi service unavailable")
	ErrInternalServer     = errors.New("internal server error")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

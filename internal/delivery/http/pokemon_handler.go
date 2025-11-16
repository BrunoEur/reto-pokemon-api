package delivery

import (
	"net/http"
	"strconv"

	"reto-pokemon-api/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PokemonHandler struct {
	pokemonUseCase domain.PokemonUseCase
	validator      *validator.Validate
}

func NewPokemonHandler(pokemonUseCase domain.PokemonUseCase) *PokemonHandler {
	return &PokemonHandler{
		pokemonUseCase: pokemonUseCase,
		validator:      validator.New(),
	}
}

func (h *PokemonHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Pokemon API is running",
		"version": "1.0.0",
	})
}

func (h *PokemonHandler) GetPokemon(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.sendError(c, http.StatusBadRequest, "Pokemon ID is required", nil)
		return
	}

	pokemon, err := h.pokemonUseCase.GetPokemonByID(id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, pokemon)
}

func (h *PokemonHandler) GetPokemonByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		h.sendError(c, http.StatusBadRequest, "Pokemon name is required", nil)
		return
	}

	pokemon, err := h.pokemonUseCase.GetPokemonByName(name)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, pokemon)
}

func (h *PokemonHandler) GetAllPokemon(c *gin.Context) {
	filter := domain.PokemonFilter{
		Limit:  h.parseIntQuery(c, "limit", 0),
		Offset: h.parseIntQuery(c, "offset", 0),
	}
	
	if isFavoriteStr := c.Query("is_favorite"); isFavoriteStr != "" {
		if isFavorite, err := strconv.ParseBool(isFavoriteStr); err == nil {
			filter.IsFavorite = &isFavorite
		}
	}
	
	pokemon, err := h.pokemonUseCase.GetPokemonAll(filter)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, pokemon)
}

func (h *PokemonHandler) parseIntQuery(c *gin.Context, key string, defaultValue int) int {
	if value := c.Query(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func (h *PokemonHandler) handleError(c *gin.Context, err error) {
	switch err {
	case domain.ErrPokemonNotFound:
		h.sendError(c, http.StatusNotFound, "Pokemon not found", err)
	case domain.ErrInvalidPokemonData:
		h.sendError(c, http.StatusBadRequest, "Invalid pokemon data", err)
	case domain.ErrPokeAPIUnavailable:
		h.sendError(c, http.StatusServiceUnavailable, "PokeAPI service unavailable", err)
	default:
		h.sendError(c, http.StatusInternalServerError, "Internal server error", err)
	}
}

func (h *PokemonHandler) sendError(c *gin.Context, code int, message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg = err.Error()
	}

	c.JSON(code, domain.ErrorResponse{
		Error:   message,
		Message: errorMsg,
		Code:    code,
	})
}

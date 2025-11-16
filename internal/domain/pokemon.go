package domain

import "time"

type Pokemon struct {
	ID         int       `json:"id"`
	Name       string    `json:"name" validate:"required"`
	Height     int       `json:"height"`
	Weight     int       `json:"weight"`
	BaseExp    int       `json:"base_experience"`
	Types      []Type    `json:"types"`
	Abilities  []Ability `json:"abilities"`
	Sprites    Sprite    `json:"sprites"`
	Stats      []Stat    `json:"stats"`
	IsFavorite bool      `json:"is_favorite"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	PokeAPIID  int       `json:"pokeapi_id"`
}

type PokemonList struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Pokemons *[]Pokemon `json:"pokemons"`
}

type Type struct {
	Slot int      `json:"slot"`
	Type TypeInfo `json:"type"`
}

type TypeInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Ability struct {
	IsHidden bool        `json:"is_hidden"`
	Slot     int         `json:"slot"`
	Ability  AbilityInfo `json:"ability"`
}

type AbilityInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Sprite struct {
	FrontDefault string `json:"front_default"`
	FrontShiny   string `json:"front_shiny"`
	BackDefault  string `json:"back_default"`
	BackShiny    string `json:"back_shiny"`
}

type Stat struct {
	BaseStat int      `json:"base_stat"`
	Effort   int      `json:"effort"`
	Stat     StatInfo `json:"stat"`
}

type StatInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonFilter struct {
	Name       string `json:"name,omitempty"`
	Type       string `json:"type,omitempty"`
	IsFavorite *bool  `json:"is_favorite,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Offset     int    `json:"offset,omitempty"`
}

type PokeAPIRepository interface {
	GetPokemonByID(id int) (*Pokemon, error)
	GetPokemonByName(name string) (*Pokemon, error)
	GetPokemonAll(filter PokemonFilter) (*PokemonList, error)
}

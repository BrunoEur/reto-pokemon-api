package application

import (
	"strconv"
	"time"

	"reto-pokemon-api/internal/domain"
)

type pokemonUseCase struct {
	pokeAPIRepo domain.PokeAPIRepository
}

func NewPokemonUseCase(pokeAPIRepo domain.PokeAPIRepository) domain.PokemonUseCase {
	return &pokemonUseCase{
		pokeAPIRepo: pokeAPIRepo,
	}
}

func (uc *pokemonUseCase) GetPokemonByID(id string) (*domain.Pokemon, error) {
	i, err := strconv.Atoi(id)
	if err != nil {
		return nil, domain.ErrInvalidPokemonData
	}

	apiPokemon, err := uc.pokeAPIRepo.GetPokemonByID(i)
	if err != nil {
		// Si hay error en la API, retornamos un pokemon básico
		return &domain.Pokemon{
			ID:         i,
			IsFavorite: false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}, nil
	}

	pokemon := &domain.Pokemon{
		ID:         i,
		Name:       apiPokemon.Name,
		Height:     apiPokemon.Height,
		Weight:     apiPokemon.Weight,
		BaseExp:    apiPokemon.BaseExp,
		Types:      apiPokemon.Types,
		Abilities:  apiPokemon.Abilities,
		Sprites:    apiPokemon.Sprites,
		Stats:      apiPokemon.Stats,
		IsFavorite: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		PokeAPIID:  apiPokemon.PokeAPIID,
	}
	return pokemon, nil
}

func (uc *pokemonUseCase) GetPokemonByName(name string) (*domain.Pokemon, error) {
	return uc.pokeAPIRepo.GetPokemonByName(name)
}

func (uc *pokemonUseCase) GetPokemonAll(filter domain.PokemonFilter) (*domain.PokemonList, error) {
	return uc.pokeAPIRepo.GetPokemonAll(filter)
}

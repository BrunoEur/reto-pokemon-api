package application

import (
	"errors"
	"reto-pokemon-api/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
)

type MockPokeAPIRepository struct {
	mock.Mock
}

func (m *MockPokeAPIRepository) GetPokemonByID(id int) (*domain.Pokemon, error) {

	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Pokemon), args.Error(1)

}

func (m *MockPokeAPIRepository) GetPokemonByName(name string) (*domain.Pokemon, error) {

	args := m.Called(name)

	if args.Get(0) == nil {

		return nil, args.Error(1)

	}

	return args.Get(0).(*domain.Pokemon), args.Error(1)

}

func (m *MockPokeAPIRepository) GetPokemonAll(filter domain.PokemonFilter) (*domain.PokemonList, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PokemonList), args.Error(1)
}

func TestPokemonUseCase_GetPokemonByID(t *testing.T) {

	mockPokeAPIRepo := new(MockPokeAPIRepository)

	useCase := NewPokemonUseCase(mockPokeAPIRepo)

	t.Run("Success", func(t *testing.T) {

		apiPokemon := &domain.Pokemon{
			Name:      "pikachu",
			Height:    4,
			Weight:    60,
			BaseExp:   112,
			PokeAPIID: 25,
			Types: []domain.Type{
				{
					Slot: 1,
					Type: domain.TypeInfo{Name: "electric", URL: "https://pokeapi.co/api/v2/type/13/"},
				},
			},
		}

		mockPokeAPIRepo.On("GetPokemonByID", 25).Return(apiPokemon, nil)

		result, err := useCase.GetPokemonByID("25")

		assert.NoError(t, err)

		assert.NotNil(t, result)
		assert.Equal(t, "pikachu", result.Name)

		assert.Equal(t, 25, result.ID)
		assert.Equal(t, 25, result.PokeAPIID)

		assert.Equal(t, 4, result.Height)

		assert.Equal(t, 60, result.Weight)

		assert.False(t, result.IsFavorite)

		mockPokeAPIRepo.AssertExpectations(t)
	})

	t.Run("Success - API error ignored", func(t *testing.T) {
		mockPokeAPIRepo.On("GetPokemonByID", 999).Return(nil, errors.New("API error"))
		result, err := useCase.GetPokemonByID("999")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockPokeAPIRepo.AssertExpectations(t)

	})
}

func TestPokemonUseCase_GetPokemonAll(t *testing.T) {

	mockPokeAPIRepo := new(MockPokeAPIRepository)
	useCase := NewPokemonUseCase(mockPokeAPIRepo)

	t.Run("Success - returns pokemon list", func(t *testing.T) {
		filter := domain.PokemonFilter{
			Limit:  10,
			Offset: 0,
		}

		expectedPokemons := []domain.Pokemon{
			{
				ID:        25,
				Name:      "pikachu",
				Height:    4,
				Weight:    60,
				BaseExp:   112,
				PokeAPIID: 25,
				Types: []domain.Type{
					{
						Slot: 1,
						Type: domain.TypeInfo{Name: "electric", URL: "https://pokeapi.co/api/v2/type/13/"},
					},
				},
			},
			{
				ID:        4,
				Name:      "charmander",
				Height:    6,
				Weight:    85,
				BaseExp:   62,
				PokeAPIID: 4,
				Types: []domain.Type{
					{
						Slot: 1,
						Type: domain.TypeInfo{Name: "fire", URL: "https://pokeapi.co/api/v2/type/10/"},
					},
				},
			},
		}

		expectedList := &domain.PokemonList{
			Count:    2,
			Next:     "",
			Previous: "",
			Pokemons: &expectedPokemons,
		}

		mockPokeAPIRepo.On("GetPokemonAll", filter).Return(expectedList, nil)

		result, err := useCase.GetPokemonAll(filter)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.Count)
		assert.NotNil(t, result.Pokemons)
		assert.Len(t, *result.Pokemons, 2)
		assert.Equal(t, "pikachu", (*result.Pokemons)[0].Name)
		assert.Equal(t, "charmander", (*result.Pokemons)[1].Name)

		mockPokeAPIRepo.AssertExpectations(t)
	})

	t.Run("Error - API returns error", func(t *testing.T) {
		filter := domain.PokemonFilter{}

		mockPokeAPIRepo.On("GetPokemonAll", filter).Return(nil, domain.ErrPokeAPIUnavailable)

		result, err := useCase.GetPokemonAll(filter)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrPokeAPIUnavailable, err)

		mockPokeAPIRepo.AssertExpectations(t)
	})
}

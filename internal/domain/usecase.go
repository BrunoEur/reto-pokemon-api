package domain

type PokemonUseCase interface {
	GetPokemonByID(id string) (*Pokemon, error)
	GetPokemonByName(name string) (*Pokemon, error)
	GetPokemonAll(filter PokemonFilter) (*PokemonList, error)
}

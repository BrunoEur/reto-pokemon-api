package infrastructure

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"reto-pokemon-api/internal/domain"
)

type pokeAPIRepository struct {
	client  *http.Client
	baseURL string
	cache   *Cache
}

func NewPokeAPIRepository() domain.PokeAPIRepository {
	baseURL := os.Getenv("POKEAPI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://pokeapi.co/api/v2"
	}

	cachettlEnv := os.Getenv("CACHE_TTL")
	if cachettlEnv == "0" {
		cachettlEnv = "60"
	}
	
	cacheTTL := 60 * time.Minute
	if ttlEnv := cachettlEnv; ttlEnv != "" {
		if ttlMinutes, err := strconv.Atoi(ttlEnv); err == nil {
			cacheTTL = time.Duration(ttlMinutes) * time.Minute
			log.Printf("Cache TTL configured to %d minutes", ttlMinutes)
		}
	}
	
	log.Printf("Initializing cache with TTL: %v", cacheTTL)
	
	return &pokeAPIRepository{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
		cache:   NewCache(cacheTTL),
	}
}

type PokeAPIResponse struct {
	ID             int              `json:"id"`
	Name           string           `json:"name"`
	Height         int              `json:"height"`
	Weight         int              `json:"weight"`
	BaseExperience int              `json:"base_experience"`
	Types          []PokeAPIType    `json:"types"`
	Abilities      []PokeAPIAbility `json:"abilities"`
	Sprites        PokeAPISprites   `json:"sprites"`
	Stats          []PokeAPIStat    `json:"stats"`
}

type PokeAPIResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type PokeAPIResponseList struct {
	Count    int             `json:"count"`
	Next     string          `json:"next"`
	Previous string          `json:"previous"`
	Results  []PokeAPIResult `json:"results"`
}

type PokeAPIType struct {
	Slot int `json:"slot"`
	Type struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"type"`
}

type PokeAPIAbility struct {
	IsHidden bool `json:"is_hidden"`
	Slot     int  `json:"slot"`
	Ability  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"ability"`
}

type PokeAPISprites struct {
	FrontDefault string `json:"front_default"`
	FrontShiny   string `json:"front_shiny"`
	BackDefault  string `json:"back_default"`
	BackShiny    string `json:"back_shiny"`
}

type PokeAPIStat struct {
	BaseStat int `json:"base_stat"`
	Effort   int `json:"effort"`
	Stat     struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"stat"`
}

func (r *pokeAPIRepository) GetPokemonByID(id int) (*domain.Pokemon, error) {

	cacheKey := fmt.Sprintf("pokemon:id:%d", id)
	if cached, found := r.cache.Get(cacheKey); found {
		log.Printf("Cache HIT for pokemon ID: %d", id)
		return cached.(*domain.Pokemon), nil
	}
	log.Printf("Cache MISS for pokemon ID: %d", id)
	url := fmt.Sprintf("%s/pokemon/%d", r.baseURL, id)
	pokemon, err := r.fetchPokemon(url)
	
	if err == nil && pokemon != nil {
		r.cache.Set(cacheKey, pokemon)
	}
	
	return pokemon, err
}

func (r *pokeAPIRepository) GetPokemonByName(name string) (*domain.Pokemon, error) {

	cacheKey := fmt.Sprintf("pokemon:name:%s", name)
	if cached, found := r.cache.Get(cacheKey); found {
		log.Printf("Cache HIT for pokemon name: %s", name)
		return cached.(*domain.Pokemon), nil
	}
	
	log.Printf("Cache MISS for pokemon name: %s", name)
	url := fmt.Sprintf("%s/pokemon/%s", r.baseURL, name)
	pokemon, err := r.fetchPokemon(url)
	
	if err == nil && pokemon != nil {
		r.cache.Set(cacheKey, pokemon)
	}
	
	return pokemon, err
}

func (r *pokeAPIRepository) GetPokemonAll(filter domain.PokemonFilter) (*domain.PokemonList, error) {
	offset := 0
	limitset := 20
	
	if filter.Offset > 0 {
		offset = filter.Offset
	}
	if filter.Limit > 0 {
		limitset = filter.Limit
	}
	
	cacheKey := fmt.Sprintf("pokemon:list:offset:%d:limit:%d", offset, limitset)
	if cached, found := r.cache.Get(cacheKey); found {
		log.Printf("Cache HIT for pokemon list (offset: %d, limit: %d)", offset, limitset)
		return cached.(*domain.PokemonList), nil
	}
	
	log.Printf("Cache MISS for pokemon list (offset: %d, limit: %d)", offset, limitset)
	url := fmt.Sprintf("%s/pokemon?offset=%d&limit=%d", r.baseURL, offset, limitset)
	log.Println("url:", url)
	
	pokemonList, err := r.fetchPokemonAll(url)
	
	if err == nil && pokemonList != nil {
		r.cache.Set(cacheKey, pokemonList)
	}
	
	return pokemonList, err
}

func (r *pokeAPIRepository) fetchPokemon(url string) (*domain.Pokemon, error) {
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pokemon from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrPokemonNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var pokeAPIResp PokeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&pokeAPIResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return r.mapToDomainPokemon(&pokeAPIResp), nil
}

func (r *pokeAPIRepository) fetchPokemonAll(url string) (*domain.PokemonList, error) {
	resultPokemon := &[]domain.Pokemon{}
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pokemon from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrPokemonNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var PokeAPIResponseList PokeAPIResponseList
	if err := json.NewDecoder(resp.Body).Decode(&PokeAPIResponseList); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	for _, t := range PokeAPIResponseList.Results {
		pokemon, err := r.fetchPokemon(t.URL)
		if err != nil {
			return nil, err
		}
		*resultPokemon = append(*resultPokemon, *pokemon)
	}

	return &domain.PokemonList{
		Count:    PokeAPIResponseList.Count,
		Next:     PokeAPIResponseList.Next,
		Previous: PokeAPIResponseList.Previous,
		Pokemons: resultPokemon,
	}, nil
}

func (r *pokeAPIRepository) mapToDomainPokemon(apiPokemon *PokeAPIResponse) *domain.Pokemon {
	now := time.Now()

	types := make([]domain.Type, len(apiPokemon.Types))
	for i, t := range apiPokemon.Types {
		types[i] = domain.Type{
			Slot: t.Slot,
			Type: domain.TypeInfo{
				Name: t.Type.Name,
				URL:  t.Type.URL,
			},
		}
	}

	abilities := make([]domain.Ability, len(apiPokemon.Abilities))
	for i, a := range apiPokemon.Abilities {
		abilities[i] = domain.Ability{
			IsHidden: a.IsHidden,
			Slot:     a.Slot,
			Ability: domain.AbilityInfo{
				Name: a.Ability.Name,
				URL:  a.Ability.URL,
			},
		}
	}

	stats := make([]domain.Stat, len(apiPokemon.Stats))
	for i, s := range apiPokemon.Stats {
		stats[i] = domain.Stat{
			BaseStat: s.BaseStat,
			Effort:   s.Effort,
			Stat: domain.StatInfo{
				Name: s.Stat.Name,
				URL:  s.Stat.URL,
			},
		}
	}

	return &domain.Pokemon{
		ID: apiPokemon.ID,
		Name:      apiPokemon.Name,
		Height:    apiPokemon.Height,
		Weight:    apiPokemon.Weight,
		BaseExp:   apiPokemon.BaseExperience,
		Types:     types,
		Abilities: abilities,
		Sprites: domain.Sprite{
			FrontDefault: apiPokemon.Sprites.FrontDefault,
			FrontShiny:   apiPokemon.Sprites.FrontShiny,
			BackDefault:  apiPokemon.Sprites.BackDefault,
			BackShiny:    apiPokemon.Sprites.BackShiny,
		},
		Stats:     stats,
		CreatedAt: now,
		UpdatedAt: now,
		PokeAPIID: apiPokemon.ID,
	}
}

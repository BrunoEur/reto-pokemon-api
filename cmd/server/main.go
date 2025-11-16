package main

import (
	"log"
	"os"

	"reto-pokemon-api/internal/application"
	delivery "reto-pokemon-api/internal/delivery/http"
	"reto-pokemon-api/internal/infrastructure"
)

func main() {
	pokeAPIRepo := infrastructure.NewPokeAPIRepository()

	pokemonUseCase := application.NewPokemonUseCase(pokeAPIRepo)

	pokemonHandler := delivery.NewPokemonHandler(pokemonUseCase)

	router := delivery.SetupRoutes(pokemonHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Pokemon API server on port %s", port)
	log.Printf("Health check available at: http://localhost:%s/health", port)
	log.Printf("API documentation available at: http://localhost:%s/api/v1/pokemon", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

package delivery

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(pokemonHandler *PokemonHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

	router.GET("/health", pokemonHandler.HealthCheck)

	v1 := router.Group("/api/v1")
	{
		pokemon := v1.Group("/pokemon")
		{
			pokemon.GET("", pokemonHandler.GetAllPokemon)
			pokemon.GET("/:id", pokemonHandler.GetPokemon)
			pokemon.GET("/name/:name", pokemonHandler.GetPokemonByName)
		}
	}

	return router
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

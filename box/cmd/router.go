package main

import (
	"crypto/rsa"

	"pob/box/cmd/di"
	"pob/box/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(c *di.Container, publicKey *rsa.PublicKey) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.TraceMiddleware())

	auth := router.Group("/", middleware.AuthMiddleware(publicKey))
	{
		box := auth.Group("/boxes")
		{
			box.GET("", c.Box.GetBoxes)
			box.POST("", c.Box.CreateBox)
			box.PATCH("/:id", c.Box.UpdateBoxName)
			box.DELETE("/:id", c.Box.DeleteBox)
			box.GET("/:id/pokemon", c.Box.GetBoxPokemon)
			box.POST("/:id/pokemon", c.Box.AddBoxPokemon)
			box.PATCH("/:id/pokemon/:pid", c.Box.UpdateBoxPokemon)
			box.DELETE("/:id/pokemon/:pid", c.Box.DeleteBoxPokemon)
		}

		party := auth.Group("/parties")
		{
			party.GET("", c.Party.GetParties)
			party.POST("", c.Party.CreateParty)
			party.PATCH("/:id", c.Party.UpdatePartyName)
			party.DELETE("/:id", c.Party.DeleteParty)
			party.GET("/:id/pokemon", c.Party.GetPartyPokemon)
			party.PUT("/:id/pokemon", c.Party.SetPartyPokemon)
		}
	}

	dex := router.Group("/dex")
	{
		dex.GET("/pokemon", c.Dex.GetPokemonList)
		dex.GET("/pokemon/:id", c.Dex.GetPokemon)
	}

	return router
}

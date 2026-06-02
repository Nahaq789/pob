package handler

import (
	"net/http"
	"strconv"

	gen "pob/box/proto"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type DexHandler struct {
	dex gen.DexServiceClient
}

func NewDexHandler(dex gen.DexServiceClient) *DexHandler {
	return &DexHandler{dex: dex}
}

func (h *DexHandler) GetPokemonList(c *gin.Context) {
	res, err := h.dex.GetPokemonList(c.Request.Context(), &gen.GetPokemonListRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": res.Pokemon})
}

func (h *DexHandler) GetPokemon(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "invalid id"})
		return
	}

	ctx := c.Request.Context()
	eg, egCtx := errgroup.WithContext(ctx)

	var pokemon *gen.PokemonResponse
	var moves *gen.LearnableMovesResponse

	eg.Go(func() error {
		var e error
		pokemon, e = h.dex.GetPokemon(egCtx, &gen.GetPokemonRequest{PokemonId: int32(id)})
		return e
	})
	eg.Go(func() error {
		var e error
		moves, e = h.dex.GetLearnableMoves(egCtx, &gen.GetLearnableMovesRequest{PokemonId: int32(id)})
		return e
	})

	if err := eg.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": gin.H{
		"pokemon": pokemon,
		"moves":   moves.Moves,
	}})
}

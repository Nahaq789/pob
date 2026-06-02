package handler

import (
	"log/slog"
	"net/http"

	"pob/box/internal/service"
	"pob/box/internal/service/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PartyHandler struct {
	partySvc        *service.PartyService
	partyPokemonSvc *service.PartyPokemonService
}

func NewPartyHandler(partySvc *service.PartyService, partyPokemonSvc *service.PartyPokemonService) *PartyHandler {
	return &PartyHandler{partySvc: partySvc, partyPokemonSvc: partyPokemonSvc}
}

func (h *PartyHandler) GetParties(c *gin.Context) {
	ctx := c.Request.Context()
	userId, err := uuid.Parse(getUserId(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "invalid user_id"})
		return
	}

	parties, err := h.partySvc.GetParties(ctx, userId)
	if err != nil {
		slog.ErrorContext(ctx, "GetParties failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": parties})
}

func (h *PartyHandler) CreateParty(c *gin.Context) {
	ctx := c.Request.Context()
	userId, err := uuid.Parse(getUserId(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "invalid user_id"})
		return
	}

	var req dto.CreatePartyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		return
	}

	party, err := h.partySvc.Create(ctx, userId, req.Name)
	if err != nil {
		slog.ErrorContext(ctx, "Create party failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": party})
}

func (h *PartyHandler) UpdatePartyName(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "invalid id"})
		return
	}

	var req dto.UpdatePartyNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		return
	}

	if err := h.partySvc.UpdateName(ctx, id, req.Name); err != nil {
		slog.ErrorContext(ctx, "UpdateName failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "updated"})
}

func (h *PartyHandler) DeleteParty(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "invalid id"})
		return
	}

	if err := h.partySvc.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Delete party failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "deleted"})
}

func (h *PartyHandler) GetPartyPokemon(c *gin.Context) {
	ctx := c.Request.Context()
	partyId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "invalid id"})
		return
	}

	pokemon, err := h.partyPokemonSvc.GetByPartyId(ctx, partyId)
	if err != nil {
		slog.ErrorContext(ctx, "GetByPartyId failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": pokemon})
}

func (h *PartyHandler) SetPartyPokemon(c *gin.Context) {
	ctx := c.Request.Context()
	partyId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "invalid id"})
		return
	}

	var req dto.SetPartyPokemonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		return
	}

	for _, p := range req.Pokemon {
		if p.Slot < 1 || p.Slot > 6 {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "slot must be between 1 and 6"})
			return
		}
	}

	if err := h.partyPokemonSvc.SetPokemon(ctx, partyId, req.Pokemon); err != nil {
		slog.ErrorContext(ctx, "SetPokemon failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "updated"})
}

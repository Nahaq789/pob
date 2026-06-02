package handler

import (
	"log/slog"
	"net/http"

	"pob/box/internal/service"
	"pob/box/internal/service/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BoxHandler struct {
	boxSvc        *service.BoxService
	boxPokemonSvc *service.BoxPokemonService
}

func NewBoxHandler(boxSvc *service.BoxService, boxPokemonSvc *service.BoxPokemonService) *BoxHandler {
	return &BoxHandler{boxSvc: boxSvc, boxPokemonSvc: boxPokemonSvc}
}

func (h *BoxHandler) GetBoxes(c *gin.Context) {
	ctx := c.Request.Context()
	userId, err := uuid.Parse(getUserId(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "invalid user_id"})
		return
	}

	boxes, err := h.boxSvc.GetBoxes(ctx, userId)
	if err != nil {
		slog.ErrorContext(ctx, "GetBoxes failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": boxes})
}

func (h *BoxHandler) CreateBox(c *gin.Context) {
	ctx := c.Request.Context()
	userId, err := uuid.Parse(getUserId(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "invalid user_id"})
		return
	}

	var req dto.CreateBoxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		return
	}

	box, err := h.boxSvc.Create(ctx, userId, req.Name)
	if err != nil {
		slog.ErrorContext(ctx, "Create box failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": box})
}

func (h *BoxHandler) UpdateBoxName(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "invalid id"})
		return
	}

	var req dto.UpdateBoxNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		return
	}

	if err := h.boxSvc.UpdateName(ctx, id, req.Name); err != nil {
		slog.ErrorContext(ctx, "UpdateName failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "updated"})
}

func (h *BoxHandler) DeleteBox(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "invalid id"})
		return
	}

	if err := h.boxSvc.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Delete box failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "deleted"})
}

func (h *BoxHandler) GetBoxPokemon(c *gin.Context) {
	ctx := c.Request.Context()
	boxId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "invalid id"})
		return
	}

	pokemon, err := h.boxPokemonSvc.GetByBoxId(ctx, boxId)
	if err != nil {
		slog.ErrorContext(ctx, "GetByBoxId failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": pokemon})
}

func (h *BoxHandler) AddBoxPokemon(c *gin.Context) {
	ctx := c.Request.Context()
	boxId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "invalid id"})
		return
	}

	var req dto.AddBoxPokemonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		return
	}

	bp, err := h.boxPokemonSvc.Add(ctx, boxId, req)
	if err != nil {
		slog.ErrorContext(ctx, "Add box_pokemon failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": bp})
}

func (h *BoxHandler) UpdateBoxPokemon(c *gin.Context) {
	ctx := c.Request.Context()
	if _, err := uuid.Parse(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "invalid id"})
		return
	}

	var req dto.UpdateBoxPokemonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		return
	}
	req.BoxPokemonId = c.Param("pid")

	bp, err := h.boxPokemonSvc.Update(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "Update box_pokemon failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": bp})
}

func (h *BoxHandler) DeleteBoxPokemon(c *gin.Context) {
	ctx := c.Request.Context()
	pid, err := uuid.Parse(c.Param("pid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "invalid pid"})
		return
	}

	if err := h.boxPokemonSvc.Delete(ctx, pid); err != nil {
		slog.ErrorContext(ctx, "Delete box_pokemon failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "deleted"})
}

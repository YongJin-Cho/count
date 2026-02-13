package inbound

import (
	"count-processing-service/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CountValueHandler struct {
	useCase domain.CountValueUseCase
}

func NewCountValueHandler(useCase domain.CountValueUseCase) *CountValueHandler {
	return &CountValueHandler{useCase: useCase}
}

func (h *CountValueHandler) RegisterRoutes(r *gin.Engine) {
	internal := r.Group("/api/v1/internal/counts")
	{
		internal.POST("", h.Initialize)
		internal.GET("", h.GetMultiple)
		internal.GET("/:itemId", h.GetSingle)
		internal.DELETE("/:itemId", h.Delete)
	}

	external := r.Group("/api/v1/counts/:itemId")
	{
		external.POST("/increase", h.Increase)
		external.POST("/decrease", h.Decrease)
		external.POST("/reset", h.Reset)
	}
}

type UpdateRequest struct {
	Amount int `json:"amount"`
}

func (h *CountValueHandler) Increase(c *gin.Context) {
	itemID := c.Param("itemId")
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Default to 1 if body is empty or malformed
		req.Amount = 1
	}
	if req.Amount == 0 {
		req.Amount = 1
	}

	count, err := h.useCase.Increase(c.Request.Context(), itemID, req.Amount)
	if err == domain.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"itemId": count.ItemID, "value": count.CurrentValue})
}

func (h *CountValueHandler) Decrease(c *gin.Context) {
	itemID := c.Param("itemId")
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Default to 1
		req.Amount = 1
	}
	if req.Amount == 0 {
		req.Amount = 1
	}

	count, err := h.useCase.Decrease(c.Request.Context(), itemID, req.Amount)
	if err == domain.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"itemId": count.ItemID, "value": count.CurrentValue})
}

func (h *CountValueHandler) Reset(c *gin.Context) {
	itemID := c.Param("itemId")
	count, err := h.useCase.Reset(c.Request.Context(), itemID)
	if err == domain.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"itemId": count.ItemID, "value": count.CurrentValue})
}

type InitializeRequest struct {
	ItemID       string `json:"itemId" binding:"required"`
	InitialValue int    `json:"initialValue"`
}

func (h *CountValueHandler) Initialize(c *gin.Context) {
	var req InitializeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := h.useCase.Initialize(c.Request.Context(), req.ItemID, req.InitialValue)
	if err == domain.ErrAlreadyExists {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, count)
}

func (h *CountValueHandler) GetMultiple(c *gin.Context) {
	itemIDs := c.QueryArray("itemIds")
	counts, err := h.useCase.GetMultiple(c.Request.Context(), itemIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"counts": counts})
}

func (h *CountValueHandler) GetSingle(c *gin.Context) {
	itemID := c.Param("itemId")
	count, err := h.useCase.Get(c.Request.Context(), itemID)
	if err == domain.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, count)
}

func (h *CountValueHandler) Delete(c *gin.Context) {
	itemID := c.Param("itemId")
	err := h.useCase.Delete(c.Request.Context(), itemID)
	if err == domain.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "count value not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})
}

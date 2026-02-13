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

	external := r.Group("/api/v1/counts")
	{
		external.GET("/values", h.GetAllExternal)
		external.GET("/:itemId/value", h.GetSingleExternal)
		external.POST("/:itemId/increase", h.Increase)
		external.POST("/:itemId/decrease", h.Decrease)
		external.POST("/:itemId/reset", h.Reset)
	}
}

func (h *CountValueHandler) GetSingleExternal(c *gin.Context) {
	itemID := c.Param("itemId")
	count, err := h.useCase.Get(c.Request.Context(), itemID)
	if err == domain.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"message": "Count item " + itemID + " not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"itemId":       count.ItemID,
		"currentValue": count.CurrentValue,
	})
}

func (h *CountValueHandler) GetAllExternal(c *gin.Context) {
	counts, err := h.useCase.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	type externalCountValue struct {
		ItemID       string `json:"itemId"`
		CurrentValue int    `json:"currentValue"`
	}
	var resp []externalCountValue
	for _, count := range counts {
		resp = append(resp, externalCountValue{
			ItemID:       count.ItemID,
			CurrentValue: count.CurrentValue,
		})
	}

	c.JSON(http.StatusOK, resp)
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

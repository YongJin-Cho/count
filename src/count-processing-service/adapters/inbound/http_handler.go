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
	group := r.Group("/api/v1/internal/counts")
	{
		group.POST("", h.Initialize)
		group.GET("", h.GetMultiple)
		group.GET("/:itemId", h.GetSingle)
		group.DELETE("/:itemId", h.Delete)
	}
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

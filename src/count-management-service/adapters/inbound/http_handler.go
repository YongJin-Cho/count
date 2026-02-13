package inbound

import (
	"count-management-service/domain"
	"count-management-service/ports"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HTTPHandler struct {
	service ports.CountItemService
}

func NewHTTPHandler(service ports.CountItemService) *HTTPHandler {
	return &HTTPHandler{service: service}
}

func (h *HTTPHandler) mapDomainError(err error) int {
	switch err {
	case domain.ErrEmptyName:
		return http.StatusBadRequest
	case domain.ErrDuplicateName:
		return http.StatusConflict
	case domain.ErrItemNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// UI API (HTMX)

func (h *HTTPHandler) RegisterItemUI(c *gin.Context) {
	name := c.PostForm("name")
	description := c.PostForm("description")

	item, err := h.service.RegisterItem(c.Request.Context(), name, description)
	if err != nil {
		status := h.mapDomainError(err)
		c.HTML(status, "error_message.html", gin.H{"Message": err.Error()})
		return
	}

	c.HTML(http.StatusCreated, "count_item_row.html", item)
}

func (h *HTTPHandler) ListItemsUI(c *gin.Context) {
	items, err := h.service.ListItemWithValues(c.Request.Context())
	if err != nil {
		status := h.mapDomainError(err)
		c.HTML(status, "error_message.html", gin.H{"Message": err.Error()})
		return
	}

	if len(items) == 0 {
		c.String(http.StatusOK, "<tr><td colspan='4'>No items found.</td></tr>")
		return
	}

	for _, item := range items {
		c.HTML(http.StatusOK, "count_item_row.html", item)
	}
}

func (h *HTTPHandler) GetItemValueUI(c *gin.Context) {
	id := c.Param("id")
	value, err := h.service.GetItemValue(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrItemNotFound {
			c.String(http.StatusNotFound, "<span class=\"text-red-500\">Item not found</span>")
			return
		}
		c.String(http.StatusInternalServerError, "<span class=\"text-red-500\">Error</span>")
		return
	}

	c.String(http.StatusOK, "%d", value)
}

func (h *HTTPHandler) UpdateItemUI(c *gin.Context) {
	id := c.Param("id")
	name := c.PostForm("name")
	description := c.PostForm("description")

	_, err := h.service.UpdateItem(c.Request.Context(), id, name, description)
	if err != nil {
		status := h.mapDomainError(err)
		c.HTML(status, "error_message.html", gin.H{"Message": err.Error()})
		return
	}

	c.Header("HX-Redirect", "/ui/count-items")
	c.Status(http.StatusOK)
}

func (h *HTTPHandler) DeleteItemUI(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeleteItem(c.Request.Context(), id)
	if err != nil {
		status := h.mapDomainError(err)
		c.HTML(status, "error_message.html", gin.H{"Message": err.Error()})
		return
	}

	c.Status(http.StatusOK) // HTMX will remove the row if target is the row and swap is outerHTML
}

// External API (JSON)

func (h *HTTPHandler) ListItemsAPI(c *gin.Context) {
	items, err := h.service.ListItem(c.Request.Context())
	if err != nil {
		status := h.mapDomainError(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *HTTPHandler) RegisterItemAPI(c *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.service.RegisterItem(c.Request.Context(), req.Name, req.Description)
	if err != nil {
		status := h.mapDomainError(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *HTTPHandler) UpdateItemAPI(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.service.UpdateItem(c.Request.Context(), id, req.Name, req.Description)
	if err != nil {
		status := h.mapDomainError(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *HTTPHandler) DeleteItemAPI(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeleteItem(c.Request.Context(), id)
	if err != nil {
		status := h.mapDomainError(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

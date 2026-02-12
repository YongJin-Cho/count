package collector

import (
	"count-api-service/internal/common/auth"
	"count-api-service/internal/common/event"
	"count-api-service/internal/common/model"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Repository interface {
	FindAll(filter string, limit int, offset int) ([]model.CountItem, error)
	FindById(id string) (*model.CountItem, error)
	CountTotal(filter string) (int, error)
	Create(item model.CountItem) error
	UpdateValue(id string, value int) error
}

type CollectorHandler struct {
	authProvider *auth.AuthProvider
	publisher    event.Publisher
	repo         Repository
}

func NewCollectorHandler(ap *auth.AuthProvider, pub event.Publisher, repo Repository) *CollectorHandler {
	return &CollectorHandler{
		authProvider: ap,
		publisher:    pub,
		repo:         repo,
	}
}

func (h *CollectorHandler) CollectCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Authentication
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	valid, err := h.authProvider.ValidateToken(token)
	if err != nil || !valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	authorized, err := h.authProvider.IsAuthorized(token, "collect")
	if err != nil || !authorized {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Forbidden"})
		return
	}

	// 2. Parsing
	var req model.CountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Validation failed", "message": "Invalid JSON"})
		return
	}

	// 3. Validation
	if err := req.Validate(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Validation failed", "message": err.Error()})
		return
	}

	// 4. Event Emission
	h.publisher.Publish(event.CountCollectedEvent{
		ExternalID: req.ExternalID,
		Count:      *req.Count,
		Timestamp:  time.Now().Format(time.RFC3339),
	})

	// 5. Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *CollectorHandler) GetCounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Authentication
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	valid, err := h.authProvider.ValidateToken(token)
	if err != nil || !valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	authorized, err := h.authProvider.IsAuthorized(token, "query")
	if err != nil || !authorized {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Forbidden"})
		return
	}

	// 2. Parameter Parsing
	query := r.URL.Query()
	externalID := query.Get("external_id")

	limit := 10
	if l := query.Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val >= 1 {
			limit = val
		}
	}

	offset := 0
	if o := query.Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			offset = val
		}
	}

	// 3. Data Retrieval
	total, err := h.repo.CountTotal(externalID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	counts, err := h.repo.FindAll(externalID, limit, offset)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	// 4. Response Construction
	resp := model.CountListResponse{
		TotalCount: total,
		Counts:     counts,
	}
	if resp.Counts == nil {
		resp.Counts = []model.CountItem{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Internal Management Methods

func (h *CollectorHandler) CreateCount(sourceID string, initialValue int) error {
	// Validate source_id format (^[a-z0-9-]+$)
	re := regexp.MustCompile("^[a-z0-9-]+$")
	if !re.MatchString(sourceID) {
		return fmt.Errorf("invalid source_id format: %s", sourceID)
	}

	// Check if already exists
	_, err := h.repo.FindById(sourceID)
	if err == nil {
		return fmt.Errorf("source_id already exists: %s", sourceID)
	}

	return h.repo.Create(model.CountItem{
		ExternalID: sourceID,
		Count:      initialValue,
	})
}

func (h *CollectorHandler) IncrementCount(sourceID string) (*model.CountItem, error) {
	item, err := h.repo.FindById(sourceID)
	if err != nil {
		return nil, err
	}

	newValue := item.Count + 1
	err = h.repo.UpdateValue(sourceID, newValue)
	if err != nil {
		return nil, err
	}

	return h.repo.FindById(sourceID)
}

func (h *CollectorHandler) DecrementCount(sourceID string) (*model.CountItem, error) {
	item, err := h.repo.FindById(sourceID)
	if err != nil {
		return nil, err
	}

	newValue := item.Count - 1
	if newValue < 0 {
		newValue = 0
	}

	err = h.repo.UpdateValue(sourceID, newValue)
	if err != nil {
		return nil, err
	}

	return h.repo.FindById(sourceID)
}

func (h *CollectorHandler) UpdateCount(sourceID string, value int) error {
	if value < 0 {
		return fmt.Errorf("invalid value: %d", value)
	}
	return h.repo.UpdateValue(sourceID, value)
}

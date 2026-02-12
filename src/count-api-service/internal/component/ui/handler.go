package ui

import (
	"count-api-service/internal/common/auth"
	"count-api-service/internal/common/model"
	"count-api-service/internal/component/collector"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type Storage interface {
	FindAll(filter string, limit int, offset int) ([]model.CountItem, error)
	FindById(id string) (*model.CountItem, error)
}

type UIHandler struct {
	collector *collector.CollectorHandler
	storage   Storage
	templates *template.Template
	auth      *auth.AuthProvider
}

func NewUIHandler(coll *collector.CollectorHandler, stor Storage, templateDir string, ap *auth.AuthProvider) (*UIHandler, error) {
	tmpl, err := template.ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &UIHandler{
		collector: coll,
		storage:   stor,
		templates: tmpl,
		auth:      ap,
	}, nil
}

func (h *UIHandler) authenticate(r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return fmt.Errorf("Unauthorized")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	valid, err := h.auth.ValidateToken(token)
	if err != nil || !valid {
		return fmt.Errorf("Unauthorized")
	}
	return nil
}

func (h *UIHandler) GetCountList(w http.ResponseWriter, r *http.Request) {
	if err := h.authenticate(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	counts, err := h.storage.FindAll("", 100, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, count := range counts {
		if err := h.templates.ExecuteTemplate(w, "count_row.html", count); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *UIHandler) GetCreateForm(w http.ResponseWriter, r *http.Request) {
	if err := h.authenticate(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err := h.templates.ExecuteTemplate(w, "create_form.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UIHandler) CreateCount(w http.ResponseWriter, r *http.Request) {
	if err := h.authenticate(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	sourceID := r.FormValue("source_id")
	initialValueStr := r.FormValue("initial_value")
	initialValue, _ := strconv.Atoi(initialValueStr)

	err := h.collector.CreateCount(sourceID, initialValue)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Header().Set("HX-Retarget", "#count-create-error-msg")
		fmt.Fprintf(w, `<div class="text-red-500">%s</div>`, err.Error())
		return
	}

	item, err := h.storage.FindById(sourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "count-list-changed")
	if err := h.templates.ExecuteTemplate(w, "count_row.html", item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UIHandler) GetCountRow(w http.ResponseWriter, r *http.Request) {
	if err := h.authenticate(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	sourceID := strings.TrimPrefix(r.URL.Path, "/ui/counts/")
	item, err := h.storage.FindById(sourceID)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "count_row.html", item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UIHandler) GetEditForm(w http.ResponseWriter, r *http.Request) {
	if err := h.authenticate(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	sourceID := strings.TrimPrefix(r.URL.Path, "/ui/counts/")
	sourceID = strings.TrimSuffix(sourceID, "/edit")

	item, err := h.storage.FindById(sourceID)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "edit_row.html", item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UIHandler) UpdateCount(w http.ResponseWriter, r *http.Request) {
	if err := h.authenticate(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	sourceID := strings.TrimPrefix(r.URL.Path, "/ui/counts/")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	valueStr := r.FormValue("value")
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		http.Error(w, "Invalid value", http.StatusBadRequest)
		return
	}

	err = h.collector.UpdateCount(sourceID, value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item, err := h.storage.FindById(sourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "count_row.html", item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UIHandler) IncrementCount(w http.ResponseWriter, r *http.Request) {
	if err := h.authenticate(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	sourceID := strings.TrimPrefix(r.URL.Path, "/ui/counts/")
	sourceID = strings.TrimSuffix(sourceID, "/increment")

	item, err := h.collector.IncrementCount(sourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "count_row.html", item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UIHandler) DecrementCount(w http.ResponseWriter, r *http.Request) {
	if err := h.authenticate(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	sourceID := strings.TrimPrefix(r.URL.Path, "/ui/counts/")
	sourceID = strings.TrimSuffix(sourceID, "/decrement")

	item, err := h.collector.DecrementCount(sourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "count_row.html", item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

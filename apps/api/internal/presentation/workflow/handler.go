// internal/handler/workflow_definition_handler.go
package workflow

import (
	"encoding/json"
	"net/http"
	"strings"

	"seo-backend/internal/domain/workflow"
)

type WorkflowDefinitionHandler struct {
	service workflow.WorkflowDefinitionService
}

func NewWorkflowDefinitionHandler(service workflow.WorkflowDefinitionService) *WorkflowDefinitionHandler {
	return &WorkflowDefinitionHandler{service: service}
}

func extractID(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	return parts[len(parts)-1]
}

func (h *WorkflowDefinitionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)

	wf, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wf)
}

func (h *WorkflowDefinitionHandler) GetByProductID(w http.ResponseWriter, r *http.Request) {
	productID := extractID(r.URL.Path)

	workflows, err := h.service.GetByProductID(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflows)
}

func (h *WorkflowDefinitionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var wf workflow.WorkflowDefinition
	if err := json.NewDecoder(r.Body).Decode(&wf); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &wf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(wf)
}

func (h *WorkflowDefinitionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Update(r.Context(), id, updates); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "updated"})
}

func (h *WorkflowDefinitionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)

	if err := h.service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "deleted"})
}

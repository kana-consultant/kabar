// internal/handler/workflow/workflow_node_handler.go
package workflow_node

import (
	"encoding/json"
	"net/http"

	"seo-backend/internal/domain/workflow_node"

	"github.com/go-chi/chi/v5"
)

type WorkflowNodeHandler struct {
	service workflow_node.WorkflowNodeService
}

func NewWorkflowNodeHandler(service workflow_node.WorkflowNodeService) *WorkflowNodeHandler {
	return &WorkflowNodeHandler{service: service}
}

func (h *WorkflowNodeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	node, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(node)
}

func (h *WorkflowNodeHandler) GetByWorkflowID(w http.ResponseWriter, r *http.Request) {
	workflowID := chi.URLParam(r, "id")

	nodes, err := h.service.GetByWorkflowID(r.Context(), workflowID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}

func (h *WorkflowNodeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var node workflow_node.WorkflowNode
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &node); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(node)
}

func (h *WorkflowNodeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

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

func (h *WorkflowNodeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "deleted"})
}

func (h *WorkflowNodeHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	workflowID := chi.URLParam(r, "id")

	var nodeIDs []string
	if err := json.NewDecoder(r.Body).Decode(&nodeIDs); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.ReorderNodes(r.Context(), workflowID, nodeIDs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "reordered"})
}

package product

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"seo-backend/internal/domain/product"
	auth "seo-backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	service product.ProductService
}

func NewProductHandler(service product.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("encode error: %v", err)
	}
}

func (h *ProductHandler) writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch err.Error() {
	case "product not found":
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})

	case "product id is required":
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Product ID is required"})

	case "no updates provided":
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "No updates provided"})

	case "product name is required":
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Product name is required"})

	case "platform is required":
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Platform is required"})

	case "API endpoint is required":
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "API endpoint is required"})

	case "API key is required":
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "API key is required"})

	case "cannot update status of active product":
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cannot update status of active product"})

	case "cannot update platform of active product":
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cannot update platform of active product"})

	case "cannot delete active product, deactivate it first":
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cannot delete active product, deactivate it first"})

	case "cannot disconnect active product":
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cannot disconnect active product"})

	default:
		log.Printf("unexpected error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
	}
}

// =======================
// CREATE PRODUCT
// =======================
// @Summary Create new product
// @Tags products
// @Accept json
// @Produce json
// @Param request body product.ProductRequest true "Product data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [post]
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Println("========== HTTP CREATE PRODUCT ==========")
	log.Printf("Method: %s | URL: %s\n", r.Method, r.URL.Path)

	// Parse request body
	var req product.ProductRequest

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed read body: %v\n", err)

		h.writeJSON(w, map[string]string{
			"error": "Failed read request body",
		}, http.StatusBadRequest)

		return
	}

	log.Printf("RAW BODY: %s\n", string(bodyBytes))

	// balikin lagi body ke reader
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	log.Println("Decoding request body...")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed decode request body: %v\n", err)
		h.writeJSON(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
		return
	}

	// FULL REQUEST
	log.Printf("FULL REQUEST: %+v\n", req)

	// DETAIL PRODUCT
	log.Println("========== PRODUCT ==========")
	log.Printf("Name: %v\n", req.Name)
	log.Printf("APIEndpoint: %v\n", req.APIEndpoint)
	log.Printf("APIKey: %v\n", req.APIKey)

	// DETAIL ADAPTER CONFIG
	log.Println("========== ADAPTER CONFIG ==========")

	if req.AdapterConfig != nil {
		log.Printf("EndpointPath: %v\n", req.AdapterConfig.EndpointPath)
		log.Printf("HTTPMethod: %v\n", req.AdapterConfig.HTTPMethod)
		log.Printf("CustomHeaders: %+v\n", req.AdapterConfig.CustomHeaders)
		log.Printf("FieldMapping: %v\n", req.AdapterConfig.FieldMapping)
		log.Printf("ResponseMapping: %+v\n", req.AdapterConfig.ResponseMapping)
		log.Printf("MetaConfig: %v\n", req.AdapterConfig.MetaConfig)
		log.Printf("SitemapConfig: %v\n", req.AdapterConfig.SitemapConfig)
		log.Printf("TimeoutSeconds: %v\n", req.AdapterConfig.TimeoutSeconds)
		log.Printf("RetryCount: %v\n", req.AdapterConfig.RetryCount)
	} else {
		log.Println("AdapterConfig is NIL")
	}

	// DEBUG: TAMBAHKAN INI UNTUK ADAPTERCONFIGS
	log.Println("========== ADAPTERCONFIGS (ARRAY) ==========")

	// DEBUG: TAMPILKAN JUGA WORKFLOWS
	log.Println("========== WORKFLOWS ==========")
	if len(req.Workflows) > 0 {
		log.Printf("Number of workflows: %d\n", len(req.Workflows))
		for i, wf := range req.Workflows {
			log.Printf("Workflow[%d]: ID=%s, Name=%s\n", i, wf.ID, wf.Name)
		}
	} else {
		log.Println("Workflows is EMPTY or NIL")
	}

	// Validasi request
	if req.Name == "" {
		h.writeJSON(w, map[string]string{"error": "Product name is required"}, http.StatusBadRequest)
		return
	}

	if req.APIEndpoint == "" {
		h.writeJSON(w, map[string]string{"error": "API endpoint is required"}, http.StatusBadRequest)
		return
	}

	if req.APIKey == "" {
		h.writeJSON(w, map[string]string{"error": "API key is required"}, http.StatusBadRequest)
		return
	}

	// Get user context from auth middleware
	log.Println("Getting user context from middleware...")

	userCtx := auth.GetUserContext(r)

	log.Printf(
		"User Context => UserID: %s | TeamID: %s | Role: %s\n",
		userCtx.GetUserID(),
		userCtx.GetTeamID(),
		userCtx.GetRole(),
	)

	if userCtx.GetUserID() == "" {
		h.writeJSON(w, map[string]string{"error": "Unauthorized: User ID not found"}, http.StatusUnauthorized)
		return
	}

	// Call application service
	log.Println("Calling ProductService.CreateProduct...")

	productID, err := h.service.CreateProduct(ctx, req, userCtx)
	if err != nil {
		log.Printf("CreateProduct failed: %v\n", err)
		h.writeError(w, err)
		return
	}

	log.Printf("Product created successfully. ProductID: %s\n", productID)

	// Return response
	response := map[string]string{
		"id":      productID,
		"message": "Product created successfully",
	}

	log.Printf("Response: %+v\n", response)

	h.writeJSON(w, response, http.StatusCreated)

	log.Println("========== END CREATE PRODUCT ==========")
}

// =======================
// GET PRODUCT BY ID
// =======================
// @Summary Get product by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} product.Product
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [get]
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	useCtx := auth.GetUserContext(r)

	// Call application service
	product, err := h.service.GetByID(ctx, id, useCtx)
	if err != nil {
		h.writeError(w, err)
		return
	}

	if product == nil {
		h.writeJSON(w, map[string]string{"error": "Product not found"}, http.StatusNotFound)
		return
	}

	h.writeJSON(w, product, http.StatusOK)
}

// =======================
// GET ALL PRODUCTS
// =======================
// @Summary Get all products
// @Tags products
// @Produce json
// @Success 200 {array} product.Product
// @Failure 500 {object} map[string]string
// @Router /products [get]
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)

	// Call application service (masih perlu implementasi di service)
	products, err := h.service.GetAllProducts(ctx, userCtx)
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, products, http.StatusOK)
}

// =======================
// UPDATE PRODUCT
// =======================
// @Summary Update product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body map[string]interface{} true "Update fields"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [put]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	// Parse request body ke UpdateProductRequest struct
	var req product.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSON(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
		return
	}

	// Get user context
	userCtx := auth.GetUserContext(r)

	// Call application service
	if err := h.service.UpdateProduct(ctx, id, req, userCtx); err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, map[string]string{
		"id":      id,
		"message": "Product updated successfully",
	}, http.StatusOK)
}

// =======================
// DELETE PRODUCT
// =======================
// @Summary Delete product
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	// Call application service
	if err := h.service.DeleteProduct(ctx, id); err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// =======================
// UPDATE CONNECTION STATUS
// =======================
// @Summary Update product connection status
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body object true "Connection status" example({"isConnected": true})
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id}/connection [patch]
func (h *ProductHandler) UpdateConnectionStatus(
	w http.ResponseWriter,
	r *http.Request,
) {

	log.Println("========== UPDATE CONNECTION STATUS ==========")

	ctx := r.Context()

	id := chi.URLParam(r, "id")

	log.Println("[INFO] Product ID:", id)
	log.Println("[INFO] Method:", r.Method)
	log.Println("[INFO] URL:", r.URL.Path)
	log.Println("[INFO] Remote Address:", r.RemoteAddr)

	// =========================
	// PARSE REQUEST BODY
	// =========================
	var req struct {
		IsConnected bool `json:"isConnected"`
	}

	if err := h.service.UpdateConnectionStatus(
		ctx,
		id,
		req.IsConnected,
	); err != nil {

		log.Println("[ERROR] Failed to update connection status:", err)

		h.writeError(w, err)

		return
	}

	log.Println("[SUCCESS] Connection status updated successfully")

	// =========================
	// RESPONSE
	// =========================
	response := map[string]string{
		"id":      id,
		"message": "Connection status updated successfully",
	}

	log.Println("[INFO] Response:")
	log.Printf("%#v\n", response)

	h.writeJSON(
		w,
		response,
		http.StatusOK,
	)
}

// =======================
// HELPER FUNCTIONS FOR QUERY BUILDING
// =======================
func buildQueryFromRequest(r *http.Request) string {
	// Build dynamic query based on query parameters
	// Example: ?status=active&platform=shopify
	query := "SELECT * FROM products WHERE 1=1"

	if status := r.URL.Query().Get("status"); status != "" {
		query += " AND status = '" + status + "'"
	}

	if platform := r.URL.Query().Get("platform"); platform != "" {
		query += " AND platform = '" + platform + "'"
	}

	// Add pagination
	if limit := r.URL.Query().Get("limit"); limit != "" {
		query += " LIMIT " + limit
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		query += " OFFSET " + offset
	}

	return query
}

func buildArgsFromRequest(r *http.Request) []interface{} {
	// Build args for parameterized query
	args := make([]interface{}, 0)

	if status := r.URL.Query().Get("status"); status != "" {
		args = append(args, status)
	}

	if platform := r.URL.Query().Get("platform"); platform != "" {
		args = append(args, platform)
	}

	return args
}

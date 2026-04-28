package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"seo-backend/internal/builder"
	"seo-backend/internal/database"
	"seo-backend/internal/middleware/auth"
	"seo-backend/internal/models"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	teamID := auth.GetTeamID(ctx)
	userRole := auth.GetUserRole(ctx)
	userID := auth.GetUserID(ctx)

	log.Printf("GetAll products - role: %s, teamID: %s, userID: %s", userRole, teamID, userID)

	// Build query dengan generic builder
	qb := builder.NewQueryBuilder("products")

	query := qb.Select(
		"id", "name", "platform", "api_endpoint", "status", "sync_status", "last_sync",
		"created_by", "team_id", "user_id", "created_at", "updated_at",
	).OrderBy("created_at DESC")

	// Logic filter berdasarkan role
	switch userRole {
	case "super_admin":
		// Super admin: melihat SEMUA product
		log.Printf("Super admin - melihat semua product")

	case "admin":
		// Admin: hanya melihat product dalam team yang sama
		if teamID != "" {
			query = query.WhereEq("team_id", teamID)
			log.Printf("Admin - melihat product dalam team: %s", teamID)
		} else {
			// Admin tanpa team hanya melihat product miliknya sendiri
			query = query.WhereEq("user_id", userID)
			log.Printf("Admin tanpa team - hanya melihat product sendiri")
		}

	default:
		// Role lain (manager, editor, viewer): hanya melihat product milik sendiri
		query = query.WhereEq("user_id", userID)
		log.Printf("User role %s - hanya melihat product sendiri", userRole)
	}

	// Filter tambahan
	platform := r.URL.Query().Get("platform")
	if platform != "" {
		query = query.WhereEq("platform", platform)
	}

	status := r.URL.Query().Get("status")
	if status != "" {
		query = query.WhereEq("status", status)
	}

	syncStatus := r.URL.Query().Get("sync_status")
	if syncStatus != "" {
		query = query.WhereEq("sync_status", syncStatus)
	}

	sqlQuery, args, err := query.Build()
	if err != nil {
		log.Printf("Failed to build query: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch products: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Query: %s | Args: %v", sqlQuery, args)

	var rows *sql.Rows
	rows, err = database.GetDB().Query(sqlQuery, args...)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch products: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Platform, &p.APIEndpoint, &p.Status, &p.SyncStatus,
			&p.LastSync, &p.CreatedBy, &p.TeamID, &p.UserID, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		// Get adapter config
		configQb := builder.NewQueryBuilder("adapter_configs")
		configQuery, configArgs, err := configQb.Select(
			"id", "product_id", "endpoint_path", "http_method", "custom_headers", "field_mapping",
			"timeout_seconds", "retry_count", "created_at", "updated_at",
		).WhereEq("product_id", p.ID).Build()

		if err == nil {
			var adapterConfig models.AdapterConfig
			var customHeadersJSON []byte
			err = database.GetDB().QueryRow(configQuery, configArgs...).Scan(
				&adapterConfig.ID, &adapterConfig.ProductID, &adapterConfig.EndpointPath,
				&adapterConfig.HTTPMethod, &customHeadersJSON, &adapterConfig.FieldMapping,
				&adapterConfig.TimeoutSeconds, &adapterConfig.RetryCount,
				&adapterConfig.CreatedAt, &adapterConfig.UpdatedAt)
			if err == nil {
				json.Unmarshal(customHeadersJSON, &adapterConfig.CustomHeaders)
				p.AdapterConfig = &adapterConfig
			}
		}

		products = append(products, p)
	}

	log.Printf("Returning %d products", len(products))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teamID := auth.GetTeamID(ctx)
	userRole := auth.GetUserRole(ctx)
	userID := auth.GetUserID(ctx)

	productID := chi.URLParam(r, "id")

	log.Printf("[PRODUCT GET] id=%s user=%s role=%s team=%s",
		productID, userID, userRole, teamID,
	)

	// =========================
	// BUILD QUERY
	// =========================
	qb := builder.NewQueryBuilder("products")

	sqlQuery, args, err := qb.Select(
		"id", "name", "platform", "api_key_encrypted", "api_endpoint",
		"status", "sync_status", "last_sync",
		"created_by", "team_id", "user_id",
		"created_at", "updated_at",
	).WhereEq("id", productID).Build()

	if err != nil {
		log.Printf("[PRODUCT QUERY BUILD ERROR] %v", err)
		http.Error(w, "Failed to build query", http.StatusInternalServerError)
		return
	}

	log.Printf("[PRODUCT SQL] %s", sqlQuery)
	log.Printf("[PRODUCT ARGS] %#v", args)

	// =========================
	// SCAN
	// =========================
	var product models.Product

	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(
		&product.ID,
		&product.Name,
		&product.Platform,
		&product.APIKeyEncrypted,
		&product.APIEndpoint,
		&product.Status,
		&product.SyncStatus,
		&product.LastSync,
		&product.CreatedBy,
		&product.TeamID,
		&product.UserID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		log.Printf("[PRODUCT SCAN ERROR] %v", err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	log.Printf("[PRODUCT RAW RESULT] %+v", product)
	log.Printf("[PRODUCT API KEY RAW] %+v", product.APIKeyEncrypted)

	// =========================
	// AUTH CHECK
	// =========================
	switch userRole {

	case "super_admin":
		log.Printf("[AUTH] super_admin access granted product=%s", productID)

	case "admin":
		if product.TeamID != nil && teamID != "" && *product.TeamID != teamID {
			log.Printf("[AUTH DENIED] admin team mismatch product_team=%v user_team=%s",
				product.TeamID, teamID,
			)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		log.Printf("[AUTH] admin access granted product=%s", productID)

	default:
		if product.UserID == nil || *product.UserID != userID {
			log.Printf("[AUTH DENIED] user mismatch product_user=%v user=%s",
				product.UserID, userID,
			)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		log.Printf("[AUTH] user access granted product=%s", productID)
	}

	// =========================
	// ADAPTER CONFIG
	// =========================
	configQb := builder.NewQueryBuilder("adapter_configs")

	configQuery, configArgs, err := configQb.Select(
		"id", "product_id", "endpoint_path", "http_method",
		"custom_headers", "field_mapping",
		"timeout_seconds", "retry_count",
		"created_at", "updated_at",
	).WhereEq("product_id", productID).Build()

	log.Printf("[ADAPTER SQL] %s", configQuery)

	if err == nil {
		var adapterConfig models.AdapterConfig
		var customHeadersJSON []byte

		err = database.GetDB().QueryRow(configQuery, configArgs...).Scan(
			&adapterConfig.ID,
			&adapterConfig.ProductID,
			&adapterConfig.EndpointPath,
			&adapterConfig.HTTPMethod,
			&customHeadersJSON,
			&adapterConfig.FieldMapping,
			&adapterConfig.TimeoutSeconds,
			&adapterConfig.RetryCount,
			&adapterConfig.CreatedAt,
			&adapterConfig.UpdatedAt,
		)

		if err != nil {
			log.Printf("[ADAPTER SCAN ERROR] %v", err)
		} else {
			err = json.Unmarshal(customHeadersJSON, &adapterConfig.CustomHeaders)
			if err != nil {
				log.Printf("[ADAPTER JSON UNMARSHAL ERROR] %v", err)
			}

			product.AdapterConfig = &adapterConfig
			log.Printf("[ADAPTER SUCCESS LOADED] product_id=%s", productID)
		}
	}

	// =========================
	// RESPONSE
	// =========================
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// Create product
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Println("Create product called")

	var req models.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Request: %+v", req)

	// Validate required fields
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	if req.Platform == "" {
		http.Error(w, "Platform is required", http.StatusBadRequest)
		return
	}
	if req.APIEndpoint == "" {
		http.Error(w, "API Endpoint is required", http.StatusBadRequest)
		return
	}

	tx, err := database.GetDB().Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Get user ID from context
	userID := auth.GetUserID(r.Context())
	if userID == "" {
		userID = "00000000-0000-0000-0000-000000000000"
	}
	teamID := auth.GetTeamID(r.Context())
	log.Printf("UserID: %s, TeamID: %s", userID, teamID)

	var teamIDPtr interface{}
	if teamID != "" {
		teamIDPtr = teamID
	} else {
		teamIDPtr = nil
	}

	// Insert product using builder
	qb := builder.NewQueryBuilder("products")
	sqlQuery, args, err := qb.Insert().
		Columns("name", "platform", "api_endpoint", "api_key_encrypted", "status", "sync_status", "created_by", "team_id", "user_id").
		Values(req.Name, req.Platform, req.APIEndpoint, req.APIKey, "pending", "idle", userID, teamIDPtr, userID).
		Returning("id").
		Build()

	if err != nil {
		log.Printf("Failed to build insert query: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create product: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Insert query: %s | Args: %v", sqlQuery, args)

	var productID string
	err = tx.QueryRow(sqlQuery, args...).Scan(&productID)
	if err != nil {
		log.Printf("Failed to insert product: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create product: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Product created with ID: %s", productID)

	// Insert adapter config if provided
	if req.AdapterConfig != nil {
		customHeadersJSON, _ := json.Marshal(req.AdapterConfig.CustomHeaders)
		timeoutSeconds := req.AdapterConfig.TimeoutSeconds
		if timeoutSeconds == 0 {
			timeoutSeconds = 30
		}
		retryCount := req.AdapterConfig.RetryCount
		if retryCount == 0 {
			retryCount = 3
		}

		configQb := builder.NewQueryBuilder("adapter_configs")
		configQuery, configArgs, err := configQb.Insert().
			Columns("product_id", "endpoint_path", "http_method", "custom_headers", "field_mapping", "timeout_seconds", "retry_count").
			Values(productID, req.AdapterConfig.EndpointPath, req.AdapterConfig.HTTPMethod, customHeadersJSON,
				req.AdapterConfig.FieldMapping, timeoutSeconds, retryCount).
			Build()

		if err != nil {
			log.Printf("Failed to build adapter config insert query: %v", err)
			http.Error(w, "Failed to create adapter config", http.StatusInternalServerError)
			return
		}

		log.Printf("Adapter config insert query: %s | Args: %v", configQuery, configArgs)

		_, err = tx.Exec(configQuery, configArgs...)
		if err != nil {
			log.Printf("Failed to insert adapter config: %v", err)
			http.Error(w, "Failed to create adapter config", http.StatusInternalServerError)
			return
		}
		log.Println("Adapter config created")
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// Get created product
	productQb := builder.NewQueryBuilder("products")
	productQuery, productArgs, err := productQb.Select(
		"id", "name", "platform", "api_endpoint", "status", "sync_status", "created_at", "updated_at",
	).WhereEq("id", productID).Build()

	if err == nil {
		var product models.Product
		err = database.GetDB().QueryRow(productQuery, productArgs...).Scan(
			&product.ID, &product.Name, &product.Platform, &product.APIEndpoint,
			&product.Status, &product.SyncStatus, &product.CreatedAt, &product.UpdatedAt)
		if err == nil {
			log.Printf("Product creation successful: %s", product.Name)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(product)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": productID, "message": "Product created successfully"})
}

// Update product
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	tx, err := database.GetDB().Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Update products table using builder
	qb := builder.NewQueryBuilder("products")
	updateBuilder := qb.Update()

	// Set fields
	if name, ok := updates["name"]; ok {
		updateBuilder = updateBuilder.Set("name", name)
	}
	if platform, ok := updates["platform"]; ok {
		updateBuilder = updateBuilder.Set("platform", platform)
	}
	if apiEndpoint, ok := updates["apiEndpoint"]; ok {
		updateBuilder = updateBuilder.Set("api_endpoint", apiEndpoint)
	}

	if status, ok := updates["status"]; ok {
		updateBuilder = updateBuilder.Set("status", status)
	}
	if apiKey, ok := updates["apiKey"]; ok {
		updateBuilder = updateBuilder.Set("api_key_encrypted", apiKey)
	}

	updateBuilder = updateBuilder.Set("updated_at", "NOW()")
	updateBuilder = updateBuilder.WhereEq("id", id)

	sqlQuery, args, err := updateBuilder.Build()
	if err != nil {
		log.Printf("Failed to build update query: %v", err)
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	log.Printf("Update query: %s | Args: %v", sqlQuery, args)

	if _, err := tx.Exec(sqlQuery, args...); err != nil {
		log.Printf("Failed to update product: %v", err)
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	// Update adapter_configs table
	if adapterUpdates, ok := updates["adapterConfig"].(map[string]interface{}); ok {
		configQb := builder.NewQueryBuilder("adapter_configs")
		configUpdateBuilder := configQb.Update()

		if endpointPath, ok := adapterUpdates["endpointPath"]; ok {
			configUpdateBuilder = configUpdateBuilder.Set("endpoint_path", endpointPath)
		}
		if httpMethod, ok := adapterUpdates["httpMethod"]; ok {
			configUpdateBuilder = configUpdateBuilder.Set("http_method", httpMethod)
		}
		if customHeaders, ok := adapterUpdates["customHeaders"]; ok {
			customHeadersJSON, _ := json.Marshal(customHeaders)
			configUpdateBuilder = configUpdateBuilder.Set("custom_headers", customHeadersJSON)
		}
		if fieldMapping, ok := adapterUpdates["fieldMapping"]; ok {
			configUpdateBuilder = configUpdateBuilder.Set("field_mapping", fieldMapping)
		}
		if timeoutSeconds, ok := adapterUpdates["timeoutSeconds"]; ok {
			configUpdateBuilder = configUpdateBuilder.Set("timeout_seconds", timeoutSeconds)
		}
		if retryCount, ok := adapterUpdates["retryCount"]; ok {
			configUpdateBuilder = configUpdateBuilder.Set("retry_count", retryCount)
		}

		configUpdateBuilder = configUpdateBuilder.Set("updated_at", "NOW()")
		configUpdateBuilder = configUpdateBuilder.WhereEq("product_id", id)

		configQuery, configArgs, err := configUpdateBuilder.Build()
		if err != nil {
			log.Printf("Failed to build adapter config update query: %v", err)
			http.Error(w, "Failed to update adapter config", http.StatusInternalServerError)
			return
		}

		log.Printf("Adapter config update query: %s | Args: %v", configQuery, configArgs)

		if _, err := tx.Exec(configQuery, configArgs...); err != nil {
			log.Printf("Failed to update adapter config: %v", err)
			http.Error(w, "Failed to update adapter config", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Product updated successfully"})
}

// Delete product
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	qb := builder.NewQueryBuilder("products")
	sqlQuery, args, err := qb.Delete().WhereEq("id", id).Build()
	if err != nil {
		log.Printf("Failed to build delete query: %v", err)
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	log.Printf("Delete query: %s | Args: %v", sqlQuery, args)

	_, err = database.GetDB().Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("Failed to delete product: %v", err)
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) TestConnection(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// 1. Ambil data product
	var productName, platform, apiEndpoint, apiKey string

	err := database.GetDB().QueryRow(`
		SELECT name, platform, api_endpoint, COALESCE(api_key_encrypted, '')
		FROM products WHERE id = $1
	`, id).Scan(&productName, &platform, &apiEndpoint, &apiKey)

	if err != nil {
		responseError(w, "Product not found", http.StatusNotFound)
		return
	}

	// 2. Default adapter config
	var (
		endpointPath      string
		httpMethod        string = "GET"
		customHeadersJSON []byte
		timeoutSeconds    int = 10
		retryCount        int = 1
		customHeaders     map[string]string
	)

	// 3. Ambil adapter_configs HANYA jika API key ada
	if apiKey != "" {
		err = database.GetDB().QueryRow(`
			SELECT 
				COALESCE(endpoint_path, ''), 
				COALESCE(http_method, 'GET'),
				COALESCE(custom_headers, '{}'),
				COALESCE(timeout_seconds, 10),
				COALESCE(retry_count, 1)
			FROM adapter_configs 
			WHERE product_id = $1
		`, id).Scan(
			&endpointPath,
			&httpMethod,
			&customHeadersJSON,
			&timeoutSeconds,
			&retryCount,
		)

		if err != nil {
			// fallback default jika config tidak ada
			endpointPath = ""
			httpMethod = "GET"
			timeoutSeconds = 10
			retryCount = 1
			customHeadersJSON = []byte(`{}`)
		}

		// parse headers
		json.Unmarshal(customHeadersJSON, &customHeaders)

	} else {
		// skip adapter config
		log.Printf("[INFO] API key empty, skipping adapter_configs for product %s", id)
		customHeaders = make(map[string]string)
	}

	// 4. Build URL
	fullURL := apiEndpoint + endpointPath
	log.Printf("Testing connection to: %s [%s]", fullURL, httpMethod)

	// 5. Execute request dengan retry
	var lastResp *http.Response
	var lastErr error

	for i := 0; i < retryCount; i++ {
		req, err := http.NewRequest(httpMethod, fullURL, nil)
		if err != nil {
			lastErr = err
			continue
		}

		// default headers
		req.Header.Set("Content-Type", "application/json")

		// API key header (optional)
		if apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+apiKey)
		}

		// custom headers
		for k, v := range customHeaders {
			req.Header.Set(k, v)
		}

		client := &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			if i < retryCount-1 {
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}

		lastResp = resp
		lastErr = nil
		break
	}

	if lastErr != nil {
		updateProductConnectionStatus(id, false)
		responseError(w, "Failed to connect: "+lastErr.Error(), http.StatusBadRequest)
		return
	}
	defer lastResp.Body.Close()

	// 6. Read response
	bodyBytes, _ := io.ReadAll(io.LimitReader(lastResp.Body, 1024))

	// 7. Status check
	isConnected := lastResp.StatusCode >= 200 && lastResp.StatusCode < 300
	updateProductConnectionStatus(id, isConnected)

	// 8. Response
	response := map[string]interface{}{
		"success":      isConnected,
		"status_code":  lastResp.StatusCode,
		"status_text":  lastResp.Status,
		"message":      getStatusMessage(lastResp.StatusCode),
		"product_name": productName,
		"endpoint":     fullURL,
		"method":       httpMethod,
		"response":     truncateString(string(bodyBytes), 500),
		"tested_at":    time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func updateProductConnectionStatus(productID string, isConnected bool) {
	status := "connected"
	log.Printf("CONNECT = %v", isConnected)
	if !isConnected {
		status = "pending" // atau "pending", ganti yang valid
	}

	_, err := database.GetDB().Exec(`
		UPDATE products 
		SET status = $1, 
		    last_sync = NOW(), 
		    updated_at = NOW()
		WHERE id = $2
	`, status, productID)

	if err != nil {
		log.Printf("ERROR update product status: %v", err)
	}
}

func getStatusMessage(statusCode int) string {
	switch statusCode {
	case 200, 201, 204:
		return "Connection successful"
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized - Invalid API key"
	case 403:
		return "Forbidden"
	case 404:
		return "Endpoint not found"
	default:
		return fmt.Sprintf("HTTP %d", statusCode)
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func responseError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"message": message,
	})
}

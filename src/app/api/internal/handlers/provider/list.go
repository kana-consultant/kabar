package provider

import (
	"log"
	"net/http"

	"seo-backend/internal/builder"
	"seo-backend/internal/database"
)

// GetAll returns all API providers based on user permissions
func (h *ProviderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)

	log.Println("==== GET ALL PROVIDERS DEBUG ====")
	log.Println("userRole:", userCtx.Role)
	log.Println("teamID:", userCtx.TeamID)

	// Build query
	qb := builder.NewQueryBuilder("api_providers")
	query := qb.Select(
		"id", "name", "display_name", "description",
		"base_url", "auth_type", "auth_header", "auth_prefix",
		"text_endpoint", "image_endpoint",
		"default_headers", "request_template",
		"response_text_path", "response_image_path",
		"is_active", "created_at", "updated_at",
	).OrderBy("name ASC")

	// Apply role-based filtering
	whereClause := h.buildWhereClause(userCtx.Role)
	if whereClause != "" {
		query = query.Where(whereClause)
	}

	sqlQuery, args, err := query.Build()
	if err != nil {
		log.Println("query build error:", err)
		h.writeError(w, "Failed to fetch providers", http.StatusInternalServerError)
		return
	}

	log.Println("SQL:", sqlQuery)
	log.Printf("ARGS: %#v\n", args)

	rows, err := database.GetDB().Query(sqlQuery, args...)
	if err != nil {
		log.Println("query execute error:", err)
		h.writeError(w, "Failed to fetch providers", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var providers []APIProvider
	rowCount := 0

	for rows.Next() {
		rowCount++

		provider, err := h.scanProvider(rows)
		if err != nil {
			log.Printf("scan error row %d: %v", rowCount, err)
			continue
		}

		log.Println("provider loaded:", provider.ID, provider.Name)
		providers = append(providers, *provider)
	}

	if err := rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
	}

	log.Println("total rows:", rowCount)
	log.Println("providers returned:", len(providers))
	log.Println("==== END DEBUG ====")

	h.writeJSON(w, providers, http.StatusOK)
}

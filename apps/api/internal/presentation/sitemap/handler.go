// internal/domain/sitemap/handler.go
package sitemap

import (
	"encoding/json"
	"net/http"
	"seo-backend/internal/domain/sitemap"
	"strconv"
)

type SitemapHandler struct {
	sitemapService sitemap.Service
}

func NewSitemapHandler(sitemapService sitemap.Service) *SitemapHandler {
	return &SitemapHandler{
		sitemapService: sitemapService,
	}
}

// GenerateSitemap handles GET /sitemap
func (h *SitemapHandler) GenerateSitemap(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	productID := r.URL.Query().Get("product_id")
	baseURL := r.URL.Query().Get("base_url")

	// Validasi: product_id atau base_url harus ada
	if productID == "" && baseURL == "" {
		http.Error(w, "either product_id or base_url is required", http.StatusBadRequest)
		return
	}

	// Jika product_id ada tapi base_url kosong, ambil base_url dari product
	if productID != "" && baseURL == "" {
		// TODO: Ambil base_url dari database berdasarkan product_id
		// baseURL = getBaseURLByProductID(productID)
		// Untuk sementara, return error
		http.Error(w, "base_url is required when product_id is provided", http.StatusBadRequest)
		return
	}

	includeImages, _ := strconv.ParseBool(r.URL.Query().Get("include_images"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	req := sitemap.GenerateRequest{
		ProductID:     productID,
		BaseURL:       baseURL,
		IncludeImages: includeImages,
		Limit:         limit,
	}

	resp, err := h.sitemapService.GenerateSitemap(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return as XML
	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("X-Total-URLs", strconv.Itoa(resp.TotalURLs))
	w.Header().Set("X-Generated-At", resp.GeneratedAt.Format(http.TimeFormat))
	w.Header().Set("X-Product-ID", resp.ProductID)
	w.Header().Set("X-Base-URL", resp.BaseURL)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp.SitemapXML))
}

// GetSitemapHistory handles GET /sitemap/history
func (h *SitemapHandler) GetSitemapHistory(w http.ResponseWriter, r *http.Request) {
	// Get history from service
	history, err := h.sitemapService.GetSitemapHistory(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(history)
}

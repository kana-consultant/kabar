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

// GenerateSitemap godoc
// @Summary Generate sitemap
// @Description Generate an XML sitemap for a product or base URL
// @Tags Sitemap
// @Accept json
// @Produce xml
// @Param product_id query string false "Product ID to generate sitemap for"
// @Param base_url query string true "Base URL for the sitemap" example:"https://example.com"
// @Param include_images query bool false "Include images in sitemap" default(false)
// @Success 200 {string} string "XML Sitemap"
// @Failure 400 {object} map[string]string "Bad request - missing required parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Header 200 {string} X-Total-URLs "Total number of URLs in sitemap"
// @Header 200 {string} X-Generated-At "Timestamp when sitemap was generated"
// @Header 200 {string} X-Product-ID "Product ID used for generation"
// @Header 200 {string} X-Base-URL "Base URL used for generation"
// @Security BearerAuth
// @Router /api/sitemap [get]
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

	req := sitemap.GenerateRequest{
		ProductID:     productID,
		BaseURL:       baseURL,
		IncludeImages: includeImages,
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

// GetSitemapHistory godoc
// @Summary Get sitemap generation history
// @Description Get history of all sitemap generations
// @Tags Sitemap
// @Accept json
// @Produce json
// @Success 200 {array} []sitemap.HistoryResponse "List of sitemap generation history"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/sitemap/history [get]
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

// internal/domain/sitemap/handler.go
package sitemap

import (
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
	baseURL := r.URL.Query().Get("base_url")
	if baseURL == "" {
		http.Error(w, "base_url is required", http.StatusBadRequest)
		return
	}

	includeImages, _ := strconv.ParseBool(r.URL.Query().Get("include_images"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	req := sitemap.GenerateRequest{
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp.SitemapXML))
}

// internal/domain/draft/entity.go
package draft

import "time"

type Draft struct {
	ID             string     `json:"id"`
	Title          string     `json:"title"`
	Topic          string     `json:"topic"`
	Article        string     `json:"article"`
	ImageURL       *string    `json:"image_url"`
	ImagePrompt    string     `json:"image_prompt"`
	Status         string     `json:"status"`
	ScheduledFor   *time.Time `json:"scheduled_for"`
	TargetProducts []string   `json:"target_products"`
	HasImage       bool       `json:"has_image"`
	TeamID         *string    `json:"team_id"`
	UserID         *string    `json:"user_id"`
	SeoScore       int        `json:"seo_score"`
	CreatedBy      *string    `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type Keywords struct {
	ID        string    `json:"id"`
	IDDraft   string    `json:"id_draft"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DraftData struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Topic          string   `json:"topic"`
	Article        string   `json:"article"`
	ImageURL       *string  `json:"image_url"`
	ImagePrompt    string   `json:"image_prompt"`
	TargetProducts []string `json:"target_products"`
	SEOScore       int      `json:"seo_score"`
	Keywords       []string `json:"keywords,omitempty"` // Tambahkan ini
	Excerpt        string   `json:"excerpt"`
	Slug           string   `json:"slug"`
	Status         string
	ScheduledFor   *time.Time
}

// draft/models.go
type PublishResult struct {
	Results       interface{} `json:"results"`
	SomeFailed    bool        `json:"some_failed"`
	AllFailed     bool        `json:"all_failed"`
	Status        string      `json:"status"`  // "published", "partial", "failed", "scheduled"
	Message       string      `json:"message"` // Deskripsi status
	ScheduledFor  *time.Time  `json:"scheduled_for,omitempty"`
	PublishedAt   *time.Time  `json:"published_at,omitempty"`
	TotalProducts int         `json:"total_products"`
	SuccessCount  int         `json:"success_count"`
	FailedCount   int         `json:"failed_count"`
	Errors        []string    `json:"errors,omitempty"`
}

// Jangan lupa update struct SEOScore
type SEOScore struct {
	Total       int            `json:"total"`
	MaxScore    int            `json:"max_score"`
	Details     map[string]int `json:"details"`
	Suggestions []string       `json:"suggestions"`
}

type SimilarityResult struct {
	DraftID    string  `json:"draft_id"`
	Title      string  `json:"title"`
	Similarity float64 `json:"similarity"`
}
type DraftStats struct {
	// Basic metrics
	TotalDraft        int `json:"total_draft"`
	TotalWithImage    int `json:"total_with_image"`
	TotalWithoutImage int `json:"total_without_image"`
	TotalScheduled    int `json:"total_scheduled"`
	TotalPublished    int `json:"total_published"`
	TotalWithKeywords int `json:"total_with_keywords"`
	TotalWithSEO      int `json:"total_with_seo"`

	// Derived metrics
	CompletionRate    float64 `json:"completion_rate"`
	ScheduledRate     float64 `json:"scheduled_rate"`
	ImageCoverageRate float64 `json:"image_coverage_rate"`
	SEOScoreAvg       float64 `json:"seo_score_avg"`
	KeywordsAvgCount  float64 `json:"keywords_avg_count"`

	// Breakdowns
	StatusBreakdown      map[string]int `json:"status_breakdown"`
	ProductCoverage      map[string]int `json:"product_coverage"`
	ProductStatus        map[string]int `json:"product_status"`
	TopicBreakdown       map[string]int `json:"topic_breakdown"`
	SEOScoreDistribution map[string]int `json:"seo_score_distribution"`

	// Time series
	DailyActivity     []DailyActivity `json:"daily_activity"`
	WeeklyTrend       []WeeklyTrend   `json:"weekly_trend,omitempty"`
	ScheduledUpcoming []ScheduledItem `json:"scheduled_upcoming,omitempty"`

	// Content quality metrics
	TopTopics   []TopicStats   `json:"top_topics,omitempty"`
	TopKeywords []KeywordStats `json:"top_keywords,omitempty"`

	// Cache metadata
	CacheMetadata CacheMetadata `json:"cache_metadata,omitempty"`
}

type DailyActivity struct {
	Date         string  `json:"date"`
	Count        int     `json:"count"`
	Scheduled    int     `json:"scheduled"`
	Published    int     `json:"published"`
	WithImage    int     `json:"with_image"`
	WithKeywords int     `json:"with_keywords"`
	AvgSEO       float64 `json:"avg_seo_score"`
}

type WeeklyTrend struct {
	Week      string `json:"week"`
	Created   int    `json:"created"`
	Scheduled int    `json:"scheduled"`
	Published int    `json:"published"`
}

type TopicStats struct {
	Topic  string  `json:"topic"`
	Count  int     `json:"count"`
	AvgSEO float64 `json:"avg_seo_score"`
}

type KeywordStats struct {
	Keyword string `json:"keyword"`
	Count   int    `json:"count"`
}

type ScheduledItem struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	ScheduledFor time.Time `json:"scheduled_for"`
	Products     []string  `json:"products,omitempty"`
}

type CacheMetadata struct {
	CachedAt   time.Time `json:"cached_at"`
	TTL        string    `json:"ttl"`
	Generation float64   `json:"generation_time_ms"`
}

func (s *DraftStats) CalculateDerivedMetrics() {
	if s.TotalDraft > 0 {
		s.CompletionRate = float64(s.TotalWithImage) / float64(s.TotalDraft) * 100
		s.ScheduledRate = float64(s.TotalScheduled) / float64(s.TotalDraft) * 100
		s.ImageCoverageRate = float64(s.TotalWithImage) / float64(s.TotalDraft) * 100
	}
}

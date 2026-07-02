package repositories

import (
	"context"
	"database/sql"
)

type dashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) *dashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetTotalContent(ctx context.Context, where string, args []interface{}) (int, error) {
	query := `
		SELECT COUNT(*) FROM (
			SELECT id FROM drafts WHERE ` + where + `
			UNION ALL
			SELECT id FROM drafts WHERE ` + where + ` AND status IN ('published', 'generated')  
		) AS all_content
	`

	var total int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&total)
	return total, err
}

func (r *dashboardRepository) GetTotalProducts(ctx context.Context, where string, args []interface{}) (int, error) {
	query := `SELECT COUNT(*) FROM products WHERE ` + where

	var total int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&total)
	return total, err
}

func (r *dashboardRepository) GetTotalPublished(ctx context.Context, where string, args []interface{}) (int, error) {
	query := `
		SELECT COUNT(*) FROM drafts  
		WHERE status IN ('published', 'generated') AND ` + where

	var total int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&total)
	return total, err
}

func (r *dashboardRepository) GetAverageSeoScore(ctx context.Context, where string, args []interface{}) (float64, error) {
	query := `
		SELECT COALESCE(AVG(seo_score), 0) FROM drafts where status IN ('published', 'generated') 
		AND ` + where + ` AND seo_score IS NOT NULL 
	`

	var avg float64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&avg)
	return avg, err
}

func (r *dashboardRepository) GetContentCountByPeriod(ctx context.Context, where string, args []interface{}) (int, error) {
	query := `
		SELECT COUNT(*) FROM (
			SELECT id FROM drafts WHERE ` + where + `
			UNION ALL
			SELECT id FROM drafts WHERE ` + where + ` AND status IN ('published', 'generated')
		) AS all_content
	`

	var total int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&total)
	return total, err
}

func (r *dashboardRepository) GetProductsCountByPeriod(ctx context.Context, where string, args []interface{}) (int, error) {
	query := `SELECT COUNT(*) FROM products WHERE ` + where

	var total int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&total)
	return total, err
}

func (r *dashboardRepository) GetSeoScoreByPeriod(ctx context.Context, where string, args []interface{}) (float64, error) {
	query := `
		SELECT COALESCE(AVG(seo_score), 0) FROM drafts where status IN ('published', 'generated') 
		AND ` + where + ` AND seo_score IS NOT NULL 
	`

	var avg float64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&avg)
	return avg, err
}

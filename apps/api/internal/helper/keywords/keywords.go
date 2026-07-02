package keywords

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"seo-backend/internal/domain/draft"

	"github.com/google/uuid"
)

type KeywordSource interface {
	Column() string // "id_draft"
	ID() string     // nilai UUID-nya
}

type DraftSource struct {
	DraftID string
}

func (d DraftSource) Column() string { return "id_draft" }
func (d DraftSource) ID() string     { return d.DraftID }

// keyword baru berdasarkan sumber yang diberikan (draft).
func ReplaceKeywords(ctx context.Context, tx *sql.Tx, source KeywordSource, keywordNames []string) error {
	// Hapus semua keyword yang ada
	deleteQuery := fmt.Sprintf(`DELETE FROM keywords WHERE %s = $1`, source.Column())
	if _, err := tx.ExecContext(ctx, deleteQuery, source.ID()); err != nil {
		return fmt.Errorf("failed to delete existing keywords: %w", err)
	}

	// Tidak ada keyword baru, selesai
	if len(keywordNames) == 0 {
		return nil
	}

	// Hapus duplikat
	uniqueNames := UniqueStrings(keywordNames)

	// Batch insert keyword baru
	insertQuery := fmt.Sprintf(`
        INSERT INTO keywords (id, %s, name, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())`,
		source.Column(),
	)

	stmt, err := tx.PrepareContext(ctx, insertQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	for _, name := range uniqueNames {
		if _, err = stmt.ExecContext(ctx, uuid.New().String(), source.ID(), name); err != nil {
			return fmt.Errorf("failed to insert keyword %q: %w", name, err)
		}
	}

	log.Printf("Updated keywords for %s=%s: %d keywords inserted",
		source.Column(), source.ID(), len(uniqueNames))
	return nil
}

func UniqueStrings(strs []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, s := range strs {
		if s != "" && !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	return result
}

func InsertKeywordsWithDuplicateCheck(ctx context.Context, tx *sql.Tx, draftID string, keywords []draft.Keywords) error {
	query := `INSERT INTO keywords (id, id_draft, name, created_at, updated_at) 
	          VALUES ($1, $2, $3, NOW(), NOW())`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	inserted := 0
	for _, kw := range keywords {
		if kw.ID == "" {
			kw.ID = uuid.New().String()
		}

		result, err := stmt.ExecContext(ctx, kw.ID, draftID, kw.Name)
		if err != nil {
			log.Printf("Warning: failed to insert keyword %s: %v", kw.Name, err)
			continue
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			inserted++
		}
	}

	log.Printf("Inserted %d out of %d keywords for draft %s", inserted, len(keywords), draftID)
	return nil
}

func UpdateKeywords(
	ctx context.Context,
	tx *sql.Tx,
	source KeywordSource,
	kw interface{},
) error {
	var names []string

	switch kw := kw.(type) {
	case []string:
		names = kw

	case []draft.Keywords:
		names = make([]string, len(kw))
		for i, k := range kw {
			names[i] = k.Name
		}

	default:
		return fmt.Errorf("unsupported keywords type: %T", kw)
	}

	return ReplaceKeywords(ctx, tx, source, names)
}

func GetKeywords(
	ctx context.Context,
	db *sql.DB,
	source KeywordSource,
) ([]string, error) {

	query := fmt.Sprintf(`
	SELECT name
	FROM keywords
	WHERE %s = $1
	ORDER BY created_at ASC
	`, source.Column())

	rows, err := db.QueryContext(ctx, query, source.ID())
	if err != nil {
		return nil, fmt.Errorf("error querying keywords: %w", err)
	}
	defer rows.Close()

	var result []string

	for rows.Next() {
		var kw string
		if err := rows.Scan(&kw); err != nil {
			log.Printf("error scanning keyword row: %v", err)
			continue
		}
		result = append(result, kw)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating keyword rows: %w", err)
	}

	log.Printf("found %d keywords for %s=%s", len(result), source.Column(), source.ID())

	return result, nil
}

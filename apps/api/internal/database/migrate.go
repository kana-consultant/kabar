package database

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func RunMigrations() error {
	log.Println("Running migrations...")

	migrationPath := "internal/database/migrations/full_migration_schema.up.sql"

	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		return fmt.Errorf("migration file not found: %s", migrationPath)
	}

	sqlBytes, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	sqlContent := string(sqlBytes)

	if strings.TrimSpace(sqlContent) == "" {
		log.Println("Migration file is empty, skipping...")
		return nil
	}

	log.Printf("Migration file size: %d bytes", len(sqlBytes))

	// Better splitter untuk PostgreSQL
	statements := splitSQLStatements(sqlContent)

	successCount := 0
	errorCount := 0
	errorStatements := []int{}

	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		log.Printf("Executing statement %d...", i+1)
		log.Printf("SQL: %s", stmt)

		if _, err := GetDB().Exec(stmt); err != nil {
			errMsg := err.Error()
			// Skip error yang umum (already exists)
			if strings.Contains(errMsg, "already exists") ||
				strings.Contains(errMsg, "duplicate key") ||
				strings.Contains(errMsg, "duplicate constraint") ||
				strings.Contains(errMsg, "relation already exists") ||
				strings.Contains(errMsg, "already a constraint") {
				log.Printf("⚠️ Statement %d already exists, skipping...", i+1)
				successCount++
				continue
			}
			log.Printf("❌ Statement %d FAILED: %v", i+1, err)
			log.Printf("   SQL: %s", stmt)
			errorCount++
			errorStatements = append(errorStatements, i+1)
		} else {
			log.Printf("✅ Statement %d SUCCESS", i+1)
			successCount++
		}
	}

	log.Printf("Migrations completed: %d success, %d errors", successCount, errorCount)
	if len(errorStatements) > 0 {
		log.Printf("Failed statements: %v", errorStatements)
	}

	return nil
}

// splitSQLStatements splits SQL into statements while respecting:
// - Dollar-quoted strings ($$...$$, $tag$...$tag$)
// - Single quotes
// - Line comments (--)
// - Block comments (/* */)
func splitSQLStatements(sql string) []string {
	var statements []string
	var current strings.Builder

	inDollarQuote := false
	dollarTag := ""
	inSingleQuote := false
	inBlockComment := false
	inLineComment := false

	i := 0
	for i < len(sql) {
		ch := sql[i]
		remaining := sql[i:]

		// Handle line comments
		if !inDollarQuote && !inSingleQuote && !inBlockComment && strings.HasPrefix(remaining, "--") {
			inLineComment = true
			current.WriteString("--")
			i += 2
			continue
		}

		// End of line comment
		if inLineComment && ch == '\n' {
			inLineComment = false
			current.WriteByte(ch)
			i++
			continue
		}

		// Handle block comments
		if !inDollarQuote && !inSingleQuote && !inLineComment && strings.HasPrefix(remaining, "/*") {
			inBlockComment = true
			current.WriteString("/*")
			i += 2
			continue
		}

		if inBlockComment && strings.HasPrefix(remaining, "*/") {
			inBlockComment = false
			current.WriteString("*/")
			i += 2
			continue
		}

		// Handle dollar quotes
		if !inSingleQuote && !inBlockComment && !inLineComment {
			if !inDollarQuote && strings.HasPrefix(remaining, "$$") {
				inDollarQuote = true
				dollarTag = "$$"
				current.WriteString("$$")
				i += 2
				continue
			} else if !inDollarQuote && strings.HasPrefix(remaining, "$") {
				// Find dollar tag name: $tag$
				end := strings.Index(remaining[1:], "$")
				if end > 0 {
					tag := remaining[:end+2]
					inDollarQuote = true
					dollarTag = tag
					current.WriteString(tag)
					i += len(tag)
					continue
				}
			} else if inDollarQuote && strings.HasPrefix(remaining, dollarTag) {
				inDollarQuote = false
				current.WriteString(dollarTag)
				i += len(dollarTag)
				continue
			}
		}

		// Handle single quotes
		if !inDollarQuote && !inBlockComment && !inLineComment && ch == '\'' {
			inSingleQuote = !inSingleQuote
		}

		// Check for statement separator (semicolon)
		if ch == ';' && !inDollarQuote && !inSingleQuote && !inBlockComment && !inLineComment {
			stmt := strings.TrimSpace(current.String())
			if stmt != "" {
				statements = append(statements, stmt)
			}
			current.Reset()
			i++
			continue
		}

		current.WriteByte(ch)
		i++
	}

	// Last statement
	if stmt := strings.TrimSpace(current.String()); stmt != "" {
		statements = append(statements, stmt)
	}

	return statements
}

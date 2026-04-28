// internal/builder/query.go
package builder

import (
	sq "github.com/Masterminds/squirrel"
)

// InitDB inisialisasi squirrel dengan placeholder untuk PostgreSQL
func InitDB() {
	sq.StatementBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

// StatementBuilder untuk query (pakai $1, $2, $3)
var StatementBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// QueryBuilder generic query builder untuk semua tabel
type QueryBuilder struct {
	table string
}

// NewQueryBuilder membuat query builder baru untuk tabel tertentu
func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{table: table}
}

// SelectBuilder builder untuk SELECT query
func (q *QueryBuilder) Select(columns ...string) *SelectBuilder {
	return &SelectBuilder{
		builder: StatementBuilder.Select(columns...).From(q.table),
	}
}

// InsertBuilder builder untuk INSERT query
func (q *QueryBuilder) Insert() *InsertBuilder {
	return &InsertBuilder{
		builder: StatementBuilder.Insert(q.table),
	}
}

// UpdateBuilder builder untuk UPDATE query
func (q *QueryBuilder) Update() *UpdateBuilder {
	return &UpdateBuilder{
		builder: StatementBuilder.Update(q.table),
	}
}

// DeleteBuilder builder untuk DELETE query
func (q *QueryBuilder) Delete() *DeleteBuilder {
	return &DeleteBuilder{
		builder: StatementBuilder.Delete(q.table),
	}
}

// ========== SELECT BUILDER ==========

type SelectBuilder struct {
	builder sq.SelectBuilder
}

func (s *SelectBuilder) Where(condition string, args ...interface{}) *SelectBuilder {
	s.builder = s.builder.Where(condition, args...)
	return s
}

func (s *SelectBuilder) WhereEq(field string, value interface{}) *SelectBuilder {
	if value != nil && value != "" {
		s.builder = s.builder.Where(sq.Eq{field: value})
	}
	return s
}

func (s *SelectBuilder) WhereLike(field, value string) *SelectBuilder {
	if value != "" {
		s.builder = s.builder.Where(sq.ILike{field: "%" + value + "%"})
	}
	return s
}

func (s *SelectBuilder) WhereOr(conditions ...sq.Sqlizer) *SelectBuilder {
	s.builder = s.builder.Where(sq.Or(conditions))
	return s
}

func (s *SelectBuilder) WhereIn(field string, values []interface{}) *SelectBuilder {
	if len(values) > 0 {
		s.builder = s.builder.Where(sq.Eq{field: values})
	}
	return s
}

func (s *SelectBuilder) OrderBy(orderBy string) *SelectBuilder {
	s.builder = s.builder.OrderBy(orderBy)
	return s
}

func (s *SelectBuilder) Limit(limit uint64) *SelectBuilder {
	if limit > 0 {
		s.builder = s.builder.Limit(limit)
	}

	return s
}

func (s *SelectBuilder) Offset(offset uint64) *SelectBuilder {
	if offset > 0 {
		s.builder = s.builder.Offset(offset)
	}
	return s
}

func (s *SelectBuilder) Build() (string, []interface{}, error) {
	return s.builder.ToSql()
}

// ========== INSERT BUILDER ==========

type InsertBuilder struct {
	builder sq.InsertBuilder
}

func (i *InsertBuilder) Columns(columns ...string) *InsertBuilder {
	i.builder = i.builder.Columns(columns...)
	return i
}

func (i *InsertBuilder) Values(values ...interface{}) *InsertBuilder {
	i.builder = i.builder.Values(values...)
	return i
}

func (i *InsertBuilder) SetMap(data map[string]interface{}) *InsertBuilder {
	i.builder = i.builder.SetMap(data)
	return i
}

func (i *InsertBuilder) Returning(column string) *InsertBuilder {
	i.builder = i.builder.Suffix("RETURNING " + column)
	return i
}

func (i *InsertBuilder) Build() (string, []interface{}, error) {
	return i.builder.ToSql()
}

// ========== UPDATE BUILDER ==========

type UpdateBuilder struct {
	builder sq.UpdateBuilder
}

func (u *UpdateBuilder) Set(column string, value interface{}) *UpdateBuilder {
	u.builder = u.builder.Set(column, value)
	return u
}

func (u *UpdateBuilder) SetMap(data map[string]interface{}) *UpdateBuilder {
	u.builder = u.builder.SetMap(data)
	return u
}

func (u *UpdateBuilder) Where(condition string, args ...interface{}) *UpdateBuilder {
	u.builder = u.builder.Where(condition, args...)
	return u
}

func (u *UpdateBuilder) WhereEq(field string, value interface{}) *UpdateBuilder {
	if value != nil && value != "" {
		u.builder = u.builder.Where(sq.Eq{field: value})
	}
	return u
}

func (u *UpdateBuilder) Build() (string, []interface{}, error) {
	return u.builder.ToSql()
}

// ========== DELETE BUILDER ==========

type DeleteBuilder struct {
	builder sq.DeleteBuilder
}

func (d *DeleteBuilder) Where(condition string, args ...interface{}) *DeleteBuilder {
	d.builder = d.builder.Where(condition, args...)
	return d
}

func (d *DeleteBuilder) WhereEq(field string, value interface{}) *DeleteBuilder {
	if value != nil && value != "" {
		d.builder = d.builder.Where(sq.Eq{field: value})
	}
	return d
}

func (d *DeleteBuilder) Build() (string, []interface{}, error) {
	return d.builder.ToSql()
}

package model

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// ModelBase is the base model for all models
type ModelBase struct {
	ID        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Save inserts or updates the model using upsert logic
func (m *ModelBase) Save(db *sql.DB, tableName string, data map[string]interface{}) error {
	var keys []string
	var values []interface{}
	var placeholders []string
	var updates []string

	idx := 1
	for k, v := range data {
		keys = append(keys, k)
		values = append(values, v)
		placeholders = append(placeholders, fmt.Sprintf("?")) // MySQL uses '?' for placeholders
		if k != "id" {                                        // Assuming 'id' is the primary key and not updated
			updates = append(updates, fmt.Sprintf("%s = VALUES(%s)", k, k))
		}
		idx++
	}

	query := fmt.Sprintf(`
        INSERT INTO %s (%s)
        VALUES (%s)
        ON DUPLICATE KEY UPDATE
        %s
    `, tableName, strings.Join(keys, ", "), strings.Join(placeholders, ", "), strings.Join(updates, ", "))

	_, err := db.Exec(query, values...)
	return err
}

// FindByID retrieves a model by its ID from the database
func (m *ModelBase) FindByID(db *sql.DB, tableName string, id int) error {
	v := reflect.ValueOf(m).Elem()
	t := v.Type()

	var columns []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		column := strings.ToLower(field.Name)
		columns = append(columns, column)
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = ?", strings.Join(columns, ", "), tableName)
	row := db.QueryRow(query, id)

	var values []interface{}
	for i := 0; i < t.NumField(); i++ {
		values = append(values, v.Field(i).Addr().Interface())
	}

	return row.Scan(values...)
}

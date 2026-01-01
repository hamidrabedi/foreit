package advanced

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// DatabaseHistoryStore stores history in a database
type DatabaseHistoryStore struct {
	db        *sql.DB
	tableName string
}

// NewDatabaseHistoryStore creates a new database history store
func NewDatabaseHistoryStore(db *sql.DB, tableName string) *DatabaseHistoryStore {
	if tableName == "" {
		tableName = "admin_history"
	}
	return &DatabaseHistoryStore{
		db:        db,
		tableName: tableName,
	}
}

// EnsureTable creates the history table if it doesn't exist
func (dhs *DatabaseHistoryStore) EnsureTable(ctx context.Context) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id BIGSERIAL PRIMARY KEY,
			object_type VARCHAR(255) NOT NULL,
			object_id BIGINT NOT NULL,
			action VARCHAR(50) NOT NULL,
			user_id TEXT,
			user_name VARCHAR(255),
			changes JSONB,
			message TEXT,
			timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
			INDEX idx_object (object_type, object_id),
			INDEX idx_user (user_id),
			INDEX idx_timestamp (timestamp)
		)
	`, dhs.tableName)

	_, err := dhs.db.ExecContext(ctx, query)
	return err
}

// LogAction logs an action to the database
func (dhs *DatabaseHistoryStore) LogAction(ctx context.Context, entry *HistoryEntry) error {
	// Serialize changes to JSON
	changesJSON, err := json.Marshal(entry.Changes)
	if err != nil {
		return fmt.Errorf("failed to marshal changes: %w", err)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (object_type, object_id, action, user_id, user_name, changes, message, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`, dhs.tableName)

	// Convert user_id to string for storage
	userIDStr := fmt.Sprintf("%v", entry.UserID)
	
	err = dhs.db.QueryRowContext(
		ctx,
		query,
		entry.ObjectType,
		entry.ObjectID,
		entry.Action,
		userIDStr,
		entry.UserName,
		changesJSON,
		entry.Message,
		entry.Timestamp,
	).Scan(&entry.ID)

	if err != nil {
		return fmt.Errorf("failed to insert history entry: %w", err)
	}

	return nil
}

// GetObjectHistory gets history for a specific object
func (dhs *DatabaseHistoryStore) GetObjectHistory(ctx context.Context, objectType string, objectID int64) ([]*HistoryEntry, error) {
	query := fmt.Sprintf(`
		SELECT id, object_type, object_id, action, user_id, user_name, changes, message, timestamp
		FROM %s
		WHERE object_type = $1 AND object_id = $2
		ORDER BY timestamp DESC
	`, dhs.tableName)

	rows, err := dhs.db.QueryContext(ctx, query, objectType, objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query history: %w", err)
	}
	defer rows.Close()

	return dhs.scanEntries(rows)
}

// GetUserHistory gets history for a specific user
func (dhs *DatabaseHistoryStore) GetUserHistory(ctx context.Context, userID interface{}, limit int) ([]*HistoryEntry, error) {
	query := fmt.Sprintf(`
		SELECT id, object_type, object_id, action, user_id, user_name, changes, message, timestamp
		FROM %s
		WHERE user_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`, dhs.tableName)

	userIDStr := fmt.Sprintf("%v", userID)
	rows, err := dhs.db.QueryContext(ctx, query, userIDStr, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query user history: %w", err)
	}
	defer rows.Close()

	return dhs.scanEntries(rows)
}

// GetRecentHistory gets recent history entries
func (dhs *DatabaseHistoryStore) GetRecentHistory(ctx context.Context, limit int) ([]*HistoryEntry, error) {
	query := fmt.Sprintf(`
		SELECT id, object_type, object_id, action, user_id, user_name, changes, message, timestamp
		FROM %s
		ORDER BY timestamp DESC
		LIMIT $1
	`, dhs.tableName)

	rows, err := dhs.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent history: %w", err)
	}
	defer rows.Close()

	return dhs.scanEntries(rows)
}

// scanEntries scans rows into history entries
func (dhs *DatabaseHistoryStore) scanEntries(rows *sql.Rows) ([]*HistoryEntry, error) {
	entries := make([]*HistoryEntry, 0)

	for rows.Next() {
		entry := &HistoryEntry{}
		var changesJSON []byte
		var userIDStr string

		err := rows.Scan(
			&entry.ID,
			&entry.ObjectType,
			&entry.ObjectID,
			&entry.Action,
			&userIDStr,
			&entry.UserName,
			&changesJSON,
			&entry.Message,
			&entry.Timestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history entry: %w", err)
		}

		entry.UserID = userIDStr

		// Deserialize changes from JSON
		if len(changesJSON) > 0 {
			if err := json.Unmarshal(changesJSON, &entry.Changes); err != nil {
				return nil, fmt.Errorf("failed to unmarshal changes: %w", err)
			}
		}

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return entries, nil
}

// DeleteOldEntries deletes history entries older than the specified duration
func (dhs *DatabaseHistoryStore) DeleteOldEntries(ctx context.Context, olderThan time.Duration) (int64, error) {
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE timestamp < $1
	`, dhs.tableName)

	cutoffTime := time.Now().Add(-olderThan)
	result, err := dhs.db.ExecContext(ctx, query, cutoffTime)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old entries: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// GetStatistics gets statistics about history entries
func (dhs *DatabaseHistoryStore) GetStatistics(ctx context.Context) (map[string]int64, error) {
	query := fmt.Sprintf(`
		SELECT 
			action,
			COUNT(*) as count
		FROM %s
		GROUP BY action
	`, dhs.tableName)

	rows, err := dhs.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query statistics: %w", err)
	}
	defer rows.Close()

	stats := make(map[string]int64)
	for rows.Next() {
		var action string
		var count int64
		if err := rows.Scan(&action, &count); err != nil {
			return nil, fmt.Errorf("failed to scan statistics: %w", err)
		}
		stats[action] = count
	}

	return stats, nil
}

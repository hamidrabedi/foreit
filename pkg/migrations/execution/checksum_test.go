package execution

import (
	"testing"
)


func TestValidateChecksum(t *testing.T) {
	sql := "CREATE TABLE users (id BIGINT PRIMARY KEY);"
	checksum := CalculateChecksum(sql)

	tests := []struct {
		name    string
		sql     string
		checksum string
		wantErr bool
	}{
		{
			name:    "valid checksum",
			sql:     sql,
			checksum: checksum,
			wantErr: false,
		},
		{
			name:    "invalid checksum",
			sql:     sql,
			checksum: "invalid_checksum",
			wantErr: true,
		},
		{
			name:    "modified SQL",
			sql:     "CREATE TABLE users (id BIGINT);", // Modified SQL
			checksum: checksum,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateChecksum(tt.sql, tt.checksum)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateChecksum() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}



package execution

import (
	"testing"
)

func TestCalculateChecksum(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want string
	}{
		{
			name: "simple SQL",
			sql:  "CREATE TABLE users (id BIGINT PRIMARY KEY);",
			want: "", // We'll check it's not empty
		},
		{
			name: "empty SQL",
			sql:  "",
			want: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", // SHA256 of empty string
		},
		{
			name: "multi-line SQL",
			sql: `CREATE TABLE users (
				id BIGINT PRIMARY KEY,
				email VARCHAR(255) NOT NULL
			);`,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checksum := CalculateChecksum(tt.sql)
			if tt.want != "" && checksum != tt.want {
				t.Errorf("Expected checksum '%s', got '%s'", tt.want, checksum)
			}
			if checksum == "" {
				t.Error("Expected non-empty checksum")
			}
			if len(checksum) != 64 { // SHA256 produces 64 hex characters
				t.Errorf("Expected checksum to be 64 characters, got %d", len(checksum))
			}
		})
	}
}

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

func TestChecksumConsistency(t *testing.T) {
	sql := "CREATE TABLE users (id BIGINT PRIMARY KEY, email VARCHAR(255));"
	
	// Calculate checksum multiple times
	checksum1 := CalculateChecksum(sql)
	checksum2 := CalculateChecksum(sql)
	
	if checksum1 != checksum2 {
		t.Errorf("Checksums should be consistent: got %s and %s", checksum1, checksum2)
	}
	
	// Validate with both checksums
	if err := ValidateChecksum(sql, checksum1); err != nil {
		t.Errorf("Checksum validation failed: %v", err)
	}
	if err := ValidateChecksum(sql, checksum2); err != nil {
		t.Errorf("Checksum validation failed: %v", err)
	}
}


package cmd

import "testing"

func TestSanitizePackageName(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"users", "users"},
		{"user-service", "userservice"},
		{"123-service", "pkg123service"},
		{"", "pkg"},
		{"!!!", "pkg"},
	}

	for _, tc := range tests {
		if got := sanitizePackageName(tc.in); got != tc.want {
			t.Fatalf("sanitizePackageName(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestExportedIdentifier(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"users", "Users"},
		{"user-service", "UserService"},
		{"9service", "X9service"},
		{"", "Generated"},
	}

	for _, tc := range tests {
		if got := exportedIdentifier(tc.in); got != tc.want {
			t.Fatalf("exportedIdentifier(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestValidateScaffoldName(t *testing.T) {
	if err := validateScaffoldName("users", "module name"); err != nil {
		t.Fatalf("expected valid name, got %v", err)
	}
	if err := validateScaffoldName("user_service-1", "module name"); err != nil {
		t.Fatalf("expected valid name, got %v", err)
	}

	invalid := []string{"", "../x", "a/b", `a\\b`, "my app", "a!"}
	for _, v := range invalid {
		if err := validateScaffoldName(v, "module name"); err == nil {
			t.Fatalf("expected error for %q", v)
		}
	}
}

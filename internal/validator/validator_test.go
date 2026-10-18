package validator

import (
	"testing"
)

func TestNotBlank(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty string", "", false},
		{"whitespace only", "   \t\n  ", false},
		{"valid text", "hello", true},
		{"valid with spaces", "  hello world  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NotBlank(tt.input); got != tt.want {
				t.Errorf("NotBlank(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestMaxRunes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		max   int
		want  bool
	}{
		{"within limit", "hello", 10, true},
		{"exact limit", "hello", 5, true},
		{"exceeds limit", "hello world", 5, false},
		{"utf8 characters count", "привет", 6, true},
		{"utf8 characters exceed", "привет мир", 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxRunes(tt.input, tt.max); got != tt.want {
				t.Errorf("MaxRunes(%q, %d) = %v, want %v", tt.input, tt.max, got, tt.want)
			}
		})
	}
}

func TestPermittedInt(t *testing.T) {
	tests := []struct {
		name      string
		value     int
		permitted []int
		want      bool
	}{
		{"permitted value exists", 7, []int{1, 7, 365}, true},
		{"permitted value first", 1, []int{1, 7, 365}, true},
		{"permitted value last", 365, []int{1, 7, 365}, true},
		{"not permitted value", 30, []int{1, 7, 365}, false},
		{"empty permitted list", 5, []int{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PermittedInt(tt.value, tt.permitted...); got != tt.want {
				t.Errorf("PermittedInt(%d, %v) = %v, want %v", tt.value, tt.permitted, got, tt.want)
			}
		})
	}
}

func TestValidator(t *testing.T) {
	v := &Validator{}

	if !v.Valid() {
		t.Errorf("empty validator should be valid")
	}

	v.AddFieldError("title", "Title is required")
	if v.Valid() {
		t.Errorf("validator with field error should be invalid")
	}

	// Adding duplicate error for same field shouldn't overwrite the first error
	v.AddFieldError("title", "Second title error")
	if v.FieldErrors["title"] != "Title is required" {
		t.Errorf("got %q, want first error %q", v.FieldErrors["title"], "Title is required")
	}

	v.AddNonFieldError("General form error")
	if len(v.NonFieldErrors) != 1 {
		t.Errorf("got %d non-field errors, want 1", len(v.NonFieldErrors))
	}

	v2 := &Validator{}
	v2.CheckField(false, "content", "Content required")
	v2.CheckField(true, "expires", "Should not be added")

	if v2.FieldErrors["content"] != "Content required" {
		t.Errorf("expected content error, got %v", v2.FieldErrors["content"])
	}
	if _, exists := v2.FieldErrors["expires"]; exists {
		t.Errorf("did not expect error for expires")
	}
}

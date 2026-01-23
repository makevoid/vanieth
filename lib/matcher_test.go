package lib

import (
	"testing"
)

func TestParseMatchArg_ValidHexPrefix(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		ignoreCase bool
		wantPrefix string
	}{
		{"simple hex", "abc", false, "0xabc"},
		{"uppercase hex", "ABC", false, "0xABC"},
		{"mixed case hex", "AbC123", false, "0xAbC123"},
		{"with 0x prefix", "0xdef", false, "0xdef"},
		{"single char", "a", false, "0xa"},
		{"all digits", "123456", false, "0x123456"},
		{"ignore case converts to lower", "ABC", true, "0xabc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseMatchArg(tt.input, tt.ignoreCase)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !result.IsPrefix {
				t.Error("expected IsPrefix to be true")
			}
			if result.Prefix != tt.wantPrefix {
				t.Errorf("got prefix %q, want %q", result.Prefix, tt.wantPrefix)
			}
			if result.Regex != nil {
				t.Error("expected Regex to be nil for prefix match")
			}
		})
	}
}

func TestParseMatchArg_ValidRegex(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		ignoreCase bool
	}{
		{"dot wildcard", "a.b", false},
		{"star quantifier", "a*b", false},
		{"plus quantifier", "a+b", false},
		{"question mark", "a?b", false},
		{"character class", "[abc]", false},
		{"alternation", "abc|def", false},
		{"anchor end", "abc$", false},
		{"anchor start in pattern", "^abc", false},
		{"complex pattern", ".*ABC", false},
		{"ends with pattern", ".*DEF$", false},
		{"quantifier range", "a{2,4}", false},
		{"escape sequence", `a\.b`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseMatchArg(tt.input, tt.ignoreCase)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.IsPrefix {
				t.Error("expected IsPrefix to be false for regex")
			}
			if result.Regex == nil {
				t.Error("expected Regex to be non-nil")
			}
		})
	}
}

func TestParseMatchArg_InvalidCharacters(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantCharErr rune
	}{
		{"contains 'o'", "foo", 'o'},
		{"contains 'g'", "deadbeefg", 'g'},
		{"contains 'z'", "abcz", 'z'},
		{"contains 'x' not at start", "axbc", 'x'},
		{"contains space", "ab c", ' '},
		{"contains 'O' uppercase", "ABC0O", 'O'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseMatchArg(tt.input, false)
			if err == nil {
				t.Fatal("expected error for invalid character")
			}
			if result != nil {
				t.Error("expected nil result on error")
			}
			// Check error message contains the invalid character
			expectedMsg := "invalid character '" + string(tt.wantCharErr) + "'"
			if !contains(err.Error(), expectedMsg) {
				t.Errorf("error %q should mention invalid character %q", err.Error(), string(tt.wantCharErr))
			}
		})
	}
}

func TestParseMatchArg_EmptyInput(t *testing.T) {
	result, err := ParseMatchArg("", false)
	if err == nil {
		t.Fatal("expected error for empty input")
	}
	if result != nil {
		t.Error("expected nil result on error")
	}
}

func TestParseMatchArg_OnlyPrefix(t *testing.T) {
	result, err := ParseMatchArg("0x", false)
	if err == nil {
		t.Fatal("expected error for '0x' only input")
	}
	if result != nil {
		t.Error("expected nil result on error")
	}
}

// ============================================================================
// Top 5 Real-World Use Case Tests
// ============================================================================

// UseCase 1: Simple vanity prefix - most common use case
// Users want addresses starting with a memorable hex pattern like "dead", "cafe", "1337"
func TestUseCase_SimpleVanityPrefix(t *testing.T) {
	testCases := []struct {
		pattern    string
		wantPrefix string
	}{
		{"dead", "0xdead"},
		{"cafe", "0xcafe"},
		{"1337", "0x1337"},
		{"beef", "0xbeef"},
		{"ace", "0xace"},
	}

	for _, tc := range testCases {
		t.Run(tc.pattern, func(t *testing.T) {
			result, err := ParseMatchArg(tc.pattern, false)
			if err != nil {
				t.Fatalf("common vanity prefix %q should be valid: %v", tc.pattern, err)
			}
			if !result.IsPrefix {
				t.Error("simple hex patterns should use prefix matching for performance")
			}
			if result.Prefix != tc.wantPrefix {
				t.Errorf("got %q, want %q", result.Prefix, tc.wantPrefix)
			}
		})
	}
}

// UseCase 2: Case-insensitive search - users often don't care about checksum case
func TestUseCase_CaseInsensitiveSearch(t *testing.T) {
	// With ignore case, "DEAD" should become "0xdead"
	result, err := ParseMatchArg("DEAD", true)
	if err != nil {
		t.Fatalf("case insensitive search should work: %v", err)
	}
	if result.Prefix != "0xdead" {
		t.Errorf("ignore case should lowercase: got %q, want %q", result.Prefix, "0xdead")
	}

	// Mixed case should also work
	result, err = ParseMatchArg("DeAdBeEf", true)
	if err != nil {
		t.Fatalf("mixed case with ignore should work: %v", err)
	}
	if result.Prefix != "0xdeadbeef" {
		t.Errorf("got %q, want %q", result.Prefix, "0xdeadbeef")
	}
}

// UseCase 3: Regex patterns for flexible matching
// Users want to find patterns anywhere in address or complex patterns
func TestUseCase_RegexPatterns(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
	}{
		{"match anywhere", ".*dead"},
		{"ends with pattern", ".*beef$"},
		{"repeating chars", "0+"},
		{"alternation", "dead|beef"},
		{"character class", "[abc]{4}"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseMatchArg(tc.pattern, false)
			if err != nil {
				t.Fatalf("regex pattern %q should be valid: %v", tc.pattern, err)
			}
			if result.IsPrefix {
				t.Error("regex pattern should not use prefix matching")
			}
			if result.Regex == nil {
				t.Error("regex should be compiled")
			}
		})
	}
}

// UseCase 4: Typo detection - catch common mistakes with non-hex characters
// This is the main issue the PR #5 was trying to solve
func TestUseCase_TypoDetection(t *testing.T) {
	invalidPatterns := []struct {
		pattern     string
		invalidChar rune
		description string
	}{
		{"hello", 'h', "common word with non-hex 'h'"},
		{"test", 't', "common word with non-hex 't'"},
		{"good", 'g', "common word with non-hex 'g'"},
		{"zero", 'z', "word with 'z' instead of '0'"},
		{"0xgood", 'g', "looks like hex but contains 'g'"},
	}

	for _, tc := range invalidPatterns {
		t.Run(tc.description, func(t *testing.T) {
			result, err := ParseMatchArg(tc.pattern, false)
			if err == nil {
				t.Fatalf("pattern %q should fail: %s", tc.pattern, tc.description)
			}
			if result != nil {
				t.Error("should return nil result on error")
			}
			expectedMsg := "invalid character '" + string(tc.invalidChar) + "'"
			if !contains(err.Error(), expectedMsg) {
				t.Errorf("error should mention '%c': got %q", tc.invalidChar, err.Error())
			}
		})
	}
}

// UseCase 5: Edge cases users might encounter
func TestUseCase_EdgeCases(t *testing.T) {
	t.Run("full address length prefix", func(t *testing.T) {
		// 40 hex chars (full address without 0x)
		fullAddr := "1234567890abcdef1234567890abcdef12345678"
		result, err := ParseMatchArg(fullAddr, false)
		if err != nil {
			t.Fatalf("full length hex should be valid: %v", err)
		}
		if !result.IsPrefix {
			t.Error("should use prefix matching")
		}
	})

	t.Run("single character search", func(t *testing.T) {
		result, err := ParseMatchArg("a", false)
		if err != nil {
			t.Fatalf("single char should be valid: %v", err)
		}
		if result.Prefix != "0xa" {
			t.Errorf("got %q, want %q", result.Prefix, "0xa")
		}
	})

	t.Run("regex with valid hex only chars", func(t *testing.T) {
		// This pattern uses regex metachar but only hex letters
		result, err := ParseMatchArg("dead.*beef", false)
		if err != nil {
			t.Fatalf("regex with hex chars should work: %v", err)
		}
		if result.Regex == nil {
			t.Error("should compile regex")
		}
	})

	t.Run("0x prefix is stripped", func(t *testing.T) {
		result, err := ParseMatchArg("0xabc", false)
		if err != nil {
			t.Fatalf("0x prefix should be handled: %v", err)
		}
		if result.Prefix != "0xabc" {
			t.Errorf("got %q, want %q", result.Prefix, "0xabc")
		}
	})

	t.Run("invalid regex pattern", func(t *testing.T) {
		// Unclosed bracket is invalid regex
		_, err := ParseMatchArg("[abc", false)
		if err == nil {
			t.Fatal("invalid regex should return error")
		}
		if !contains(err.Error(), "invalid regex") {
			t.Errorf("error should mention invalid regex: %q", err.Error())
		}
	})
}

// helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

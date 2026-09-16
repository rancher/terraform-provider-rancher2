package main

import (
	"go/parser"
	"go/token"
	"os"
	"strings"
	"testing"
)

func TestMatchesKeywords(t *testing.T) {
	keywords := []string{"password", "token", "credential", "pemencodedprivatekey", "secret", "private"}

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"exact password", "password", true},
		{"exact token", "token", true},
		{"prefix", "my_password", true},
		{"suffix", "token_key", true},
		{"mixed case", "PeMeNcoDeDpRiVaTeKeY", true},
		{"no match", "username", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesKeywords(tt.input, keywords)
			if got != tt.expected {
				t.Errorf("matchesKeywords(%q) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestLoadKeywords(t *testing.T) {
	kwData := `
# Comment line
custom_password
custom_secret
`
	tmpFile, err := os.CreateTemp(t.TempDir(), "test-keywords-*")
	if err != nil {
		t.Fatalf("failed to create temp keywords file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(kwData)); err != nil {
		t.Fatalf("failed to write temp keywords file: %v", err)
	}
	tmpFile.Close()

	keywords, err := loadKeywords(tmpFile.Name())
	if err != nil {
		t.Fatalf("loadKeywords failed: %v", err)
	}

	if len(keywords) != 2 {
		t.Errorf("expected 2 keywords, got %d: %v", len(keywords), keywords)
	}

	if keywords[0] != "custom_password" || keywords[1] != "custom_secret" {
		t.Errorf("unexpected keywords returned: %v", keywords)
	}
}

func TestLoadIgnoreList(t *testing.T) {
	ignoreData := `
# Comment line
global_ignore

test.go:file_local_ignore
test.tf:local_tf_ignore
`
	tmpFile, err := os.CreateTemp(t.TempDir(), "test-ignore-*")
	if err != nil {
		t.Fatalf("failed to create temp ignore file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(ignoreData)); err != nil {
		t.Fatalf("failed to write temp ignore file: %v", err)
	}
	tmpFile.Close()

	list, err := loadIgnoreList(tmpFile.Name())
	if err != nil {
		t.Fatalf("loadIgnoreList failed: %v", err)
	}

	if !list.IsIgnored("any.go", "global_ignore") {
		t.Errorf("expected global_ignore to be ignored globally")
	}

	if !list.IsIgnored("test.go", "file_local_ignore") {
		t.Errorf("expected file_local_ignore to be ignored for test.go")
	}

	if list.IsIgnored("other.go", "file_local_ignore") {
		t.Errorf("did not expect file_local_ignore to be ignored for other.go")
	}

	if !list.IsIgnored("test.tf", "local_tf_ignore") {
		t.Errorf("expected local_tf_ignore to be ignored for test.tf")
	}
}

func TestCheckGoFile(t *testing.T) {
	src := `package test
import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func testSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"unmarked_password": {
			Type:     schema.TypeString,
			Required: true,
		},
		"marked_password": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"ignored_password": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("failed to parse test source: %v", err)
	}

	ignoreList := &IgnoreList{
		global: map[string]bool{"ignored_password": true},
		local:  make(map[string]map[string]bool),
	}

	keywords := []string{"password"}

	findings := checkGoFile(fset, file, "test.go", ignoreList, keywords)

	expectedFindings := 1
	if len(findings) != expectedFindings {
		t.Errorf("expected %d findings, got %d: %v", expectedFindings, len(findings), findings)
	} else if got := findings[0]; !strings.Contains(got, "unmarked_password") {
		t.Errorf("expected finding for unmarked_password, got %q", got)
	}
}

func TestAnalyzeTfContent(t *testing.T) {
	tfSrc := `
variable "normal" {
  type = string
}

variable "unmarked_password" {
  type = string
}

variable "ignored_password" {
  type = string
}

variable "marked_token" {
  type      = string
  sensitive = true
}

output "unmarked_credential" {
  value = "foo"
}

output "marked_secret" {
  value     = "bar"
  sensitive = true
}

variable "commented_sensitive_password" {
  type = string
  # sensitive = true
}

variable "heredoc_password" {
  type    = string
  default = <<EOF
this is a password
EOF
}

variable "heredoc_indented_password" {
  type    = string
  default = <<-EOF
this is a password
EOF
}

variable "url_comment_password" {
  type        = string
  description = "http://example.com" // This should not be treated as a comment start
}

variable "braces_in_quotes_password" {
  type        = string
  description = "The {password} of the user" // Braces inside quotes should be ignored
}

variable "single_line_unmarked_password" { type = string }
variable "single_line_marked_password" { type = string; sensitive = true }

variable "marked_password_after_heredoc" {
  type = map(string)
  default = { foo = <<EOF
password
EOF
  }
  sensitive = true
}
`
	ignoreList := &IgnoreList{
		global: map[string]bool{"ignored_password": true},
		local:  make(map[string]map[string]bool),
	}

	keywords := []string{"password", "token", "credential", "pemencodedprivatekey", "secret", "private"}

	findings, err := analyzeTfContent(tfSrc, "test.tf", ignoreList, keywords)
	if err != nil {
		t.Fatalf("analyzeTfContent failed: %v", err)
	}

	expectedFindings := 8
	if len(findings) != expectedFindings {
		t.Errorf("expected %d findings, got %d: %v", expectedFindings, len(findings), findings)
	}

	var foundUnmarkedPassword, foundUnmarkedCredential, foundCommentedSensitive, foundHeredoc, foundHeredocIndented, foundUrlComment, foundBracesInQuotes, foundSingleLineUnmarked bool
	for _, finding := range findings {
		if strings.Contains(finding, "unmarked_password") {
			foundUnmarkedPassword = true
		}
		if strings.Contains(finding, "unmarked_credential") {
			foundUnmarkedCredential = true
		}
		if strings.Contains(finding, "commented_sensitive_password") {
			foundCommentedSensitive = true
		}
		if strings.Contains(finding, "heredoc_password") {
			foundHeredoc = true
		}
		if strings.Contains(finding, "heredoc_indented_password") {
			foundHeredocIndented = true
		}
		if strings.Contains(finding, "url_comment_password") {
			foundUrlComment = true
		}
		if strings.Contains(finding, "braces_in_quotes_password") {
			foundBracesInQuotes = true
		}
		if strings.Contains(finding, "single_line_unmarked_password") {
			foundSingleLineUnmarked = true
		}
		if strings.Contains(finding, "ignored_password") {
			t.Errorf("unexpected finding for ignored_password")
		}
		if strings.Contains(finding, "single_line_marked_password") {
			t.Errorf("unexpected finding for single_line_marked_password")
		}
		if strings.Contains(finding, "marked_password_after_heredoc") {
			t.Errorf("unexpected finding for marked_password_after_heredoc")
		}
	}

	if !foundUnmarkedPassword {
		t.Errorf("expected finding for variable unmarked_password")
	}
	if !foundUnmarkedCredential {
		t.Errorf("expected finding for output unmarked_credential")
	}
	if !foundCommentedSensitive {
		t.Errorf("expected finding for variable commented_sensitive_password")
	}
	if !foundHeredoc {
		t.Errorf("expected finding for variable heredoc_password")
	}
	if !foundHeredocIndented {
		t.Errorf("expected finding for variable heredoc_indented_password")
	}
	if !foundUrlComment {
		t.Errorf("expected finding for variable url_comment_password")
	}
	if !foundBracesInQuotes {
		t.Errorf("expected finding for variable braces_in_quotes_password")
	}
	if !foundSingleLineUnmarked {
		t.Errorf("expected finding for variable single_line_unmarked_password")
	}
}

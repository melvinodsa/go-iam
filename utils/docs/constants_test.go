package docs

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntoDescription(t *testing.T) {
	// Test that the constant is defined and not empty
	assert.NotEmpty(t, intoDescription)

	// Test that it contains expected content
	assert.Contains(t, intoDescription, "Go IAM")
	assert.Contains(t, intoDescription, "Identity and Access Management")
	assert.Contains(t, strings.ToLower(intoDescription), "golang")
	assert.Contains(t, intoDescription, "multi-tenant")
	assert.Contains(t, intoDescription, "authentication")
	assert.Contains(t, intoDescription, "authorization")

	// Test that it contains resource links
	assert.Contains(t, intoDescription, "go-iam-ui")
	assert.Contains(t, intoDescription, "go-iam-docker")
	assert.Contains(t, intoDescription, "github.com/melvinodsa")

	// Test that it contains architecture information
	assert.Contains(t, intoDescription, "Multi-Tenant Architecture")
	assert.Contains(t, intoDescription, "Authentication and Authorization")
	assert.Contains(t, intoDescription, "Resources and Roles")

	// Test that it contains example sections
	assert.Contains(t, intoDescription, "Examples")
	assert.Contains(t, intoDescription, "@ui/my-super-app")

	// Test that it's properly formatted markdown
	assert.Contains(t, intoDescription, "# Documentation")
	assert.Contains(t, intoDescription, "## Resources")
	assert.Contains(t, intoDescription, "## How it works")
	assert.Contains(t, intoDescription, "### Multi-Tenant Architecture")
	assert.Contains(t, intoDescription, "### Authentication and Authorization")
	assert.Contains(t, intoDescription, "### Resources and Roles")

	// Test that it contains markdown elements
	assert.Contains(t, intoDescription, "![Go IAM Logo]")
	assert.Contains(t, intoDescription, "![Auth process]")
	assert.Contains(t, intoDescription, "<details>")
	assert.Contains(t, intoDescription, "<summary>")

	// Test that links are properly formatted
	assert.Contains(t, intoDescription, "[go-iam-ui](https://github.com/melvinodsa/go-iam-ui)")
	assert.Contains(t, intoDescription, "[go-iam-docker](https://github.com/melvinodsa/go-iam-docker)")
	assert.Contains(t, intoDescription, "[go-iam](https://github.com/melvinodsa/go-iam)")

	// Test that it contains JWT information
	assert.Contains(t, intoDescription, "JWT")

	// Test that it's a reasonable length (not too short, not ridiculously long)
	assert.Greater(t, len(intoDescription), 1000, "Description should be substantial")
	assert.Less(t, len(intoDescription), 10000, "Description should not be excessively long")
}

func TestIntoDescriptionStructure(t *testing.T) {
	lines := strings.Split(intoDescription, "\n")

	// Should have multiple lines
	assert.Greater(t, len(lines), 20, "Description should have multiple lines")

	// Should start with main heading
	assert.True(t, strings.HasPrefix(strings.TrimSpace(lines[0]), "# Documentation for Go IAM APIs"))

	// Should contain various markdown elements
	var hasLinks, hasList, hasImages bool

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "> ") {
			hasList = true
		}
		if strings.Contains(trimmed, "](") {
			hasLinks = true
		}
		if strings.Contains(trimmed, "![") {
			hasImages = true
		}
	}

	assert.True(t, hasLinks, "Description should contain markdown links")
	assert.True(t, hasList, "Description should contain list items")
	assert.True(t, hasImages, "Description should contain images")
}

func TestIntoDescriptionContent(t *testing.T) {
	// Test specific content sections
	sections := []string{
		"Documentation for Go IAM APIs",
		"Resources",
		"How it works",
		"Multi-Tenant Architecture",
		"Authentication and Authorization",
		"Resources and Roles",
		"Examples",
	}

	for _, section := range sections {
		assert.Contains(t, intoDescription, section, "Description should contain section: %s", section)
	}

	// Test technical terms
	technicalTerms := []string{
		"IAM", "JWT", "OAuth2", "API", "stateless",
		"multi-tenant", "authentication", "authorization",
		"provider", "client", "role", "resource", "policies",
	}

	for _, term := range technicalTerms {
		assert.Contains(t, strings.ToLower(intoDescription), strings.ToLower(term),
			"Description should contain technical term: %s", term)
	}
}

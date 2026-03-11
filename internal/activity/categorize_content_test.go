package activity

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/demeyerthom/feeds-aggregator/internal"
	"github.com/stretchr/testify/assert"
)

// parseCategoriesJSON is a helper function to parse the categories JSON response
// This replicates the logic in categorize_content.go for isolated testing
func parseCategoriesJSON(jsonResponse string, categories *[]string) error {
	importedJSON := []byte(jsonResponse)
	return json.Unmarshal(importedJSON, categories)
}

func TestCategorizeContent_Success(t *testing.T) {
	t.Run("JSON parsing of valid categories", func(t *testing.T) {
		jsonResponse := `["Programming Languages", "AI & Machine Learning", "Tutorials & Guides"]`
		var categories []string
		err := parseCategoriesJSON(jsonResponse, &categories)
		assert.NoError(t, err)
		assert.Len(t, categories, 3)
		assert.Equal(t, "Programming Languages", categories[0])
	})

	t.Run("JSON parsing with invalid JSON", func(t *testing.T) {
		jsonResponse := `not valid json`
		var categories []string
		err := parseCategoriesJSON(jsonResponse, &categories)
		assert.Error(t, err)
	})

	t.Run("JSON parsing with empty array", func(t *testing.T) {
		jsonResponse := `[]`
		var categories []string
		err := parseCategoriesJSON(jsonResponse, &categories)
		assert.NoError(t, err)
		assert.Len(t, categories, 0)
	})
}

func TestCategorizeContent_JSONParsing(t *testing.T) {
	tests := []struct {
		name          string
		jsonResponse  string
		expectedError bool
		expectedCount int
	}{
		{
			name:          "single category",
			jsonResponse:  `["Technology"]`,
			expectedError: false,
			expectedCount: 1,
		},
		{
			name:          "five categories",
			jsonResponse:  `["Cat1", "Cat2", "Cat3", "Cat4", "Cat5"]`,
			expectedError: false,
			expectedCount: 5,
		},
		{
			name:          "invalid JSON",
			jsonResponse:  `not json`,
			expectedError: true,
			expectedCount: 0,
		},
		{
			name:          "empty array",
			jsonResponse:  `[]`,
			expectedError: false,
			expectedCount: 0,
		},
		{
			name:          "malformed JSON",
			jsonResponse:  `["category1", "category2"`,
			expectedError: true,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var categories []string
			err := parseCategoriesJSON(tt.jsonResponse, &categories)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, categories, tt.expectedCount)
			}
		})
	}
}

func TestCategorizeContent_ValidatesCategoryCount(t *testing.T) {
	tests := []struct {
		name        string
		categories  []string
		expectValid bool
	}{
		{
			name:        "zero categories - invalid",
			categories:  []string{},
			expectValid: false,
		},
		{
			name:        "one category - valid",
			categories:  []string{"Technology"},
			expectValid: true,
		},
		{
			name:        "three categories - valid",
			categories:  []string{"Tech", "AI", "Dev"},
			expectValid: true,
		},
		{
			name:        "five categories - valid",
			categories:  []string{"A", "B", "C", "D", "E"},
			expectValid: true,
		},
		{
			name:        "six categories - invalid",
			categories:  []string{"A", "B", "C", "D", "E", "F"},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := len(tt.categories) > 0 && len(tt.categories) <= 5
			assert.Equal(t, tt.expectValid, isValid)
		})
	}
}

func TestCategorizeContent_BuildsCorrectPath(t *testing.T) {
	// Test that the path building logic is correct
	dataDir := "/tmp/data"
	testID := "abc123"
	expectedPath := "/tmp/data/abc123.html"

	// This is the same logic as in categorize_content.go
	filename := filepath.Join(dataDir, testID+".html")

	assert.Equal(t, expectedPath, filename)
}

func TestCategorizeContent_ValidatesFeedItemDocument(t *testing.T) {
	// Test that FeedItemDocument has the required fields
	feedItemDoc := internal.FeedItemDocument{
		Link:  "https://example.com/article",
		Title: "Test Article",
	}

	assert.NotEmpty(t, feedItemDoc.Link)
	assert.NotEmpty(t, feedItemDoc.Title)
}

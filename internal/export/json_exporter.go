package export

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jamesBoder/BeerInfoApp.git/internal/models"
)

// JSONExporter struct implementing Exporter interface
type JSONExporter struct{}

// NewJSONExporter constructor
func NewJSONExporter() *JSONExporter {
	return &JSONExporter{}
}

// Export method to export beers to a JSON file
func (je *JSONExporter) Export(beers []models.Beer, filename string) error {
	// create or truncate the output file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	// create a JSON encoder with indentation
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	// encode the beers slice to JSON
	if err := encoder.Encode(beers); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// FileExtension method to return the file extension for JSON files
func (je *JSONExporter) FileExtension() string {
	return "json"
}

// FormatName method to return the format name
func (je *JSONExporter) FormatName() string {
	return "JSON"
}

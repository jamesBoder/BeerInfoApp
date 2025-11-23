package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/jamesBoder/BeerInfoApp.git/internal/models"
)

// empty CSVExporter struct implementing Exporter interface
type CSVExporter struct{}

// *CSVExporter constructor
func NewCSVExporter() *CSVExporter {
	return &CSVExporter{}
}

// Export method to export beer data to a CSV file
func (e *CSVExporter) Export(beers []models.Beer, filename string) error {
	// create or open the CSV file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	// create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write the header row
	header := []string{"SKU", "Name", "Brewery", "Region", "Country", "ABV", "IBU", "Category", "Rating", "Food Pairing", "SubCategory_1", "SubCategory_2", "Description"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// write beer data rows
	for _, beer := range beers {
		// try to parse numeric fields, fall back to the original string on error
		abvField := beer.Abv
		if v, err := strconv.ParseFloat(beer.Abv, 64); err == nil {
			abvField = fmt.Sprintf("%.2f", v)
		}

		ibuField := beer.Ibu
		if v, err := strconv.ParseFloat(beer.Ibu, 64); err == nil {
			ibuField = fmt.Sprintf("%.2f", v)
		}

		row := []string{
			beer.Sku,
			beer.Name,
			beer.Brewery,
			beer.Region,
			beer.Country,
			abvField,
			ibuField,
			beer.Category,
			beer.Rating,
			beer.FoodPairing,
			beer.SubCategory_1,
			beer.SubCategory_2,
			beer.Description,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

// FileExtension method that returns "csv"
func (e *CSVExporter) FileExtension() string {
	return "csv"
}

// FormatName method that returns "CSV"
func (e *CSVExporter) FormatName() string {
	return "CSV"
}

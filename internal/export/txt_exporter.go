package export

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jamesBoder/BeerInfoApp.git/internal/models"
)

// empty TXTExporter struct implementing Exporter interface
type TXTExporter struct{}

// constructor function for TXTExporter
func NewTXTExporter() *TXTExporter {
	// return pointer to new TXTExporter struct
	return &TXTExporter{}
}

// Export method to export beers to a TXT file
func (te *TXTExporter) Export(beers []models.Beer, filename string) error {
	// create or truncate the output file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// write beer data to the file
	for _, beer := range beers {
		// parse numeric fields to string
		abvField := beer.Abv
		if v, err := strconv.ParseFloat(beer.Abv, 64); err == nil {
			abvField = fmt.Sprintf("%.2f", v)
		}

		ibuField := beer.Ibu
		if v, err := strconv.ParseFloat(beer.Ibu, 64); err == nil {
			ibuField = fmt.Sprintf("%.2f", v)
		}

		beerInfo := fmt.Sprintf("Sku: %s\nName: %s\nBrewery: %s\nRegion: %s\nCountry: %s\nABV: %s\nIBU: %s\nCategory: %s\nRating: %s\nFood Pairing: %s\nSubCategory 1: %s\nSubCategory 2: %s\nDescription: %s\n\n",
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
		)
		_, err := file.WriteString(beerInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

// FileExtension method to return the file extension for TXT files
func (te *TXTExporter) FileExtension() string {
	return "txt"
}

// FormatName method to return the format name
func (te *TXTExporter) FormatName() string {
	return "TXT"
}

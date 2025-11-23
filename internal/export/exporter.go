package export

import (
	"fmt"
	"strings"
	"time"

	"github.com/jamesBoder/BeerInfoApp.git/internal/models"
)

// Exporter interface for exporting data
type Exporter interface {
	Export(beers []models.Beer, filename string) error
	FileExtension() string
	FormatName() string
}

// create a helper function to generate filename
func GenerateFilename(username, dataType, format string) string {
	// generate timestamped filenames
	timestamp := time.Now().Format("20060102_150405")
	safeUsername := strings.ToLower(strings.ReplaceAll(username, " ", "_"))
	return fmt.Sprintf("%s_%s_%s.%s", safeUsername, dataType, timestamp, format)
}

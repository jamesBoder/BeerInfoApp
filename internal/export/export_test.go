package export

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamesBoder/BeerInfoApp.git/internal/models"
)

// testBeers returns a small slice of beers for use across all exporter tests.
func testBeers() []models.Beer {
	return []models.Beer{
		{
			Sku:           "SKU001",
			Name:          "Berkshire IPA",
			Brewery:       "Berkshire Brewing",
			Region:        "New England",
			Country:       "USA",
			Abv:           "6.50",
			Ibu:           "55.00",
			Category:      "IPA",
			Rating:        "4.2",
			FoodPairing:   "Spicy foods",
			SubCategory_1: "American IPA",
			SubCategory_2: "",
			Description:   "A classic New England IPA.",
		},
		{
			Sku:           "SKU002",
			Name:          "Guinness Stout",
			Brewery:       "Guinness",
			Region:        "Dublin",
			Country:       "Ireland",
			Abv:           "4.20",
			Ibu:           "45.00",
			Category:      "Stout",
			Rating:        "4.5",
			FoodPairing:   "Oysters",
			SubCategory_1: "Dry Stout",
			SubCategory_2: "",
			Description:   "The iconic Irish dry stout.",
		},
	}
}

// ─── GenerateFilename ────────────────────────────────────────────────────────

func TestGenerateFilename_Format(t *testing.T) {
	name := GenerateFilename("james", "favs", "json")

	if !strings.HasPrefix(name, "james_favs_") {
		t.Errorf("expected prefix 'james_favs_', got %q", name)
	}
	if !strings.HasSuffix(name, ".json") {
		t.Errorf("expected suffix '.json', got %q", name)
	}
}

func TestGenerateFilename_SanitizesUsername(t *testing.T) {
	name := GenerateFilename("James Macean", "favs", "csv")

	if strings.Contains(name, " ") {
		t.Errorf("expected no spaces in filename, got %q", name)
	}
	if !strings.HasPrefix(name, "james_macean_") {
		t.Errorf("expected lowercase username with underscores, got %q", name)
	}
}

// ─── JSONExporter ────────────────────────────────────────────────────────────

func TestJSONExporter_Metadata(t *testing.T) {
	e := NewJSONExporter()
	if e.FileExtension() != "json" {
		t.Errorf("expected 'json', got %q", e.FileExtension())
	}
	if e.FormatName() != "JSON" {
		t.Errorf("expected 'JSON', got %q", e.FormatName())
	}
}

func TestJSONExporter_Export_WritesValidJSON(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "out.json")

	e := NewJSONExporter()
	if err := e.Export(testBeers(), filename); err != nil {
		t.Fatalf("Export returned unexpected error: %v", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	var beers []models.Beer
	if err := json.Unmarshal(data, &beers); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(beers) != 2 {
		t.Errorf("expected 2 beers, got %d", len(beers))
	}
	if beers[0].Name != "Berkshire IPA" {
		t.Errorf("expected first beer 'Berkshire IPA', got %q", beers[0].Name)
	}
}

func TestJSONExporter_Export_EmptySlice(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "empty.json")

	e := NewJSONExporter()
	if err := e.Export([]models.Beer{}, filename); err != nil {
		t.Fatalf("Export returned unexpected error: %v", err)
	}

	data, _ := os.ReadFile(filename)
	var beers []models.Beer
	if err := json.Unmarshal(data, &beers); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(beers) != 0 {
		t.Errorf("expected empty JSON array, got %d beers", len(beers))
	}
}

func TestJSONExporter_Export_InvalidPath(t *testing.T) {
	e := NewJSONExporter()
	err := e.Export(testBeers(), "/nonexistent/path/out.json")
	if err == nil {
		t.Error("expected error for invalid file path, got nil")
	}
}

// ─── CSVExporter ─────────────────────────────────────────────────────────────

func TestCSVExporter_Metadata(t *testing.T) {
	e := NewCSVExporter()
	if e.FileExtension() != "csv" {
		t.Errorf("expected 'csv', got %q", e.FileExtension())
	}
	if e.FormatName() != "CSV" {
		t.Errorf("expected 'CSV', got %q", e.FormatName())
	}
}

func TestCSVExporter_Export_WritesHeaderAndRows(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "out.csv")

	e := NewCSVExporter()
	if err := e.Export(testBeers(), filename); err != nil {
		t.Fatalf("Export returned unexpected error: %v", err)
	}

	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("failed to open output file: %v", err)
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatalf("output is not valid CSV: %v", err)
	}

	// header row + 2 data rows
	if len(records) != 3 {
		t.Fatalf("expected 3 rows (header + 2 beers), got %d", len(records))
	}

	header := records[0]
	if header[0] != "SKU" || header[1] != "Name" {
		t.Errorf("unexpected header: %v", header)
	}

	// Name column is index 1
	if records[1][1] != "Berkshire IPA" {
		t.Errorf("expected first data row name 'Berkshire IPA', got %q", records[1][1])
	}
}

func TestCSVExporter_Export_EmptySlice(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "empty.csv")

	e := NewCSVExporter()
	if err := e.Export([]models.Beer{}, filename); err != nil {
		t.Fatalf("Export returned unexpected error: %v", err)
	}

	f, _ := os.Open(filename)
	defer f.Close()

	records, _ := csv.NewReader(f).ReadAll()
	// only the header row should be present
	if len(records) != 1 {
		t.Errorf("expected only header row for empty beer slice, got %d rows", len(records))
	}
}

func TestCSVExporter_Export_InvalidPath(t *testing.T) {
	e := NewCSVExporter()
	err := e.Export(testBeers(), "/nonexistent/path/out.csv")
	if err == nil {
		t.Error("expected error for invalid file path, got nil")
	}
}

// ─── TXTExporter ─────────────────────────────────────────────────────────────

func TestTXTExporter_Metadata(t *testing.T) {
	e := NewTXTExporter()
	if e.FileExtension() != "txt" {
		t.Errorf("expected 'txt', got %q", e.FileExtension())
	}
	if e.FormatName() != "TXT" {
		t.Errorf("expected 'TXT', got %q", e.FormatName())
	}
}

func TestTXTExporter_Export_ContainsBeerFields(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "out.txt")

	e := NewTXTExporter()
	if err := e.Export(testBeers(), filename); err != nil {
		t.Fatalf("Export returned unexpected error: %v", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	content := string(data)
	checks := []string{"Berkshire IPA", "Berkshire Brewing", "Guinness Stout", "Guinness"}
	for _, want := range checks {
		if !strings.Contains(content, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestTXTExporter_Export_EmptySlice(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "empty.txt")

	e := NewTXTExporter()
	if err := e.Export([]models.Beer{}, filename); err != nil {
		t.Fatalf("Export returned unexpected error: %v", err)
	}

	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("expected file to be created: %v", err)
	}
	if info.Size() != 0 {
		t.Errorf("expected empty file for zero beers, got %d bytes", info.Size())
	}
}

func TestTXTExporter_Export_InvalidPath(t *testing.T) {
	e := NewTXTExporter()
	err := e.Export(testBeers(), "/nonexistent/path/out.txt")
	if err == nil {
		t.Error("expected error for invalid file path, got nil")
	}
}

// ─── Exporter interface compliance ───────────────────────────────────────────

func TestExporterInterface_AllImplemented(t *testing.T) {
	exporters := []Exporter{
		NewJSONExporter(),
		NewCSVExporter(),
		NewTXTExporter(),
	}
	for _, e := range exporters {
		if e.FileExtension() == "" {
			t.Errorf("%T: FileExtension() returned empty string", e)
		}
		if e.FormatName() == "" {
			t.Errorf("%T: FormatName() returned empty string", e)
		}
	}
}

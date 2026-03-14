package storage

import (
	"os"
	"testing"
	"time"

	"github.com/jamesBoder/BeerInfoApp.git/internal/models"
)

// testBeers returns a reusable slice of beers for test fixtures.
func testBeers() []models.Beer {
	return []models.Beer{
		{Sku: "SKU001", Name: "Berkshire IPA", Brewery: "Berkshire Brewing", Country: "USA", Abv: "6.5"},
		{Sku: "SKU002", Name: "Guinness Stout", Brewery: "Guinness", Country: "Ireland", Abv: "4.2"},
	}
}

// ─── getUserFavoritesFilename ─────────────────────────────────────────────────

func TestGetUserFavoritesFilename_Normal(t *testing.T) {
	got := getUserFavoritesFilename("james")
	want := "favorites_james.json"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestGetUserFavoritesFilename_EmptyUsername(t *testing.T) {
	got := getUserFavoritesFilename("")
	want := "favorites_guest.json"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestGetUserFavoritesFilename_UppercaseAndSpaces(t *testing.T) {
	got := getUserFavoritesFilename("James Macean")
	want := "favorites_james_macean.json"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

// ─── Favorites round-trip ─────────────────────────────────────────────────────

func TestSaveAndLoadFavorites(t *testing.T) {
	t.Chdir(t.TempDir())

	favs := models.Favorites{Beers: testBeers()}

	if err := SaveFavorites("testuser", favs); err != nil {
		t.Fatalf("SaveFavorites returned error: %v", err)
	}

	loaded, err := LoadFavorites("testuser")
	if err != nil {
		t.Fatalf("LoadFavorites returned error: %v", err)
	}

	if len(loaded.Beers) != len(favs.Beers) {
		t.Fatalf("expected %d beers, got %d", len(favs.Beers), len(loaded.Beers))
	}
	for i, b := range favs.Beers {
		if loaded.Beers[i].Name != b.Name {
			t.Errorf("beer[%d]: expected Name %q, got %q", i, b.Name, loaded.Beers[i].Name)
		}
		if loaded.Beers[i].Sku != b.Sku {
			t.Errorf("beer[%d]: expected Sku %q, got %q", i, b.Sku, loaded.Beers[i].Sku)
		}
	}
}

func TestSaveAndLoadFavorites_EmptyList(t *testing.T) {
	t.Chdir(t.TempDir())

	favs := models.Favorites{Beers: []models.Beer{}}
	if err := SaveFavorites("emptyuser", favs); err != nil {
		t.Fatalf("SaveFavorites returned error: %v", err)
	}

	loaded, err := LoadFavorites("emptyuser")
	if err != nil {
		t.Fatalf("LoadFavorites returned error: %v", err)
	}
	if len(loaded.Beers) != 0 {
		t.Errorf("expected 0 beers, got %d", len(loaded.Beers))
	}
}

func TestSaveAndLoadFavorites_OverwritesExisting(t *testing.T) {
	t.Chdir(t.TempDir())

	first := models.Favorites{Beers: testBeers()}
	SaveFavorites("overuser", first)

	second := models.Favorites{Beers: []models.Beer{
		{Sku: "SKU999", Name: "New Beer", Brewery: "New Brewery"},
	}}
	if err := SaveFavorites("overuser", second); err != nil {
		t.Fatalf("second SaveFavorites returned error: %v", err)
	}

	loaded, err := LoadFavorites("overuser")
	if err != nil {
		t.Fatalf("LoadFavorites returned error: %v", err)
	}
	if len(loaded.Beers) != 1 || loaded.Beers[0].Name != "New Beer" {
		t.Errorf("expected only 'New Beer', got %+v", loaded.Beers)
	}
}

func TestLoadFavorites_FileNotFound(t *testing.T) {
	t.Chdir(t.TempDir())

	_, err := LoadFavorites("nonexistentuser123")
	if err == nil {
		t.Error("expected error when loading non-existent favorites, got nil")
	}
	if !os.IsNotExist(err) {
		t.Errorf("expected os.IsNotExist error, got: %v", err)
	}
}

// ─── Search history round-trip ───────────────────────────────────────────────

func TestSaveAndLoadSearchHistory(t *testing.T) {
	t.Chdir(t.TempDir())

	history := []models.SearchHistoryEntry{
		{
			Term:      "IPA",
			Timestamp: time.Now().UTC().Truncate(time.Second),
			Results:   testBeers(),
		},
		{
			Term:      "Stout",
			Timestamp: time.Now().UTC().Truncate(time.Second),
			Results:   []models.Beer{testBeers()[1]},
		},
	}

	if err := SaveSearchHistory("histuser", history); err != nil {
		t.Fatalf("SaveSearchHistory returned error: %v", err)
	}

	loaded, err := LoadSearchHistory("histuser")
	if err != nil {
		t.Fatalf("LoadSearchHistory returned error: %v", err)
	}

	if len(loaded) != len(history) {
		t.Fatalf("expected %d entries, got %d", len(history), len(loaded))
	}
	for i, entry := range history {
		if loaded[i].Term != entry.Term {
			t.Errorf("entry[%d]: expected Term %q, got %q", i, entry.Term, loaded[i].Term)
		}
		if len(loaded[i].Results) != len(entry.Results) {
			t.Errorf("entry[%d]: expected %d results, got %d", i, len(entry.Results), len(loaded[i].Results))
		}
	}
}

func TestSaveAndLoadSearchHistory_EmptyHistory(t *testing.T) {
	t.Chdir(t.TempDir())

	if err := SaveSearchHistory("emptyhistuser", []models.SearchHistoryEntry{}); err != nil {
		t.Fatalf("SaveSearchHistory returned error: %v", err)
	}

	loaded, err := LoadSearchHistory("emptyhistuser")
	if err != nil {
		t.Fatalf("LoadSearchHistory returned error: %v", err)
	}
	if len(loaded) != 0 {
		t.Errorf("expected 0 entries, got %d", len(loaded))
	}
}

func TestSaveAndLoadSearchHistory_OverwritesExisting(t *testing.T) {
	t.Chdir(t.TempDir())

	first := []models.SearchHistoryEntry{{Term: "IPA", Timestamp: time.Now()}}
	SaveSearchHistory("overhistuser", first)

	second := []models.SearchHistoryEntry{
		{Term: "Lager", Timestamp: time.Now()},
		{Term: "Porter", Timestamp: time.Now()},
	}
	if err := SaveSearchHistory("overhistuser", second); err != nil {
		t.Fatalf("second SaveSearchHistory returned error: %v", err)
	}

	loaded, err := LoadSearchHistory("overhistuser")
	if err != nil {
		t.Fatalf("LoadSearchHistory returned error: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 entries after overwrite, got %d", len(loaded))
	}
	if loaded[0].Term != "Lager" {
		t.Errorf("expected first entry 'Lager', got %q", loaded[0].Term)
	}
}

func TestLoadSearchHistory_FileNotFound(t *testing.T) {
	t.Chdir(t.TempDir())

	_, err := LoadSearchHistory("nonexistenthistuser123")
	if err == nil {
		t.Error("expected error when loading non-existent history, got nil")
	}
	if !os.IsNotExist(err) {
		t.Errorf("expected os.IsNotExist error, got: %v", err)
	}
}

// ─── Username normalization on filenames ──────────────────────────────────────

func TestSaveFavorites_UsernameNormalization(t *testing.T) {
	t.Chdir(t.TempDir())

	favs := models.Favorites{Beers: testBeers()}

	// save with mixed case + spaces
	if err := SaveFavorites("James Macean", favs); err != nil {
		t.Fatalf("SaveFavorites returned error: %v", err)
	}

	// file should be on disk as favorites_james_macean.json
	if _, err := os.Stat("favorites_james_macean.json"); os.IsNotExist(err) {
		t.Error("expected normalized filename 'favorites_james_macean.json' to exist on disk")
	}
}

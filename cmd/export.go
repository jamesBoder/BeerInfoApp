package cmd

import (
	"fmt"
	"strings"

	"github.com/jamesBoder/BeerInfoApp.git/internal/export"
	"github.com/jamesBoder/BeerInfoApp.git/internal/storage"
	"github.com/jamesBoder/BeerInfoApp.git/internal/ui"
)

// ExportFavsCommand - implements cli.Command
type ExportFavsCommand struct {
	ctx *CommandContext
}

// NewExportFavsCommand returns a instance of ExportFavsCommand
func NewExportFavsCommand(ctx *CommandContext) *ExportFavsCommand {
	return &ExportFavsCommand{ctx: ctx}
}

// name method
func (c *ExportFavsCommand) Name() string {
	return "export favs"
}

// Description method
func (c *ExportFavsCommand) Description() string {
	return "Export favorite beers to a file [JSON/CSV/TXT]"
}

// usage method
func (c *ExportFavsCommand) Usage() string {
	return "export favs <filename>"
}

// Execute method
func (c *ExportFavsCommand) Execute(args []string) error {
	// validate format argument
	if len(args) == 0 {
		fmt.Println(ui.Error("Error: Export format is required (json/csv/txt)"))
		fmt.Println(ui.Tip("\nUsage: export favs <format>"))
		fmt.Println(ui.Example("\nExamples:"))
		fmt.Println(ui.Info("   export favs json"))
		fmt.Println(ui.Info("   export favs csv"))
		fmt.Println(ui.Info("   export favs txt"))
		return nil
	}

	// load favorites from storage
	favorites, err := storage.LoadFavorites(c.ctx.State.Config.User)
	if err != nil {
		fmt.Println(ui.Error("error loading favorites", err))
		return nil
	}

	// get appropriate exporter using factory
	format := strings.ToLower(args[0])
	var exporter export.Exporter

	switch format {
	case "json":
		exporter = export.NewJSONExporter()
	case "csv":
		exporter = export.NewCSVExporter()
	case "txt":
		exporter = export.NewTXTExporter()
	default:
		fmt.Println(ui.Error("Error: Unsupported export format %q. Use 'json', 'csv', or 'txt'.", format))
		return nil
	}

	// generate filename
	filename := export.GenerateFilename(c.ctx.State.Config.User, "favorites", format)

	// call the exporter.Export() method
	err = exporter.Export(favorites.Beers, filename)
	if err != nil {
		fmt.Println(ui.Error("Error: Failed to export favorites: %v", err))
		return nil
	}

	// display success msg with file location
	fmt.Println(ui.Success("Favorites exported successfully to %s", filename))
	return nil
}

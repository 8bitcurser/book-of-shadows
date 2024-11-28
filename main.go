package main

import (
	"book-of-shadows/models"
	"book-of-shadows/views"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"path/filepath"
)

func handleHome(c echo.Context) error {
	return views.Home().Render(c.Request().Context(), c.Response().Writer)
}

func handleGenerate(c echo.Context) error {
	mode := models.Classic
	if c.QueryParam("mode") == "pulp" {
		mode = models.Pulp
	}

	investigator := models.NewInvestigator(mode)

	// Set content type to ensure proper rendering
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return views.CharacterSheet(investigator).Render(c.Request().Context(), c.Response().Writer)
}

func handleGetJSON(c echo.Context) error {
	mode := models.Classic
	if c.QueryParam("mode") == "pulp" {
		mode = models.Pulp
	}

	investigator := models.NewInvestigator(mode)
	return c.JSON(http.StatusOK, investigator)
}

func handleExportPDF(c echo.Context) error {
	//mode := models.Classic
	//if c.QueryParam("mode") == "pulp" {
	//	mode = models.Pulp
	//}

	investigator := models.NewInvestigator(models.Pulp)
	PDFExport(
		"/Users/tommyboy/Projects/book-of-shadows/modernSheet.pdf",
		"/Users/tommyboy/Projects/book-of-shadows/filledModernSheet.pdf",
		investigator)

	return c.File("filledModernSheet.pdf")
}

func main() {
	// Update the PYTHONPATH in your Go code
	os.Setenv("PYTHONPATH", filepath.Join(os.Getenv("PYTHONPATH"), "./scripts"))
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Static files
	e.Static("/static", "static")

	// Routes
	e.GET("/", handleHome)
	e.GET("/api/generate", handleGenerate)
	e.POST("/api/export-pdf", handleExportPDF)
	e.POST("/api/get-json", handleGetJSON)

	e.Logger.Fatal(e.Start(":8080"))
}

package main

import (
	"book-of-shadows/models"
	"book-of-shadows/views"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func handleHome(c echo.Context) error {
	return views.Home().Render(c.Request().Context(), c.Response().Writer)
}

func handleGenerate(c echo.Context) error {
	mode := models.Pulp
	if c.QueryParam("mode") == "classic" {
		mode = models.Classic
	}

	investigator := models.NewInvestigator(mode)

	// Set content type to ensure proper rendering
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return views.CharacterSheet(investigator).Render(c.Request().Context(), c.Response().Writer)
}

func handleGetJSON(c echo.Context) error {
	mode := models.Pulp
	if c.QueryParam("mode") == "classic" {
		mode = models.Classic
	}

	investigator := models.NewInvestigator(mode)
	return c.JSON(http.StatusOK, investigator)
}

func handleExportPDF(c echo.Context) error {
	mode := models.Pulp
	if c.QueryParam("mode") == "classic" {
		mode = models.Classic
	}

	investigator := models.NewInvestigator(mode)
	investigatorPDF := "./static/" + investigator.Name + ".pdf"
	err := PDFExport(
		"./static/modernSheet.pdf",
		investigatorPDF,
		investigator)
	if err != nil {
		return err
	}
	return c.File(investigatorPDF)
}

func main() {
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
	e.GET("/api/export-pdf", handleExportPDF)
	e.POST("/api/get-json", handleGetJSON)

	e.Logger.Fatal(e.Start(":8080"))
}

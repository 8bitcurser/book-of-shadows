package main

//func main() {
//	investigator := NewInvestigator(Pulp)
//	err := PDFExport(investigator)
//	if err != nil {
//		panic(err)
//	}
//}

import (
	"book-of-shadows/models"
	"book-of-shadows/views"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

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
	mode := models.Classic
	if c.QueryParam("mode") == "pulp" {
		mode = models.Pulp
	}

	investigator := models.NewInvestigator(mode)
	err := PDFExport(investigator)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate PDF",
		})
	}

	return c.File("character_sheet.pdf") // Adjust filename as needed
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
	e.POST("/api/export-pdf", handleExportPDF)
	e.POST("/api/get-json", handleGetJSON)

	e.Logger.Fatal(e.Start(":8080"))
}

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"book-of-shadows/internal/errors"
	"book-of-shadows/models"
)

// ProcessingOptions defines the options for PDF processing
type ProcessingOptions struct {
	InputPath  string            `json:"input_path"`
	OutputPath string            `json:"output_path"`
	Metadata   map[string]string `json:"metadata"`
}

// PdfProcessor handles PDF generation via Python script
type PdfProcessor struct {
	pythonPath string
	scriptPath string
}

// NewPdfProcessor creates a new PDF processor
func NewPdfProcessor() *PdfProcessor {
	return &PdfProcessor{
		pythonPath: "python3",
		scriptPath: filepath.Join("scripts", "exporter.py"),
	}
}

// ProcessPdf processes a PDF with the given options
func (p *PdfProcessor) ProcessPdf(options ProcessingOptions) error {
	optionsJSON, err := json.Marshal(options)
	if err != nil {
		return fmt.Errorf("failed to marshal options: %w", err)
	}

	cmd := exec.Command(p.pythonPath, p.scriptPath, string(optionsJSON))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run Python script: %w", err)
	}

	return nil
}

// ExportPDF exports an investigator as a PDF
func (h *Handler) ExportPDF(w http.ResponseWriter, r *http.Request) {
	// Extract ID from context
	params := r.Context().Value("params").([]string)
	if len(params) == 0 {
		h.respondError(w, errors.NewHTTPError(http.StatusBadRequest, "Missing investigator ID", nil))
		return
	}
	id := params[0]

	// Get investigator
	investigator, err := h.store.GetInvestigator(r, id)
	if err != nil {
		h.respondError(w, err)
		return
	}

	// Convert to PDF map
	data := convertInvestigatorToMap(investigator)

	// Generate filename
	fileName := fmt.Sprintf("%s.pdf", strings.ReplaceAll(data["Investigators_Name"], " ", "_"))
	outputPath := filepath.Join(".", "static", fileName)

	// Process PDF
	processor := NewPdfProcessor()
	options := ProcessingOptions{
		InputPath:  filepath.Join(".", "static", "modernSheet.pdf"),
		OutputPath: outputPath,
		Metadata:   data,
	}

	if err := processor.ProcessPdf(options); err != nil {
		h.logger.Printf("PDF generation failed: %v", err)
		h.respondError(w, errors.NewHTTPError(http.StatusInternalServerError, "Error generating PDF", err))
		return
	}

	// Serve the file
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	http.ServeFile(w, r, outputPath)

	// Clean up the generated file
	go func() {
		if err := os.Remove(outputPath); err != nil {
			h.logger.Printf("Failed to remove temporary PDF: %v", err)
		}
	}()
}

// convertInvestigatorToMap converts an investigator to a map for PDF export
func convertInvestigatorToMap(investigator *models.Investigator) map[string]string {
	data := make(map[string]string)

	// Handle Attributes
	for key, attr := range investigator.Attributes {
		data[attr.Name] = strconv.Itoa(attr.Value)
		data[attr.Name+"_half"] = strconv.Itoa(attr.Value / 2)
		data[attr.Name+"_fifth"] = strconv.Itoa(attr.Value / 5)

		// Add Starting/Max values for HP, Magic, and Sanity
		// Note: StartingValue is not serialized (json:"-"), so we use MaxValue
		if key == models.AttrHitPoints {
			data["StartingHP"] = strconv.Itoa(attr.MaxValue)
		} else if key == models.AttrMagicPoints {
			data["StartingMagic"] = strconv.Itoa(attr.MaxValue)
		} else if key == models.AttrSanity {
			data["StartingSanity"] = strconv.Itoa(attr.Value) // Starting sanity = current sanity at creation
			data["MaxSanity"] = strconv.Itoa(attr.MaxValue)
		}
	}

	// Handle Skills
	for _, skill := range investigator.Skills {
		if skill.Base == 1 {
			continue
		}
		formField := skill.FormName
		if skill.NeedsFormDef == 1 {
			data["SkillDef_"+formField] = skill.Name
		}
		if skill.IsSelected {
			data["Skill_"+formField+"_Chk"] = "1"
		}
		data["Skill_"+formField] = strconv.Itoa(skill.Value)
		data["Skill_"+formField+"_half"] = strconv.Itoa(skill.Value / 2)
		data["Skill_"+formField+"_fifth"] = strconv.Itoa(skill.Value / 5)
	}

	// Handle other fields
	data["Investigators_Name"] = investigator.Name
	data["Occupation"] = investigator.Occupation.Name
	data["Age"] = strconv.Itoa(investigator.Age)
	data["Residence"] = investigator.Residence
	data["Birthplace"] = investigator.Birthplace
	data["MOV"] = strconv.Itoa(investigator.Move)
	data["DamageBonus"] = investigator.DamageBonus
	data["Build"] = investigator.Build
	data["Archetype"] = investigator.Archetype.Name

	// Handle talents
	var talents strings.Builder
	for i, talent := range investigator.Talents {
		if i > 0 {
			talents.WriteString(", ")
		}
		talents.WriteString(talent.Name)
	}
	data["Pulp Talents"] = talents.String()

	return data
}
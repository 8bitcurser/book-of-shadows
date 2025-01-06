package main

import (
	"book-of-shadows/models"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type ProcessingOptions struct {
	InputPath  string            `json:"input_path"`
	OutputPath string            `json:"output_path"`
	Metadata   map[string]string `json:"metadata"`
}

type PdfProcessor struct {
	pythonPath string
	scriptPath string
}

func NewPdfProcessor() (*PdfProcessor, error) {
	// Find Python executable
	pythonPath := filepath.Join("python3")

	// Set script path
	scriptPath := filepath.Join("scripts", "exporter.py")

	return &PdfProcessor{
		pythonPath: pythonPath,
		scriptPath: scriptPath,
	}, nil
}

func (p *PdfProcessor) ProcessPdf(options ProcessingOptions) error {
	// Convert options to JSON string
	optionsJSON, err := json.Marshal(options)
	if err != nil {
		return fmt.Errorf("failed to marshal options: %v", err)
	}

	// Run Python script with JSON data as argument
	cmd := exec.Command(p.pythonPath, p.scriptPath, string(optionsJSON))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run Python script: %v", err)
	}

	return nil
}

func PDFExport(input, output string, payload map[string]string) error {
	data := payload
	processor, err := NewPdfProcessor()
	if err != nil {
		return fmt.Errorf("failed to initialize processor: %v", err)
	}

	options := ProcessingOptions{
		InputPath:  input,
		OutputPath: output,
		Metadata:   data,
	}

	if err := processor.ProcessPdf(options); err != nil {
		return fmt.Errorf("failed to process PDF: %v", err)
	}

	return nil
}

func ConvertInvestigatorToMap(investigator *models.Investigator) map[string]string {
	data := make(map[string]string)
	// Handle Attributes
	for key, attr := range investigator.Attributes {
		attrName := investigator.Attributes[key]
		data[attrName.Name] = strconv.Itoa(attr.Value)
		data[attrName.Name+"_half"] = strconv.Itoa(attr.Value / 2)
		data[attrName.Name+"_fifth"] = strconv.Itoa(attr.Value / 5)
	}

	// Handle Skills
	for _, skill := range investigator.Skills {
		formField := skill.FormName
		if skill.NeedsFormDef == 1 {
			data["SkillDef_"+formField] = skill.Name
		}
		if skill.IsSelected == true {
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
	talents := ""
	for _, talent := range investigator.Talents {
		talents += talent.Name + ", "
	}
	data["Pulp Talents"] = talents
	return data
}

package main

import (
	"book-of-shadows/models"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/fjson"
	"github.com/unidoc/unipdf/v3/model"
	"log"
	"os"
	"strconv"
)

func init() {
	// Load the .env file.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the API key from the environment variable
	apiKey := os.Getenv("UNIDOC_LICENSE_API_KEY")
	if apiKey == "" {
		log.Fatal("UNIDOC_LICENSE_API_KEY not set in .env file")
	}

	// Set the metered license
	err = license.SetMeteredKey(apiKey)
	if err != nil {
		log.Fatalf("Error setting metered key: %v", err)
	}
}

type PDFField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func ConvertInvestigatorToPDFFields(investigator *models.Investigator) ([]byte, error) {
	pdfFields := []PDFField{}

	// Handle Attributes
	for key, attr := range investigator.Attributes {
		attrName := investigator.Attributes[key]
		pdfFields = append(pdfFields, PDFField{
			Name:  attrName.Name,
			Value: strconv.Itoa(attr.Value),
		})
		// Add fields for half and fifth values
		pdfFields = append(pdfFields, PDFField{
			Name:  attrName.Name + "_half",
			Value: strconv.Itoa(attr.Value / 2),
		})
		pdfFields = append(pdfFields, PDFField{
			Name:  attrName.Name + "_fifth",
			Value: strconv.Itoa(attr.Value / 5),
		})
	}

	// Handle Skills
	for name, skill := range investigator.Skills {
		if name == "Dodge_Copy" {
			pdfFields = append(pdfFields, PDFField{
				Name:  name,
				Value: strconv.Itoa(skill.Value),
			})
			pdfFields = append(pdfFields, PDFField{
				Name:  name + "_half",
				Value: strconv.Itoa(skill.Value / 2),
			})
			pdfFields = append(pdfFields, PDFField{
				Name:  name + "_fifth",
				Value: strconv.Itoa(skill.Value / 5),
			})
		} else if name == "FastTalk" {
			pdfFields = append(pdfFields, PDFField{
				Name:  "Skill_" + name,
				Value: strconv.Itoa(skill.Value),
			})
			pdfFields = append(pdfFields, PDFField{
				Name:  "Skill_" + name + " _half",
				Value: strconv.Itoa(skill.Value / 2),
			})
			pdfFields = append(pdfFields, PDFField{
				Name:  "Skill_" + name + " _fifth",
				Value: strconv.Itoa(skill.Value / 5),
			})
		} else {
			pdfFields = append(pdfFields, PDFField{
				Name:  "Skill_" + name,
				Value: strconv.Itoa(skill.Value),
			})
			// Add fields for half and fifth values
			pdfFields = append(pdfFields, PDFField{
				Name:  "Skill_" + name + "_half",
				Value: strconv.Itoa(skill.Value / 2),
			})
			pdfFields = append(pdfFields, PDFField{
				Name:  "Skill_" + name + "_fifth",
				Value: strconv.Itoa(skill.Value / 5),
			})
		}
	}

	// Handle other fields
	pdfFields = append(pdfFields, PDFField{Name: "Investigators_Name", Value: investigator.Name})
	pdfFields = append(pdfFields, PDFField{Name: "Occupation", Value: investigator.Occupation.Name})
	pdfFields = append(pdfFields, PDFField{Name: "Age", Value: strconv.Itoa(investigator.Age)})
	pdfFields = append(pdfFields, PDFField{Name: "Residence", Value: investigator.Residence})
	pdfFields = append(pdfFields, PDFField{Name: "Birthplace", Value: investigator.Birthplace})
	pdfFields = append(pdfFields, PDFField{Name: "MOV", Value: strconv.Itoa(investigator.Move)})
	pdfFields = append(pdfFields, PDFField{Name: "DamageBonus", Value: investigator.DamageBonus})
	pdfFields = append(pdfFields, PDFField{Name: "Build", Value: investigator.Build})
	pdfFields = append(pdfFields, PDFField{Name: "Archetype", Value: investigator.Archetype.Name})
	talents := ""
	for _, talent := range investigator.Talents {
		talents += talent.Name + ", "
	}
	pdfFields = append(pdfFields, PDFField{Name: "Pulp Talents", Value: talents})
	// Add more fields as needed
	return json.MarshalIndent(pdfFields, "", "  ")
}

func PDFExport(investigator *models.Investigator) error {
	// Convert investigator to JSON and save to a file
	_, err := json.MarshalIndent(investigator, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling investigator to JSON: %v", err)
	}
	pdfFieldsJSON, err := ConvertInvestigatorToPDFFields(investigator)
	if err != nil {
		fmt.Printf("Error converting to PDF fields: %v\n", err)
		return err
	}
	//
	jsonPath := "investigator_data.json"
	err = os.WriteFile(jsonPath, pdfFieldsJSON, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON to file: %v", err)
	}

	inputPath := "modernSheet.pdf"
	fdata, err := fjson.LoadFromPDFFile(inputPath)
	if err != nil {
		return fmt.Errorf("error loading input file: %v", err)
	}
	_, err = fdata.JSON()
	if err != nil {
		return fmt.Errorf("error marshalling input file: %v", err)
	}
	outputPath := "filled_modernSheet.pdf"

	err = fillFields(inputPath, jsonPath, outputPath)
	if err != nil {
		return fmt.Errorf("error filling fields: %v", err)
	}

	fmt.Printf("Success, output written to %s\n", outputPath)
	return nil
}

func fillFields(inputPath, jsonPath, outputPath string) error {
	_, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return err
	}

	fdata, err := fjson.LoadFromJSONFile(jsonPath)
	if err != nil {
		fmt.Printf("Error loading JSON: %v\n", err)
		return err
	}

	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}
	fieldAppearance := annotator.FieldAppearance{OnlyIfMissing: true, RegenerateTextFields: true}
	// Populate the form data.
	err = pdfReader.AcroForm.FillWithAppearance(fdata, fieldAppearance)
	if err != nil {
		return err
	}

	opt := &model.ReaderToWriterOpts{
		SkipAcroForm: false,
	}

	pdfWriter, err := pdfReader.ToWriter(opt)
	if err != nil {
		return err
	}

	if err := pdfWriter.WriteToFile(outputPath); err != nil {
		return err
	}
	return nil
}

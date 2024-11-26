package main

func main() {
	investigator := NewInvestigator(Pulp)
	err := PDFExport(investigator)
	if err != nil {
		panic(err)
	}
}

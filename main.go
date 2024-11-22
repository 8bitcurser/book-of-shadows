package main

type Era int

const (
	Twenties Era = iota
	Modern
)

type Occupation struct {
	Name   string
	Skills []string
}

type Archetype struct {
	Name   string
	Skills []string
}

type ProfilePic struct {
	FilePath string
	FileName string
}

type Skill struct {
	Name         string
	Abbreviation string
	Default      int
	Value        int
	Era          []Era
	Base         int
}

type Attribute struct {
	Name            string
	StartingValue   int
	Value           int
	MaxValue        int
	AssociatedSkill Skill
}

type Investigator struct {
	Era              Era
	Name             string `json:"name"`
	Residence        string `json:"residence"`
	Birthplace       string `json:"birthplace"`
	Age              int    `json:"age"`
	ProfilePic       ProfilePic
	Occupation       Occupation
	Archetype        Archetype
	Insane           bool `json:"insane"`
	TemporaryInsane  bool `json:"temporary_insane"`
	IndefiniteInsane bool `json:"indefinite_insane"`
	MajorWound       bool `json:"major_wound"`
	Unconscious      bool `json:"unconscious"`
	Dying            bool `json:"dying"`
	Strength         Attribute
	Constitution     Attribute
	Dexterity        Attribute
	Intelligence     Attribute
	Size             Attribute
	Power            Attribute
	Appearance       Attribute
	Education        Attribute
	HitPoints        Attribute
	MagicPoints      Attribute
	Luck             Attribute
	Sanity           Attribute
	Skills           map[string]Skill
	Move             int
	Build            string
	DamageBonus      string
}

func (*Investigator) NewInvestigator() Investigator {
	return Investigator{
		Era:        1,
		Name:       "John Doe",
		Residence:  "Boston",
		Birthplace: "Dallas TX",
		Age:        37,
		ProfilePic: ProfilePic{"/sample/path/env", "profile"},
		Occupation: Occupation{
			Name:   "Adventurer",
			Skills: []string{"Firearms (Handgun)", "Archaeology"},
		},
		Archetype: Archetype{
			Name:   "Indiana Jones",
			Skills: []string{"History"},
		},
		Insane:           false,
		TemporaryInsane:  false,
		IndefiniteInsane: false,
		MajorWound:       false,
		Unconscious:      false,
		Dying:            false,
		Strength:         Attribute{},
		Constitution:     Attribute{},
		Dexterity:        Attribute{},
		Intelligence:     Attribute{},
		Size:             Attribute{},
		Power:            Attribute{},
		Appearance:       Attribute{},
		Education:        Attribute{},
		HitPoints:        Attribute{},
		MagicPoints:      Attribute{},
		Sanity:           Attribute{},
		Skills:           BaseModernSkills,
		Move:             2,
		Build:            "Big",
		DamageBonus:      "1D4",
	}
}

//func handleConnection(w http.ResponseWriter, r *http.Request) {
//
//	// Log the path that was requested
//	fmt.Printf("Method: %s\nPath: %s\n", r.Method, r.URL.Path)
//
//	// Send response
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(map[string]string{"message": "Got it"})
//}
//
//func main() {
//	fmt.Println("Server listening on :8080")
//	http.HandleFunc("/", handleConnection)
//	if err := http.ListenAndServe(":8080", nil); err != nil {
//		panic(err)
//	}
//}

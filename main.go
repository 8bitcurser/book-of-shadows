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
	Name          string
	StartingValue int
	Value         int
	MaxValue      int
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
	STR              Attribute
	CON              Attribute
	DEX              Attribute
	INT              Attribute
	SIZ              Attribute
	PWR              Attribute
	APP              Attribute
	EDU              Attribute
	HP               Attribute
	MP               Attribute
	LCK              Attribute
	SAN              Attribute
	Skills           map[string]Skill
	Move             int
	Build            string
	DamageBonus      string
}

func (*Investigator) NewInvestigator() *Investigator {
	inv := Investigator{
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
		STR: Attribute{
			Name:          "Strength",
			StartingValue: 0,
			Value:         0,
			MaxValue:      0,
		},
		CON: Attribute{
			Name:          "Constitution",
			StartingValue: 0,
			Value:         0,
			MaxValue:      0,
		},
		DEX: Attribute{
			Name:          "Dexterity",
			StartingValue: 0,
			Value:         0,
			MaxValue:      0,
		},
		INT: Attribute{
			Name:          "Intelligence",
			StartingValue: 0,
			Value:         0,
			MaxValue:      0,
		},
		SIZ: Attribute{
			Name:          "Size",
			StartingValue: 0,
			Value:         0,
			MaxValue:      0,
		},
		PWR: Attribute{
			Name:          "Power",
			StartingValue: 0,
			Value:         0,
			MaxValue:      0,
		},
		APP: Attribute{
			Name:          "Appearance",
			StartingValue: 0,
			Value:         0,
			MaxValue:      0,
		},
		EDU: Attribute{
			Name:          "Education",
			StartingValue: 0,
			Value:         0,
			MaxValue:      0,
		},
		HP: Attribute{
			Name:          "HitPoints",
			StartingValue: 0,
			Value:         0,
			MaxValue:      0,
		},
		MP: Attribute{
			Name:          "MagicPoints",
			StartingValue: 0,
			Value:         0,
			MaxValue:      0,
		},
		LCK: Attribute{
			Name:          "Luck",
			StartingValue: 0,
			Value:         0,
			MaxValue:      0,
		},
		SAN: Attribute{
			Name:          "Sanity",
			StartingValue: 0,
			Value:         0,
			MaxValue:      0,
		},
		Skills:      BaseModernSkills,
		Move:        2,
		Build:       "Big",
		DamageBonus: "1D4",
	}
	inv.Skills["Dodge"] = Skill{
		Name:         "Dodge",
		Abbreviation: "Dodge",
		Default:      inv.DEX.Value / 2,
		Value:        inv.DEX.Value / 2,
	}
	inv.Skills["Idea"] = Skill{
		Name:         "Idea",
		Abbreviation: "Idea",
		Default:      inv.INT.Value / 2,
		Value:        inv.INT.Value / 2,
	}
	inv.Skills["Know"] = Skill{
		Name:         "Know",
		Abbreviation: "Know",
		Default:      inv.EDU.Value / 2,
		Value:        inv.EDU.Value / 2,
	}
	inv.Skills["Language(Own)"] = Skill{
		Name:         "Language(Own)",
		Abbreviation: "Language(Own)",
		Default:      inv.EDU.Value,
		Value:        inv.EDU.Value,
	}

	return &inv
}

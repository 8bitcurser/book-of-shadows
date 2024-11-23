package main

import (
	"fmt"
	"math/rand"
)

type Era int
type GameMode int
type TalentType int

const (
	Twenties Era = iota
	Modern
)

const (
	Classic GameMode = iota
	Pulp
)

const (
	Physical TalentType = iota
	Mental
	Combat
	Miscellaneous
)

type Occupation struct {
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

type Talent struct {
	Name        string
	Description string
	Type        TalentType
}

type Archetype struct {
	Name                 string
	Skills               []string
	BonusPoints          int
	CoreCharacteristic   []Attribute
	SuggestedOccupations []Occupation
	Talents              []Talent
	AmountOfTalents      int
}

func rollD6() int {
	return rand.Intn(6) + 1
}

func (i *Investigator) SetHP() {
	rawHP := i.CON.Value + i.SIZ.Value
	divider := 10
	if i.GameMode == Pulp {
		divider = 5
	}
	hp := rawHP / divider
	i.HP.Value = hp
	i.HP.MaxValue = hp
	i.HP.StartingValue = hp

}

type buildDamageRange struct {
	maxValue    int
	damageBonus string
	build       string
}

var buildDamageTable = []buildDamageRange{
	{64, "-2", "-2"},
	{84, "-1", "-1"},
	{124, "None", ""},
	{164, "+1D4", "+1"},
	{204, "+1D6", "+2"},
	{284, "+2D6", "+3"},
	{364, "+3D6", "+4"},
	{444, "+4D6", "+5"},
	{524, "+5D6", "+6"},
}

func (i *Investigator) SetBuildAndDMG() {
	compoundValue := i.STR.Value + i.SIZ.Value

	// Handle the special case for very high values first
	if compoundValue > 524 {
		extraSteps := (compoundValue - 524) / 80
		i.DamageBonus = fmt.Sprintf("+%dD6", 5+extraSteps)
		i.Build = fmt.Sprintf("+%d", 6+extraSteps)
		return
	}

	// Find the appropriate range in the table
	for _, r := range buildDamageTable {
		if compoundValue <= r.maxValue {
			i.DamageBonus = r.damageBonus
			i.Build = r.build
			return
		}
	}
}

func (i *Investigator) SetMovement() {
	// if both are not greater or lesser than size it means one of the two is
	if i.DEX.Value < i.SIZ.Value && i.STR.Value < i.SIZ.Value {
		i.Move = 7
	} else if i.STR.Value > i.SIZ.Value && i.DEX.Value > i.SIZ.Value {
		i.Move = 9
	} else {
		i.Move = 9
	}
}

func (a *Attribute) Initialize() {
	rolled := 0
	if a.Name == "Size" || a.Name == "Intelligence" || a.Name == "Education" {
		rolled = (rollD6() + rollD6() + 6) * 5
	} else {
		rolled = (rollD6() + rollD6() + rollD6()) * 5
	}

	a.Value = rolled
	a.StartingValue = rolled
}

type Investigator struct {
	Era              Era
	GameMode         GameMode
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
	POW              Attribute
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

func NewInvestigator(mode GameMode) *Investigator {
	inv := Investigator{
		Era:        1,
		GameMode:   mode,
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
		POW: Attribute{
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

	inv.STR.Initialize()
	inv.DEX.Initialize()
	inv.CON.Initialize()
	inv.EDU.Initialize()
	inv.INT.Initialize()
	inv.SIZ.Initialize()
	inv.APP.Initialize()
	inv.LCK.Initialize()
	inv.POW.Initialize()
	// allow re roll
	if inv.LCK.Value < 45 {
		inv.LCK.Initialize()
	}
	inv.SAN.Value = inv.POW.Value
	inv.SAN.StartingValue = inv.POW.StartingValue
	inv.SetHP()
	inv.SetMovement()
	inv.SetBuildAndDMG()
	inv.MP.Value = inv.POW.Value / 5

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

func main() {
	investigator := NewInvestigator(Pulp)
	fmt.Println(investigator)
}

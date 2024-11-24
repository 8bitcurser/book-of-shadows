package main

import (
	"fmt"
	"math/rand"
)

type Era int
type GameMode int

const (
	Twenties Era = iota
	Modern
)

const (
	Classic GameMode = iota
	Pulp
)

type ProfilePic struct {
	FilePath string
	FileName string
}

func rollD6() int {
	return rand.Intn(6) + 1
}

func coreRoll() int {
	return (rollD6() + 13) * 5
}

func (i *Investigator) SetHP() {
	rawHP := i.Attributes[AttrConstitution].Value + i.Attributes[AttrSize].Value
	divider := 10
	if i.GameMode == Pulp {
		divider = 5
	}
	hp := rawHP / divider
	HP := i.Attributes[AttrHitPoints]
	HP.Value = hp
	HP.MaxValue = hp
	HP.StartingValue = hp

}

func (i *Investigator) PickRandomTalents() {
	// ToDo: Need to support archetype talent class or specific talent suggestion
	for j := 0; j < i.Archetype.AmountOfTalents; j++ {
		talentName := TalentsList[rand.Intn(len(TalentsList))]
		i.Talents = append(i.Talents, Talents[talentName])
	}

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
	compoundValue := i.Attributes[AttrStrength].Value + i.Attributes[AttrSize].Value

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
	if i.Attributes[AttrDexterity].Value < i.Attributes[AttrSize].Value && i.Attributes[AttrStrength].Value < i.Attributes[AttrSize].Value {
		i.Move = 7
	} else if i.Attributes[AttrStrength].Value > i.Attributes[AttrSize].Value && i.Attributes[AttrDexterity].Value > i.Attributes[AttrSize].Value {
		i.Move = 9
	} else {
		i.Move = 9
	}
}

func (i *Investigator) InitializeAttributes() {
	// Create a map of all attributes
	// Create a lookup map for core characteristics for O(1) lookup
	coreCharacteristics := make(map[string]bool)
	if i.Archetype != nil {
		// there is one chore characteristic per each character
		pickedCore := rand.Intn(len(i.Archetype.CoreCharacteristic))
		coreCharacteristics[i.Archetype.CoreCharacteristic[pickedCore]] = true
	}
	// Initialize each attribute
	isPulp := i.GameMode == Pulp // or however you check for pulp mode

	for _, attr := range i.Attributes {

		// An attribute is core if we're in pulp mode AND it's in core characteristics
		isCore := isPulp && coreCharacteristics[attr.Name]

		// Initialize the attribute
		attr.Initialize(isCore)
	}

}

func (i *Investigator) CalculateSKillPoints() int {
	formula := i.Occupation.SkillPoints
	points := 0
	for _, skillAttr := range formula.BaseAttributes {
		attr := i.Attributes[skillAttr.Name]
		points += attr.Value * skillAttr.Multiplier
	}
	if len(formula.Options) > 0 {
		picked := rand.Intn(len(formula.Options))
		optional := formula.Options[picked]
		attrOptional := i.Attributes[optional.Name]
		points += attrOptional.Value * optional.Multiplier
	}
	return points
}

type Investigator struct {
	Era              Era
	GameMode         GameMode
	Name             string `json:"name"`
	Residence        string `json:"residence"`
	Birthplace       string `json:"birthplace"`
	Age              int    `json:"age"`
	ProfilePic       ProfilePic
	Occupation       *Occupation
	Archetype        *Archetype
	Insane           bool                 `json:"insane"`
	TemporaryInsane  bool                 `json:"temporary_insane"`
	IndefiniteInsane bool                 `json:"indefinite_insane"`
	MajorWound       bool                 `json:"major_wound"`
	Unconscious      bool                 `json:"unconscious"`
	Dying            bool                 `json:"dying"`
	Attributes       map[string]Attribute `json:"attributes"`
	Skills           map[string]Skill
	Move             int
	Build            string
	DamageBonus      string
	Talents          []Talent
}

func NewInvestigator(mode GameMode) *Investigator {
	inv := Investigator{
		Era:              1,
		GameMode:         mode,
		Name:             "John Doe",
		Residence:        "Boston",
		Birthplace:       "Dallas TX",
		Age:              37,
		ProfilePic:       ProfilePic{"/sample/path/env", "profile"},
		Insane:           false,
		TemporaryInsane:  false,
		IndefiniteInsane: false,
		MajorWound:       false,
		Unconscious:      false,
		Dying:            false,
		Attributes: map[string]Attribute{
			AttrStrength: {
				Name:          AttrStrength,
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrConstitution: {
				Name:          AttrConstitution,
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrDexterity: {
				Name:          AttrDexterity,
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrIntelligence: {
				Name:          AttrIntelligence,
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrSize: {
				Name:          AttrSize,
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrPower: {
				Name:          AttrPower,
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrAppearance: {
				Name:          AttrAppearance,
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrEducation: {
				Name:          AttrEducation,
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrHitPoints: {
				Name:          AttrHitPoints,
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrMagicPoints: {
				Name:          AttrMagicPoints,
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrLuck: {
				Name:          AttrLuck,
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrSanity: {
				Name:          AttrSanity,
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
		},
		Skills:      BaseModernSkills,
		Move:        2,
		Build:       "Big",
		DamageBonus: "1D4",
	}
	// assign archetype
	if mode == Pulp {
		inv.Archetype = PickRandomArchetype()
	}
	inv.PickRandomTalents()
	// assign occupation
	// ToDo
	// Initialize Attributes
	inv.InitializeAttributes()
	LCK := inv.Attributes[AttrLuck]
	SAN := inv.Attributes[AttrSanity]
	POW := inv.Attributes[AttrPower]
	MP := inv.Attributes[AttrMagicPoints]
	DEX := inv.Attributes[AttrDexterity]
	EDU := inv.Attributes[AttrEducation]
	INT := inv.Attributes[AttrIntelligence]
	LCK.Initialize(false)
	// allow re roll
	if LCK.Value < 45 {
		LCK.Initialize(false)
	}

	SAN.Value = POW.Value
	SAN.StartingValue = POW.StartingValue
	inv.SetHP()
	inv.SetMovement()
	inv.SetBuildAndDMG()
	MP.Value = POW.Value / 5

	inv.Skills["Dodge"] = Skill{
		Name:         "Dodge",
		Abbreviation: "Dodge",
		Default:      DEX.Value / 2,
		Value:        DEX.Value / 2,
	}
	inv.Skills["Idea"] = Skill{
		Name:         "Idea",
		Abbreviation: "Idea",
		Default:      INT.Value / 2,
		Value:        INT.Value / 2,
	}
	inv.Skills["Know"] = Skill{
		Name:         "Know",
		Abbreviation: "Know",
		Default:      EDU.Value / 2,
		Value:        EDU.Value / 2,
	}
	inv.Skills["Language(Own)"] = Skill{
		Name:         "Language(Own)",
		Abbreviation: "Language(Own)",
		Default:      EDU.Value,
		Value:        EDU.Value,
	}
	// assign points
	return &inv
}

// ToDO: Need to support Occupation Assignment + Skill Points Assignament based on Occ & Archetype & Free
func main() {
	investigator := NewInvestigator(Pulp)
	fmt.Println(investigator)
}

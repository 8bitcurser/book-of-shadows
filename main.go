package main

import (
	"encoding/json"
	"fmt"
	"math"
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
	FilePath string `json:"path"`
	FileName string `json:"-"`
}

func (pp *ProfilePic) String() string {
	return fmt.Sprintf("%v", pp.FilePath)
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
	i.Attributes[AttrHitPoints] = HP

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
		attribute := i.Attributes[attr.Name]
		attribute.Initialize(isCore)
		i.Attributes[attr.Name] = attribute
	}

}

func (i *Investigator) CalculateOccupationSkillPoints() int {
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

func (i *Investigator) AssignSkillPoints(assignablePoints int, skills []string) {
	skillLimit := 90
	if i.GameMode == Pulp {
		skillLimit = 95
	}

	for assignablePoints > 0 {
		skillPicked := rand.Intn(len(skills))
		skill := i.Skills[skills[skillPicked]]
		pointsToAssign := 0
		CR := i.Skills["Credit Rating"]
		if CR.Value < i.Occupation.CreditRating.Min {
			creditPointsBase := i.Occupation.CreditRating.Min - CR.Value
			assignablePoints -= creditPointsBase
			CR.Value = creditPointsBase
			i.Skills["Credit Rating"] = CR
		}
		maxPointForSkill := skillLimit
		if skill.Value <= skillLimit && skill.Name != "Cthulhu Mythos" {
			if skill.Name == "Credit Rating" {
				maxPointForSkill = int(math.Min(float64(skillLimit), float64(i.Occupation.CreditRating.Max)))
				pointsToAssign = rand.Intn(maxPointForSkill) + 1
			} else {
				maxPointForSkill = skillLimit - skill.Value
				pointsToAssign = rand.Intn(maxPointForSkill) + 1
			}
			skill.Value += pointsToAssign
			assignablePoints -= pointsToAssign
			i.Skills[skill.Name] = skill

		}
	}

}

func (i *Investigator) GetSkills() {
	filteredSkills := map[string]Skill{}
	for name, skill := range Skills {
		for _, era := range Skills[name].Era {
			if era == i.Era {
				filteredSkills[name] = skill
			}
		}
	}
	i.Skills = filteredSkills
}

func (i *Investigator) AssignOccupation() {
	occupation := rand.Intn(len(OccupationsList))
	if i.GameMode == Pulp && i.Archetype != nil {
		if len(i.Archetype.SuggestedOccupations) > 0 {
			occupation = rand.Intn(len(i.Archetype.SuggestedOccupations))
		}
	}

	pickedOccupation := Occupations[OccupationsList[occupation]]
	i.Occupation = &pickedOccupation
}

func (i *Investigator) ToJSON() (string, error) {
	bytes, err := json.Marshal(i)
	if err != nil {
		return "", fmt.Errorf("error marshaling investigator: %v", err)
	}
	return string(bytes), nil
}

type Investigator struct {
	Era              Era                  `json:"era"`
	GameMode         GameMode             `json:"game-mode"`
	Name             string               `json:"name"`
	Residence        string               `json:"residence"`
	Birthplace       string               `json:"birthplace"`
	Age              int                  `json:"age"`
	ProfilePic       ProfilePic           `json:"profile-pic"`
	Occupation       *Occupation          `json:"occupation"`
	Archetype        *Archetype           `json:"archetype"`
	Insane           bool                 `json:"insane"`
	TemporaryInsane  bool                 `json:"temporary_insane"`
	IndefiniteInsane bool                 `json:"indefinite_insane"`
	MajorWound       bool                 `json:"major_wound"`
	Unconscious      bool                 `json:"unconscious"`
	Dying            bool                 `json:"dying"`
	Attributes       map[string]Attribute `json:"attributes"`
	Skills           map[string]Skill     `json:"skills"`
	Move             int                  `json:"move"`
	Build            string               `json:"build"`
	DamageBonus      string               `json:"damage-bonus"`
	Talents          []Talent             `json:"talents"`
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
		Skills:      map[string]Skill{},
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
	inv.GetSkills()
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
	inv.AssignOccupation()
	// assign points
	occupationPoints := inv.CalculateOccupationSkillPoints()
	if inv.GameMode == Pulp {
		archetypePoints := inv.Archetype.BonusPoints
		inv.AssignSkillPoints(archetypePoints, inv.Archetype.Skills)
	}
	inv.AssignSkillPoints(occupationPoints, inv.Occupation.Skills)
	var skillsList []string
	for _, v := range inv.Skills {
		skillsList = append(skillsList, v.Name)
	}
	inv.AssignSkillPoints(INT.Value*2, skillsList)
	return &inv
}

// ToDO: Need to support Occupation Assignment
func main() {
	investigator := NewInvestigator(Pulp)
	fmt.Println(investigator.ToJSON())
}

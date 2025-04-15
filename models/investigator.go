package models

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strings"
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
	// if both are not greater or lesser than size it means one of the two is.
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

	for key, _ := range i.Attributes {

		// An attribute is core if we're in pulp mode AND it's in core characteristics
		isCore := isPulp && coreCharacteristics[key]

		// Initialize the attribute
		attribute := i.Attributes[key]
		attribute.Initialize(isCore)
		i.Attributes[key] = attribute
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

func (i *Investigator) AssignSkillPoints(assignablePoints int, skills []string) int {
	skillLimit := 90
	if i.GameMode == Pulp {
		skillLimit = 95
	}
	CR := i.Skills["Credit Rating"]
	if CR.Value < i.Occupation.CreditRating.Min {
		creditPointsBase := i.Occupation.CreditRating.Min - CR.Value
		assignablePoints -= creditPointsBase
		CR.Value = creditPointsBase
		i.Skills["Credit Rating"] = CR
	}
	for assignablePoints > 0 {
		skillPicked := rand.Intn(len(skills))
		skillName := skills[skillPicked]
		skill, ok := i.Skills[skillName]
		if !ok || skill.Base == 1 {
			continue
		}
		pointsToAssign := rand.Intn(50) + 5

		if assignablePoints < pointsToAssign || assignablePoints-pointsToAssign < 0 {
			pointsToAssign = assignablePoints
		}
		assignablePoints -= pointsToAssign
		if skill.Name == "Credit Rating" {
			skillLimit = int(math.Min(float64(skillLimit), float64(i.Occupation.CreditRating.Max)))
		}
		if skill.Value > skillLimit || skill.Name == "Cthulhu Mythos" || (skill.Value+pointsToAssign) > skillLimit {
			continue
		}

		skill.Value += pointsToAssign

		i.Skills[skillName] = skill
		fmt.Printf("Assigned skill points: %v - Amount: %d\n", skillName, pointsToAssign)
	}
	return assignablePoints
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

func (i *Investigator) addMissingSkills(skills *[]string) {
	// Add base skills as Any Skills on empty call
	if len(*skills) == 0 {
		for _, skill := range i.Skills {
			if skill.Base == 1 {
				newName := fmt.Sprintf("%s(%s)", skill.Name, "Any")
				i.Skills[newName] = Skill{
					Name:         newName,
					Abbreviation: newName,
					FormName:     skill.FormName,
					Default:      skill.Default,
					Value:        skill.Value,
					Era:          []Era{Twenties, Modern},
					Base:         0,
					NeedsFormDef: 1,
				}
			}
		}
	}

	for _, occ := range *skills {
		_, ok := i.Skills[occ]
		if ok {
			continue
		}

		if strings.Contains(occ, "(") {
			// Capture prefix example ArtCraft(Painting) ArtCraft
			prefix, _, _ := strings.Cut(occ, "(")
			matches := make([]Skill, 0)
			// find all the skills that have that category
			for _, v := range i.Skills {
				if v.Category == prefix && v.Base == 0 {
					matches = append(matches, v)
				}
			}
			// if several matched pick one
			if len(matches) > 0 {
				matchPick := rand.Intn(len(matches))
				skillMatched := matches[matchPick]
				skill, _ := i.Skills[skillMatched.Name]
				occ = skill.Name

			} else {
				baseSkill := Skill{}
				for s, v := range i.Skills {
					if s == prefix && v.Base == 1 {
						baseSkill = i.Skills[v.Name]
						break
					}
				}
				if baseSkill.Name != "" {
					// handle form names 1, 2, 3 ... etc
					i.Skills[occ] = Skill{
						Name:         occ,
						Abbreviation: occ,
						FormName:     baseSkill.FormName,
						Default:      baseSkill.Default,
						Value:        baseSkill.Value,
						Era:          baseSkill.Era,
						Base:         0,
						Category:     prefix,
						NeedsFormDef: 1,
					}
				} else {
					i.Skills[occ] = Skill{
						Name:         occ,
						Abbreviation: occ,
						FormName:     "Custom1",
						Default:      1,
						Value:        1,
						Era:          []Era{Twenties, Modern},
						Base:         0,
						NeedsFormDef: 1,
					}
				}
			}

		} else {
			matches := make([]Skill, 0)
			for _, v := range i.Skills {
				if strings.HasPrefix(v.Name, occ) {
					matches = append(matches, v)
				}
			}

			// Add the custom ones
			if len(matches) == 0 {
				i.Skills[occ] = Skill{
					Name:         occ,
					Abbreviation: occ,
					FormName:     "Custom1",
					Default:      1,
					Value:        1,
					Era:          []Era{Twenties, Modern},
					Base:         0,
					NeedsFormDef: 1,
				}
			}
		}
	}

}

func (i *Investigator) ToJSON() ([]byte, error) {
	bytes, err := json.Marshal(i)

	if err != nil {
		return []byte(""), fmt.Errorf("error marshaling investigator: %v", err)
	}
	return bytes, nil
}

type Investigator struct {
	ID                         string               `json:"id"`
	Era                        Era                  `json:"-"`
	GameMode                   GameMode             `json:"-"`
	Name                       string               `json:"Investigators_Name"`
	Residence                  string               `json:"Residence"`
	Birthplace                 string               `json:"Birthplace"`
	Age                        int                  `json:"Age"`
	ProfilePic                 ProfilePic           `json:"Portrait"`
	Occupation                 *Occupation          `json:"Occupation"`
	Archetype                  *Archetype           `json:"Archetype"`
	Insane                     bool                 `json:"insane"`
	TemporaryInsane            bool                 `json:"TempInsanity_Chk Off"`
	IndefiniteInsane           bool                 `json:"IndefInsanity_Chk"`
	MajorWound                 bool                 `json:"MajorWound_Chk"`
	Unconscious                bool                 `json:"Unconscious_Chk"`
	Dying                      bool                 `json:"Dying_Chk"`
	Attributes                 map[string]Attribute `json:"Attributes"`
	Skills                     map[string]Skill     `json:"Skills"`
	Move                       int                  `json:"MOV"`
	Build                      string               `json:"Build"`
	DamageBonus                string               `json:"DamageBonus"`
	Talents                    []Talent             `json:"Pulp-Talents"`
	OccupationPoints           int                  `json:"OccupationPoints"`
	ArchetypePoints            int                  `json:"ArchetypePoints"`
	FreePoints                 int                  `json:"FreePoints"`
	UnassignedOccupationPoints int                  `json:"UnassignedOccupationPoints"`
	UnassignedArchetypePoints  int                  `json:"UnassignedArchetypePoints"`
	UnassignedFreePoints       int                  `json:"UnassignedFreePoints"`
}

func RandomInvestigator(mode GameMode) *Investigator {
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
				Name:          "STR",
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrConstitution: {
				Name:          "CON",
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrDexterity: {
				Name:          "DEX",
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrIntelligence: {
				Name:          "INT",
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrSize: {
				Name:          "SIZ",
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrPower: {
				Name:          "POW",
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrAppearance: {
				Name:          "APP",
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrEducation: {
				Name:          "EDU",
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrHitPoints: {
				Name:          "CurrentHP",
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrMagicPoints: {
				Name:          "CurrentMagic",
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrLuck: {
				Name:          "CurrentLuck",
				StartingValue: 0,
				Value:         0,
				MaxValue:      0,
			},
			AttrSanity: {
				Name:          "CurrentSanity",
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
		inv.PickRandomTalents()
	}
	// assign occupation
	inv.AssignOccupation()
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
	inv.Attributes[AttrSanity] = SAN
	inv.SetHP()
	inv.SetMovement()
	inv.SetBuildAndDMG()
	MP.Value = POW.Value / 5
	inv.Attributes[AttrMagicPoints] = MP
	inv.GetSkills()

	inv.Skills["Dodge_Copy"] = Skill{
		Name:         "Dodge_Copy",
		Abbreviation: "Dodge",
		FormName:     "Dodge_Copy",
		Default:      DEX.Value / 2,
		Value:        (DEX.Value / 2),
	}
	inv.Skills["Dodge"] = Skill{
		Name:         "Dodge",
		Abbreviation: "Dodge",
		FormName:     "Dodge",
		Default:      DEX.Value / 2,
		Value:        DEX.Value / 2,
	}
	inv.Skills["Language(Own)"] = Skill{
		Name:         "Language(Own)",
		Abbreviation: "Language(Own)",
		FormName:     "OwnLanguage",
		Default:      EDU.Value,
		Value:        EDU.Value,
	}
	inv.addMissingSkills(&[]string{})
	// assign points
	occupationPoints := inv.CalculateOccupationSkillPoints()
	inv.UnassignedOccupationPoints = occupationPoints
	inv.OccupationPoints = occupationPoints

	inv.ArchetypePoints = inv.Archetype.BonusPoints
	inv.addMissingSkills(&inv.Archetype.Skills)
	sparePoints := inv.AssignSkillPoints(inv.ArchetypePoints, inv.Archetype.Skills)
	inv.UnassignedArchetypePoints = sparePoints

	occupationSkills := inv.GetOccupationSkills()
	inv.addMissingSkills(occupationSkills)

	sparePoints = inv.AssignSkillPoints(occupationPoints, *occupationSkills)
	inv.UnassignedOccupationPoints = sparePoints
	var skillsList []string
	for s, v := range inv.Skills {
		if v.Name != "Cthulhu Mythos" && v.Name != "Dodge_Copy" {
			skillsList = append(skillsList, s)
		}
	}
	inv.FreePoints = INT.Value * 2
	sparePoints = inv.AssignSkillPoints(inv.FreePoints, skillsList)
	inv.UnassignedFreePoints = sparePoints
	return &inv
}

func InvestigatorBaseCreate(data map[string]any) *Investigator {
	archetype := Archetypes[data["archetype"].(string)]
	occupation := Occupations[data["occupation"].(string)]
	inv := Investigator{
		Era:              1,
		GameMode:         Pulp,
		Name:             data["name"].(string),
		Residence:        data["residence"].(string),
		Birthplace:       data["birthplace"].(string),
		Age:              data["age"].(int),
		ProfilePic:       ProfilePic{"/sample/path/env", "profile"},
		Insane:           false,
		TemporaryInsane:  false,
		IndefiniteInsane: false,
		MajorWound:       false,
		Unconscious:      false,
		Dying:            false,
		Attributes:       map[string]Attribute{},
		Skills:           map[string]Skill{},
		Move:             2,
		Build:            "Big",
		DamageBonus:      "1D4",
		Archetype:        &archetype,
		Occupation:       &occupation,
	}
	inv.GetSkills()
	inv.addMissingSkills(&[]string{})
	inv.addMissingSkills(&inv.Archetype.Skills)
	occupationSkills := inv.GetOccupationSkills()
	inv.addMissingSkills(occupationSkills)

	return &inv
}

func (i *Investigator) InvestigatorUpdateAttributes(data map[string]int) {
	i.Attributes = map[string]Attribute{
		AttrStrength: {
			Name:          "STR",
			StartingValue: 0,
			Value:         data["STR"],
			MaxValue:      0,
		},
		AttrConstitution: {
			Name:          "CON",
			StartingValue: 0,
			Value:         data["CON"],
			MaxValue:      0,
		},
		AttrDexterity: {
			Name:          "DEX",
			StartingValue: 0,
			Value:         data["DEX"],
			MaxValue:      0,
		},
		AttrIntelligence: {
			Name:          "INT",
			StartingValue: 0,
			Value:         data["INT"],
			MaxValue:      0,
		},
		AttrSize: {
			Name:          "SIZ",
			StartingValue: 0,
			Value:         data["SIZ"],
			MaxValue:      0,
		},
		AttrPower: {
			Name:          "POW",
			StartingValue: 0,
			Value:         data["POW"],
			MaxValue:      0,
		},
		AttrAppearance: {
			Name:          "APP",
			StartingValue: 0,
			Value:         data["APP"],
			MaxValue:      0,
		},
		AttrEducation: {
			Name:          "EDU",
			StartingValue: 0,
			Value:         data["EDU"],
			MaxValue:      0,
		},
		AttrHitPoints: {
			Name:          "CurrentHP",
			StartingValue: 0,
			Value:         0,
			MaxValue:      0,
		},
		AttrMagicPoints: {
			Name:          "CurrentMagic",
			StartingValue: 0,
			Value:         0,
			MaxValue:      0,
		},
		AttrLuck: {
			Name:          "CurrentLuck",
			StartingValue: data["LCK"],
			Value:         data["LCK"],
			MaxValue:      0,
		},
		AttrSanity: {
			Name:          "CurrentSanity",
			StartingValue: 0,
			Value:         0,
			MaxValue:      0,
		},
	}
	SAN := i.Attributes[AttrSanity]
	POW := i.Attributes[AttrPower]
	MP := i.Attributes[AttrMagicPoints]
	DEX := i.Attributes[AttrDexterity]
	EDU := i.Attributes[AttrEducation]
	INT := i.Attributes[AttrIntelligence]
	SAN.Value = POW.Value
	SAN.StartingValue = POW.StartingValue
	i.Attributes[AttrSanity] = SAN
	i.SetHP()
	i.SetMovement()
	i.SetBuildAndDMG()
	MP.Value = POW.Value / 5
	i.Attributes[AttrMagicPoints] = MP
	i.Skills["Dodge_Copy"] = Skill{
		Name:         "Dodge_Copy",
		Abbreviation: "Dodge",
		FormName:     "Dodge_Copy",
		Default:      DEX.Value / 2,
		Value:        DEX.Value / 2,
	}
	i.Skills["Dodge"] = Skill{
		Name:         "Dodge",
		Abbreviation: "Dodge",
		FormName:     "Dodge",
		Default:      DEX.Value / 2,
		Value:        DEX.Value / 2,
	}
	i.Skills["Language(Own)"] = Skill{
		Name:         "Language(Own)",
		Abbreviation: "Language(Own)",
		FormName:     "OwnLanguage",
		Default:      EDU.Value,
		Value:        EDU.Value,
	}
	occupationPoints := i.CalculateOccupationSkillPoints()
	i.OccupationPoints = occupationPoints
	i.ArchetypePoints = i.Archetype.BonusPoints
	i.FreePoints = INT.Value * 2
}

func (i *Investigator) GetOccupationSkills() *[]string {
	occupationSkills := make([]string, 0)

	for _, skillReq := range i.Occupation.SkillRequirements {
		if skillReq.Type == "required" {
			occupationSkills = append(occupationSkills, skillReq.Skill)
		} else {
			picked := make([]int, 0)
			for i := 0; i < skillReq.SkillChoice.NumRequired; i++ {
				choice := rand.Intn(len(skillReq.SkillChoice.Skills))
				if slices.Contains(picked, choice) {
					continue
				} else {
					picked = append(picked, choice)
					occupationSkills = append(occupationSkills, skillReq.SkillChoice.Skills[choice])
				}
			}
		}
	}
	return &occupationSkills
}

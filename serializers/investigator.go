package serializers

import (
	"book-of-shadows/models"
	"encoding/json"
	"strconv"
	"strings"
)

type InvestigatorSerializer struct {
	// Basic Info
	Name                    string `json:"Investigators_Name"`
	Age                     string `json:"Age"`
	Residence               string `json:"Residence"`
	Birthplace              string `json:"Birthplace"`
	Occupation              string `json:"Occupation"`
	Archetype               string `json:"Archetype"`
	MOV                     string `json:"MOV"`
	Build                   string `json:"Build"`
	DamageBonus             string `json:"DamageBonus"`
	PulpTalents             string `json:"Pulp Talents"`
	PulpTalentsDescriptions string `json:"Pulp Talents Descriptions"`

	// Status
	Insane           string `json:"insane"`
	TemporaryInsane  string `json:"TempInsanity_Chk Off"`
	IndefiniteInsane string `json:"IndefInsanity_Chk"`
	MajorWound       string `json:"MajorWound_Chk"`
	Unconscious      string `json:"Unconscious_Chk"`
	Dying            string `json:"Dying_Chk"`

	// Attributes with their derived values
	STR             string `json:"STR"`
	STR_half        string `json:"STR_half"`
	STR_fifth       string `json:"STR_fifth"`
	DEX             string `json:"DEX"`
	DEX_half        string `json:"DEX_half"`
	DEX_fifth       string `json:"DEX_fifth"`
	POW             string `json:"POW"`
	POW_half        string `json:"POW_half"`
	POW_fifth       string `json:"POW_fifth"`
	CON             string `json:"CON"`
	CON_half        string `json:"CON_half"`
	CON_fifth       string `json:"CON_fifth"`
	APP             string `json:"APP"`
	APP_half        string `json:"APP_half"`
	APP_fifth       string `json:"APP_fifth"`
	EDU             string `json:"EDU"`
	EDU_half        string `json:"EDU_half"`
	EDU_fifth       string `json:"EDU_fifth"`
	SIZ             string `json:"SIZ"`
	SIZ_half        string `json:"SIZ_half"`
	SIZ_fifth       string `json:"SIZ_fifth"`
	INT             string `json:"INT"`
	INT_half        string `json:"INT_half"`
	INT_fifth       string `json:"INT_fifth"`
	CurrentHP       string `json:"CurrentHP"`
	CurrentHP_half  string `json:"CurrentHP_half"`
	CurrentHP_fifth string `json:"CurrentHP_fifth"`
	CurrentMagic    string `json:"CurrentMagic"`
	CurrentSanity   string `json:"CurrentSanity"`
	CurrentLuck     string `json:"CurrentLuck"`

	// All skills will be handled dynamically in the conversion methods
	// using map[string]string to store all Skill_* fields
	Skills map[string]string `json:"-"`
}

// Helper method to convert string to boolean
func strToBool(s string) bool {
	return strings.ToLower(s) == "true"
}

// Helper method to convert string to int with fallback
func strToInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}

// UnmarshalJSON implements custom JSON unmarshaling
func (s *InvestigatorSerializer) UnmarshalJSON(data []byte) error {
	// First, unmarshal into a map to handle dynamic skill fields
	var rawData map[string]string

	if err := json.Unmarshal(data, &rawData); err != nil {
		return err
	}

	// Initialize the skills map
	s.Skills = make(map[string]string)

	// Process each field
	for key, value := range rawData {
		switch {
		case strings.HasPrefix(key, "Skill_"):
			// Store all skill-related fields
			s.Skills[key] = value
		default:
			// Handle other fields using reflection or manual assignment
			switch key {
			case "Investigators_Name":
				s.Name = value
			case "Age":
				s.Age = value
			case "Residence":
				s.Residence = value
			case "Birthplace":
				s.Birthplace = value
			case "Occupation":
				s.Occupation = value
			case "Archetype":
				s.Archetype = value
			case "MOV":
				s.MOV = value
			case "Build":
				s.Build = value
			case "DamageBonus":
				s.DamageBonus = value
			case "Pulp Talents":
				s.PulpTalents = value
			case "Pulp Talents Descriptions":
				s.PulpTalentsDescriptions = value
			case "STR":
				s.STR = value
			case "DEX":
				s.DEX = value
			case "APP":
				s.APP = value
			case "INT":
				s.INT = value
			case "EDU":
				s.EDU = value
			case "CON":
				s.CON = value
			case "POW":
				s.POW = value
			case "SIZ":
				s.SIZ = value
			case "CurrentLuck":
				s.CurrentLuck = value
			case "CurrentHP":
				s.CurrentHP = value
			case "CurrentMagic":
				s.CurrentMagic = value
			case "CurrentSanity":
				s.CurrentSanity = value
				// Add other fields as needed
			}
		}
	}

	return nil
}

// ToInvestigator converts the serializer to an Investigator domain model
func (s *InvestigatorSerializer) ToInvestigator() *models.Investigator {
	inv := &models.Investigator{
		Name:             s.Name,
		Residence:        s.Residence,
		Birthplace:       s.Birthplace,
		Age:              strToInt(s.Age),
		Move:             strToInt(s.MOV),
		Build:            s.Build,
		DamageBonus:      s.DamageBonus,
		Occupation:       &models.Occupation{Name: s.Occupation},
		Archetype:        &models.Archetype{Name: s.Archetype},
		Insane:           strToBool(s.Insane),
		TemporaryInsane:  strToBool(s.TemporaryInsane),
		IndefiniteInsane: strToBool(s.IndefiniteInsane),
		MajorWound:       strToBool(s.MajorWound),
		Unconscious:      strToBool(s.Unconscious),
		Dying:            strToBool(s.Dying),

		// Initialize maps
		Attributes: make(map[string]models.Attribute),
		Skills:     make(map[string]models.Skill),
	}

	// Convert attributes
	attributeMap := map[string]string{
		"STR":                s.STR,
		"DEX":                s.DEX,
		"POW":                s.POW,
		"CON":                s.CON,
		"APP":                s.APP,
		"EDU":                s.EDU,
		"SIZ":                s.SIZ,
		"INT":                s.INT,
		"CurrentHitPoints":   s.CurrentHP,
		"CurrentMagicPoints": s.CurrentMagic,
		"CurrentSanity":      s.CurrentSanity,
		"CurrentLuck":        s.CurrentLuck,
	}

	for key, value := range attributeMap {
		if strings.HasPrefix(key, "Current") {
			inv.Attributes[strings.Replace(key, "Current", "", 1)] = models.Attribute{Name: key, Value: strToInt(value)}
		} else {
			inv.Attributes[key] = models.Attribute{Name: key, Value: strToInt(value)}
		}
	}

	// Convert skills
	for key, value := range s.Skills {
		if strings.HasPrefix(key, "Skill_") && !strings.HasSuffix(key, "_half") && !strings.HasSuffix(key, "_fifth") {
			skillName := strings.TrimPrefix(key, "Skill_")
			inv.Skills[skillName] = models.Skill{
				Name:  skillName,
				Value: strToInt(value),
			}
		}
	}
	// Convert Pulp Talents
	if s.PulpTalents != "" {
		talents := strings.Split(strings.TrimSuffix(s.PulpTalents, ", "), ", ")
		description := strings.Split(strings.TrimSuffix(s.PulpTalentsDescriptions, "~ "), "~ ")
		for i, t := range talents {
			if t != "" {
				inv.Talents = append(inv.Talents, models.Talent{Name: t, Description: description[i]})
			}
		}
	}

	return inv
}

// FromJSON creates an Investigator from JSON data
func FromJSON(data []byte) (*models.Investigator, error) {
	var serializer InvestigatorSerializer

	if err := json.Unmarshal(data, &serializer); err != nil {
		return nil, err
	}
	return serializer.ToInvestigator(), nil
}

var validUpdateSections = []string{
	"combat",
	"personalInfo",
	"skills",
	"skill_check",
	"skill_name",
}

type UpdateRequestSerializer struct {
	Section string `json:"section"`
	Field   string `json:"field"`
	Value   any    `json:"value"`
}

func (u *UpdateRequestSerializer) IsValid() bool {
	valid := false
	for _, section := range validUpdateSections {
		if u.Section == section {
			valid = true
			break
		}
	}
	return valid
}

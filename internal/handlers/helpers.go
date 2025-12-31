package handlers

import (
	"fmt"
	"strconv"

	"book-of-shadows/internal/errors"
	"book-of-shadows/models"
)

// UpdateRequest represents an update request for an investigator
type UpdateRequest struct {
	Section string      `json:"section"`
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
}

// processInvestigatorPayload processes and validates the payload for creating an investigator
func (h *Handler) processInvestigatorPayload(payload map[string]interface{}) map[string]interface{} {
	processed := make(map[string]interface{})

	// Define fields that need type conversion
	intFields := []string{"age", "Age"}

	for key, val := range payload {
		// Check if this field needs int conversion
		needsIntConversion := false
		for _, field := range intFields {
			if key == field {
				needsIntConversion = true
				break
			}
		}

		if needsIntConversion {
			// Convert string to int if needed
			switch v := val.(type) {
			case string:
				if intVal, err := strconv.Atoi(v); err == nil {
					processed[key] = intVal
				} else {
					processed[key] = 0
				}
			case float64:
				processed[key] = int(v)
			case int:
				processed[key] = v
			default:
				processed[key] = 0
			}
		} else {
			processed[key] = val
		}
	}

	return processed
}

// applyInvestigatorUpdate applies an update to an investigator based on the section and field
func (h *Handler) applyInvestigatorUpdate(inv *models.Investigator, req *UpdateRequest) error {
	switch req.Section {
	case "attributes":
		return h.updateAttribute(inv, req.Field, req.Value)
	case "skills":
		return h.updateSkill(inv, req.Field, req.Value)
	case "stats":
		return h.updateStats(inv, req.Field, req.Value)
	case "personalInfo":
		return h.updatePersonalInfo(inv, req.Field, req.Value)
	case "combat":
		return h.updateCombat(inv, req.Field, req.Value)
	case "skill_check":
		return h.toggleSkillCheck(inv, req.Field)
	case "skill_prio":
		return h.toggleSkillPriority(inv, req.Field)
	case "skill_name":
		return h.updateSkillName(inv, req.Field, req.Value)
	case "talents":
		return h.updateTalent(inv, req.Field, req.Value)
	case "phobias":
		return h.updatePhobia(inv, req.Field, req.Value)
	case "manias":
		return h.updateMania(inv, req.Field, req.Value)
	default:
		return errors.NewHTTPError(400, "Unknown section", nil)
	}
}

// Mapping from internal attribute names to PDF field names
var attrNameToPdfField = map[string]string{
	models.AttrStrength:     "STR",
	models.AttrConstitution: "CON",
	models.AttrSize:         "SIZ",
	models.AttrDexterity:    "DEX",
	models.AttrAppearance:   "APP",
	models.AttrIntelligence: "INT",
	models.AttrPower:        "POW",
	models.AttrEducation:    "EDU",
}

// updateAttribute updates an investigator's attribute
func (h *Handler) updateAttribute(inv *models.Investigator, field string, value interface{}) error {
	intValue, err := toInt(value)
	if err != nil {
		return errors.NewValidationError(field, "must be a number")
	}

	attr, exists := inv.Attributes[field]
	if !exists {
		// Create new attribute if it doesn't exist
		// Use the PDF field name (STR, CON, etc.) for the Name field
		pdfName := field
		if mapped, ok := attrNameToPdfField[field]; ok {
			pdfName = mapped
		}
		inv.Attributes[field] = models.Attribute{
			Name:  pdfName,
			Value: intValue,
		}
	} else {
		attr.Value = intValue
		inv.Attributes[field] = attr
	}

	// Recalculate dependent values
	h.recalculateDependentAttributes(inv)

	return nil
}

// updateSkill updates an investigator's skill value
func (h *Handler) updateSkill(inv *models.Investigator, field string, value interface{}) error {
	intValue, err := toInt(value)
	if err != nil {
		return errors.NewValidationError(field, "must be a number")
	}

	skill, exists := inv.Skills[field]
	if !exists {
		return errors.ErrInvalidSkill
	}

	skill.Value = intValue
	inv.Skills[field] = skill

	return nil
}

// updateStats toggles investigator status conditions
func (h *Handler) updateStats(inv *models.Investigator, field string, _ interface{}) error {
	switch field {
	case "TemporaryInsane":
		inv.TemporaryInsane = !inv.TemporaryInsane
	case "IndefiniteInsane":
		inv.IndefiniteInsane = !inv.IndefiniteInsane
	case "MajorWound":
		inv.MajorWound = !inv.MajorWound
	case "Unconscious":
		inv.Unconscious = !inv.Unconscious
	case "Dying":
		inv.Dying = !inv.Dying
	default:
		return errors.NewValidationError(field, "unknown status field")
	}
	return nil
}

// updatePersonalInfo updates investigator personal information
func (h *Handler) updatePersonalInfo(inv *models.Investigator, field string, value interface{}) error {
	switch field {
	case "Name":
		strVal, ok := value.(string)
		if !ok {
			return errors.NewValidationError(field, "must be a string")
		}
		inv.Name = strVal

	case "Age":
		intValue, err := toInt(value)
		if err != nil {
			// Try parsing as string
			if strVal, ok := value.(string); ok {
				if age, err := strconv.Atoi(strVal); err == nil {
					inv.Age = age
				} else {
					return errors.NewValidationError(field, "must be a number")
				}
			} else {
				return errors.NewValidationError(field, "must be a number")
			}
		} else {
			inv.Age = intValue
		}

	case "Residence":
		strVal, ok := value.(string)
		if !ok {
			return errors.NewValidationError(field, "must be a string")
		}
		inv.Residence = strVal

	case "Birthplace":
		strVal, ok := value.(string)
		if !ok {
			return errors.NewValidationError(field, "must be a string")
		}
		inv.Birthplace = strVal

	default:
		return errors.NewValidationError(field, "unsupported personal info field")
	}

	return nil
}

// updateCombat updates combat-related attributes
func (h *Handler) updateCombat(inv *models.Investigator, field string, value interface{}) error {
	intValue, err := toInt(value)
	if err != nil {
		return errors.NewValidationError(field, "must be a number")
	}

	attr, exists := inv.Attributes[field]
	if !exists {
		return errors.ErrInvalidAttribute
	}

	attr.Value = intValue
	inv.Attributes[field] = attr

	return nil
}

// toggleSkillCheck toggles the selection status of a skill
func (h *Handler) toggleSkillCheck(inv *models.Investigator, field string) error {
	skill, exists := inv.Skills[field]
	if !exists {
		return errors.ErrInvalidSkill
	}

	skill.IsSelected = !skill.IsSelected
	inv.Skills[field] = skill

	return nil
}

// toggleSkillPriority toggles the priority status of a skill
func (h *Handler) toggleSkillPriority(inv *models.Investigator, field string) error {
	skill, exists := inv.Skills[field]
	if !exists {
		return errors.ErrInvalidSkill
	}

	skill.IsPriority = !skill.IsPriority
	inv.Skills[field] = skill

	return nil
}

// updateSkillName renames a skill
func (h *Handler) updateSkillName(inv *models.Investigator, field string, value interface{}) error {
	newName, ok := value.(string)
	if !ok {
		return errors.NewValidationError(field, "must be a string")
	}

	skill, exists := inv.Skills[field]
	if !exists {
		return errors.ErrInvalidSkill
	}

	// Don't update if name hasn't changed
	if skill.Name == newName {
		return nil
	}

	// Update skill name and re-index
	skill.Name = newName
	inv.Skills[newName] = skill
	delete(inv.Skills, field)

	return nil
}

// updateTalent adds or removes a talent from an investigator
func (h *Handler) updateTalent(inv *models.Investigator, talentName string, value interface{}) error {
	// Check if talent exists in the global list
	talent, exists := models.Talents[talentName]
	if !exists {
		return errors.NewValidationError(talentName, "unknown talent")
	}

	// Determine if we're adding or removing
	shouldAdd, ok := value.(bool)
	if !ok {
		return errors.NewValidationError(talentName, "value must be boolean")
	}

	if shouldAdd {
		// Check if already at max talents
		if len(inv.Talents) >= inv.Archetype.AmountOfTalents {
			return errors.NewValidationError(talentName, "maximum talents already selected")
		}

		// Check if talent already exists
		for _, t := range inv.Talents {
			if t.Name == talentName {
				return nil // Already has this talent
			}
		}

		// Add the talent
		inv.Talents = append(inv.Talents, talent)
	} else {
		// Remove the talent
		newTalents := make([]models.Talent, 0, len(inv.Talents))
		for _, t := range inv.Talents {
			if t.Name != talentName {
				newTalents = append(newTalents, t)
			}
		}
		inv.Talents = newTalents
	}

	return nil
}

// updatePhobia adds or removes a phobia from an investigator
func (h *Handler) updatePhobia(inv *models.Investigator, phobiaName string, value interface{}) error {
	// Check if phobia exists in the global list
	phobia, exists := models.Phobias[phobiaName]
	if !exists {
		return errors.NewValidationError(phobiaName, "unknown phobia")
	}

	// Determine if we're adding or removing
	shouldAdd, ok := value.(bool)
	if !ok {
		return errors.NewValidationError(phobiaName, "value must be boolean")
	}

	if shouldAdd {
		// Check if phobia already exists
		for _, p := range inv.Phobias {
			if p.Name == phobiaName {
				return nil // Already has this phobia
			}
		}

		// Add the phobia
		inv.Phobias = append(inv.Phobias, phobia)
	} else {
		// Remove the phobia
		newPhobias := make([]models.Phobia, 0, len(inv.Phobias))
		for _, p := range inv.Phobias {
			if p.Name != phobiaName {
				newPhobias = append(newPhobias, p)
			}
		}
		inv.Phobias = newPhobias
	}

	return nil
}

// updateMania adds or removes a mania from an investigator
func (h *Handler) updateMania(inv *models.Investigator, maniaName string, value interface{}) error {
	// Check if mania exists in the global list
	mania, exists := models.Manias[maniaName]
	if !exists {
		return errors.NewValidationError(maniaName, "unknown mania")
	}

	// Determine if we're adding or removing
	shouldAdd, ok := value.(bool)
	if !ok {
		return errors.NewValidationError(maniaName, "value must be boolean")
	}

	if shouldAdd {
		// Check if mania already exists
		for _, m := range inv.Manias {
			if m.Name == maniaName {
				return nil // Already has this mania
			}
		}

		// Add the mania
		inv.Manias = append(inv.Manias, mania)
	} else {
		// Remove the mania
		newManias := make([]models.Mania, 0, len(inv.Manias))
		for _, m := range inv.Manias {
			if m.Name != maniaName {
				newManias = append(newManias, m)
			}
		}
		inv.Manias = newManias
	}

	return nil
}

// recalculateDependentAttributes recalculates attributes that depend on other attributes
func (h *Handler) recalculateDependentAttributes(inv *models.Investigator) {
	// Recalculate occupation points
	occupationPoints := inv.CalculateOccupationSkillPoints()
	inv.UnassignedOccupationPoints = occupationPoints
	inv.OccupationPoints = occupationPoints

	// Recalculate free points
	if intel, exists := inv.Attributes["Intelligence"]; exists {
		intPoints := intel.Value * 2
		inv.FreePoints = intPoints
		inv.UnassignedFreePoints = intPoints
	}

	// Update Sanity based on Power
	if power, exists := inv.Attributes["Power"]; exists {
		inv.Attributes["Sanity"] = models.Attribute{
			Name:          "CurrentSanity",
			Value:         power.Value,
			StartingValue: power.StartingValue,
			MaxValue:      99,
		}

		// Update Magic Points
		inv.Attributes["MagicPoints"] = models.Attribute{
			Name:          "CurrentMagic",
			Value:         power.Value / 5,
			StartingValue: power.Value / 5,
			MaxValue:      power.Value / 5,
		}
	}

	// Recalculate HP, Movement, Build & Damage
	inv.SetHP()
	inv.SetMovement()
	inv.SetBuildAndDMG()
}

// toInt converts an interface{} to int
func toInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case float32:
		return int(v), nil
	case int64:
		return int(v), nil
	case int32:
		return int(v), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int", value)
	}
}
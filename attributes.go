package main

import "fmt"

const (
	AttrStrength     = "Strength"
	AttrConstitution = "Constitution"
	AttrDexterity    = "Dexterity"
	AttrIntelligence = "Intelligence"
	AttrSize         = "Size"
	AttrPower        = "Power"
	AttrAppearance   = "Appearance"
	AttrEducation    = "Education"
	AttrHitPoints    = "HitPoints"
	AttrMagicPoints  = "MagicPoints"
	AttrLuck         = "Luck"
	AttrSanity       = "Sanity"
)

type Attribute struct {
	Name          string
	StartingValue int
	Value         int
	MaxValue      int
}

func (a *Attribute) String() string {
	return fmt.Sprintf("%v: %v", a.Name, a.Value)
}

func (a *Attribute) Initialize(isCore bool) {
	rolled := 0
	if a.Name == "Size" || a.Name == "Intelligence" || a.Name == "Education" {
		if isCore {
			rolled = coreRoll()
		} else {
			rolled = (rollD6() + rollD6() + 6) * 5
		}

	} else {

		if isCore {
			rolled = coreRoll()
		} else {
			rolled = (rollD6() + rollD6() + rollD6()) * 5
		}
	}
	a.Value = rolled
	a.StartingValue = rolled
}

package main

var Talents = map[string]Talent{
	"Keen Vision": Talent{
		Name:        "Keen Vision",
		Description: "Gain a bonus die to Spot Hidden Rolls",
		Type:        Physical,
	},
	"Quick Healer": Talent{
		Name:        "Quick Healer",
		Description: "Natural healing is increased to +3 hit points per day",
		Type:        Physical,
	},
	"Night Vision": Talent{
		Name: "Keen Vision",
		Description: "In darkness, reduce the difficulty level of Spot Hidden rolls and ignore penalty die " +
			"for shooting in the dark",
		Type: Physical,
	},
	"Endurance": Talent{
		Name:        "Endurance",
		Description: "Gain a bonus die when making CON rolls (including to determine MOV rate for chases)",
		Type:        Physical,
	},
	"Power Lifter": Talent{
		Name:        "Power Lifter",
		Description: "Gain a bonus die when making STR rolls to lift objects or people",
		Type:        Physical,
	},
	"Iron Liver": Talent{
		Name: "Iron Liver",
		Description: "May spend 5 Luck to avoid the effects of drinking excessive amounts of alcohol" +
			" (negating penalty applied to skill rolls)",
		Type: Physical,
	},
	"Stout Constitution": Talent{
		Name:        "Stout Constitution",
		Description: "May spend 10 Luck to reduce poison or disease damage and effect by half.",
		Type:        Physical,
	},
	"Tough Guy": Talent{
		Name: "Tough Guy",
		Description: "Soaks up damage, may spend 10 Luck points to shrug off up to 5 hit points worth" +
			" of damage taken in one combat round.",
		Type: Physical,
	},
	"Keen Hearing": Talent{
		Name:        "Keen Hearing",
		Description: "Gain a bonus die to Listen rolls",
		Type:        Physical,
	},
	"Smooth Talker": Talent{
		Name:        "Smooth Talker",
		Description: "Gain a bonus die to Charm Rolls",
		Type:        Physical,
	},
}

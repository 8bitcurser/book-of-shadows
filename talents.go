package main

var Talents = map[string]Talent{
	"Keen Vision": {
		Name:        "Keen Vision",
		Description: "Gain a bonus die to Spot Hidden Rolls",
		Type:        Physical,
	},
	"Quick Healer": {
		Name:        "Quick Healer",
		Description: "Natural healing is increased to +3 hit points per day",
		Type:        Physical,
	},
	"Night Vision": {
		Name: "Keen Vision",
		Description: "In darkness, reduce the difficulty level of Spot Hidden rolls and ignore penalty die " +
			"for shooting in the dark",
		Type: Physical,
	},
	"Endurance": {
		Name:        "Endurance",
		Description: "Gain a bonus die when making CON rolls (including to determine MOV rate for chases)",
		Type:        Physical,
	},
	"Power Lifter": {
		Name:        "Power Lifter",
		Description: "Gain a bonus die when making STR rolls to lift objects or people",
		Type:        Physical,
	},
	"Iron Liver": {
		Name: "Iron Liver",
		Description: "May spend 5 Luck to avoid the effects of drinking excessive amounts of alcohol" +
			" (negating penalty applied to skill rolls)",
		Type: Physical,
	},
	"Stout Constitution": {
		Name:        "Stout Constitution",
		Description: "May spend 10 Luck to reduce poison or disease damage and effect by half.",
		Type:        Physical,
	},
	"Tough Guy": {
		Name: "Tough Guy",
		Description: "Soaks up damage, may spend 10 Luck points to shrug off up to 5 hit points worth" +
			" of damage taken in one combat round.",
		Type: Physical,
	},
	"Keen Hearing": {
		Name:        "Keen Hearing",
		Description: "Gain a bonus die to Listen rolls",
		Type:        Physical,
	},
	"Smooth Talker": {
		Name:        "Smooth Talker",
		Description: "Gain a bonus die to Charm Rolls",
		Type:        Physical,
	},
	"Hardened": {
		Name:        "Hardened",
		Description: "Ignores Sanity point loss from attacking other humans, viewing horrific injuries, or the deceased",
		Type:        Mental,
	},
	"Resilient": {
		Name:        "Resilient",
		Description: "May spend Luck points to shrug-off points of Sanity loss, on a one-for-one basis",
		Type:        Mental,
	},
	"Strong Willed": {
		Name:        "Strong Willed",
		Description: "Gains a bonus die when making POW rolls",
		Type:        Mental,
	},
	"Quick Study": {
		Name:        "Quick Study",
		Description: "Halve the time required for Initial and Full Reading of Mythos tomes, as well as other books",
		Type:        Mental,
	},
	"Linguist": {
		Name:        "Linguist",
		Description: "Able to determine what language is being spoken (or what is written); gains a bonus die to Language rolls",
		Type:        Mental,
	},
	"Arcane Insight": {
		Name:        "Arcane Insight",
		Description: "Halve the time required to learn spells and gains bonus die to spell casting rolls",
		Type:        Mental,
	},
	"Photographic Memory": {
		Name:        "Photographic Memory",
		Description: "Can remember many details; gains a bonus die when making Know rolls",
		Type:        Mental,
	},
	"Lore": {
		Name:        "Lore",
		Description: "Has knowledge of a lore specialization skill (e.g. Dream Lore, Vampire Lore, Werewolf Lore, etc.). Note that occupational and/or personal interest skill points should be invested in this skill",
		Type:        Mental,
	},
	"Psychic Power": {
		Name:        "Psychic Power",
		Description: "May choose one psychic power (Clairvoyance, Divination, Medium, Psychometry, or Telekinesis). Note that occupational and/or personal interest skill points should be invested in this skill",
		Type:        Mental,
	},
	"Sharp Witted": {
		Name:        "Sharp Witted",
		Description: "Able to collate facts quickly; gain a bonus die when making Intelligence (but not Idea) rolls",
		Type:        Mental,
	},

	"Alert": {
		Name:        "Alert",
		Description: "Never surprised in combat",
		Type:        Combat,
	},
	"Heavy Hitter": {
		Name:        "Heavy Hitter",
		Description: "May spend 10 Luck points to add an additional damage die when dealing out melee combat (die type depends on the weapon being used, e.g. 1D3 for unarmed combat, 1D6 for a sword, etc.)",
		Type:        Combat,
	},
	"Fast Load": {
		Name:        "Fast Load",
		Description: "Choose a Firearm specialism; ignore penalty die for loading and firing in the same round",
		Type:        Combat,
	},
	"Nimble": {
		Name:        "Nimble",
		Description: "Does not lose next action when \"diving for cover\" versus firearms",
		Type:        Combat,
	},
	"Beady Eye": {
		Name:        "Beady Eye",
		Description: "Does not suffer penalty die when \"aiming\" at a small target (Build -2), and may also fire into melee without a penalty die",
		Type:        Combat,
	},
	"Outmaneuver": {
		Name:        "Outmaneuver",
		Description: "Character is considered to have one point higher Build when initiating a combat maneuver (e.g. Build 1 becomes Build 2 when comparing their Build to the target in a maneuver, reducing the likelihood of suffering a penalty on their Fighting roll)",
		Type:        Combat,
	},
	"Rapid Attack": {
		Name:        "Rapid Attack",
		Description: "May spend 10 Luck points to gain one further melee attack in a single combat round",
		Type:        Combat,
	},
	"Fleet Footed": {
		Name:        "Fleet Footed",
		Description: "May spend 10 Luck to avoid being \"outnumbered\" in melee combat for one combat encounter",
		Type:        Combat,
	},
	"Quick Draw": {
		Name:        "Quick Draw",
		Description: "Does not need to have their firearm \"readied\" to gain +50 DEX when determining position in the DEX order for combat",
		Type:        Combat,
	},
	"Rapid Fire": {
		Name:        "Rapid Fire",
		Description: "Ignores penalty die for multiple handgun shots",
		Type:        Combat,
	},
	"Scary": {
		Name:        "Scary",
		Description: "Reduces difficulty by one level or gains bonus die (at the Keeper's discretion) to Intimidate rolls",
		Type:        Miscellaneous,
	},
	"Gadget": {
		Name:        "Gadget",
		Description: "Starts game with one weird science gadget (see Weird Science, page 86)",
		Type:        Miscellaneous,
	},
	"Lucky": {
		Name:        "Lucky",
		Description: "Regains an additional +1D10 Luck points when Luck Recovery rolls are made",
		Type:        Miscellaneous,
	},
	"Mythos Knowledge": {
		Name:        "Mythos Knowledge",
		Description: "Begins the game with a Cthulhu Mythos Skill of 10 points",
		Type:        Miscellaneous,
	},
	"Weird Science": {
		Name:        "Weird Science",
		Description: "May build and repair weird science devices (see Weird Science, page 86)",
		Type:        Miscellaneous,
	},
	"Shadow": {
		Name:        "Shadow",
		Description: "Reduces difficulty by one level or gains bonus die (at the Keeper's discretion) to Stealth rolls, and if currently unseen is able to make two surprise attacks before their location is discovered",
		Type:        Miscellaneous,
	},
	"Handy": {
		Name:        "Handy",
		Description: "Reduces difficulty by one level or gains bonus die (at the Keeper's discretion) when making Electrical Repair, Mechanical Repair, and Operate Heavy Machinery rolls",
		Type:        Miscellaneous,
	},
	"Animal Companion": {
		Name:        "Animal Companion",
		Description: "Starts game with a faithful animal companion (e.g. dog, cat, parrot) and gains a bonus die when making Animal Handling rolls",
		Type:        Miscellaneous,
	},
	"Master of Disguise": {
		Name:        "Master of Disguise",
		Description: "May spend 10 Luck points to gain a bonus die to Disguise or Art/Craft (Acting) rolls; includes ventriloquism (able to throw voice over long distances so it appears that the sound is emanating from somewhere other than the hero). Note that if someone is trying to detect the disguise their Spot Hidden or Psychology roll's difficulty is raised to Hard",
		Type:        Miscellaneous,
	},
	"Resourceful": {
		Name:        "Resourceful",
		Description: "Always seems to have what they need to hand; may spend 10 Luck points (rather than make Luck roll) to find a certain useful piece of equipment (e.g. a flashlight, length of rope, a weapon, etc.) in their current location",
		Type:        Miscellaneous,
	},
}

package models

type BaseSkillAttribute struct {
	Name       string
	Multiplier int
}

type SkillPointFormula struct {
	BaseAttributes []BaseSkillAttribute // For handling multiple base attributes
	Options        []BaseSkillAttribute // Optional OR cases
}

type SkillChoice struct {
	NumRequired int      // Number of skills that must be chosen from the group
	Skills      []string // List of skills to choose from
}

// SkillRequirement represents either a required skill or a choice between skills
type SkillRequirement struct {
	Type        string      // "required" or "choice"
	Skill       string      // Used when Type is "required"
	SkillChoice SkillChoice // Used when Type is "choice"
}

type Occupation struct {
	Name              string
	SkillRequirements []SkillRequirement
	SuggestedContacts string
	SkillPoints       SkillPointFormula
	CreditRating      struct {
		Min int
		Max int
	}
}

func (o *Occupation) String() string {
	return o.Name
}

var Occupations = map[string]Occupation{
	"Archaeologist": {
		Name: "Archaeologist",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Appraise"},
			{Type: "required", Skill: "Archaeology"},
			{Type: "required", Skill: "History"},
			{Type: "required", Skill: "Library Use"},
			{Type: "required", Skill: "Spot Hidden"},
			{Type: "required", Skill: "Mechanical Repair"},
			{Type: "required", Skill: "Other Language"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Navigate", "Science"},
				},
			},
		},
		SuggestedContacts: "patrons, museums, universities",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{10, 40},
	},
	"Artist": {
		Name: "Artist",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Art/Craft"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"History", "Natural World"},
				},
			},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Other Language"},
			{Type: "required", Skill: "Psychology"},
			{Type: "required", Skill: "Spot Hidden"},
		},
		SuggestedContacts: "art galleries, critics, wealthy patrons, the advertising industry",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrDexterity, Multiplier: 2},
				{Name: AttrPower, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 50},
	},
	"Author": {
		Name: "Author",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Art (Literature)"},
			{Type: "required", Skill: "History"},
			{Type: "required", Skill: "Library Use"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Natural World", "Occult"},
				},
			},
			{Type: "required", Skill: "Other Language"},
			{Type: "required", Skill: "Own Language"},
			{Type: "required", Skill: "Psychology"},
		},
		SuggestedContacts: "publishers, critics, historians, etc",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 30},
	},
	"Aviator": {
		Name: "Aviator",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Accounting"},
			{Type: "required", Skill: "Electrical Repair"},
			{Type: "required", Skill: "Listen"},
			{Type: "required", Skill: "Mechanical Repair"},
			{Type: "required", Skill: "Navigate"},
			{Type: "required", Skill: "Pilot (Aircraft)"},
			{Type: "required", Skill: "Spot Hidden"},
		},
		SuggestedContacts: "old military contacts, other pilots, airfield mechanics, businessmen",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
				{Name: AttrDexterity, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{30, 60},
	},
	"Bank Robber": {
		Name: "Bank Robber",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Drive Auto"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Electrical Repair", "Mechanical Repair"},
				},
			},
			{Type: "required", Skill: "Fighting"},
			{Type: "required", Skill: "Firearms"},
			{Type: "required", Skill: "Intimidate"},
			{Type: "required", Skill: "Locksmith"},
			{Type: "required", Skill: "Operate Heavy Machinery"},
		},
		SuggestedContacts: "other gang members (current and retired), criminal freelancers, organized crime",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrStrength, Multiplier: 2},
				{Name: AttrDexterity, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{5, 75},
	},
	"Bartender/Waitress": {
		Name: "Bartender/Waitress",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Accounting"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 2,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Fighting (Brawl)"},
			{Type: "required", Skill: "Listen"},
			{Type: "required", Skill: "Psychology"},
			{Type: "required", Skill: "Spot Hidden"},
		},
		SuggestedContacts: "regular customers, possibly organized crime",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
				{Name: AttrAppearance, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{8, 25},
	},
	"Beat Cop": {
		Name: "Beat Cop",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Fighting (Brawl)"},
			{Type: "required", Skill: "Firearms"},
			{Type: "required", Skill: "First Aid"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Law"},
			{Type: "required", Skill: "Psychology"},
			{Type: "required", Skill: "Spot Hidden"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Drive Automobile", "Ride"},
				},
			},
		},
		SuggestedContacts: "law enforcement, local businesses and residents, street level crime, organized crime",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrDexterity, Multiplier: 2},
				{Name: AttrStrength, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 30},
	},
	"Big Game Hunter": {
		Name: "Big Game Hunter",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Firearms"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Listen", "Spot Hidden"},
				},
			},
			{Type: "required", Skill: "Natural World"},
			{Type: "required", Skill: "Navigate"},
			{Type: "required", Skill: "Other Language"},
			{Type: "required", Skill: "Survival"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Science (Biology)", "Science (Botany)"},
				},
			},
			{Type: "required", Skill: "Stealth"},
			{Type: "required", Skill: "Track"},
		},
		SuggestedContacts: "foreign government officials, game wardens, past (usually wealthy) clients, black-market gangs and traders, zoo owners",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrDexterity, Multiplier: 2},
				{Name: AttrStrength, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{20, 50},
	},
	"Bounty Hunter": {
		Name: "Bounty Hunter",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Drive Auto"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Mechanical Repair", "Electrical Repair"},
				},
			},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Fighting", "Firearms"},
				},
			},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Fast Talk", "Charm", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Law"},
			{Type: "required", Skill: "Psychology"},
			{Type: "required", Skill: "Track"},
			{Type: "required", Skill: "Stealth"},
		},
		SuggestedContacts: "bail bondsmen, local police, criminal informants",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrDexterity, Multiplier: 2},
				{Name: AttrStrength, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 30},
	},
	"Boxer/Wrestler": {
		Name: "Boxer/Wrestler",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Dodge"},
			{Type: "required", Skill: "Fighting (Brawl)"},
			{Type: "required", Skill: "Intimidate"},
			{Type: "required", Skill: "Jump"},
			{Type: "required", Skill: "Psychology"},
			{Type: "required", Skill: "Spot Hidden"},
		},
		SuggestedContacts: "sports promoters, journalists, organized crime, professional trainers",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
				{Name: AttrStrength, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 60},
	},
	"Butler": {
		Name: "Butler",
		SkillRequirements: []SkillRequirement{
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Accounting", "Appraise"},
				},
			},
			{Type: "required", Skill: "Art/Craft"}, // Note: Any craft type allowed
			{Type: "required", Skill: "First Aid"},
			{Type: "required", Skill: "Listen"},
			{Type: "required", Skill: "Other Language"},
			{Type: "required", Skill: "Psychology"},
			{Type: "required", Skill: "Spot Hidden"},
		},
		SuggestedContacts: "waiting staff of other households, local businesses, and household suppliers",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 40},
	},
	"Cat Burglar": {
		Name: "Cat Burglar",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Appraise"},
			{Type: "required", Skill: "Climb"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Electrical Repair", "Mechanical Repair"},
				},
			},
			{Type: "required", Skill: "Listen"},
			{Type: "required", Skill: "Locksmith"},
			{Type: "required", Skill: "Sleight of Hand"},
			{Type: "required", Skill: "Stealth"},
			{Type: "required", Skill: "Spot Hidden"},
		},
		SuggestedContacts: "fences, other burglars",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
				{Name: AttrDexterity, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{5, 40},
	},
	"Chauffeur": {
		Name: "Chauffeur",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Drive Auto"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 2,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Listen"},
			{Type: "required", Skill: "Mechanical Repair"},
			{Type: "required", Skill: "Navigate"},
			{Type: "required", Skill: "Spot Hidden"},
		},
		SuggestedContacts: "successful business people (criminals included), political representatives",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
				{Name: AttrDexterity, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{10, 40},
	},
	"Confidence Trickster": {
		Name: "Confidence Trickster",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Appraise"},
			{Type: "required", Skill: "Art/Craft (Acting)"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Law", "Other Language"},
				},
			},
			{Type: "required", Skill: "Listen"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 2,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Psychology"},
			{Type: "required", Skill: "Sleight of Hand"},
		},
		SuggestedContacts: "other confidence artists, freelance criminals",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
				{Name: AttrAppearance, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{10, 65},
	},

	"Criminal": {
		Name: "Criminal",
		SkillRequirements: []SkillRequirement{
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Art/Craft", "Disguise"},
				},
			},
			{Type: "required", Skill: "Appraise"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate"},
				},
			},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Fighting", "Firearms"},
				},
			},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Locksmith", "Mechanical Repair"},
				},
			},
			{Type: "required", Skill: "Stealth"},
			{Type: "required", Skill: "Psychology"},
			{Type: "required", Skill: "Spot Hidden"},
		},
		SuggestedContacts: "other criminals, organized crime, law enforcement, street thugs, private detectives",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrDexterity, Multiplier: 2},
				{Name: AttrAppearance, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{5, 65},
	},

	"Cult Leader": {
		Name: "Cult Leader",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Accounting"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 2,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Occult"},
			{Type: "required", Skill: "Psychology"},
			{Type: "required", Skill: "Spot Hidden"},
		},
		SuggestedContacts: "while the majority of followers will be \"regular\" people, the more charismatic the leader, the greater the possibility of celebrity followers, such as movie stars and rich widows",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
				{Name: AttrAppearance, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{30, 60},
	},
	"Dilettante": {
		Name: "Dilettante",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Art/Craft"},
			{Type: "required", Skill: "Firearms"},
			{Type: "required", Skill: "Other Language"},
			{Type: "required", Skill: "Ride"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
		},
		SuggestedContacts: "variable, but usually people of a similar background and tastes, fraternal organizations, bohemian circles, high society at large",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
				{Name: AttrAppearance, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{50, 99},
	},
	"Doctor of Medicine": {
		Name: "Doctor of Medicine",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "First Aid"},
			{Type: "required", Skill: "Medicine"},
			{Type: "required", Skill: "Other Language (Latin)"},
			{Type: "required", Skill: "Psychology"},
			{Type: "required", Skill: "Science (Biology)"},
			{Type: "required", Skill: "Science (Pharmacy)"},
		},
		SuggestedContacts: "other physicians, medical workers, patients, and ex-patients",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{30, 80},
	},
	"Drifter": {
		Name: "Drifter",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Climb"},
			{Type: "required", Skill: "Jump"},
			{Type: "required", Skill: "Listen"},
			{Type: "required", Skill: "Navigate"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Stealth"},
		},
		SuggestedContacts: "other hobos, a few friendly railroad guards, soft touches in numerous towns",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrAppearance, Multiplier: 2},
				{Name: AttrDexterity, Multiplier: 2},
				{Name: AttrStrength, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{0, 5},
	},
	"Elected Official": {
		Name: "Elected Official",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Charm"},
			{Type: "required", Skill: "History"},
			{Type: "required", Skill: "Intimidate"},
			{Type: "required", Skill: "Fast Talk"},
			{Type: "required", Skill: "Listen"},
			{Type: "required", Skill: "Own Language"},
			{Type: "required", Skill: "Persuade"},
			{Type: "required", Skill: "Psychology"},
		},
		SuggestedContacts: "political operatives, government, news media, business, foreign governments, possibly organized crime",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
				{Name: AttrAppearance, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{50, 90},
	},

	"Engineer": {
		Name: "Engineer",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Art/Craft (Technical Drawing)"},
			{Type: "required", Skill: "Electrical Repair"},
			{Type: "required", Skill: "Library Use"},
			{Type: "required", Skill: "Mechanical Repair"},
			{Type: "required", Skill: "Operate Heavy Machinery"},
			{Type: "required", Skill: "Science (Chemistry)"},
			{Type: "required", Skill: "Science (Physics)"},
		},
		SuggestedContacts: "business or military workers, local government, architects",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{30, 60},
	},

	"Entertainer": {
		Name: "Entertainer",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Art/Craft"}, // Acting, Singer, Comedian, etc.
			{Type: "required", Skill: "Disguise"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 2,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Listen"},
			{Type: "required", Skill: "Psychology"},
		},
		SuggestedContacts: "Vaudeville, theater, film industry, entertainment critics, organized crime, and television (for modern-day)",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
				{Name: AttrAppearance, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 70},
	},

	"Exorcist": {
		Name: "Exorcist",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Anthropology"},
			{Type: "required", Skill: "History"},
			{Type: "required", Skill: "Library Use"},
			{Type: "required", Skill: "Listen"},
			{Type: "required", Skill: "Occult"},
			{Type: "required", Skill: "Other Language"},
			{Type: "required", Skill: "Psychology"},
		},
		SuggestedContacts: "Religious organizations",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{25, 55},
	},

	"Explorer": {
		Name: "Explorer",
		SkillRequirements: []SkillRequirement{
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Climb", "Swim"},
				},
			},
			{Type: "required", Skill: "Firearms"},
			{Type: "required", Skill: "History"},
			{Type: "required", Skill: "Jump"},
			{Type: "required", Skill: "Natural World"},
			{Type: "required", Skill: "Navigate"},
			{Type: "required", Skill: "Other Language"},
			{Type: "required", Skill: "Survival"},
		},
		SuggestedContacts: "major libraries, universities, museums, wealthy patrons, other explorers, publishers, foreign government officials, local tribespeople",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrAppearance, Multiplier: 2},
				{Name: AttrDexterity, Multiplier: 2},
				{Name: AttrStrength, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{55, 80},
	},
	"Federal Agent": {
		Name: "Federal Agent",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Drive Auto"},
			{Type: "required", Skill: "Fighting (Brawl)"},
			{Type: "required", Skill: "Firearms"},
			{Type: "required", Skill: "Law"},
			{Type: "required", Skill: "Persuade"},
			{Type: "required", Skill: "Stealth"},
			{Type: "required", Skill: "Spot Hidden"},
		},
		SuggestedContacts: "federal agencies, law enforcement, organized crime",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{20, 40},
	},

	"Gambler": {
		Name: "Gambler",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Accounting"},
			{Type: "required", Skill: "Art/Craft (Acting)"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 2,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Listen"},
			{Type: "required", Skill: "Psychology"},
			{Type: "required", Skill: "Sleight of Hand"},
			{Type: "required", Skill: "Spot Hidden"},
		},
		SuggestedContacts: "bookies, organized crime, street scene",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrAppearance, Multiplier: 2},
				{Name: AttrDexterity, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{8, 50},
	},

	"Gangster, Boss": {
		Name: "Gangster, Boss",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Fighting"},
			{Type: "required", Skill: "Firearms"},
			{Type: "required", Skill: "Law"},
			{Type: "required", Skill: "Listen"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 2,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Psychology"},
			{Type: "required", Skill: "Spot Hidden"},
		},
		SuggestedContacts: "organized crime, street-level crime, police, city government, politicians, judges, unions, lawyers, businesses and residents of the same ethnic community",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
				{Name: AttrAppearance, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{60, 95},
	},

	"Gangster, Underling": {
		Name: "Gangster, Underling",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Drive Auto"},
			{Type: "required", Skill: "Fighting"},
			{Type: "required", Skill: "Firearms"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 2,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Psychology"},
		},
		SuggestedContacts: "street-level crime, police, businesses and residents of the same ethnic community",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrDexterity, Multiplier: 2},
				{Name: AttrStrength, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 20},
	},

	"Laborer": {
		Name: "Laborer",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Drive Auto"},
			{Type: "required", Skill: "Electrical Repair"},
			{Type: "required", Skill: "Fighting"},
			{Type: "required", Skill: "First Aid"},
			{Type: "required", Skill: "Mechanical Repair"},
			{Type: "required", Skill: "Operate Heavy Machinery"},
			{Type: "required", Skill: "Throw"},
		},
		SuggestedContacts: "other workers and supervisors within their industry",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrDexterity, Multiplier: 2},
				{Name: AttrStrength, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{5, 20},
	},

	"Librarian": {
		Name: "Librarian",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Accounting"},
			{Type: "required", Skill: "Library Use"},
			{Type: "required", Skill: "Other Language"},
			{Type: "required", Skill: "Own Language"},
		},
		SuggestedContacts: "booksellers, community groups, specialist researchers",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 35},
	},

	"Mechanic": {
		Name: "Mechanic",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Art/Craft"}, // Carpentry, Welding, Plumbing, etc.
			{Type: "required", Skill: "Climb"},
			{Type: "required", Skill: "Drive Auto"},
			{Type: "required", Skill: "Electrical Repair"},
			{Type: "required", Skill: "Mechanical Repair"},
			{Type: "required", Skill: "Operate Heavy Machinery"},
		},
		SuggestedContacts: "Union members, trade-relevant specialists",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 40},
	},

	"Military Officer": {
		Name: "Military Officer",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Accounting"},
			{Type: "required", Skill: "Firearms"},
			{Type: "required", Skill: "Navigate"},
			{Type: "required", Skill: "First Aid"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 2,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Psychology"},
		},
		SuggestedContacts: "military, federal government",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrDexterity, Multiplier: 2},
				{Name: AttrStrength, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{20, 70},
	},

	"Missionary": {
		Name: "Missionary",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Art/Craft"},
			{Type: "required", Skill: "First Aid"},
			{Type: "required", Skill: "Mechanical Repair"},
			{Type: "required", Skill: "Medicine"},
			{Type: "required", Skill: "Natural World"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
		},
		SuggestedContacts: "church hierarchy, foreign officials",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
				{Name: AttrAppearance, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{0, 30},
	},
	"Musician": {
		Name: "Musician",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Art/Craft (Instrument)"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Listen"},
			{Type: "required", Skill: "Psychology"},
		},
		SuggestedContacts: "club owners, musicians' union, organized crime, street-level criminals",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrAppearance, Multiplier: 2},
				{Name: AttrDexterity, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 30},
	},

	"Nurse": {
		Name: "Nurse",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "First Aid"},
			{Type: "required", Skill: "Listen"},
			{Type: "required", Skill: "Medicine"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Psychology"},
			{Type: "required", Skill: "Science (Biology)"},
			{Type: "required", Skill: "Science (Chemistry)"},
			{Type: "required", Skill: "Spot Hidden"},
		},
		SuggestedContacts: "hospital workers, physicians, community workers",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 30},
	},

	"Occultist": {
		Name: "Occultist",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Anthropology"},
			{Type: "required", Skill: "History"},
			{Type: "required", Skill: "Library Use"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Occult"},
			{Type: "required", Skill: "Other Language"},
			{Type: "required", Skill: "Science (Astronomy)"},
		},
		SuggestedContacts: "libraries, occult societies or fraternities, other occultists",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{10, 80},
	},

	"Parapsychologist": {
		Name: "Parapsychologist",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Anthropology"},
			{Type: "required", Skill: "Art/Craft (Photography)"},
			{Type: "required", Skill: "History"},
			{Type: "required", Skill: "Library Use"},
			{Type: "required", Skill: "Occult"},
			{Type: "required", Skill: "Other Language"},
			{Type: "required", Skill: "Psychology"},
		},
		SuggestedContacts: "universities, parapsychological societies, clients",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 30},
	},

	"Photographer": {
		Name: "Photographer",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Art/Craft (Photography)"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Psychology"},
			{Type: "required", Skill: "Science (Chemistry)"},
			{Type: "required", Skill: "Stealth"},
			{Type: "required", Skill: "Spot Hidden"},
		},
		SuggestedContacts: "advertising industry, local clients (including political organizations and newspapers)",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 30},
	},

	"Police Detective": {
		Name: "Police Detective",
		SkillRequirements: []SkillRequirement{
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Art/Craft (Acting)", "Disguise"},
				},
			},
			{Type: "required", Skill: "Firearms"},
			{Type: "required", Skill: "Law"},
			{Type: "required", Skill: "Listen"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Psychology"},
			{Type: "required", Skill: "Spot Hidden"},
		},
		SuggestedContacts: "law enforcement, street-level crime, coroner's office, judiciary, organized crime",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrDexterity, Multiplier: 2},
				{Name: AttrStrength, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{20, 50},
	},

	"Priest": {
		Name: "Priest",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Accounting"},
			{Type: "required", Skill: "History"},
			{Type: "required", Skill: "Library Use"},
			{Type: "required", Skill: "Listen"},
			{Type: "required", Skill: "Other Language"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Psychology"},
		},
		SuggestedContacts: "church hierarchy, local congregations, community leaders",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 60},
	},

	"Private Investigator": {
		Name: "Private Investigator",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Art/Craft (Photography)"},
			{Type: "required", Skill: "Disguise"},
			{Type: "required", Skill: "Law"},
			{Type: "required", Skill: "Library Use"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Psychology"},
			{Type: "required", Skill: "Spot Hidden"},
		},
		SuggestedContacts: "law enforcement, clients",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrDexterity, Multiplier: 2},
				{Name: AttrStrength, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 30},
	},

	"Professor": {
		Name: "Professor",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Library Use"},
			{Type: "required", Skill: "Other Language"},
			{Type: "required", Skill: "Own Language"},
			{Type: "required", Skill: "Psychology"},
		},
		SuggestedContacts: "scholars, universities, libraries",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{20, 70},
	},

	"Ranger": {
		Name: "Ranger",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Firearms"},
			{Type: "required", Skill: "First Aid"},
			{Type: "required", Skill: "Listen"},
			{Type: "required", Skill: "Natural World"},
			{Type: "required", Skill: "Navigate"},
			{Type: "required", Skill: "Spot Hidden"},
			{Type: "required", Skill: "Survival"},
			{Type: "required", Skill: "Track"},
		},
		SuggestedContacts: "local people and native folk, traders",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrDexterity, Multiplier: 2},
				{Name: AttrStrength, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{5, 20},
	},

	"Reporter": {
		Name: "Reporter",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Art/Craft (Acting)"},
			{Type: "required", Skill: "History"},
			{Type: "required", Skill: "Listen"},
			{Type: "required", Skill: "Own Language"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Psychology"},
			{Type: "required", Skill: "Stealth"},
			{Type: "required", Skill: "Spot Hidden"},
		},
		SuggestedContacts: "news and media industries, political organizations and government, business, law enforcement, street criminals, high and low society",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 30},
	},

	"Sailor": {
		Name: "Sailor",
		SkillRequirements: []SkillRequirement{
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Electrical Repair", "Mechanical Repair"},
				},
			},
			{Type: "required", Skill: "Fighting"},
			{Type: "required", Skill: "Firearms"},
			{Type: "required", Skill: "First Aid"},
			{Type: "required", Skill: "Navigate"},
			{Type: "required", Skill: "Pilot (Boat)"},
			{Type: "required", Skill: "Survival (Sea)"},
			{Type: "required", Skill: "Swim"},
		},
		SuggestedContacts: "military, veterans' associations",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrDexterity, Multiplier: 2},
				{Name: AttrStrength, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 30},
	},
	"Scientist": {
		Name: "Scientist",
		SkillRequirements: []SkillRequirement{
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Computer Use", "Library Use"},
				},
			},
			{Type: "required", Skill: "Other Language"},
			{Type: "required", Skill: "Own Language"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Spot Hidden"},
		},
		SuggestedContacts: "other scientists and academics, universities, their employers and former employers",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 50},
	},

	"Secretary": {
		Name: "Secretary",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Accounting"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Art/Craft (Typing)", "Art/Craft (Short Hand)"},
				},
			},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 2,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Own Language"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Library Use", "Computer Use"},
				},
			},
			{Type: "required", Skill: "Psychology"},
		},
		SuggestedContacts: "other office workers, senior executives in client firms",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrDexterity, Multiplier: 2},
				{Name: AttrAppearance, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 30},
	},

	"Soldier": {
		Name: "Soldier",
		SkillRequirements: []SkillRequirement{
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Climb", "Swim"},
				},
			},
			{Type: "required", Skill: "Dodge"},
			{Type: "required", Skill: "Fighting"},
			{Type: "required", Skill: "Firearms"},
			{Type: "required", Skill: "Stealth"},
			{Type: "required", Skill: "Survival"},
		},
		SuggestedContacts: "military, veterans' associations",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrDexterity, Multiplier: 2},
				{Name: AttrStrength, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{9, 30},
	},

	"Spy": {
		Name: "Spy",
		SkillRequirements: []SkillRequirement{
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Art/Craft (Acting)", "Disguise"},
				},
			},
			{Type: "required", Skill: "Firearms"},
			{Type: "required", Skill: "Listen"},
			{Type: "required", Skill: "Other Language"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Psychology"},
			{Type: "required", Skill: "Sleight of Hand"},
			{Type: "required", Skill: "Stealth"},
		},
		SuggestedContacts: "generally only the person the spy reports to, other connections developed while undercover",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrAppearance, Multiplier: 2},
				{Name: AttrDexterity, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{20, 60},
	},

	"Street Punk": {
		Name: "Street Punk",
		SkillRequirements: []SkillRequirement{
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Fighting"},
			{Type: "required", Skill: "Firearms"},
			{Type: "required", Skill: "Jump"},
			{Type: "required", Skill: "Sleight of Hand"},
			{Type: "required", Skill: "Stealth"},
			{Type: "required", Skill: "Throw"},
		},
		SuggestedContacts: "petty criminals, other punks, the local fence, maybe the local gangster, certainly the local police",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrDexterity, Multiplier: 2},
				{Name: AttrStrength, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{3, 10},
	},

	"Student/Intern": {
		Name: "Student/Intern",
		SkillRequirements: []SkillRequirement{
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Own Language", "Other Language"},
				},
			},
			{Type: "required", Skill: "Library Use"},
			{Type: "required", Skill: "Listen"},
		},
		SuggestedContacts: "academics and other students, while interns may also know business people",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{5, 10},
	},

	"Tribe Member": {
		Name: "Tribe Member",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Climb"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 1,
					Skills:      []string{"Fighting", "Throw"},
				},
			},
			{Type: "required", Skill: "Listen"},
			{Type: "required", Skill: "Natural World"},
			{Type: "required", Skill: "Occult"},
			{Type: "required", Skill: "Spot Hidden"},
			{Type: "required", Skill: "Swim"},
			{Type: "required", Skill: "Survival"},
		},
		SuggestedContacts: "fellow tribe members",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrDexterity, Multiplier: 2},
				{Name: AttrStrength, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{0, 15},
	},

	"Union Activist": {
		Name: "Union Activist",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "Accounting"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 2,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Fighting (Brawl)"},
			{Type: "required", Skill: "Law"},
			{Type: "required", Skill: "Listen"},
			{Type: "required", Skill: "Operate Heavy Machinery"},
			{Type: "required", Skill: "Psychology"},
		},
		SuggestedContacts: "other labor leaders and activists, political friends, possibly organized crime. In the 1920s, also socialists, communists, and subversive anarchists",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{5, 30},
	},

	"Yogi": {
		Name: "Yogi",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "First Aid"},
			{Type: "required", Skill: "History"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 2,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Natural World"},
			{Type: "required", Skill: "Occult"},
			{Type: "required", Skill: "Other Language"},
		},
		SuggestedContacts: "tribespeople, occult or spiritual fraternities, wealthy patrons",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 4},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{6, 60},
	},

	"Zealot": {
		Name: "Zealot",
		SkillRequirements: []SkillRequirement{
			{Type: "required", Skill: "History"},
			{
				Type: "choice",
				SkillChoice: SkillChoice{
					NumRequired: 2,
					Skills:      []string{"Charm", "Fast Talk", "Intimidate", "Persuade"},
				},
			},
			{Type: "required", Skill: "Psychology"},
			{Type: "required", Skill: "Stealth"},
		},
		SuggestedContacts: "religious or fraternal groups, news media",
		SkillPoints: SkillPointFormula{
			BaseAttributes: []BaseSkillAttribute{
				{Name: AttrEducation, Multiplier: 2},
			},
			Options: []BaseSkillAttribute{
				{Name: AttrAppearance, Multiplier: 2},
				{Name: AttrPower, Multiplier: 2},
			},
		},
		CreditRating: struct {
			Min int
			Max int
		}{0, 30},
	},
}

var OccupationsList = func() []string {
	keys := make([]string, 0, len(Occupations))
	for k := range Occupations {
		keys = append(keys, k)
	}
	return keys
}()

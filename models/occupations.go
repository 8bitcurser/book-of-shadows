package models

type BaseSkillAttribute struct {
	Name       string
	Multiplier int
}

type SkillPointFormula struct {
	BaseAttributes []BaseSkillAttribute // For handling multiple base attributes
	Options        []BaseSkillAttribute // Optional OR cases
}

type Occupation struct {
	Name              string            `json:"name"`
	Skills            []string          `json:"-"`
	SuggestedContacts string            `json:"-"`
	SkillPoints       SkillPointFormula `json:"-"`
	CreditRating      struct {
		Min int
		Max int
	} `json:"-"`
}

func (o *Occupation) String() string {
	return o.Name
}

var Occupations = map[string]Occupation{
	"Archaeologist": {
		Name: "Archaeologist",
		Skills: []string{
			"Appraise",
			"Archaeology",
			"History",
			"Other Language (any)",
			"Library Use",
			"Spot Hidden",
			"Mechanical Repair",
			"Navigate or Science (e.g. chemistry, physics, geology, etc.)",
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
		Skills: []string{
			"Art/Craft (any)",
			"History or Natural World",
			"Charm, Fast Talk, Intimidate, or Persuade", // interpersonal skill choice
			"Other Language",
			"Psychology",
			"Spot Hidden",
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
		Skills: []string{
			"Art (Literature)",
			"History",
			"Library Use",
			"Natural World or Occult",
			"Other Language",
			"Own Language",
			"Psychology",
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
		Skills: []string{
			"Accounting",
			"Electrical Repair",
			"Listen",
			"Mechanical Repair",
			"Navigate",
			"Pilot (Aircraft)",
			"Spot Hidden",
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
		Skills: []string{
			"Drive Auto",
			"Electrical or Mechanical Repair",
			"Fighting",
			"Firearms",
			"Intimidate",
			"Locksmith",
			"Operate Heavy Machinery",
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
		Skills: []string{
			"Accounting",
			"Charm, Fast Talk, Intimidate, or Persuade", // Two of these
			"Fighting (Brawl)",
			"Listen",
			"Psychology",
			"Spot Hidden",
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
		Skills: []string{
			"Fighting (Brawl)",
			"Firearms",
			"First Aid",
			"Charm, Fast Talk, Intimidate, or Persuade",
			"Law",
			"Psychology",
			"Spot Hidden",
			"Drive Automobile or Ride",
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
		Skills: []string{
			"Firearms",
			"Listen or Spot Hidden",
			"Natural World",
			"Navigate",
			"Other Language",
			"Survival (any)",
			"Science (Biology or Botany)",
			"Stealth",
			"Track",
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
		Skills: []string{
			"Drive Auto",
			"Mechanical or Electrical Repair",
			"Fighting or Firearms",
			"Fast Talk, Charm, Intimidate, or Persuade",
			"Law",
			"Psychology",
			"Track",
			"Stealth",
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
		Skills: []string{
			"Dodge",
			"Fighting (Brawl)",
			"Intimidate",
			"Jump",
			"Psychology",
			"Spot Hidden",
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
		Skills: []string{
			"Accounting or Appraise",
			"Art/Craft (any, e.g. Cook, Tailor, Barber)",
			"First Aid",
			"Listen",
			"Other Language",
			"Psychology",
			"Spot Hidden",
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
		Skills: []string{
			"Appraise",
			"Climb",
			"Electrical or Mechanical Repair",
			"Listen",
			"Locksmith",
			"Sleight of Hand",
			"Stealth",
			"Spot Hidden",
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
		Skills: []string{
			"Drive Auto",
			"Charm, Fast Talk, Intimidate, or Persuade", // Two of these
			"Listen",
			"Mechanical Repair",
			"Navigate",
			"Spot Hidden",
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
		Skills: []string{
			"Appraise",
			"Art/Craft (Acting)",
			"Law or Other Language",
			"Listen",
			"Charm, Fast Talk, Intimidate, or Persuade", // Two of these
			"Psychology",
			"Sleight of Hand",
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
		Skills: []string{
			"Art/Craft (any) or Disguise",
			"Appraise",
			"Charm, Fast Talk or Intimidate",
			"Fighting or Firearms",
			"Locksmith or Mechanical Repair",
			"Stealth",
			"Psychology",
			"Spot Hidden",
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
		Skills: []string{
			"Accounting",
			"Charm, Fast Talk, Intimidate, or Persuade", // Two of these
			"Occult",
			"Psychology",
			"Spot Hidden",
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
		Skills: []string{
			"Art/Craft (Any)",
			"Firearms",
			"Other Language",
			"Ride",
			"Charm, Fast Talk, Intimidate, or Persuade", // One of these
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
		Skills: []string{
			"First Aid",
			"Medicine",
			"Other Language (Latin)",
			"Psychology",
			"Science (Biology and Pharmacy)",
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
		Skills: []string{
			"Climb",
			"Jump",
			"Listen",
			"Navigate",
			"Charm, Fast Talk, Intimidate, or Persuade", // One of these
			"Stealth",
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
		Skills: []string{
			"Charm",
			"History",
			"Intimidate",
			"Fast Talk",
			"Listen",
			"Own Language",
			"Persuade",
			"Psychology",
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
		Skills: []string{
			"Art/Craft (Technical Drawing)",
			"Electrical Repair",
			"Library Use",
			"Mechanical Repair",
			"Operate Heavy Machine",
			"Science (Chemistry and Physics)",
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
		Skills: []string{
			"Art/Craft (e.g. Acting, Singer, Comedian, etc.)",
			"Disguise",
			"Charm, Fast Talk, Intimidate, or Persuade", // Two of these
			"Listen",
			"Psychology",
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
		Skills: []string{
			"Anthropology",
			"History",
			"Library Use",
			"Listen",
			"Occult",
			"Other Language",
			"Psychology",
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
		Skills: []string{
			"Climb or Swim",
			"Firearms",
			"History",
			"Jump",
			"Natural World",
			"Navigate",
			"Other Language",
			"Survival",
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
		Skills: []string{
			"Drive Auto",
			"Fighting (Brawl)",
			"Firearms",
			"Law",
			"Persuade",
			"Stealth",
			"Spot Hidden",
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
		Skills: []string{
			"Accounting",
			"Art/Craft (Acting)",
			"Charm, Fast Talk, Intimidate, or Persuade", // Two of these
			"Listen",
			"Psychology",
			"Sleight of Hand",
			"Spot Hidden",
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
		Skills: []string{
			"Fighting",
			"Firearms",
			"Law",
			"Listen",
			"Charm, Fast Talk, Intimidate, or Persuade", // Two of these
			"Psychology",
			"Spot Hidden",
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
		Skills: []string{
			"Drive Auto",
			"Fighting",
			"Firearms",
			"Charm, Fast Talk, Intimidate, or Persuade", // Two of these
			"Psychology",
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
		Skills: []string{
			"Drive Auto",
			"Electrical Repair",
			"Fighting",
			"First Aid",
			"Mechanical Repair",
			"Operate Heavy Machinery",
			"Throw",
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
		Skills: []string{
			"Accounting",
			"Library use",
			"Other Language",
			"Own Language",
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
		Skills: []string{
			"Art/Craft (Carpentry, Welding, Plumbing, etc.)",
			"Climb",
			"Drive Auto",
			"Electrical Repair",
			"Mechanical Repair",
			"Operate Heavy Machinery",
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
		Skills: []string{
			"Accounting",
			"Firearms",
			"Navigate",
			"First Aid",
			"Charm, Fast Talk, Intimidate, or Persuade", // Two of these
			"Psychology",
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
		Skills: []string{
			"Art/Craft (any)",
			"First Aid",
			"Mechanical Repair",
			"Medicine",
			"Natural World",
			"Charm, Fast Talk, Intimidate, or Persuade", // One of these
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
		Skills: []string{
			"Art/Craft (Instrument)",
			"Charm, Fast Talk, Intimidate, or Persuade", // One of these
			"Listen",
			"Psychology",
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
		Skills: []string{
			"First Aid",
			"Listen",
			"Medicine",
			"Charm, Fast Talk, Intimidate, or Persuade", // One of these
			"Psychology",
			"Science (Biology and Chemistry)",
			"Spot Hidden",
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
		Skills: []string{
			"Anthropology",
			"History",
			"Library Use",
			"Charm, Fast Talk, Intimidate, or Persuade", // One of these
			"Occult",
			"Other Language",
			"Science (Astronomy)",
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
		Skills: []string{
			"Anthropology",
			"Art/Craft (Photography)",
			"History",
			"Library Use",
			"Occult",
			"Other Language",
			"Psychology",
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
		Skills: []string{
			"Art/Craft (Photography)",
			"Charm, Fast Talk, Intimidate, or Persuade", // One of these
			"Psychology",
			"Science (Chemistry)",
			"Stealth",
			"Spot Hidden",
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
		Skills: []string{
			"Art/Craft (Acting) or Disguise",
			"Firearms",
			"Law",
			"Listen",
			"Charm, Fast Talk, Intimidate, or Persuade", // One of these
			"Psychology",
			"Spot Hidden",
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
		Skills: []string{
			"Accounting",
			"History",
			"Library Use",
			"Listen",
			"Other Language",
			"Charm, Fast Talk, Intimidate, or Persuade", // One of these
			"Psychology",
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
		Skills: []string{
			"Art/Craft (Photography)",
			"Disguise",
			"Law",
			"Library Use",
			"Charm, Fast Talk, Intimidate, or Persuade", // One of these
			"Psychology",
			"Spot Hidden",
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
		Skills: []string{
			"Library Use",
			"Other Language",
			"Own Language",
			"Psychology",
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
		Skills: []string{
			"Firearms",
			"First Aid",
			"Listen",
			"Natural World",
			"Navigate",
			"Spot Hidden",
			"Survival (any)",
			"Track",
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
		Skills: []string{
			"Art/Craft (Acting)",
			"History",
			"Listen",
			"Own Language",
			"Charm, Fast Talk, Intimidate, or Persuade", // One of these
			"Psychology",
			"Stealth",
			"Spot Hidden",
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
		Skills: []string{
			"Electrical or Mechanical Repair",
			"Fighting",
			"Firearms",
			"First Aid",
			"Navigate",
			"Pilot (Boat)",
			"Survival (Sea)",
			"Swim",
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
		Skills: []string{
			"Computer Use or Library Use",
			"Other Language",
			"Own Language",
			"Charm, Fast Talk, Intimidate, or Persuade", // One of these
			"Spot Hidden",
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
		Skills: []string{
			"Accounting",
			"Art/Craft (Typing or Short Hand)",
			"Charm, Fast Talk, Intimidate, or Persuade", // Two of these
			"Own Language",
			"Library Use or Computer Use",
			"Psychology",
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
		Skills: []string{
			"Climb or Swim",
			"Dodge",
			"Fighting",
			"Firearms",
			"Stealth",
			"Survival",
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
		Skills: []string{
			"Art/Craft (Acting) or Disguise",
			"Firearms",
			"Listen",
			"Other Language",
			"Charm, Fast Talk, Intimidate, or Persuade", // One of these
			"Psychology",
			"Sleight of Hand",
			"Stealth",
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
		Skills: []string{
			"Charm, Fast Talk, Intimidate, or Persuade", // One of these
			"Fighting",
			"Firearms",
			"Jump",
			"Sleight of Hand",
			"Stealth",
			"Throw",
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
		Skills: []string{
			"Language (Own or Other)",
			"Library Use",
			"Listen",
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
		Skills: []string{
			"Climb",
			"Fighting or Throw",
			"Listen",
			"Natural World",
			"Occult",
			"Spot Hidden",
			"Swim",
			"Survival (any)",
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
		Skills: []string{
			"Accounting",
			"Charm, Fast Talk, Intimidate, or Persuade", // Two of these
			"Fighting (Brawl)",
			"Law",
			"Listen",
			"Operate Heavy Machinery",
			"Psychology",
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
		Skills: []string{
			"First Aid",
			"History",
			"Charm, Fast Talk, Intimidate, or Persuade", // Two of these
			"Natural World",
			"Occult",
			"Other Language",
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
		Skills: []string{
			"History",
			"Charm, Fast Talk, Intimidate, or Persuade", // Two of these
			"Psychology",
			"Stealth",
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

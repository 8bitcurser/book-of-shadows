package models

import (
	"math/rand"
)

type SpecialArchetypeRules struct {
	RecommendedTalents []string
	RequiredTalents    []string
	Notes              string
}

type Archetype struct {
	Name                  string                `json:"name"`
	Skills                []string              `json:"-"`
	BonusPoints           int                   `json:"-"`
	CoreCharacteristic    []string              `json:"-"`
	SuggestedOccupations  []string              `json:"-"`
	AmountOfTalents       int                   `json:"-"`
	Description           string                `json:"description"`
	SuggestedTraits       string                `json:"-"`
	SpecialArchetypeRules SpecialArchetypeRules `json:"-"`
}

func (a *Archetype) String() string {
	return a.Name
}

func PickRandomArchetype() *Archetype {
	archetypeName := ArchetypesList[rand.Intn(len(ArchetypesList))]
	archetype := Archetypes[archetypeName]
	return &archetype
}

var Archetypes = map[string]Archetype{
	"Adventurer": {
		Name: "Adventurer",
		Skills: []string{
			"Climb", "Diving", "Drive Auto", "First Aid",
			"Fighting (any)", "Firearms (any)", "Jump",
			"Language Other (any)", "Mechanical Repair",
			"Pilot (any)", "Ride", "Stealth",
			"Survival (any)", "Swim",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrDexterity, AttrAppearance}, // Note: rules say "choose either DEX or APP"
		SuggestedOccupations: []string{
			"Actor", "Archaeologist", "Athlete", "Aviator",
			"Bank Robber", "Big Game Hunter", "Cat Burglar",
			"Dilettante", "Drifter", "Gambler", "Gangster",
			"Hobo", "Investigative Journalist", "Missionary",
			"Nurse", "Photographer", "Ranger", "Sailor",
			"Soldier", "Tribe Member",
		},
		AmountOfTalents: 2,
		Description:     "A life without adventure is not worth living. The world is a big place and there is much to experience and many chances for glory. Sitting behind a desk, working a job nine to five is a death sentence for such folk. The adventurer yearns for excitement, fun, and a challenge.",
		SuggestedTraits: "easily bored, tenacious, glory hunter, egotistical",
	},
	"Beefcake": {
		Name: "Beefcake",
		Description: "Physical, muscular, and capable of handling themselves when " +
			"the chips are down. Born that way or has worked hard in the " +
			"pursuit of physical perfection. You won't find these guys and " +
			"gals in the library, but you might see their faces on a billboard. " +
			"Beefcakes come in two varieties: the caring, silent type, or the " +
			"brazen loud-mouth.",
		Skills: []string{
			"Climb", "Fighting (Brawl)", "Intimidate", "Listen",
			"Mechanical Repair", "Psychology", "Swim", "Throw",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrStrength},
		SuggestedOccupations: []string{
			"Athlete", "Beat Cop", "Bounty Hunter", "Boxer",
			"Entertainer", "Gangster", "Hired Muscle", "Hobo",
			"Itinerant Worker", "Laborer", "Mechanic", "Sailor",
			"Soldier", "Street Punk", "Tribe Member",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "domineering, brash, quiet, soft-centered, slow to anger, quick to anger",
	},
	"Cold Blooded": {
		Name: "Cold Blooded",
		Description: "A rationalist who is capable of just about anything. Cold " +
			"blooded types may follow some twisted moral code, however, " +
			"their view of humanity is cold and stark; you're either good " +
			"or bad. There are no shades of gray to navigate, just the harsh " +
			"realities of life and death. Such people make effective killers " +
			"as they have little self-doubt; they are ready to follow orders " +
			"to the letter, or pursue some personal agenda for revenge. Such " +
			"people may do anything to get the job done. They are rarely " +
			"spontaneous people; instead, they embody ruthlessness and " +
			"premeditation. Sometimes they will try to fool themselves " +
			"into believing they have a \"line\" they will not cross, when in " +
			"reality they are merciless and will go to any length to fulfill " +
			"what they see as their goal.",
		Skills: []string{
			"Art/Craft (Acting)", "Disguise", "Fighting (any)",
			"Firearms (any)", "First Aid", "History", "Intimidate",
			"Law", "Listen", "Mechanical Repair", "Psychology",
			"Stealth", "Survival (any)", "Track",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrIntelligence},
		SuggestedOccupations: []string{
			"Bank Robber", "Beat Cop", "Bounty Hunter", "Cult Leader",
			"Drifter", "Exorcist", "Federal Agent", "Gangster",
			"Gun Moll", "Hired Muscle", "Hit Man", "Professor",
			"Reporter", "Soldier", "Street Punk", "Tribe Member",
			"Zealot",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "rationalist, sees everything in black and white, ruthless, callous, brutal, pitiless, hardnosed",
	},
	"Dreamer": {
		Name: "Dreamer",
		Description: "Whether an idealist or visionary, the dreamer has a strong and " +
			"powerful mind. Such types tend to follow their own direction " +
			"in life. The dreamer looks beyond the mundane realities of " +
			"life, perhaps as a form of escapism or because they yearn for " +
			"\"what could be,\" wishing to right wrongs or improve the world " +
			"around them.",
		Skills: []string{
			"Art/Craft (any)", "Charm", "History",
			"Language Other (any)", "Library Use", "Listen",
			"Natural World", "Occult",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrPower},
		SuggestedOccupations: []string{
			"Artist", "Author", "Bartender/Waitress", "Priest",
			"Cult Leader", "Dilettante", "Drifter", "Elected Official",
			"Gambler", "Gentleman/Lady", "Hobo", "Hooker", "Librarian",
			"Musician", "Nurse", "Occultist", "Professor", "Secretary",
			"Student", "Tribe Member",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "idealist, optimist, lazy, generous, quiet, thoughtful, always late",
	},
	"Egghead": {
		Name: "Egghead",
		Description: "Everything can be broken down and analyzed in order to " +
			"understand how it works. Knowledge is a treasure and a " +
			"joy—a puzzle to explore. Where the scholar is bookish, the " +
			"egghead is practical and thoroughly enjoys getting their hands " +
			"dirty. Whether it's wires and gears, valves and computational " +
			"engines, or blood and bones, the egghead likes to figure out " +
			"what makes things tick. Perhaps an absent-minded genius " +
			"or a razor-sharp virtuoso, the egghead can easily become " +
			"absorbed in the problem before them, leaving them exposed " +
			"and unaware of what is actually happening around them. " +
			"Depending on the pulp level of your game, the egghead may " +
			"be able to invent all manner of gizmos, useful or otherwise, " +
			"see Weird Science on page 86 for details.",
		Skills: []string{
			"Anthropology", "Appraise", "Computer Use",
			"Electrical Repair", "Language Other (any)",
			"Library Use", "Mechanical Repair",
			"Operate Heavy Machinery", "Science (any)",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrIntelligence, AttrEducation}, // Note: choose either INT or EDU
		SuggestedOccupations: []string{
			"Butler", "Cult Leader", "Doctor of Medicine",
			"Engineer", "Gentleman/Lady", "Investigative Journalist",
			"Mechanic", "Priest", "Scientist",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "knowledgeable, focused, tunnel vision",
	},
	"Explorer": {
		Name: "Explorer",
		Description: "\"Don't fence me in,\" is the oft-heard cry of the explorer, who " +
			"wishes for a more authentic and fulfilling life. Strong willed, " +
			"virtually unshakeable, the explorer is ever questing for what " +
			"lies over the horizon. Possibly at one with nature, such types " +
			"are content to sleep where they fall, happily disdaining the " +
			"soft comforts of urban life. Whether hacking through jungles, " +
			"squeezing through caverns, or simply charting the hidden quarters " +
			"of the city, the explorer is often a misfit who grows restless and " +
			"annoyed by those they consider to be \"weak\" or \"cowards.\"",
		Skills: []string{
			"Animal Handling",
			"Anthropology",
			"Archaeology",
			"Climb",
			"Fighting (Brawl)",
			"First Aid",
			"Jump",
			"Language Other (any)",
			"Natural World",
			"Navigate",
			"Pilot (any)",
			"Ride",
			"Stealth",
			"Survival (any)",
			"Track",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrDexterity, AttrPower}, // Note: choose either DEX or POW
		SuggestedOccupations: []string{
			"Agency Detective",
			"Archaeologist",
			"Big Game Hunter",
			"Bounty Hunter",
			"Dilettante",
			"Explorer",
			"Get-Away Driver",
			"Gun Moll",
			"Itinerant Worker",
			"Investigative Journalist",
			"Missionary",
			"Photographer",
			"Ranger",
			"Sailor",
			"Soldier",
			"Tribe Member",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "outcast, brave, misfit, loner, bullish, strong willed, leader, restless",
	},
	"Femme Fatale": {
		Name: "Femme Fatale",
		Description: "A deadly woman or man whose outward beauty usually masks " +
			"a self-centered approach to life; one who is ever vigilant. By " +
			"constructing an alluring and glamorous persona the femme " +
			"fatale is akin to a spider. She draws others to her web in order " +
			"to possess what she desires or destroy her target. Brave and " +
			"cunning, the femme fatale is not shy of getting her hands dirty " +
			"and is a capable foe. Neither is she foolhardy, and she will wait " +
			"until her web is constructed before dealing out a sudden and " +
			"well-timed assault (be it mental or physical). A classic pulp " +
			"archetype, the femme fatale could as easily be termed homme " +
			"fatale if so desired.",
		Skills: []string{
			"Art/Craft (Acting)",
			"Appraise",
			"Charm",
			"Disguise",
			"Drive Auto",
			"Fast Talk",
			"Fighting (Brawl)",
			"Firearms (Handgun)",
			"Listen",
			"Psychology",
			"Sleight of Hand",
			"Stealth",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrAppearance, AttrIntelligence}, // Note: choose either APP or INT
		SuggestedOccupations: []string{
			"Actor",
			"Agency Detective",
			"Author",
			"Cat Burglar",
			"Confidence Trickster",
			"Dilettante",
			"Elected Official",
			"Entertainer",
			"Federal Agent",
			"Gangster",
			"Gun Moll",
			"Hit Man",
			"Hooker",
			"Investigative Journalist",
			"Musician",
			"Nurse",
			"Private Investigator",
			"Reporter",
			"Spy",
			"Zealot",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "alluring, glamorous, wicked, deceitful, cunning, focused, fraudulent",
		// Could add a special note about the Smooth Talker talent being recommended
		SpecialArchetypeRules: SpecialArchetypeRules{
			RecommendedTalents: []string{"Smooth Talker"},
			// Could potentially add a field for gender flexibility note
			Notes: "Can be played as homme fatale instead of femme fatale",
		},
	},
	"Grease Monkey": {
		Name: "Grease Monkey",
		Description: "The grease monkey is practically minded, able to make " +
			"and repair all manner of things, be they useful inventions, " +
			"machines, engines, or other devices. Grease Monkeys may " +
			"be found tinkering under the hood of a car, or playing with the " +
			"telephone exchange wires. Such types have a \"can do\" attitude, " +
			"able to make the most of what they have at hand, using their " +
			"skills and experience to wow those around them. " +
			"Depending on the pulp level of your game, the grease " +
			"monkey may be able to \"jury-rig\" all manner of gizmos, useful " +
			"or otherwise; see Weird Science on page 86 for details.",
		Skills: []string{
			"Appraise",
			"Art/Craft (any)",
			"Fighting (Brawl)",
			"Drive Auto",
			"Electrical Repair",
			"Locksmith",
			"Mechanical Repair",
			"Operate Heavy Machinery",
			"Spot Hidden",
			"Throw",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrIntelligence},
		SuggestedOccupations: []string{
			"Bartender/Waitress",
			"Butler",
			"Cat Burglar",
			"Chauffeur",
			"Drifter",
			"Engineer",
			"Get-Away Driver",
			"Hobo",
			"Itinerant Worker",
			"Mechanic",
			"Sailor",
			"Soldier",
			"Student",
			"Union Activist",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "practical, hands-on, hard working, oil-stained, capable",
		SpecialArchetypeRules: SpecialArchetypeRules{
			RecommendedTalents: []string{"Weird Science"},
			Notes:              "Can create jury-rigged inventions depending on pulp level",
		},
	},
	"Hard Boiled": {
		Name: "Hard Boiled",
		Description: "Tough and streetwise, someone who is hard boiled understands " +
			"that to catch a thief you have to think like a thief. Usually, " +
			"such a person isn't above breaking the law in order to get the " +
			"job done. They'll use whatever tools are at their disposal and " +
			"may crack a few skulls in the process. Often, at their core, they " +
			"are honest souls who wish the world wasn't so despicable and " +
			"downright nasty; however, in order to fight for justice, they " +
			"can be just as nasty as they need to be.",
		Skills: []string{
			"Art/Craft (any)",
			"Fighting (Brawl)",
			"Firearms (any)",
			"Drive Auto",
			"Fast Talk",
			"Intimidate",
			"Law",
			"Listen",
			"Locksmith",
			"Sleight of Hand",
			"Spot Hidden",
			"Stealth",
			"Throw",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrConstitution},
		SuggestedOccupations: []string{
			"Agency Detective",
			"Bank Robber",
			"Beat Cop",
			"Bounty Hunter",
			"Boxer",
			"Gangster",
			"Gun Moll",
			"Laborer",
			"Police Detective",
			"Private Investigator",
			"Ranger",
			"Union Activist",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "cynical, objective, practical, world-weary, corrupt, violent",
	},
	"Harlequin": {
		Name: "Harlequin",
		Description: "While similar to the femme fatale, the harlequin does not " +
			"like to get their hands dirty (if they can help it). Usually " +
			"possessing a magnetic personality, although not necessarily " +
			"classically beautiful, such types find enjoyment in manipulating " +
			"others to do their bidding, and often hide their own agendas " +
			"behind outright lies or subtle deceptions. Sometimes they " +
			"are committed to a cause (personal or otherwise), or act like " +
			"agents of chaos, delighting in watching how people react to " +
			"the situations they construe.",
		Skills: []string{
			"Art/Craft (Acting)",
			"Charm",
			"Climb",
			"Disguise",
			"Fast Talk",
			"Jump",
			"Language Other (any)",
			"Listen",
			"Persuade",
			"Psychology",
			"Sleight of Hand",
			"Stealth",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrAppearance},
		SuggestedOccupations: []string{
			"Actor",
			"Agency Detective",
			"Artist",
			"Bartender/Waitress",
			"Confidence Trickster",
			"Cult Leader",
			"Dilettante",
			"Elected Official",
			"Entertainer",
			"Gambler",
			"Gentleman/Lady",
			"Musician",
			"Reporter",
			"Secretary",
			"Union Activist",
			"Zealot",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "calculating, cunning, two-faced, manipulative, chaotic, wild, flamboyant",
	},
	"Hunter": {
		Name: "Hunter",
		Description: "Maybe it's the thrill of the chase, the prize at the end, or just " +
			"because they have an innate drive to master their environment, " +
			"the hunter is relentless in pursuing their prey. Calm and " +
			"calculated, the hunter is willing to wait for the most opportune " +
			"moment, despising the reckless behavior of the unwary.",
		Skills: []string{
			"Animal Handling",
			"Fighting (any)",
			"Firearms (Rifle and/or Handgun)",
			"First Aid",
			"Listen",
			"Natural World",
			"Navigate",
			"Spot Hidden",
			"Stealth",
			"Survival (any)",
			"Swim",
			"Track",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrIntelligence, AttrConstitution},
		SuggestedOccupations: []string{
			"Agency Detective",
			"Bank Robber",
			"Beat Cop",
			"Bounty Hunter",
			"Boxer",
			"Gangster",
			"Gun Moll",
			"Laborer",
			"Police Detective",
			"Private Investigator",
			"Ranger",
			"Union Activist",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "relentless, cunning, patient, driven, calm, quiet",
	},
	"Mystic": {
		Name: "Mystic",
		Description: "A seeker of the hidden, explorer of the unseen realm; the " +
			"mystic quests for secrets and the fundamental truth of " +
			"existence. They may be book-learned academics, shamanistic " +
			"healers, circus diviners, or visionaries, but all pursue knowledge " +
			"and the experience of forces outside of the natural order, be it " +
			"for personal gain or the betterment of mankind. " +
			"With the Keeper's permission, a mystic is able to tap into " +
			"supernatural powers beyond the ken of average folk. Often " +
			"they have been persecuted and hunted, hiding their \"gifts\" " +
			"from those who would call them \"witch,\" while others are " +
			"considered charlatans and little more than sideshow freaks. " +
			"Such heroes must take the Psychic talent, allowing them to " +
			"invest skill points in one or more psychic skills (see Psychic " +
			"Powers, page 83).",
		Skills: []string{
			"Art/Craft (any)",
			"Science (Astronomy)",
			"Disguise",
			"History",
			"Hypnosis",
			"Language Other (any)",
			"Natural World",
			"Occult",
			"Psychology",
			"Sleight of Hand",
			"Stealth",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrPower},
		SuggestedOccupations: []string{
			"Artist",
			"Cult Leader",
			"Dilettante",
			"Exorcist",
			"Entertainer",
			"Occultist",
			"Parapsychologist",
			"Tribe Member",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "collector, knowledgeable, responsible, calculating, opportunist, shrewd, studious, risk taker, wise",
		SpecialArchetypeRules: SpecialArchetypeRules{
			RequiredTalents: []string{"Psychic"},
			Notes:           "Must invest skill points in chosen psychic skill(s)",
		},
	},
	"Rogue": {
		Name: "Rogue",
		Description: "The rogue disobeys rules of society, openly questioning the " +
			"status quo and mocking those in authority. They delight " +
			"in being non-conformists, acting on impulse and deriding " +
			"conventional behavior. Laws are there to be broken or skirted " +
			"around. Most rogues are not necessarily criminals or anarchists " +
			"intent on spreading chaos, but rather they find amusement " +
			"in pulling off stunts that will confound others. They are often " +
			"sophisticated, governed by their own unique moral codes, " +
			"loveable and careless.",
		Skills: []string{
			"Appraise",
			"Art/Craft (any)",
			"Charm",
			"Disguise",
			"Fast Talk",
			"Law",
			"Locksmith",
			"Psychology",
			"Read Lips",
			"Spot Hidden",
			"Stealth",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrDexterity, AttrAppearance},
		SuggestedOccupations: []string{
			"Artist",
			"Bank Robber",
			"Cat Burglar",
			"Confidence Trickster",
			"Dilettante",
			"Entertainer",
			"Gambler",
			"Get-Away Driver",
			"Spy",
			"Student",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "charming, disarming, self-absorbed, crafty, shrewd, scheming",
	},
	"Outsider": {
		Name: "Outsider",
		Description: "The outsider stands apart from the rest of society, either " +
			"figuratively or literally. Such people may be alien to the " +
			"environment in which they find themselves, perhaps from a " +
			"different country or culture, or they are part of the society but " +
			"find themselves at odds with it. The outsider is usually on some " +
			"form of journey, physically or spiritually, and must complete " +
			"their objective before they can return to, or at last feel part of, " +
			"the greater whole. Often the outsider will have distinct skills " +
			"or a different way of approaching things; utilizing forgotten, " +
			"secret, or alien knowledge.",
		Skills: []string{
			"Art/Craft (any)",
			"Animal Handling",
			"Fighting (any)",
			"First Aid",
			"Intimidate",
			"Language Other (any)",
			"Listen",
			"Medicine",
			"Navigation",
			"Stealth",
			"Survival (any)",
			"Track",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrIntelligence, AttrConstitution}, // Choose either INT or CON
		SuggestedOccupations: []string{
			"Artist",
			"Drifter",
			"Explorer",
			"Hired Muscle",
			"Itinerant Worker",
			"Laborer",
			"Nurse",
			"Occultist",
			"Ranger",
			"Tribe Member",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "cold, quiet, detached, indifferent, brutal",
		SpecialArchetypeRules: SpecialArchetypeRules{
			Notes: "Character should have a defined objective or journey that sets them apart from society",
		},
	},
	"Bon Vivant": {
		Name: "Bon Vivant",
		Description: "A bon vivant is \"one who lives well,\" but that doesn't necessarily " +
			"mean they are rich. While many are accustomed to wealth, the " +
			"bon vivant is someone who could be said to enjoy life to the " +
			"fullest and damn the consequences! Why wait until tomorrow " +
			"when you can start living life today? Enjoying food and drink, " +
			"as well as other pleasurable pursuits, is the key to a lifestyle " +
			"where excess is the norm. Whether poor or rich, such a person " +
			"puts little thought to saving for a rainy day, preferring to be " +
			"the center of attention and a friend to all.",
		Skills: []string{
			"Appraise",
			"Art/Craft (any)",
			"Charm",
			"Fast Talk",
			"Language Other (any)",
			"Listen",
			"Spot Hidden",
			"Psychology",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrSize}, // Unique core characteristic
		SuggestedOccupations: []string{
			"Actor",
			"Artist",
			"Butler",
			"Confidence Trickster",
			"Cult Leader",
			"Dilettante",
			"Elected Official",
			"Entertainer",
			"Gambler",
			"Gun Moll",
			"Gentleman/Lady",
			"Military Officer",
			"Musician",
			"Priest",
			"Professor",
			"Zealot",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "excessive, greedy, hoarder, collector, name-dropper, boastful, attention seeking, kind, generous",
		SpecialArchetypeRules: SpecialArchetypeRules{
			Notes: "Character should demonstrate a lifestyle of excess regardless of wealth level",
		},
	},

	"Scholar": {
		Name: "Scholar",
		Description: "Uses intelligence and analysis to understand the world around " +
			"them. Normally quite happy sitting in the library with a book " +
			"(rather than actually facing the realities of life). A seeker of " +
			"knowledge, the scholar is not particularly action orientated; " +
			"however, when it comes to the crunch, he or she might be the " +
			"only person who knows what to do.",
		Skills: []string{
			"Accounting",
			"Anthropology",
			"Cryptography",
			"History",
			"Language Other (any)",
			"Library Use",
			"Medicine",
			"Natural World",
			"Occult",
			"Science (any)",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrEducation},
		SuggestedOccupations: []string{
			"Archaeologist",
			"Author",
			"Doctor of Medicine",
			"Librarian",
			"Parapsychologist",
			"Professor",
			"Scientist",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "studious, bookish, superiority complex, condescending, loner, fussy, speaks too quickly, pensive",
		SpecialArchetypeRules: SpecialArchetypeRules{
			Notes: "Always begins the game as a non-believer of the Mythos (see Chapter 9: Sanity, Call of Cthulhu Rulebook)",
		},
	},

	"Seeker": {
		Name: "Seeker",
		Description: "Puzzles and riddles enthrall the seeker, who uses intelligence " +
			"and reasoning to uncover mysteries and solve problems. They " +
			"look for and enjoy mental challenges, always focused on " +
			"finding the truth, no matter the consequences or tribulations " +
			"they must face.",
		Skills: []string{
			"Accounting",
			"Appraise",
			"Disguise",
			"History",
			"Law",
			"Library Use",
			"Listen",
			"Occult",
			"Psychology",
			"Science (any)",
			"Spot Hidden",
			"Stealth",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrIntelligence},
		SuggestedOccupations: []string{
			"Agency Detective",
			"Author",
			"Beat Cop",
			"Federal Agent",
			"Investigative Journalist",
			"Occultist",
			"Parapsychologist",
			"Police Detective",
			"Reporter",
			"Spy",
			"Student",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "risk taker, tunnel vision, deceitful, boastful, driven",
	},
	"Sidekick": {
		Name: "Sidekick",
		Description: "The sidekick embodies aspects of the steadfast, rogue, and " +
			"thrill seeker archetypes. Usually, a younger person who has " +
			"yet to live up to their full potential, someone who seeks to " +
			"learn from a mentor type figure, or those content not to be " +
			"the center of attention. Alternatively, the sidekick wishes to " +
			"belong, to be the hero but is overshadowed by their peers or " +
			"mentor. Subordinate sidekicks can at times struggle against " +
			"their (usually) self-imposed restraints, venturing off on flights " +
			"of fancy that mostly just get them into trouble. Sidekicks " +
			"usually possess a strong moral code of duty and responsibility.",
		Skills: []string{
			"Animal Handling",
			"Climb",
			"Electrical Repair",
			"Fast Talk",
			"First Aid",
			"Jump",
			"Library Use",
			"Listen",
			"Navigate",
			"Photography",
			"Science (any)",
			"Stealth",
			"Track",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrDexterity, AttrConstitution}, // Choose either DEX or CON
		SuggestedOccupations: []string{
			"Author",
			"Bartender/Waitress",
			"Beat Cop",
			"Butler",
			"Chauffeur",
			"Doctor of Medicine",
			"Federal Agent",
			"Get-Away Driver",
			"Gun Moll",
			"Hobo",
			"Hooker",
			"Laborer",
			"Librarian",
			"Nurse",
			"Photographer",
			"Scientist",
			"Secretary",
			"Street Punk",
			"Student",
			"Tribe Member",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "helpful, resourceful, loyal, accident-prone, questioning, inquisitive, plucky",
		SpecialArchetypeRules: SpecialArchetypeRules{
			Notes: "Character should have a defined relationship with a mentor or peer figure they look up to",
		},
	},
	"Steadfast": {
		Name: "Steadfast",
		Description: "Moral righteousness runs thickly in the blood of the steadfast. " +
			"They protect the weak, put the interests of the others before " +
			"themselves, and would willingly sacrifice their life for another's " +
			"safety. Whether they follow a clear spiritual or religious path " +
			"or some internal moral code, they do not stoop to the depths " +
			"of others, fighting with honor and acting as role models to " +
			"those around them. Whatever else they fight for, they also " +
			"fight for justice.",
		Skills: []string{
			"Accounting",
			"Drive Auto",
			"Fighting (any)",
			"Firearms (Handgun)",
			"First Aid",
			"History",
			"Intimidate",
			"Law",
			"Natural World",
			"Navigate",
			"Persuade",
			"Psychology",
			"Ride",
			"Spot Hidden",
			"Survival (any)",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrConstitution},
		SuggestedOccupations: []string{
			"Athlete", "Beat Cop", "Butler", "Priest",
			"Chauffeur", "Doctor of Medicine", "Elected Official",
			"Exorcist", "Federal Agent", "Gentleman/Lady",
			"Missionary", "Nurse", "Police Detective",
			"Private Detective", "Reporter", "Sailor",
			"Soldier", "Tribe Member",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "unwavering, loyal, resolute, committed, dedicated, firm but fair, faithful",
	},

	"Swashbuckler": {
		Name: "Swashbuckler",
		Description: "Passionate and idealistic souls who are always looking to rescue " +
			"damsels in distress. Gallant and heroic, the swashbuckler is " +
			"action-orientated and fights fairly, disdaining the use of " +
			"firearms as the tools of cowards. Most likely boastful, noisy, " +
			"and joyous, even when in the direst of situations. A romantic " +
			"at heart, a swashbuckler possesses a strong code of honor but " +
			"is prone to reckless behavior",
		Skills: []string{
			"Art/Craft (any)",
			"Charm",
			"Climb",
			"Fighting (any)",
			"Jump",
			"Language Other (any)",
			"Mechanical Repair",
			"Navigate",
			"Pilot (any)",
			"Stealth",
			"Swim",
			"Throw",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrDexterity, AttrAppearance},
		SuggestedOccupations: []string{
			"Actor", "Artist", "Aviator", "Big Game Hunter",
			"Bounty Hunter", "Dilettante", "Entertainer",
			"Gentleman/Lady", "Investigative Journalist",
			"Military Officer", "Missionary", "Private Detective",
			"Ranger", "Sailor", "Soldier", "Spy",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "boastful, gallant, action-orientated, romantic, passionate, highly-strung",
		SpecialArchetypeRules: SpecialArchetypeRules{
			Notes: "Disdains the use of firearms, preferring melee combat",
		},
	},

	"Thrill seeker": {
		Name: "Thrill Seeker",
		Description: "Some people are like moths to a flame. For them, the easy life " +
			"is no life at all, and they must seek out adventure and danger " +
			"in order to feel alive. The stakes are never high enough for " +
			"thrill seekers, who are always ready to bet large in order to " +
			"feel the rush of adrenaline pumping through their veins. Such " +
			"daredevils are drawn to high-octane sports and activities, and " +
			"for them, a mountain is a challenge to master. Foolhardy to a " +
			"fault, they cannot understand why no one else is prepared to " +
			"take the same risks as they do.",
		Skills: []string{
			"Art/Craft (any)",
			"Charm",
			"Climb",
			"Diving",
			"Drive Auto",
			"Fast Talk",
			"Jump",
			"Mechanical Repair",
			"Navigate",
			"Pilot (any)",
			"Ride",
			"Stealth",
			"Survival (any)",
			"Swim",
			"Throw",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrDexterity, AttrPower},
		SuggestedOccupations: []string{
			"Actor", "Athlete", "Aviator", "Bank Robber",
			"Bounty Hunter", "Cat Burglar", "Dilettante",
			"Entertainer", "Explorer", "Gambler", "Gangster",
			"Get-Away Driver", "Gun Moll", "Gentleman/Lady",
			"Hooker", "Investigative Journalist", "Missionary",
			"Musician", "Occultist", "Parapsychologist",
			"Ranger", "Sailor", "Soldier", "Spy", "Union Activist",
			"Zealot",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "daredevil, risk taker, manic, exhibitionist, braggart, trouble maker",
	},

	"Two-Fisted": {
		Name: "Two-Fisted",
		Description: "\"Live fast, die hard\" is the motto of the two-fisted. Such " +
			"individuals are storehouses of energy, strong, tough, and very " +
			"capable. Such types are inclined to resolve disputes with their " +
			"fists rather than words. Usually hard-drinking and hard-" +
			"talking, they like getting straight to the point and dislike pomp " +
			"and ceremony. They do not suffer fools gladly. The two-fisted " +
			"seem to live life in a hurry, quick to anger, contemptuous of " +
			"authority, and ready to play as dirty as the next guy.",
		Skills: []string{
			"Drive Auto",
			"Fighting (Brawl)",
			"Firearms (any)",
			"Intimidate",
			"Listen",
			"Mechanical Repair",
			"Spot Hidden",
			"Swim",
			"Throw",
		},
		BonusPoints:        100,
		CoreCharacteristic: []string{AttrStrength, AttrSize},
		SuggestedOccupations: []string{
			"Agency Detective", "Bank Robber", "Beat Cop",
			"Boxer", "Gangster", "Gun Moll", "Hired Muscle",
			"Hit Man", "Hooker", "Laborer", "Mechanic",
			"Nurse", "Police Detective", "Ranger", "Reporter",
			"Sailor", "Soldier", "Street Punk", "Tribe Member",
			"Union Activist",
		},
		AmountOfTalents: 2,
		SuggestedTraits: "tough, capable, determined, quick to anger, violent, dirty, corrupt, underhand",
	},
}
var ArchetypesList = func() []string {
	keys := make([]string, 0, len(Archetypes))
	for k := range Archetypes {
		keys = append(keys, k)
	}
	return keys
}()

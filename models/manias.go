package models

// Mania represents an obsession or compulsion acquired during gameplay
type Mania struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

func (m Mania) String() string {
	return m.Name
}

// Manias is the master list of all available manias from the Call of Cthulhu rulebook
var Manias = map[string]Mania{
	"Ablutomania": {
		Name:        "Ablutomania",
		Description: "Compulsion for washing oneself.",
	},
	"Aboulomania": {
		Name:        "Aboulomania",
		Description: "Pathological indecisiveness.",
	},
	"Achluomania": {
		Name:        "Achluomania",
		Description: "An excessive liking for darkness.",
	},
	"Acromania": {
		Name:        "Acromania",
		Description: "Compulsion for high places.",
	},
	"Agathomania": {
		Name:        "Agathomania",
		Description: "Pathological kindness.",
	},
	"Agromania": {
		Name:        "Agromania",
		Description: "Intense desire to be in open spaces.",
	},
	"Aichmomania": {
		Name:        "Aichmomania",
		Description: "Obsession with sharp or pointed objects.",
	},
	"Ailuromania": {
		Name:        "Ailuromania",
		Description: "Abnormal fondness for cats.",
	},
	"Algomania": {
		Name:        "Algomania",
		Description: "Obsession with pain.",
	},
	"Alliomania": {
		Name:        "Alliomania",
		Description: "Obsession with garlic.",
	},
	"Amaxomania": {
		Name:        "Amaxomania",
		Description: "Obsession with being in vehicles.",
	},
	"Amenomania": {
		Name:        "Amenomania",
		Description: "Irrational cheerfulness.",
	},
	"Anthomania": {
		Name:        "Anthomania",
		Description: "Obsession with flowers.",
	},
	"Arithmomania": {
		Name:        "Arithmomania",
		Description: "Obsessive preoccupation with numbers.",
	},
	"Asoticamania": {
		Name:        "Asoticamania",
		Description: "Impulsive or reckless spending.",
	},
	"Automania": {
		Name:        "Automania",
		Description: "An excessive liking for solitude.",
	},
	"Balletomania": {
		Name:        "Balletomania",
		Description: "Abnormal fondness for ballet.",
	},
	"Bibliokleptomania": {
		Name:        "Bibliokleptomania",
		Description: "Compulsion for stealing books.",
	},
	"Bibliomania": {
		Name:        "Bibliomania",
		Description: "Obsession with books and/or reading.",
	},
	"Bruxomania": {
		Name:        "Bruxomania",
		Description: "Compulsion for grinding teeth.",
	},
	"Cacodemomania": {
		Name:        "Cacodemomania",
		Description: "Pathological belief that one is inhabited by an evil spirit.",
	},
	"Callomania": {
		Name:        "Callomania",
		Description: "Obsession with one's own beauty.",
	},
	"Cartacoethes": {
		Name:        "Cartacoethes",
		Description: "Uncontrollable compulsion to see maps everywhere.",
	},
	"Catapedamania": {
		Name:        "Catapedamania",
		Description: "Obsession with jumping from high places.",
	},
	"Cheimatomania": {
		Name:        "Cheimatomania",
		Description: "Abnormal desire for cold and/or cold things.",
	},
	"Choreomania": {
		Name:        "Choreomania",
		Description: "Dancing mania or uncontrollable frenzy.",
	},
	"Clinomania": {
		Name:        "Clinomania",
		Description: "Excessive desire to stay in bed.",
	},
	"Coimetromania": {
		Name:        "Coimetromania",
		Description: "Obsession with cemeteries.",
	},
	"Coloromania": {
		Name:        "Coloromania",
		Description: "Obsession with a specific color.",
	},
	"Coulromania": {
		Name:        "Coulromania",
		Description: "Obsession with clowns.",
	},
	"Countermania": {
		Name:        "Countermania",
		Description: "Compulsion to experience fearful situations.",
	},
	"Dacnomania": {
		Name:        "Dacnomania",
		Description: "Obsession with killing.",
	},
	"Demonomania": {
		Name:        "Demonomania",
		Description: "Pathological belief that one is possessed by demons.",
	},
	"Dermatillomania": {
		Name:        "Dermatillomania",
		Description: "Compulsion for picking at one's skin.",
	},
	"Dikemania": {
		Name:        "Dikemania",
		Description: "Obsession to see justice done.",
	},
	"Dipsomania": {
		Name:        "Dipsomania",
		Description: "Abnormal craving for alcohol.",
	},
	"Doramania": {
		Name:        "Doramania",
		Description: "Obsession with owning furs.",
	},
	"Doromania": {
		Name:        "Doromania",
		Description: "Obsession with giving gifts.",
	},
	"Drapetomania": {
		Name:        "Drapetomania",
		Description: "Compulsion for running away.",
	},
	"Ecdemiomania": {
		Name:        "Ecdemiomania",
		Description: "Compulsion for wandering.",
	},
	"Egomania": {
		Name:        "Egomania",
		Description: "Irrational self-centered attitude or self-worship.",
	},
	"Empleomania": {
		Name:        "Empleomania",
		Description: "Insatiable urge to hold office.",
	},
	"Enosimania": {
		Name:        "Enosimania",
		Description: "Pathological belief that one has sinned.",
	},
	"Epistemomania": {
		Name:        "Epistemomania",
		Description: "Obsession for acquiring knowledge.",
	},
	"Eremiomania": {
		Name:        "Eremiomania",
		Description: "Compulsion for stillness.",
	},
	"Etheromania": {
		Name:        "Etheromania",
		Description: "Craving for ether.",
	},
	"Gamomania": {
		Name:        "Gamomania",
		Description: "Obsession with issuing odd marriage proposals.",
	},
	"Geliomania": {
		Name:        "Geliomania",
		Description: "Uncontrollable compulsion to laugh.",
	},
	"Goetomania": {
		Name:        "Goetomania",
		Description: "Obsession with witches and witchcraft.",
	},
	"Graphomania": {
		Name:        "Graphomania",
		Description: "Obsession with writing everything down.",
	},
	"Gymnomania": {
		Name:        "Gymnomania",
		Description: "Compulsion with nudity.",
	},
	"Habromania": {
		Name:        "Habromania",
		Description: "Abnormal tendency to create pleasant delusions (in spite of reality).",
	},
	"Helminthomania": {
		Name:        "Helminthomania",
		Description: "An excessive liking for worms.",
	},
	"Hoplomania": {
		Name:        "Hoplomania",
		Description: "Obsession with firearms.",
	},
	"Hydromania": {
		Name:        "Hydromania",
		Description: "Irrational craving for water.",
	},
	"Ichthyomania": {
		Name:        "Ichthyomania",
		Description: "Obsession with fish.",
	},
	"Iconomania": {
		Name:        "Iconomania",
		Description: "Obsession with icons or portraits.",
	},
	"Idolomania": {
		Name:        "Idolomania",
		Description: "Obsession or devotion to an idol.",
	},
	"Infomania": {
		Name:        "Infomania",
		Description: "Excessive devotion to accumulating facts.",
	},
	"Klazomania": {
		Name:        "Klazomania",
		Description: "Irrational compulsion to shout.",
	},
	"Kleptomania": {
		Name:        "Kleptomania",
		Description: "Irrational compulsion for stealing.",
	},
	"Ligyromania": {
		Name:        "Ligyromania",
		Description: "Uncontrollable compulsion to make loud or shrill noises.",
	},
	"Linonomania": {
		Name:        "Linonomania",
		Description: "Obsession with string.",
	},
	"Lotterymania": {
		Name:        "Lotterymania",
		Description: "An extreme desire to take part in lotteries.",
	},
	"Lypemania": {
		Name:        "Lypemania",
		Description: "An abnormal tendency toward deep melancholy.",
	},
	"Megalithomania": {
		Name:        "Megalithomania",
		Description: "Abnormal tendency to compose bizarre ideas when in the presence of stone circles/standing stones.",
	},
	"Melomania": {
		Name:        "Melomania",
		Description: "Obsession with music or a specific tune.",
	},
	"Metromania": {
		Name:        "Metromania",
		Description: "Insatiable desire for writing verse.",
	},
	"Misomania": {
		Name:        "Misomania",
		Description: "Hatred of everything, obsession of hating some subject or group.",
	},
	"Monomania": {
		Name:        "Monomania",
		Description: "Abnormal obsession with a single thought or idea.",
	},
	"Mythomania": {
		Name:        "Mythomania",
		Description: "Lying or exaggerating to an abnormal extent.",
	},
	"Nosomania": {
		Name:        "Nosomania",
		Description: "Delusion of suffering from an imagined disease.",
	},
	"Notomania": {
		Name:        "Notomania",
		Description: "Compulsion to record everything (e.g. photograph).",
	},
	"Onomamania": {
		Name:        "Onomamania",
		Description: "Obsession with names (people, places, things).",
	},
	"Onomatomania": {
		Name:        "Onomatomania",
		Description: "Irresistible desire to repeat certain words.",
	},
	"Onychotillomania": {
		Name:        "Onychotillomania",
		Description: "Compulsive picking at the fingernails.",
	},
	"Opsomania": {
		Name:        "Opsomania",
		Description: "Abnormal love for one kind of food.",
	},
	"Paramania": {
		Name:        "Paramania",
		Description: "An abnormal pleasure in complaining.",
	},
	"Personamania": {
		Name:        "Personamania",
		Description: "Compulsion to wear masks.",
	},
	"Phasmomania": {
		Name:        "Phasmomania",
		Description: "Obsession with ghosts.",
	},
	"Phonomania": {
		Name:        "Phonomania",
		Description: "Pathological tendency to murder.",
	},
	"Photomania": {
		Name:        "Photomania",
		Description: "Pathological desire for light.",
	},
	"Planomania": {
		Name:        "Planomania",
		Description: "Abnormal desire to disobey social norms.",
	},
	"Plutomania": {
		Name:        "Plutomania",
		Description: "Obsessive desire for wealth.",
	},
	"Pseudomania": {
		Name:        "Pseudomania",
		Description: "Irrational compulsion for lying.",
	},
	"Pyromania": {
		Name:        "Pyromania",
		Description: "Compulsion for starting fires.",
	},
	"Question-Asking Mania": {
		Name:        "Question-Asking Mania",
		Description: "Compulsive urge to ask questions.",
	},
	"Rhinotillexomania": {
		Name:        "Rhinotillexomania",
		Description: "Compulsive nose picking.",
	},
	"Scribbleomania": {
		Name:        "Scribbleomania",
		Description: "Obsession with scribbling/doodling.",
	},
	"Siderodromomania": {
		Name:        "Siderodromomania",
		Description: "Intense fascination with trains and railroad travel.",
	},
	"Sophomania": {
		Name:        "Sophomania",
		Description: "The delusion that one is incredibly intelligent.",
	},
	"Technomania": {
		Name:        "Technomania",
		Description: "Obsession with new technology.",
	},
	"Thanatomania": {
		Name:        "Thanatomania",
		Description: "Belief that one is cursed by death magic.",
	},
	"Theomania": {
		Name:        "Theomania",
		Description: "Belief that he or she is a god.",
	},
	"Titillomania": {
		Name:        "Titillomania",
		Description: "Compulsion for scratching oneself.",
	},
	"Tomomania": {
		Name:        "Tomomania",
		Description: "Irrational predilection for performing surgery.",
	},
	"Trichotillomania": {
		Name:        "Trichotillomania",
		Description: "Craving for pulling out own hair.",
	},
	"Typhlomania": {
		Name:        "Typhlomania",
		Description: "Pathological blindness.",
	},
	"Xenomania": {
		Name:        "Xenomania",
		Description: "Obsession with foreign things.",
	},
	"Zoomania": {
		Name:        "Zoomania",
		Description: "Insane fondness for animals.",
	},
}

// ManiasList provides a slice of all mania names for easy iteration
var ManiasList = func() []string {
	keys := make([]string, 0, len(Manias))
	for k := range Manias {
		keys = append(keys, k)
	}
	return keys
}()

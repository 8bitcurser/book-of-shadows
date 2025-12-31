package models

// Phobia represents an irrational fear acquired during gameplay
type Phobia struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

func (p Phobia) String() string {
	return p.Name
}

// Phobias is the master list of all available phobias from the Call of Cthulhu rulebook
var Phobias = map[string]Phobia{
	"Ablutophobia": {
		Name:        "Ablutophobia",
		Description: "Fear of washing or bathing.",
	},
	"Acrophobia": {
		Name:        "Acrophobia",
		Description: "Fear of heights.",
	},
	"Aerophobia": {
		Name:        "Aerophobia",
		Description: "Fear of flying.",
	},
	"Agoraphobia": {
		Name:        "Agoraphobia",
		Description: "Fear of open, public (crowded) places.",
	},
	"Alektorophobia": {
		Name:        "Alektorophobia",
		Description: "Fear of chickens.",
	},
	"Alliumphobia": {
		Name:        "Alliumphobia",
		Description: "Fear of garlic.",
	},
	"Amaxophobia": {
		Name:        "Amaxophobia",
		Description: "Fear of being in or riding in vehicles.",
	},
	"Ancraophobia": {
		Name:        "Ancraophobia",
		Description: "Fear of wind.",
	},
	"Androphobia": {
		Name:        "Androphobia",
		Description: "Fear of men.",
	},
	"Anglophobia": {
		Name:        "Anglophobia",
		Description: "Fear of England or English culture, etc.",
	},
	"Anthrophobia": {
		Name:        "Anthrophobia",
		Description: "Fear of flowers.",
	},
	"Apotemnophobia": {
		Name:        "Apotemnophobia",
		Description: "Fear of people with amputations.",
	},
	"Arachnophobia": {
		Name:        "Arachnophobia",
		Description: "Fear of spiders.",
	},
	"Astraphobia": {
		Name:        "Astraphobia",
		Description: "Fear of lightning.",
	},
	"Atephobia": {
		Name:        "Atephobia",
		Description: "Fear of ruin or ruins.",
	},
	"Aulophobia": {
		Name:        "Aulophobia",
		Description: "Fear of flutes.",
	},
	"Bacteriophobia": {
		Name:        "Bacteriophobia",
		Description: "Fear of bacteria.",
	},
	"Ballistophobia": {
		Name:        "Ballistophobia",
		Description: "Fear of missiles or bullets.",
	},
	"Basophobia": {
		Name:        "Basophobia",
		Description: "Fear of falling.",
	},
	"Bibliophobia": {
		Name:        "Bibliophobia",
		Description: "Fear of books.",
	},
	"Botanophobia": {
		Name:        "Botanophobia",
		Description: "Fear of plants.",
	},
	"Caligynephobia": {
		Name:        "Caligynephobia",
		Description: "Fear of beautiful women.",
	},
	"Cheimaphobia": {
		Name:        "Cheimaphobia",
		Description: "Fear of cold.",
	},
	"Chronomentrophobia": {
		Name:        "Chronomentrophobia",
		Description: "Fear of clocks.",
	},
	"Claustrophobia": {
		Name:        "Claustrophobia",
		Description: "Fear of confined spaces.",
	},
	"Coulrophobia": {
		Name:        "Coulrophobia",
		Description: "Fear of clowns.",
	},
	"Cynophobia": {
		Name:        "Cynophobia",
		Description: "Fear of dogs.",
	},
	"Demonophobia": {
		Name:        "Demonophobia",
		Description: "Fear of spirits or demons.",
	},
	"Demophobia": {
		Name:        "Demophobia",
		Description: "Fear of crowds.",
	},
	"Dentophobia": {
		Name:        "Dentophobia",
		Description: "Fear of dentists.",
	},
	"Disposophobia": {
		Name:        "Disposophobia",
		Description: "Fear of throwing stuff out (hoarding).",
	},
	"Doraphobia": {
		Name:        "Doraphobia",
		Description: "Fear of fur.",
	},
	"Dromophobia": {
		Name:        "Dromophobia",
		Description: "Fear of crossing streets.",
	},
	"Ecclesiophobia": {
		Name:        "Ecclesiophobia",
		Description: "Fear of church.",
	},
	"Eisoptrophobia": {
		Name:        "Eisoptrophobia",
		Description: "Fear of mirrors.",
	},
	"Enetophobia": {
		Name:        "Enetophobia",
		Description: "Fear of needles or pins.",
	},
	"Entomophobia": {
		Name:        "Entomophobia",
		Description: "Fear of insects.",
	},
	"Felinophobia": {
		Name:        "Felinophobia",
		Description: "Fear of cats.",
	},
	"Gephyrophobia": {
		Name:        "Gephyrophobia",
		Description: "Fear of crossing bridges.",
	},
	"Gerontophobia": {
		Name:        "Gerontophobia",
		Description: "Fear of old people or of growing old.",
	},
	"Gynophobia": {
		Name:        "Gynophobia",
		Description: "Fear of women.",
	},
	"Haemaphobia": {
		Name:        "Haemaphobia",
		Description: "Fear of blood.",
	},
	"Hamartophobia": {
		Name:        "Hamartophobia",
		Description: "Fear of sinning.",
	},
	"Haphophobia": {
		Name:        "Haphophobia",
		Description: "Fear of touch.",
	},
	"Herpetophobia": {
		Name:        "Herpetophobia",
		Description: "Fear of reptiles.",
	},
	"Homichlophobia": {
		Name:        "Homichlophobia",
		Description: "Fear of fog.",
	},
	"Hoplophobia": {
		Name:        "Hoplophobia",
		Description: "Fear of firearms.",
	},
	"Hydrophobia": {
		Name:        "Hydrophobia",
		Description: "Fear of water.",
	},
	"Hypnophobia": {
		Name:        "Hypnophobia",
		Description: "Fear of sleep or of being hypnotized.",
	},
	"Iatrophobia": {
		Name:        "Iatrophobia",
		Description: "Fear of doctors.",
	},
	"Ichthyophobia": {
		Name:        "Ichthyophobia",
		Description: "Fear of fish.",
	},
	"Katsaridaphobia": {
		Name:        "Katsaridaphobia",
		Description: "Fear of cockroaches.",
	},
	"Keraunophobia": {
		Name:        "Keraunophobia",
		Description: "Fear of thunder.",
	},
	"Lachanophobia": {
		Name:        "Lachanophobia",
		Description: "Fear of vegetables.",
	},
	"Ligyrophobia": {
		Name:        "Ligyrophobia",
		Description: "Fear of loud noises.",
	},
	"Limnophobia": {
		Name:        "Limnophobia",
		Description: "Fear of lakes.",
	},
	"Mechanophobia": {
		Name:        "Mechanophobia",
		Description: "Fear of machines or machinery.",
	},
	"Megalophobia": {
		Name:        "Megalophobia",
		Description: "Fear of large things.",
	},
	"Merinthophobia": {
		Name:        "Merinthophobia",
		Description: "Fear of being bound or tied up.",
	},
	"Meteorophobia": {
		Name:        "Meteorophobia",
		Description: "Fear of meteors or meteorites.",
	},
	"Monophobia": {
		Name:        "Monophobia",
		Description: "Fear of being alone.",
	},
	"Mysophobia": {
		Name:        "Mysophobia",
		Description: "Fear of dirt or contamination.",
	},
	"Myxophobia": {
		Name:        "Myxophobia",
		Description: "Fear of slime.",
	},
	"Necrophobia": {
		Name:        "Necrophobia",
		Description: "Fear of dead things.",
	},
	"Octophobia": {
		Name:        "Octophobia",
		Description: "Fear of the figure 8.",
	},
	"Odontophobia": {
		Name:        "Odontophobia",
		Description: "Fear of teeth.",
	},
	"Oneirophobia": {
		Name:        "Oneirophobia",
		Description: "Fear of dreams.",
	},
	"Onomatophobia": {
		Name:        "Onomatophobia",
		Description: "Fear of hearing a certain word or words.",
	},
	"Ophidiophobia": {
		Name:        "Ophidiophobia",
		Description: "Fear of snakes.",
	},
	"Ornithophobia": {
		Name:        "Ornithophobia",
		Description: "Fear of birds.",
	},
	"Parasitophobia": {
		Name:        "Parasitophobia",
		Description: "Fear of parasites.",
	},
	"Pediophobia": {
		Name:        "Pediophobia",
		Description: "Fear of dolls.",
	},
	"Phagophobia": {
		Name:        "Phagophobia",
		Description: "Fear of swallowing, of eating or of being eaten.",
	},
	"Pharmacophobia": {
		Name:        "Pharmacophobia",
		Description: "Fear of drugs.",
	},
	"Phasmophobia": {
		Name:        "Phasmophobia",
		Description: "Fear of ghosts.",
	},
	"Phenogophobia": {
		Name:        "Phenogophobia",
		Description: "Fear of daylight.",
	},
	"Pogonophobia": {
		Name:        "Pogonophobia",
		Description: "Fear of beards.",
	},
	"Potamophobia": {
		Name:        "Potamophobia",
		Description: "Fear of rivers.",
	},
	"Potophobia": {
		Name:        "Potophobia",
		Description: "Fear of alcohol or alcoholic beverages.",
	},
	"Pyrophobia": {
		Name:        "Pyrophobia",
		Description: "Fear of fire.",
	},
	"Rhabdophobia": {
		Name:        "Rhabdophobia",
		Description: "Fear of magic.",
	},
	"Scotophobia": {
		Name:        "Scotophobia",
		Description: "Fear of darkness or of the night.",
	},
	"Selenophobia": {
		Name:        "Selenophobia",
		Description: "Fear of the moon.",
	},
	"Siderodromophobia": {
		Name:        "Siderodromophobia",
		Description: "Fear of train travel.",
	},
	"Siderophobia": {
		Name:        "Siderophobia",
		Description: "Fear of stars.",
	},
	"Stenophobia": {
		Name:        "Stenophobia",
		Description: "Fear of narrow things or places.",
	},
	"Symmetrophobia": {
		Name:        "Symmetrophobia",
		Description: "Fear of symmetry.",
	},
	"Taphephobia": {
		Name:        "Taphephobia",
		Description: "Fear of being buried alive or of cemeteries.",
	},
	"Taurophobia": {
		Name:        "Taurophobia",
		Description: "Fear of bulls.",
	},
	"Telephonophobia": {
		Name:        "Telephonophobia",
		Description: "Fear of telephones.",
	},
	"Teratophobia": {
		Name:        "Teratophobia",
		Description: "Fear of monsters.",
	},
	"Thalassophobia": {
		Name:        "Thalassophobia",
		Description: "Fear of the sea.",
	},
	"Tomophobia": {
		Name:        "Tomophobia",
		Description: "Fear of surgical operations.",
	},
	"Triskadekaphobia": {
		Name:        "Triskadekaphobia",
		Description: "Fear of the number 13.",
	},
	"Vestiphobia": {
		Name:        "Vestiphobia",
		Description: "Fear of clothing.",
	},
	"Wiccaphobia": {
		Name:        "Wiccaphobia",
		Description: "Fear of witches and witchcraft.",
	},
	"Xanthophobia": {
		Name:        "Xanthophobia",
		Description: "Fear of the color yellow or the word \"yellow\".",
	},
	"Xenoglossophobia": {
		Name:        "Xenoglossophobia",
		Description: "Fear of foreign languages.",
	},
	"Xenophobia": {
		Name:        "Xenophobia",
		Description: "Fear of strangers or foreigners.",
	},
	"Zoophobia": {
		Name:        "Zoophobia",
		Description: "Fear of animals.",
	},
}

// PhobiasList provides a slice of all phobia names for easy iteration
var PhobiasList = func() []string {
	keys := make([]string, 0, len(Phobias))
	for k := range Phobias {
		keys = append(keys, k)
	}
	return keys
}()

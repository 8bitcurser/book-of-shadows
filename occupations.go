package main

type BaseSkillAttribute struct {
	Name       string
	Multiplier int
}

type SkillPointFormula struct {
	BaseAttributes []BaseSkillAttribute // For handling multiple base attributes
	Options        []BaseSkillAttribute // Optional OR cases
}

type Occupation struct {
	Name              string
	Skills            []string
	SuggestedContacts string
	SkillPoints       SkillPointFormula
}

package wtf

type Acronym struct {
	Acronym  string
	Meaning  string
	Location string
	Note     string
}

type Acronyms []Acronym

func (a Acronym) String() string {
	return a.Meaning
}

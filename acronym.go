package wtf

import (
	"fmt"
)

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

func (a Acronyms) String() string {
	ret := ""
	for _, el := range a {
		ret = fmt.Sprintf("%s - %s\n", ret, el.String())
	}
	return ret
}

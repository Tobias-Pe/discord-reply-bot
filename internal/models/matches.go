package models

type MessageMatch struct {
	Message      string
	IsExactMatch bool
}

// AllMatchChoices check usage before changing order
var AllMatchChoices = []string{"exact", "occurrence"}

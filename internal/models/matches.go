package models

type MessageMatch struct {
	Message      string `json:"to-be-matched"`
	IsExactMatch bool   `json:"match-type"`
}

type CompleteData struct {
	MessageMatch MessageMatch `json:"message-match"`
	Reply        []string     `json:"replies"`
}

// AllMatchChoices check usage before changing order
var AllMatchChoices = []string{"exact", "occurrence"}

package models

import "encoding/json"

type Code struct {
	Lang string
	Body string
}

func (c Code) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

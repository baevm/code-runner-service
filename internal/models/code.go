package models

import "encoding/json"

type Code struct {
	Lang      string `json:"lang"`
	Body      string `json:"body"`
	RequestId string `json:"request_id"`
}

func (c Code) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}


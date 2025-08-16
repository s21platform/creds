package model

import "encoding/json"

type CredsData struct {
	Data json.RawMessage `db:"data"`
}

type Creds struct {
	Creds   []Cred `json:"creds"`
}

type Cred struct {
	Name  string          `json:"name"`
	Value string          `json:"value"`
}


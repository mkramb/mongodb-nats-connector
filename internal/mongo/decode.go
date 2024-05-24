package mongo

import (
	"encoding/json"
)

type ChangeEventDocumentId struct {
	Value string `json:"$oid"`
}

type ChangeEventDocument struct {
	Id ChangeEventDocumentId `json:"_id"`
}

type ChangeEventResumeToken struct {
	Value string `json:"_data"`
}

type ChangeEventNs struct {
	Coll string `json:"coll"`
}

type ChangeEvent struct {
	OperationType string                 `json:"operationType"`
	FullDocument  ChangeEventDocument    `json:"fullDocument"`
	ResumeToken   ChangeEventResumeToken `json:"_id"`
	Ns            ChangeEventNs          `json:"ns"`
}

func DecodeChangeEvent(changeEvent []byte) (*ChangeEvent, error) {
	var event ChangeEvent

	err := json.Unmarshal(changeEvent, &event)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

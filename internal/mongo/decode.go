package mongo

import (
	"encoding/json"
)

type changeEventDocumentId struct {
	Value string `json:"$oid"`
}

type changeEventDocument struct {
	Id changeEventDocumentId `json:"_id"`
}

type changeEventResumeToken struct {
	Value string `json:"_data"`
}

type changeEventNs struct {
	Coll string `json:"coll"`
}

type ChangeEvent struct {
	OperationType string                 `json:"operationType"`
	FullDocument  changeEventDocument    `json:"fullDocument"`
	ResumeToken   changeEventResumeToken `json:"_id"`
	Ns            changeEventNs          `json:"ns"`
}

func DecodeChangeEvent(changeEvent []byte) (*ChangeEvent, error) {
	var event ChangeEvent

	err := json.Unmarshal(changeEvent, &event)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

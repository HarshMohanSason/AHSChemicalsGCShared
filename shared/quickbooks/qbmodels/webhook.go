package qbmodels

import "encoding/json"

type QBWebHookResp struct {
	EventNotifications []EventNotification `json:"eventNotifications"`
}

type EventNotification struct {
	RealmID         string          `json:"realmId"`
	DataChangeEvent DataChangeEvent `json:"dataChangeEvent"`
}

type DataChangeEvent struct {
	Entities []Entity `json:"entities"`
}

type Entity struct {
	Name           string `json:"name"`
	ID             string `json:"id"`
	Operation      string `json:"operation"`
	LastUpdatedUTC string `json:"lastUpdated"`
}

func (e *Entity) ToBytes() ([]byte, error) {
	return json.Marshal(e)
}

// Represents the data payload pushed by quickbooks to the pubsub
// The subscriber is made to only accept the data field. Gen2 functions 
// still send this data field instead of sending the raw json bytes when 
// pushed as a message to the topic so we need this struct to get the data 
// then have a helper function to convert it to an Entity struct
type QBWebHookPubMessage struct {
	Data []byte `json:"data"`
}

func (qbwkm *QBWebHookPubMessage) ToEntity() (*Entity, error) {
	var entity Entity
	err := json.Unmarshal(qbwkm.Data, &entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

package qbmodels

import "encoding/json"

type QBWebHookResp struct {
	EventNotifications []EventNotification `json:"eventNotifications"`
}

type EventNotification struct {
	RealmID         string         `json:"realmId"`
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

func (e *Entity) ToBytes() ([]byte, error){
	return json.Marshal(e)
}
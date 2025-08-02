package qbmodels

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
	Name           string `json:"name"`            // e.g., "Customer", "Item"
	ID             string `json:"id"`              // ID of the updated entity
	Operation      string `json:"operation"`       // e.g., "Update", "Create", "Delete"
	LastUpdatedUTC string `json:"lastUpdated"`     // e.g., "2025-07-30T10:45:12Z"
}
package shared

import "sync"

var (
	InitFirebaseOnce   sync.Once
	InitQuickBooksOnce sync.Once
	InitThirdPartyOnce sync.Once
	InitGCPOnce        sync.Once
)

// Reference represents a reference to another QuickBooks entity.
//
// Fields:
//   - Value: The unique ID of the referenced entity (required).
//   - Name:  Optional human-readable name for the referenced entity.
type Reference struct {
	Value string `json:"value"`            // Required ID of the referenced entity.
	Name  string `json:"name,omitempty"`   // Optional display name of the referenced entity.
}

// MetaData contains creation and update timestamps for QuickBooks entities.
//
// Fields:
//   - CreateTime:      ISO 8601 timestamp of when the entity was created.
//   - LastUpdatedTime: ISO 8601 timestamp of the most recent update to the entity.
type MetaData struct {
	CreateTime      string `json:"CreateTime,omitempty"`       // Timestamp of entity creation.
	LastUpdatedTime string `json:"LastUpdatedTime,omitempty"`  // Timestamp of last entity update.
}
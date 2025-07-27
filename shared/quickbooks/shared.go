package quickbooks

import "time"

type MetaData struct {
	CreateTime      time.Time `json:"CreateTime"`
	LastUpdatedTime time.Time `json:"LastUpdatedTime"`
}
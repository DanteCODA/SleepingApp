package entities

// Checkpoint struct
type Checkpoint struct {
	PageSize  int64 `json:"size,omitempty"`
	PageIndex int64 `json:"index,omitem
package model

// Status represents the status of an order.
type Status string

// Possible order statuses.
const (
	StatusNew        Status = "NEW"
	StatusRegistered Status = "REGISTERED"
	StatusProcessing Status = "PROCESSING"
	StatusInvalid    Status = "INVALID"
	StatusProcessed  Status = "PROCESSED"
)

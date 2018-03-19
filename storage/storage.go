package storage

import (
	"errors"
	"time"
)

var (
	// ErrInvalidPublicKey is returned when a public key is not in the correct format
	ErrInvalidPublicKey = errors.New("invalid public key")
)

// Partner defines the partner schema in the partners collection
type Partner struct {
	ID       string    `json:"_id"`
	Name     string    `json:"name"`
	RevShare int       `json:"revShareTotalPercentage"`
	Created  time.Time `json:"created"`
}

// DB is the contract that all databases will need to adhere to
type DB interface {
}

// PartnerC is the interface defining methods needed to interact with the partner collection
type PartnerC interface {
	GetPartner(ID string) (*Partner, error)
}

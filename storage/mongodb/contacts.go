package mongodb

import (
	"time"
)

// Contact defines the user schema
type Contact struct {
	ID          string    `json:"_id"`
	LastSeen    time.Time `json:"lastSeen"`
	Port        int       `json:"port"`
	Address     string    `json:"address"`
	UserAgent   string    `json:"userAgent"`
	Protocol    string    `json:"protocol"`
	LastTimeout time.Time `json:"lastTimeout"`
	TimeoutRate int       `json:"timeoutRate"`
}

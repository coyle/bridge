package mongodb

import "time"

// Frame defines the frame schema in the frames collection
type Frame struct {
	ID      string    `json:"_id"`
	User    string    `json:"user"`
	Shards  []string  `json:"shards"`
	Size    int       `json:"size"`
	Locked  bool      `json:"locked"`
	Created time.Time `json:"created"`
}

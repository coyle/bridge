package mongodb

// Bucket defines the bucket schema in the buckets collection
type Bucket struct {
	ID       string   `json:"_id"`
	User     string   `json:"user"`
	Created  string   `json:"created"`
	Name     string   `json:"name"`
	Pubkeys  []string `json:"pubkeys"`
	Status   string   `json:"status"`
	Transfer int      `json:"transfer"`
	Storage  int      `json:"storage"`
}

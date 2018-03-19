package mongodb

import (
	"github.com/coyle/bridge/storage"
	"github.com/globalsign/mgo/bson"
)

// GetPartner searches for a partner with the provided name from the partners collection
func (c *Client) GetPartner(name string) (*storage.Partner, error) {
	p := &storage.Partner{}
	err := c.partners.Find(bson.M{"name": name}).One(p)

	return p, err

}

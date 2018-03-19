package mongodb

import (
	"encoding/hex"
	"fmt"

	"github.com/coyle/bridge/storage"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	secp256k1 "github.com/haltingstate/secp256k1-go"
)

// PublicKey defines the PublicKey schema in the PublicKeys collection
type PublicKey struct {
	ID    string `bson:"_id" json:"_id"`
	User  string `json:"user"`
	Label string `json:"label"`
}

// CreatePublicKey instantiates a new PublicKey for a user
func (c *Client) CreatePublicKey(u *User, pubKey string) error {
	fmt.Printf("in CreatePublicKey\n")
	pk, err := hex.DecodeString(pubKey)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if valid := secp256k1.VerifyPubkey(pk); valid != 1 {
		fmt.Println(storage.ErrInvalidPublicKey)
		return storage.ErrInvalidPublicKey
	}

	// if key exists return error
	exists, err := c.PublicKeyExists(pubKey)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if exists {
		fmt.Printf("exists: %t\n", exists)
		return nil
	}

	// else save new pubkey
	pubk := &PublicKey{
		ID:   pubKey,
		User: u.ID,
	}
	fmt.Printf("struct: %#v\n", pubk)
	return c.publicKeys.Insert(pubk)
}

// GetPublickey looks up a PublicKey document with the provided key
func (c *Client) GetPublickey(key string) (*PublicKey, error) {
	pk := &PublicKey{}
	err := c.publicKeys.Find(bson.M{"_id": key}).One(pk)

	return pk, err
}

// PublicKeyExists determines if the provided key exists in the publickeys collection
func (c *Client) PublicKeyExists(key string) (bool, error) {
	c.session.SetSafe(&mgo.Safe{})
	cnt, err := c.publicKeys.FindId(key).Count()
	if cnt > 0 {
		return true, err
	}

	return false, err
}

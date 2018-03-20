package mongodb

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
)

var (
	// ErrInvalidID is returned when the user ID is not in compliance with RFC 5322
	ErrInvalidID = errors.New("invalid id format")
)

// User defines the user schema
type User struct {
	ID                string      `bson:"_id" json:"id,omitempty"`
	UUID              string      `json:"uuid,omitempty"`
	Hashpass          string      `json:"hashpass,omitempty"`
	Activated         bool        `json:"activated"`
	IsFreeTier        bool        `json:"isFreeTier"`
	Activator         string      `json:"activator,omitempty"`
	Deactivator       string      `json:"deactivator,omitempty"`
	Created           time.Time   `json:"created"`
	BytesUploaded     BytesMeta   `json:"bytesUploaded,omitempty"`
	BytesDownloaded   BytesMeta   `json:"bytesDownloaded,omitempty"`
	PaymentProcessors []string    `json:"paymentProcessors,omitempty"`
	ReferralPartner   string      `json:"referralPartner,omitempty"`
	Preferences       Preferences `json:"preferences,omitempty"`
	Resetter          string      `json:"resetter,omitempty"`
}

// Preferences contains all user preferences
type Preferences struct {
	DNT bool `json:"dnt"`
}

// BytesMeta contains metadata about the data uploaded/downloaded
type BytesMeta struct {
	LastDayBytes     int64     `json:"lastDayBytes"`
	LastDayStarted   time.Time `json:"lastDayStarted"`
	LastHourBytes    int64     `json:"lastHourBytes"`
	LastHourStarted  time.Time `json:"lastHourStarted"`
	LastMonthBytes   int64     `json:"lastMonthBytes"`
	LastMonthStarted time.Time `json:"lastMonthStarted"`
}

// CreateUser initalizes and saves a new user in the Users collection
func (c *Client) CreateUser(u User) (User, error) {
	zeroTime := time.Time{}
	if u.Created == zeroTime {
		u.Created = time.Now().UTC()
	}

	if u.UUID == "" {
		u.UUID = uuid.New().String()
	}

	if _, err := mail.ParseAddress(u.ID); err != nil {
		return User{}, ErrInvalidID
	}

	err := c.users.Insert(&u)

	return u, err
}

// GetUser queries for a user by their ID
func (c *Client) GetUser(id string) (*User, error) {
	u := &User{}
	err := c.users.Find(bson.M{"_id": id}).One(u)
	return u, err
}

// GetUserByToken queries for a user by their activator token
func (c *Client) GetUserByToken(name, token string) (*User, error) {
	u := &User{}
	err := c.users.Find(bson.M{name: token}).One(u)
	return u, err
}

// ActivateUser flips the activate flag on the user model with the provided ID
func (c *Client) ActivateUser(id string) error {
	return c.users.UpdateId(id, bson.M{"$set": bson.M{"activate": true, "activator": nil}})
}

// DeactivateUser sets the deactivator to a randomly generated hex string
func (c *Client) DeactivateUser(id string) error {
	b := make([]byte, 256)
	rand.Read(b)

	return c.users.UpdateId(id, bson.M{"$set": bson.M{"deactivator": hex.EncodeToString(b)}})
}

// ConfirmUserDeactivation sets the deactivator to a randomly generated hex string
func (c *Client) ConfirmUserDeactivation(id string) error {
	b := make([]byte, 256)
	rand.Read(b)
	return c.users.UpdateId(id, bson.M{"$set": bson.M{"deactivate": true, "activated": false, "activator": hex.EncodeToString(b)}})
}

// CreatePasswordResetToken generates a random hex string and saves to the user document
func (c *Client) CreatePasswordResetToken(id string) (string, error) {
	b := make([]byte, 256)
	rand.Read(b)
	err := c.users.UpdateId(id, bson.M{"$set": bson.M{"resetter": hex.EncodeToString(b)}})
	return string(b), err
}

// ResetPassword hashes the users new password and updates the document
func (c *Client) ResetPassword(id, p string) error {
	h := sha256.New()
	h.Write([]byte(p))

	return c.users.UpdateId(id, bson.M{"$set": bson.M{"resetter": "", "hashpass": fmt.Sprintf("%x", h.Sum(nil))}})
}

// UserToView returns a user object with private fields hidden
func UserToView(u *User) *User {
	return &User{
		UUID:              u.UUID,
		Activated:         u.Activated,
		IsFreeTier:        u.IsFreeTier,
		Created:           u.Created,
		PaymentProcessors: u.PaymentProcessors,
		ReferralPartner:   u.ReferralPartner,
		Preferences:       u.Preferences,
	}
}

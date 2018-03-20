package mongodb

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TestUser initializes a user struct with required values for testing purposes
func TestUser(activated bool) *User {
	uid := uuid.New().String()
	u := User{
		ID:              fmt.Sprintf("%s@storj.io", uid),
		UUID:            uid,
		Hashpass:        fmt.Sprintf("%x", sha256.Sum256([]byte("password"))),
		Activated:       activated,
		Created:         time.Now().UTC(),
		ReferralPartner: "CITIZEN",
	}

	if activated {
		u.Activator = uid
	} else {
		u.Deactivator = uid
	}

	return &u
}

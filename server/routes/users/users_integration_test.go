package users

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/coyle/bridge/storage/mongodb"
	secp256k1 "github.com/haltingstate/secp256k1-go"
	"github.com/stretchr/testify/assert"
)

func init() {
	waitUntilReady("http://bridge-server:8080/health")
}

var password = sha256.Sum256([]byte("password"))
var pubKey, _ = secp256k1.GenerateKeyPair()
var pubKeyString = hex.EncodeToString(pubKey)

func TestCreate(t *testing.T) {
	storageClient, err := mongodb.NewClient(os.Getenv("MONGO"))

	if err != nil {
		t.Errorf("unable to connect to Mongo: %s", err)
	}

	cases := []struct {
		name                 string
		body                 []byte
		id                   string
		expectedResponseCode int
		expectedError        bool
		expectedUser         mongodb.User
	}{
		{
			"valid user createtion request",
			[]byte(fmt.Sprintf(`{"email":"test@storj.io", "password":"%x", "pubkey": "%s"}`, password, pubKeyString)),
			"test@storj.io",
			201,
			false,
			mongodb.User{ID: "test@storj.io", Hashpass: fmt.Sprintf("%x", password)},
		},
		{
			"invalid email",
			[]byte(fmt.Sprintf(`{"email":"test+storj.io", "password":"%x", "pubkey": "%s"}`, password, pubKeyString)),
			"test+storj.io",
			400,
			true,
			mongodb.User{},
		},
	}
	for _, c := range cases {
		req, _ := http.NewRequest("POST", "http://bridge-server:8080/users", bytes.NewBuffer(c.body))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != c.expectedResponseCode {
			assert.Equal(t, c.expectedResponseCode, resp.StatusCode)
		}

		if c.expectedError {
			continue
		}

		actualUser, err := storageClient.GetUser(c.id)
		assert.NoError(t, err)
		assert.Equal(t, c.id, actualUser.ID)
		assert.Equal(t, fmt.Sprintf("%x", password), actualUser.Hashpass)
		assert.WithinDuration(t, time.Now(), actualUser.Created, 1*time.Minute)
	}
}

func TestReactivate(t *testing.T) {
	storageClient, err := mongodb.NewClient(os.Getenv("MONGO"))
	assert.NoError(t, err)

	testUser := mongodb.TestUser(false)
	_, err = storageClient.CreateUser(*testUser)
	assert.NoError(t, err)

	cases := []struct {
		name                 string
		body                 []byte
		id                   string
		expectedResponseCode int
		expectedError        bool
		expectedUser         mongodb.User
	}{
		{
			name:                 "valid user reactivation",
			body:                 []byte(fmt.Sprintf(`{"email":"%s"}`, testUser.ID)),
			expectedResponseCode: http.StatusCreated,
		},
	}

	for _, c := range cases {
		req, _ := http.NewRequest("POST", "http://bridge-server:8080/activations", bytes.NewBuffer(c.body))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != c.expectedResponseCode {
			assert.Equal(t, c.expectedResponseCode, resp.StatusCode)
		}

		if c.expectedError {
			continue
		}

		u := mongodb.User{}
		json.NewDecoder(resp.Body).Decode(&u)

		// assert the user's private fields are omitted
		assert.Empty(t, u.ID)
		assert.Empty(t, u.Hashpass)
		assert.Empty(t, u.Activator)
		assert.Empty(t, u.Deactivator)
		assert.Empty(t, u.Resetter)
		assert.Empty(t, u.BytesDownloaded)
		assert.Empty(t, u.BytesUploaded)

		// assert expected fields are there
		assert.Equal(t, testUser.UUID, u.UUID)
		assert.Equal(t, testUser.IsFreeTier, u.IsFreeTier)
		assert.Equal(t, testUser.Activated, u.Activated)
		assert.WithinDuration(t, testUser.Created, u.Created, 1*time.Second)
	}
}

func TestConfirmActivation(t *testing.T) {
	assert.NotNil(t, nil)
}

func TestDeactivation(t *testing.T) {
	assert.NotNil(t, nil)
}

func TestPasswordResetToken(t *testing.T) {
	assert.NotNil(t, nil)
}

func TestConfirmPasswordReset(t *testing.T) {
	assert.NotNil(t, nil)
}

func TestDispatchActivationEmailSwitch(t *testing.T) {
	assert.NotNil(t, nil)
}

func waitUntilReady(host string) {
	attempts := 0
	for {
		_, err := http.Get(host)
		if err != nil && attempts < 10 {
			attempts++
			time.Sleep(1 * time.Second)
		} else {
			return
		}

	}
}

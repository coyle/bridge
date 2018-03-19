package users

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/coyle/bridge/storage/mongodb"
	"github.com/go-kit/kit/log"
)

// Request contains all fields that will be used in a users request body
type Request struct {
	ReferralPartner string `json:"referralPartner"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PublicKey       string `json:"pubkey"`
	Token           string `json:"token"`
}

// User contains all configuration and methods to process user requests
type User struct {
	db     *mongodb.Client
	logger log.Logger
}

// NewServer returns a new instance of a configured User Server
func NewServer(client *mongodb.Client, logger log.Logger) *User {
	// start
	return &User{
		db:     client,
		logger: logger,
	}
}

// Create a new user
func (u *User) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("@ beginning")
	body, err := getBody(r)
	if err != nil {
		u.logger.Log("Error getting body", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("@ body\n%s\n", err)
		return
	}

	if body.PublicKey == "" {
		fmt.Println("@ no pubkey")
		u.logger.Log("No PublicKey provided", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	nuser := mongodb.User{
		ID:       body.Email,
		Hashpass: body.Password,
	}

	if body.ReferralPartner != "" {
		fmt.Println("@ no rp")
		p, err := u.db.GetPartner(body.ReferralPartner)
		if err != nil {
			u.logger.Log("Error getting partner", err)
		}

		nuser.ReferralPartner = p.ID
	}

	// do all concurrently ?
	user, err := u.db.CreateUser(nuser)
	if err != nil && err == mongodb.ErrInvalidID {
		fmt.Println("@ invalid Id")
		u.logger.Log("Error creating user", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		u.logger.Log("Error creating user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO(coyle): implement Email service
	defer u.dispatchActivationEmailSwitch(user)

	if body.PublicKey == "" {
		fmt.Println("no pub key")
		u.logger.Log("No PublicKey provided", err)
		// TODO(coyle): properly handle
		return
	}

	fmt.Println("creating public key")
	if err := u.db.CreatePublicKey(&user, body.PublicKey); err != nil {
		fmt.Printf("errors: %s\n", err)
		u.logger.Log("Error creating public key", err)
		// TODO(coyle): Should we cancel the request and remove the created user or just log?
		// looks like the node code currently removes the created user
	}

	w.WriteHeader(http.StatusCreated)
	return
}

// Reactivate a user
func (u *User) Reactivate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := getBody(r)
	if err != nil {
		// TODO(coyle): handle error
		return
	}

	user, err := u.db.GetUser(body.Email)
	if err != nil {
		// TODO(coyle): handle error
		return
	}

	if user.Activated {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO(coyle): implement Email service
	go u.dispatchActivationEmailSwitch(*user)

	w.WriteHeader(http.StatusCreated)
}

// ConfirmActivation of a user
func (u *User) ConfirmActivation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, err := u.db.GetUserByToken("activator", ps.ByName("token"))
	if err != nil {
		// TODO(coyle): handle error
	}

	if err := u.db.ActivateUser(user.ID); err != nil {
		// TODO(coyle): handle error
	}

}

// Remove a user
func (u *User) Remove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, err := u.db.GetUser(ps.ByName("id"))
	if err != nil {
		// TODO(coyle): handle error
	}

	// TODO(coyle): determine how to do check req.user._id !== user._id from node code
	// what is being passed in as req.user

	// TODO(coyle): dispatch email after mail service is completed

	if err := u.db.DeactivateUser(user.ID); err != nil {
		// TODO(coyle): handle error
	}

}

// ConfirmDeactivation of a user
func (u *User) ConfirmDeactivation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, err := u.db.GetUserByToken("activator", ps.ByName("token"))
	if err != nil {
		// TODO(coyle): handle error
		return
	}

	if err := u.db.ConfirmUserDeactivation(user.ID); err != nil {
		// TODO(coyle): handle error
	}

}

// CreatePasswordResetToken for a user
func (u *User) CreatePasswordResetToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, err := u.db.GetUser(ps.ByName("id"))
	if err != nil {
		// TODO(coyle): handle error
	}

	_, err = u.db.CreatePasswordResetToken(user.ID)
	if err != nil {
		// TODO(coyle): handle error
	}

	// TODO(coyle): dispatch email after mail service is completed
}

// ConfirmPasswordReset for a user
func (u *User) ConfirmPasswordReset(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	body, err := getBody(r)

	user, err := u.db.GetUserByToken("resetter", ps.ByName("token"))
	if err != nil {
		// TODO(coyle): handle error
	}

	if err := u.db.ResetPassword(user.ID, body.Password); err != nil {
		// TODO(coyle): handle error
	}
}

func (u *User) dispatchActivationEmailSwitch(usr mongodb.User) {
	return
}

func getBody(r *http.Request) (Request, error) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	ur := Request{}

	if err := decoder.Decode(&ur); err != nil && err != io.EOF {
		return ur, err
	}

	return ur, nil

}

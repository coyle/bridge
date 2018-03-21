package users

import (
	"encoding/json"
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
	body, err := getBody(r)
	if err != nil {
		u.logger.Log("Error getting body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if body.PublicKey == "" {
		u.logger.Log("No PublicKey provided", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	nuser := mongodb.User{
		ID:       body.Email,
		Hashpass: body.Password,
	}

	if body.ReferralPartner != "" {
		p, err := u.db.GetPartner(body.ReferralPartner)
		if err != nil {
			u.logger.Log("Error getting partner", err)
		}

		nuser.ReferralPartner = p.ID
	}

	// do all concurrently ?
	user, err := u.db.CreateUser(nuser)
	if err != nil && err == mongodb.ErrInvalidID {
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

	if err := u.db.CreatePublicKey(&user, body.PublicKey); err != nil {
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
		u.logger.Log("invalid request body", err)
		w.WriteHeader(http.StatusBadRequest)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(mongodb.UserToView(user))

}

// ConfirmActivation of a user
func (u *User) ConfirmActivation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, err := u.db.GetUserByToken("activator", ps.ByName("token"))
	if err != nil {
		u.logger.Log("Failed to Get user with provided token", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := u.db.ActivateUser(user.ID); err != nil {
		u.logger.Log("Failed to activate user", err, "ID", user.ID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	user.Activated = true
	json.NewEncoder(w).Encode(mongodb.UserToView(user))
}

// Remove a user
func (u *User) Remove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, err := u.db.GetUser(ps.ByName("id"))
	if err != nil {
		u.logger.Log("failed to get user", "ID", ps.ByName("id"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID, _, ok := r.BasicAuth()
	if !ok {
		u.logger.Log("failed to get user authentication", "ID", ps.ByName("id"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if userID != user.ID {
		u.logger.Log("auth user ID did not match request ID", "request_id", ps.ByName("id"), "auth_id", userID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO(coyle): dispatch email after mail service is completed

	if err := u.db.DeactivateUser(user.ID); err != nil {
		u.logger.Log("failed to deactivate user", "ID", ps.ByName("id"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(mongodb.UserToView(user))

}

// ConfirmDeactivation of a user
func (u *User) ConfirmDeactivation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, err := u.db.GetUserByToken("deactivator", ps.ByName("token"))
	if err != nil {
		u.logger.Log("failed to get user", "token", ps.ByName("token"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := u.db.ConfirmUserDeactivation(user.ID); err != nil {
		u.logger.Log("failed to confirm deactivate user", "ID", ps.ByName("id"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(mongodb.UserToView(user))
}

// CreatePasswordResetToken for a user
func (u *User) CreatePasswordResetToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, err := u.db.GetUser(ps.ByName("id"))
	if err != nil {
		u.logger.Log("failed to get user", "ID", ps.ByName("id"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = u.db.CreatePasswordResetToken(user.ID)
	if err != nil {
		u.logger.Log("failed to create password reset token", "ID", ps.ByName("id"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO(coyle): dispatch email after mail service is completed

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(mongodb.UserToView(user))
}

// ConfirmPasswordReset for a user
func (u *User) ConfirmPasswordReset(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	body, err := getBody(r)
	if err != nil {
		u.logger.Log("unable to get request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := u.db.GetUserByToken("resetter", ps.ByName("token"))
	if err != nil {
		u.logger.Log("failed to find user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := u.db.ResetPassword(user.ID, body.Password); err != nil {
		u.logger.Log("failed to reset password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(mongodb.UserToView(user))
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

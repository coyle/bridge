package mongodb

import (
	"github.com/globalsign/mgo"
)

// Client is the mongoDB implementation of the DB interface
type Client struct {
	session    *mgo.Session
	users      *mgo.Collection
	partners   *mgo.Collection
	publicKeys *mgo.Collection
}

// NewClient instantiates a connection to our MongoDB server
func NewClient(url string) (*Client, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	return &Client{
		session:    session,
		users:      session.DB("bridge").C("users"),
		partners:   session.DB("bridge").C("partners"),
		publicKeys: session.DB("bridge").C("publickeys"),
	}, nil

}

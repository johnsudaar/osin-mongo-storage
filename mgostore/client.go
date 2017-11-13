package mgostore

import (
	"github.com/Scalingo/osin"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type Client struct {
	ID          bson.ObjectId `bson:"_id"`
	ClientID    string        `bson:"client_id"`
	Secret      string        `bson:"secret"`
	RedirectUri string        `bson:"redirect_uri"`
	UserData    interface{}   `bson:"user_data"`
}

func (c *Client) GetId() string            { return c.ClientID }
func (c *Client) GetSecret() string        { return c.Secret }
func (c *Client) GetRedirectUri() string   { return c.RedirectUri }
func (c *Client) GetUserData() interface{} { return c.UserData }

func ToMongoClient(c osin.Client) *Client {
	return &Client{
		ID:          bson.NewObjectId(),
		ClientID:    c.GetId(),
		Secret:      c.GetSecret(),
		RedirectUri: c.GetRedirectUri(),
		UserData:    c.GetUserData(),
	}
}

func (store *MongoStorage) GetClient(id string) (osin.Client, error) {
	clients := store.session.DB(store.dbName).C(CLIENT_COL)
	client := new(Client)
	err := clients.Find(bson.M{
		"client_id": id,
	}).One(client)
	if err != nil {
		return nil, errors.Wrap(err, "fail to get client")
	}
	return client, nil
}

func (store *MongoStorage) SetClient(id string, client osin.Client) error {
	clients := store.session.DB(store.dbName).C(CLIENT_COL)
	_, err := clients.Upsert(bson.M{
		"client_id": id,
	}, ToMongoClient(client))

	if err != nil {
		return errors.Wrap(err, "fail to save client")
	}
	return nil
}

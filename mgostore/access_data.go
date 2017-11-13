package mgostore

import (
	"time"

	"github.com/Scalingo/osin"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type AccessData struct {
	ID            bson.ObjectId  `bson:"_id"`
	Client        *Client        `bson:"client"`
	AuthorizeData *AuthorizeData `bson:"authorize_data"`
	AccessData    *AccessData    `bson:"access_data"`
	AccessToken   string         `bson:"access_token"`
	RefreshToken  string         `bson:"refresh_token"`
	ExpiresIn     int32          `bson:"expires_in"`
	Scope         string         `bson:"scope"`
	RedirectUri   string         `bson:"redirect_uri"`
	CreatedAt     time.Time      `bson:"created_at"`
	UserData      interface{}    `bson:"user_data"`
}

func (d *AccessData) GetClient() osin.Client               { return d.Client }
func (d *AccessData) GetAuthorizeData() osin.AuthorizeData { return d.AuthorizeData }
func (d *AccessData) GetAccessData() osin.AccessData       { return d.AccessData }
func (d *AccessData) GetAccessToken() string               { return d.AccessToken }
func (d *AccessData) GetRefreshToken() string              { return d.RefreshToken }
func (d *AccessData) GetExpiresIn() int32                  { return d.ExpiresIn }
func (d *AccessData) GetScope() string                     { return d.Scope }
func (d *AccessData) GetRedirectUri() string               { return d.RedirectUri }
func (d *AccessData) GetCreatedAt() time.Time              { return d.CreatedAt }
func (d *AccessData) GetUserData() interface{}             { return d.UserData }

func ToMongoAccessData(d osin.AccessData) *AccessData {
	if d == nil {
		return nil
	}
	return &AccessData{
		ID:            bson.NewObjectId(),
		Client:        ToMongoClient(d.GetClient()),
		AuthorizeData: ToMongoAuthorizeData(d.GetAuthorizeData()),
		AccessData:    ToMongoAccessData(d.GetAccessData()),
		AccessToken:   d.GetAccessToken(),
		RefreshToken:  d.GetRefreshToken(),
		ExpiresIn:     d.GetExpiresIn(),
		Scope:         d.GetScope(),
		RedirectUri:   d.GetRedirectUri(),
		CreatedAt:     d.GetCreatedAt(),
		UserData:      d.GetUserData(),
	}
}

func (store *MongoStorage) SaveAccess(data osin.AccessData) error {
	accesses := store.session.DB(store.dbName).C(ACCESS_COL)

	savedData := ToMongoAccessData(data)
	if savedData.AccessData != nil && savedData.AccessData.AccessData != nil {
		savedData.AccessData.AccessData = nil
	}
	_, err := accesses.Upsert(bson.M{
		"access_token": data.GetAccessToken(),
	}, savedData)
	if err != nil {
		return errors.Wrap(err, "fail to save access data")
	}
	return err
}

func (store *MongoStorage) LoadAccess(token string) (osin.AccessData, error) {
	accesses := store.session.DB(store.dbName).C(ACCESS_COL)
	accData := new(AccessData)
	err := accesses.Find(bson.M{
		"access_token": token,
	}).One(accData)
	if err != nil {
		return nil, errors.Wrap(err, "fail to get access data")
	}
	return accData, nil
}

func (store *MongoStorage) RemoveAccess(token string) error {
	accesses := store.session.DB(store.dbName).C(ACCESS_COL)
	err := accesses.Remove(bson.M{
		"access_token": token,
	})
	if err != nil {
		return errors.Wrap(err, "fail to remove access data")
	}
	return nil
}

func (store *MongoStorage) LoadRefresh(token string) (osin.AccessData, error) {
	accesses := store.session.DB(store.dbName).C(ACCESS_COL)
	accData := new(AccessData)
	err := accesses.Find(bson.M{
		"refresh_token": token,
	}).One(accData)
	if err != nil {
		return nil, errors.Wrap(err, "fail to get refresh token")
	}
	return accData, nil
}

func (store *MongoStorage) RemoveRefresh(token string) error {
	accesses := store.session.DB(store.dbName).C(ACCESS_COL)
	err := accesses.Update(bson.M{
		"refresh_token": token,
	}, bson.M{
		"$unset": bson.M{
			"refresh_token": 1,
		},
	})

	if err != nil {
		return errors.Wrap(err, "fail to unset refresh token")
	}
	return nil
}

package mgostore

import (
	"time"

	"github.com/Scalingo/osin"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type AuthorizeData struct {
	ID                  bson.ObjectId `bson:"_id"`
	Client              *Client       `bson:"client"`
	Code                string        `bson:"code"`
	ExpiresIn           int32         `bson:"expires_in"`
	Scope               string        `bson:"scope"`
	RedirectUri         string        `bson:"redirect_uri"`
	State               string        `bson:"state"`
	CreatedAt           time.Time     `bson:"created_at"`
	UserData            interface{}   `bson:"user_data"`
	CodeChallenge       string        `bson:"code_challenge"`
	CodeChallengeMethod string        `bson:"code_challenge_method"`
}

func (d AuthorizeData) GetClient() osin.Client         { return d.Client }
func (d AuthorizeData) GetCode() string                { return d.Code }
func (d AuthorizeData) GetExpiresIn() int32            { return d.ExpiresIn }
func (d AuthorizeData) GetScope() string               { return d.Scope }
func (d AuthorizeData) GetRedirectUri() string         { return d.RedirectUri }
func (d AuthorizeData) GetState() string               { return d.State }
func (d AuthorizeData) GetCreatedAt() time.Time        { return d.CreatedAt }
func (d AuthorizeData) GetUserData() interface{}       { return d.UserData }
func (d AuthorizeData) GetCodeChallenge() string       { return d.CodeChallenge }
func (d AuthorizeData) GetCodeChallengeMethod() string { return d.CodeChallengeMethod }

func ToMongoAuthorizeData(d osin.AuthorizeData) *AuthorizeData {
	return &AuthorizeData{
		ID:                  bson.NewObjectId(),
		Client:              ToMongoClient(d.GetClient()),
		Code:                d.GetCode(),
		ExpiresIn:           d.GetExpiresIn(),
		Scope:               d.GetScope(),
		RedirectUri:         d.GetRedirectUri(),
		State:               d.GetState(),
		CreatedAt:           d.GetCreatedAt(),
		UserData:            d.GetUserData(),
		CodeChallenge:       d.GetCodeChallenge(),
		CodeChallengeMethod: d.GetCodeChallengeMethod(),
	}
}

func (store *MongoStorage) SaveAuthorize(data osin.AuthorizeData) error {
	authorizations := store.session.DB(store.dbName).C(AUTHORIZE_COL)
	_, err := authorizations.Upsert(bson.M{
		"code": data.GetCode(),
	}, ToMongoAuthorizeData(data))

	if err != nil {
		return errors.Wrap(err, "fail to save authorize data")
	}
	return nil
}

func (store *MongoStorage) LoadAuthorize(code string) (osin.AuthorizeData, error) {
	authorizations := store.session.DB(store.dbName).C(AUTHORIZE_COL)
	authData := new(AuthorizeData)
	err := authorizations.Find(bson.M{
		"code": code,
	}).One(authData)
	if err != nil {
		return nil, errors.Wrap(err, "fail to get authorize data")
	}
	return authData, nil
}

func (store *MongoStorage) RemoveAuthorize(code string) error {
	authorizations := store.session.DB(store.dbName).C(AUTHORIZE_COL)
	err := authorizations.Remove(bson.M{
		"code": code,
	})

	if err != nil {
		return errors.Wrap(err, "fail to remove authorize")
	}
	return nil
}

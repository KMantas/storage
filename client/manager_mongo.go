package client

import (
	"github.com/imdario/mergo"
	"github.com/ory/fosite"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoManager cares for the managing of the Mongo Session instance of a Client
type MongoManager struct {
	DB     *mgo.Database
	Hasher fosite.Hasher
}

// GetConcreteClient finds a Client based on ID and returns it, if found in Mongo.
func (m *MongoManager) GetConcreteClient(id string) (*Client, error) {
	collection := m.DB.C("clients").With(m.DB.Session.Copy())
	defer collection.Database.Session.Close()

	client := &Client{}
	err := collection.Find(bson.M{"_id": id}).One(client)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return client, nil
}

// GetClient returns a Client if found by an ID lookup.
func (m *MongoManager) GetClient(id string) (fosite.Client, error) {
	return m.GetConcreteClient(id)
}

// UpdateClient updates an OAuth 2.0 Client record. This is done using the equivalent of an object replace.
func (m *MongoManager) UpdateClient(c *Client) error {
	o, err := m.GetConcreteClient(c.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	// If the password isn't updated, grab it from the stored object
	if c.Secret != "" {
		h, err := m.Hasher.Hash([]byte(c.Secret))
		if err != nil {
			return errors.WithStack(err)
		}
		c.Secret = string(h)
	}

	// Otherwise, update the object with the new updates
	if err := mergo.Merge(c, o); err != nil {
		return errors.WithStack(err)
	}

	// Update Mongo reference with the updated object
	collection := m.DB.C("clients").With(m.DB.Session.Copy())
	defer collection.Database.Session.Close()
	selector := bson.M{"_id": c.ID}
	if err := collection.Update(selector, c); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

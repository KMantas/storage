package user

import (
	// Standard Library Imports
	"fmt"
	// External Imports
	"github.com/ory/fosite"
)

// User provides the specific types for storing, editing, deleting and retrieving a User record in mongo.
type User struct {
	// ID is the uniquely assigned uuid that references the user
	ID string `bson:"_id" json:"id" xml:"id"`

	// The Tenant IDs that the user has been given rights to access
	TenantIDs []string `bson:"tenantIDs,omitempty" json:"tenantIDs,omitempty" xml:"tenantIDs,omitempty"`

	// Username is used to authenticate a user
	Username string `bson:"username" json:"username" xml:"username"`

	// Password of the user - will be a hash based on your fosite selected hasher
	Password string `bson:"password" json:"-" xml:"-"`

	// Scopes contains the scopes that have been granted to
	Scopes []string `bson:"scopes" json:"scopes" xml:"scopes"`

	// FirstName stores the user's Last Name
	FirstName string `bson:"firstName" json:"firstName" xml:"firstName"`

	// LastName stores the user's Last Name
	LastName string `bson:"lastName" json:"lastName" xml:"lastName"`

	// ProfileURI is a pointer to where their profile picture lives
	ProfileURI string `bson:"profileUri" json:"profileUri,omitempty" xml:"profileUri,omitempty"`
}

// GetFullName concatenates the User's First Name and Last Name for templating purposes
func (u User) GetFullName() (fn string) {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// SetPassword takes a cleartext secret, hashes it with a hasher and sets it as the user's password
func (u *User) SetPassword(cleartext string, hasher fosite.Hasher) (err error) {
	h, err := hasher.Hash([]byte(cleartext))
	if err != nil {
		return err
	}
	u.Password = string(h)
	return nil
}

// GetHashedSecret returns the Users's Hashed Secret as a byte array
func (u *User) GetHashedSecret() []byte {
	return []byte(u.Password)
}

// Authenticate compares a cleartext string against the user's
func (u User) Authenticate(cleartext string, hasher fosite.Hasher) error {
	return hasher.Compare(u.GetHashedSecret(), []byte(cleartext))
}

// AddScopes adds multiple scopes to the given user
func (u *User) AddScopes(addScopes ...string) {
	for i := range addScopes {
		found := false
		for j := range u.Scopes {
			if addScopes[i] == u.Scopes[j] {
				found = true
				break
			}
		}
		if !found {
			u.Scopes = append(u.Scopes, addScopes[i])
		}
	}
}

// AddScopes adds multiple scopes to the given user
func (u *User) RemoveScopes(removeScopes ...string) {
	for i := range removeScopes {
		for j := range u.Scopes {
			if removeScopes[i] == u.Scopes[j] {
				copy(u.Scopes[j:], u.Scopes[j+1:])
				u.Scopes[len(u.Scopes)-1] = ""
				u.Scopes = u.Scopes[:len(u.Scopes)-1]
				break
			}
		}
	}
}

// AddTenantIDs adds a single or multiple tenantIDs to the given user
func (u *User) AddTenantIDs(addTenantIDs ...string) {
	for i := range addTenantIDs {
		found := false
		for j := range u.TenantIDs {
			if addTenantIDs[i] == u.TenantIDs[j] {
				found = true
				break
			}
		}
		if !found {
			u.TenantIDs = append(u.TenantIDs, addTenantIDs[i])
		}
	}
}

// RemoveTenants removes a single or multiple tenantIDs from the given user
func (u *User) RemoveTenantIDs(removeTenants ...string) {
	for i := range removeTenants {
		for j := range u.TenantIDs {
			if removeTenants[i] == u.TenantIDs[j] {
				copy(u.TenantIDs[j:], u.TenantIDs[j+1:])
				u.TenantIDs[len(u.TenantIDs)-1] = ""
				u.TenantIDs = u.TenantIDs[:len(u.TenantIDs)-1]
				break
			}
		}
	}
}

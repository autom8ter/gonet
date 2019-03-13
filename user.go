package gonet

import (
	"github.com/autom8ter/util"
	"github.com/go-bongo/bongo"
	"time"
)

type UserOptionFn func(u *User)
type UserValidateFn func(u *User) error

// User is a retrieved and authentiacted user.
type User struct {
	bongo.DocumentBase `bson:",inline"`
	UserID             string            `json:"user_id"`
	Nickname           string            `json:"nickname"`
	Sub                string            `json:"sub"`
	Name               string            `json:"name"`
	GivenName          string            `json:"given_name"`
	FamilyName         string            `json:"family_name"`
	Profile            string            `json:"profile"`
	Picture            string            `json:"picture"`
	Email              string            `json:"email"`
	EmailVerified      bool              `json:"email_verified"`
	Gender             string            `json:"gender"`
	Phone              string            `json:"phone"`
	Content            string            `json:"content"`
	Tags               map[string]string `json:"tags"`
	AuthKey            string            `json:"auth_key"`
	AuthToken          string            `json:"auth_token"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

func NewUser(opts ...UserOptionFn) *User {
	u := &User{}
	for _, o := range opts {
		o(u)
	}
	return u
}

func ValidateUser(opts ...UserValidateFn) (*User, error) {
	u := &User{}
	for _, o := range opts {
		if err := o(u); err != nil {
			return u, err
		}
	}
	return u, nil
}

func (o *User) JSONify() []byte {
	return util.ToPrettyJson(o)
}

func (o *User) String() string {
	return util.ToPrettyJsonString(o)
}

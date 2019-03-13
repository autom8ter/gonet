package gonet

import "github.com/go-bongo/bongo"

// User is a retrieved and authentiacted user.
type User struct {
	bongo.DocumentBase `bson:",inline"`
	Sub                string `json:"sub"`
	Name               string `json:"name"`
	GivenName          string `json:"given_name"`
	FamilyName         string `json:"family_name"`
	Profile            string `json:"profile"`
	Picture            string `json:"picture"`
	Email              string `json:"email"`
	EmailVerified      bool   `json:"email_verified"`
	Gender             string `json:"gender"`
	Phone              string `json:"phone"`
}

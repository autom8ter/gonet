package gonet

import (
	"fmt"
	"github.com/autom8ter/authzero"
	"net/http"
)

type AuthZero struct {
	*Router
	auth0 *authzero.Auth0
}

func NewAuthZeroRouter(addr string, cfg *authzero.Config) *AuthZero {
	return &AuthZero{
		Router: NewRouter(addr),
		auth0:  authzero.NewAuth0(cfg),
	}
}

func (a *AuthZero) Login(path string) {
	a.Router.Mux().HandleFunc(path, a.auth0.Login())
	fmt.Println("registered handler: ", path)
}

func (a *AuthZero) LogOut(path string, returnTo string) {
	a.Router.Mux().HandleFunc(path, a.auth0.Logout(path, returnTo))
	fmt.Printf("registered handler: %s return to: %s\n", path, returnTo)
}

func (a *AuthZero) OAuthCallBack(path, loginPath string) {
	a.Router.Mux().HandleFunc(path, a.auth0.OAuth(loginPath))
	fmt.Println("registered handler: ", path)
}

func (a *AuthZero) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return a.auth0.RequireAuth(next)
}

func (a *AuthZero) Auth0(next http.HandlerFunc) *authzero.Auth0 {
	return a.auth0
}

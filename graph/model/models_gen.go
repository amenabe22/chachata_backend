// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AuthResult struct {
	Token  string `json:"token"`
	Status bool   `json:"status"`
}

type NewUsrInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
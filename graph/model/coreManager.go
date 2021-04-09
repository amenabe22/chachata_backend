package model

import (
	"crypto/sha1"

	p "github.com/wuriyanto48/go-pbkdf2"
	"gorm.io/gorm"
)

func Authenticate(password string, email string, db *gorm.DB) (*User, bool) {
	// db := gorm.DB{}
	usrs := []*User{}
	db.First(&usrs, "email = ?", email)
	println("HERE")
	if len(usrs) == 1 {
		hashStat := CheckPasswordHash(password, usrs[0].Password)
		if hashStat {
			return usrs[0], true
		}
	}
	return nil, false

}

func HashPassword(password string) (string, error) {
	pass := p.NewPassword(sha1.New, 8, 32, 15000)
	hashed := pass.HashPassword(password)
	// bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	println(hashed.Salt)
	println(hashed.CipherText)
	return hashed.Salt, nil
}

func CheckPasswordHash(password string, hash string) bool {
	pass := p.NewPassword(sha1.New, 8, 32, 15000)
	hashed := pass.HashPassword(password)
	isValid := pass.VerifyPassword(password, hashed.CipherText, hashed.Salt)

	// err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return isValid
}

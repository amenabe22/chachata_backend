package core

import (
	"github.com/amenabe22/chachata_backend/graph/model"
	"github.com/amenabe22/chachata_backend/graph/setup"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(password string, email string) (*model.User, bool) {
	db := setup.SetupModels()
	usrs := []*model.User{}
	db.First(&usrs, "email = ?", email)
	if len(usrs) == 1 {
		hashStat := CheckPasswordHash(password, usrs[0].Password)
		if hashStat {
			return usrs[0], true
		}
	}
	return nil, false

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

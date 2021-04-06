package helpers

import (
	"math/rand"

	"github.com/amenabe22/chachata_backend/graph/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Export struct {
}

func (c Export) RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func CheckDuplicate(phone string, username string, user model.User, coredb *gorm.DB) (map[string]bool, error) {
	allUsrs := []*model.User{}
	allExceptUser := []*model.User{}
	dupStats := map[string]bool{
		"dupPhone": false,
		"dupUname": false,
	}
	coredb.Preload(clause.Associations).Find(&allUsrs)
	for _, usr := range allUsrs {
		if usr.ID != user.ID {
			allExceptUser = append(allExceptUser, usr)
		}
	}
	for _, usr := range allExceptUser {
		if phone == usr.Profile.Phone {
			dupStats["dupPhone"] = true
		}
		if username == usr.Profile.Username {
			dupStats["dupUname"] = true
		}
	}
	return dupStats, nil
}

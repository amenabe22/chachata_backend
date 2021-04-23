package model

import (
	"errors"
	"time"

	"github.com/amenabe22/chachata_backend/graph/core"
	"github.com/dgryski/trifles/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	SomeFlag  bool       `gorm:"column:some_flag;not null;default:true"`
	// main content goes here
	Email      string `gorm:"uniqueIndex,unique"`
	Password   string
	ProfileId  string
	Profile    Profile
	IsVerified bool
	Qrcode     string
	// UserDevices []Devices `gorm:"many2many:devices;" json:"devices,omitempty"`
	UserDevices []Devices `gorm:"polymorphic:Owner"`
	OwnerID     string
	OwnerType   string
}

// Profile is the model for the profile table.
type Profile struct {
	ID         string `gorm:"type:char(36);primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `sql:"index"`
	Name       string     `gorm:"column:name;size:128;not null;"`
	Username   string     `gorm:"uniqueIndex,unique"`
	Phone      string     `gorm:"uniqueIndex,unique"`
	ProfilePic string
	Complete   bool
	Progress   int
	Bio        string
}

type Devices struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	// UserID     string     `json:"user_id" gorm:"user_id"` //nolint:gofmt
	AppID      string
	DeviceName string
	OwnerID    string
	OwnerType  string
}

func AllUsrs(coredb *gorm.DB) ([]*User, error) {
	// return all users list
	usrs := []*User{}
	coredb.Preload("Profile").Preload("UserDevices").Find(&usrs)
	return usrs, nil
}

func UpdateUserProfile(coredb *gorm.DB, user User, input ProfileStarterInput) (ProfileUpdateResult, error) {
	profile := Profile{}
	otherUsrsSet := User{}
	// TODO: FIX excluded updates here
	coredb.Preload(clause.Associations).First(&otherUsrsSet).Where("id ! ?", user.ID)
	// preload the user profile object
	coredb.Preload(clause.Associations).First(&profile, "id = ?", user.ProfileId)
	// update the profile with the new coming content
	coredb.Preload(clause.Associations).Find(&user)
	duplicateData, _ := CheckDuplicate(input.Phone, input.Username, user, coredb)
	usernameDup := duplicateData["dupUname"]
	phoneDup := duplicateData["dupPhone"]
	if usernameDup == true {
		return ProfileMessages("err", "phone is duplicate"), errors.New("Username is already taken !")
	}
	if phoneDup == true {
		return ProfileMessages("err", "username is duplicate"), errors.New("Phone Number is already taken !")
	}

	coredb.Model(&profile).Updates(map[string]interface{}{
		"name": input.Name, "phone": input.Phone, "username": input.Username})

	return ProfileMessages("succ", "profile changed"), nil
}

func ProfileMessages(messageType string, message string) ProfileUpdateResult {
	errResult := ProfileUpdateResult{
		Message: message,
		Stat:    false,
	}
	successResult := ProfileUpdateResult{
		Message: message,
		Stat:    true,
	}
	if messageType == "err" {
		return errResult
	} else if messageType == "succ" {
		return successResult
	}
	return errResult
}

func CheckDuplicate(phone string, username string, user User, coredb *gorm.DB) (map[string]bool, error) {
	allUsrs := []*User{}
	allExceptUser := []*User{}
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
func AddNewUsr(input NewUsrInput, coredb *gorm.DB) (string, error) {
	// room := r.AdminChans[]
	allUsrs := []*User{}

	coredb.First(&allUsrs, "email = ?", input.Email)
	if len(allUsrs) != 0 {
		return "", errors.New("Email is already taken")
	}
	hashedPassword, perr := HashPassword(input.Password)
	if perr != nil {
		return "", perr
	}
	usrId := uuid.UUIDv4()
	qrProfile, _ := core.GenerateQrOnSignup(usrId)
	usr := User{
		ID:       usrId,
		SomeFlag: false,
		Email:    input.Email,
		Password: hashedPassword,
		Qrcode:   qrProfile,
	}

	usr.Profile = Profile{
		ID:       uuid.UUIDv4(),
		Name:     "",
		Username: uuid.UUIDv4(),
		Phone:    uuid.UUIDv4(),
		// Qrcode:   qrProfile,
	}
	usr.UserDevices = append(usr.UserDevices, Devices{
		ID:         uuid.UUIDv4(),
		DeviceName: input.DeviceInput.DeviceName,
		AppID:      input.DeviceInput.AppID,
	})
	coredb.Save(&usr)
	tokenRes, _ := EmailAuthLogin(usr.Email, usr.Password, coredb)
	return tokenRes.Token, nil
}

func EmailAuthLogin(email string, password string, coredb *gorm.DB) (*AuthResult, error) {
	usr, authStat := Authenticate(password, email, coredb)
	errModel := AuthResult{
		Token:  "",
		Status: false,
	}
	authErr := errors.New("Incorrect credentials")
	if !authStat {
		return &errModel, authErr
	}
	expiredAt := int(time.Now().Add(time.Hour * 87600).Unix())
	tokenString := GenerateJwt(usr.ID, int64(expiredAt))
	return &AuthResult{
		Token:  tokenString,
		Status: true,
	}, nil
}

func forgotPassword(email string, coredb *gorm.DB) (bool, error) {
	usrs := []*User{}
	coredb.First(&usrs, "email = ?", email)
	if len(usrs) == 0 {
		return false, errors.New("Email is unknown")
	}
	// send email to user to change password and stuff
	return false, nil
}

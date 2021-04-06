package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/amenabe22/chachata_backend/graph/chans"
	"github.com/amenabe22/chachata_backend/graph/core"
	"github.com/amenabe22/chachata_backend/graph/generated"
	"github.com/amenabe22/chachata_backend/graph/helpers"
	"github.com/amenabe22/chachata_backend/graph/model"
	"github.com/amenabe22/chachata_backend/graph/setup"
	"github.com/amenabe22/chachata_backend/middlewares"
	"github.com/dgryski/trifles/uuid"
	"gorm.io/gorm/clause"
)

func (r *mutationResolver) RemoveAllUsrs(ctx context.Context) (bool, error) {
	usrs := []*model.User{}
	coredb.Find(&usrs)
	err := coredb.Delete(&usrs).Error
	if err != nil {
		return false, errors.New("Error removing users")
	}
	return true, nil
}

func (r *mutationResolver) NewUsr(ctx context.Context, input model.NewUsrInput) (string, error) {
	// room := r.AdminChans[]
	allUsrs := []*model.User{}

	coredb.First(&allUsrs, "email = ?", input.Email)
	if len(allUsrs) != 0 {
		return "err", errors.New("Email is already taken")
	}
	message := "success"
	hashedPassword, perr := core.HashPassword(input.Password)
	if perr != nil {
		return "error", perr
	}
	usr := model.User{
		ID:       uuid.UUIDv4(),
		SomeFlag: false,
		Email:    input.Email,
		Password: hashedPassword,
	}
	// coredb.Create(&usr)
	println(usr.Profile.Name)
	println("LOOK  UP THERE")
	// err := coredb.Create(&usr).Error
	// if err != nil {
	// 	return "", nil
	// }

	// usr.Profile = model.Profile{}
	profile, _ := core.GenerateQrOnSignup(usr.ID)
	println(profile, "PROF")
	usr.Profile = model.Profile{
		ID:       uuid.UUIDv4(),
		Name:     "",
		Username: uuid.UUIDv4(),
		Phone:    uuid.UUIDv4(),
	}
	coredb.Save(&usr)
	// println(usr.Profile.Username)
	return message, nil
}

func (r *mutationResolver) EmailAuthLogin(ctx context.Context, email string, password string) (*model.AuthResult, error) {
	usr, authStat := core.Authenticate(password, email)
	errModel := model.AuthResult{
		Token:  "",
		Status: false,
	}
	authErr := errors.New("Incorrect credentials")
	if !authStat {
		return &errModel, authErr
	}
	expiredAt := int(time.Now().Add(time.Hour * 87600).Unix())
	tokenString := middlewares.GenerateJwt(usr.ID, int64(expiredAt))
	// tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	// _, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": usr.ID})
	// fmt.Printf("DEBUG: sample jwt is %s\n\n", tokenString)
	// token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiQW1lbiBBYmUifQ.WKJ7WKx9Pg76uoqhKkR-Cjb8lxigvEO-ruCroAsy9KM"
	// fmt.Println(tokenAuth.Decode(token))

	// token, err := jwt.GenerateToken(usr.ID)
	// if err != nil {
	// 	return &errModel, authErr
	// }
	return &model.AuthResult{
		Token:  tokenString,
		Status: true,
	}, nil
}

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func (r *mutationResolver) UpdateProfileStarter(ctx context.Context, uid model.ProfileStarterInput) (*model.ProfileUpdateResult, error) {
	errResult := &model.ProfileUpdateResult{
		Message: "err",
		Stat:    false,
	}
	successResult := &model.ProfileUpdateResult{
		Message: "Success",
		Stat:    true,
	}
	isValid, user, err := middlewares.GetAuthStat(coredb, ctx, "Invalid token")
	if err != nil {
		return errResult, err
	}
	// check if the token in the header is valid
	if isValid {
		profile := model.Profile{}
		otherUsrsSet := model.User{}
		// TODO: FIX excluded updates here
		coredb.Preload(clause.Associations).First(&otherUsrsSet).Where("id ! ?", user.ID)
		// preload the user profile object
		coredb.Preload(clause.Associations).First(&profile, "id = ?", user.ProfileId)
		// update the profile with the new coming content
		coredb.Preload(clause.Associations).Find(&user)
		duplicateData, _ := helpers.CheckDuplicate(uid.Phone, uid.Username, user, coredb)
		usernameDup := duplicateData["dupUname"]
		phoneDup := duplicateData["dupPhone"]
		if usernameDup == true {
			return errResult, errors.New("Username is already taken !")
		}
		if phoneDup == true {
			return errResult, errors.New("Phone is already taken !")
		}

		coredb.Model(&profile).Updates(map[string]interface{}{"name": uid.Name, "phone": uid.Phone, "username": uid.Username})
	}
	return successResult, nil
}

func (r *queryResolver) AllUsrs(ctx context.Context) ([]*model.User, error) {
	usrs := []*model.User{}
	coredb.Preload(clause.Associations).Find(&usrs)
	// coredb.Preload("User").Preload("Profile").Find(&usrs)
	// coredb.Find(&usrs)
	return usrs, nil
}

func (r *queryResolver) SecureInfo(ctx context.Context) (string, error) {
	tokenStat, user, err := middlewares.GetAuthStat(coredb, ctx, "Invalid token")
	if err != nil {
		return "", err
	}
	if tokenStat {
		println("Authentication approved !!!")
		println(user.Email)
	}
	println("hi")
	return "Hey there", nil
}

func (r *queryResolver) UserData(ctx context.Context, id string) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *subscriptionResolver) AdminsNotified(ctx context.Context) (<-chan *string, error) {
	roomId := "admins"
	r.mu.Lock()
	room := r.AdminChans[roomId]
	if room == nil {
		room = &chans.CoreAdminChannel{
			RoomId: roomId,
			Observers: map[string]struct {
				Username string
				Message  chan *string
			}{},
		}
		r.AdminChans[roomId] = room
	}
	var value helpers.Export
	r.mu.Unlock()

	id := value.RandString(8)
	events := make(chan *string, 1)

	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(room.Observers, id)
		r.mu.Unlock()
	}()

	r.mu.Lock()
	room.Observers[id] = struct {
		Username string
		Message  chan *string
	}{Username: "hey", Message: events}
	r.mu.Unlock()
	return events, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
var coredb = setup.SetupModels()

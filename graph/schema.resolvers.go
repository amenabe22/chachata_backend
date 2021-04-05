package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/amenabe22/chachata_backend/graph/chans"
	"github.com/amenabe22/chachata_backend/graph/core"
	"github.com/amenabe22/chachata_backend/graph/generated"
	"github.com/amenabe22/chachata_backend/graph/helpers"
	"github.com/amenabe22/chachata_backend/graph/jwt"
	"github.com/amenabe22/chachata_backend/graph/model"
	"github.com/amenabe22/chachata_backend/graph/setup"
	"github.com/amenabe22/chachata_backend/middlewares"
	"github.com/dgryski/trifles/uuid"
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
	err := coredb.Create(&usr).Error
	if err != nil {
		return "", nil
	}
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

	token, err := jwt.GenerateToken(usr.ID)
	if err != nil {
		return &errModel, authErr
	}
	return &model.AuthResult{
		Token:  token,
		Status: true,
	}, nil
}

func (r *queryResolver) AllUsrs(ctx context.Context) ([]*model.User, error) {
	println("NIGGAs")
	usrs := []*model.User{}
	coredb.Find(&usrs)
	return usrs, nil
}

func (r *queryResolver) SecureInfo(ctx context.Context) (string, error) {

	user := middlewares.ForContext(ctx)
	if user == nil {
		return "error", fmt.Errorf("access denied")
	}
	return "Hey there", nil

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

package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/amenabe22/chachata_backend/graph/chans"
	"github.com/amenabe22/chachata_backend/graph/generated"
	"github.com/amenabe22/chachata_backend/graph/helpers"
	"github.com/amenabe22/chachata_backend/graph/model"
)

func (r *mutationResolver) RemoveAllUsrs(ctx context.Context) (bool, error) {
	usrs := []*model.User{}
	r.Coredb.Find(&usrs)
	err := r.Coredb.Delete(&usrs).Error
	if err != nil {
		return false, errors.New("Error removing users")
	}
	return true, nil
}

func (r *mutationResolver) NewUsr(ctx context.Context, input model.NewUsrInput) (string, error) {
	token, err := model.AddNewUsr(input, r.Coredb)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *mutationResolver) EmailAuthLogin(ctx context.Context, email string, password string) (*model.AuthResult, error) {
	usr, authStat := model.Authenticate(password, email, r.Coredb)
	errModel := model.AuthResult{
		Token:  "",
		Status: false,
	}
	authErr := errors.New("Incorrect credentials")
	if !authStat {
		return &errModel, authErr
	}
	expiredAt := int(time.Now().Add(time.Hour * 87600).Unix())
	tokenString := model.GenerateJwt(usr.ID, int64(expiredAt))
	return &model.AuthResult{
		Token:  tokenString,
		Status: true,
	}, nil
}

func (r *mutationResolver) UpdateProfileStarter(ctx context.Context, input model.ProfileStarterInput) (*model.ProfileUpdateResult, error) {
	errResult := model.ProfileUpdateResult{
		Message: "err",
		Stat:    false,
	}
	isValid, user, err := model.GetAuthStat(r.Coredb, ctx, "Invalid token")
	if err != nil {
		return &errResult, err
	}
	// check if the token in the header is valid
	if isValid {
		message, err := model.UpdateUserProfile(r.Coredb, user, input)
		if err != nil {
			return nil, err
		}
		return &message, nil
	}
	return nil, nil
}

func (r *queryResolver) AllUsrs(ctx context.Context) ([]*model.User, error) {
	users, err := model.AllUsrs(r.Coredb)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *queryResolver) SecureInfo(ctx context.Context) (string, error) {
	tokenStat, user, err := model.GetAuthStat(r.Coredb, ctx, "Invalid token")
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

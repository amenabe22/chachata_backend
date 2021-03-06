package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/amenabe22/chachata_backend/graph/chans"
	"github.com/amenabe22/chachata_backend/graph/generated"
	"github.com/amenabe22/chachata_backend/graph/helpers"
	"github.com/amenabe22/chachata_backend/graph/model"
	"github.com/dgryski/trifles/uuid"
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
	room := r.NotificationChans["users"]
	// message := "Hey everyone we have a new User"
	notification := model.Notifications{ID: uuid.UUIDv4(), Type: "urgent", Text: "Hey there nigga"}
	for _, observer := range room.Observers {
		observer.Message <- &notification
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

func (r *mutationResolver) ForgotPassword(ctx context.Context) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Post(ctx context.Context, text string, username string, roomName string) (*model.Message, error) {
	r.mu.Lock()
	// fetch rooms from redis and check if they are not duplicate
	cmd := r.RedisClient.SMembers("rooms")
	if cmd.Err() != nil {
		log.Println(cmd.Err(), "CMD ERR")
		return nil, cmd.Err()
	}
	res, err := cmd.Result()
	if err != nil {
		log.Println(err, "ERROR WITH RESULT")
		return nil, err
	}
	// roomId := ""
	newRoom := true
	for _, cr := range res {
		var finRoom model.InstatntMessage
		json.Unmarshal([]byte(cr), &finRoom)
		// finRooms = append(finRooms, &finRoom)
		// println(finRoom.ID, "HERE AGAIN")
		if finRoom.Name == roomName {
			newRoom = false
			// roomId = finRoom.ID
			// println("ROOM IS NOT NEW", finRoom.ID)
		}
	}

	room := r.Rooms[roomName]
	println(roomName, room.Name, room.ID, "GO GO GOGO")
	// the below uses the above logic after the room is checked for duplicate
	// if newRoom || room == nil {
	// 	println("NEW ONE", roomId)
	// 	if roomId == "" {

	// 		room = &model.Chatroom{
	// 			ID:   uuid.UUIDv4(),
	// 			Name: roomName,
	// 			Observers: map[string]struct {
	// 				Username string
	// 				Message  chan *model.Message
	// 			}{},
	// 		}
	// 	} else {
	// 		room = &model.Chatroom{
	// 			ID:   roomId,
	// 			Name: roomName,
	// 			Observers: map[string]struct {
	// 				Username string
	// 				Message  chan *model.Message
	// 			}{},
	// 		}
	// 	}
	// 	r.Rooms[roomName] = room
	// 	// var rm model.Chatroom
	// 	// json.Unmarshal([]byte(rj), &rm)
	// 	// // unmarshal the json object
	// 	// println(rm.Name, "FO REAL")
	// }
	// else {
	// 	// make sure to keep the room ID
	// 	room.ID = roomId
	// }
	println("HOW ABOUT NOW", newRoom, room == nil)
	// println(room.ID, room.Name, "SSSS")
	r.mu.Unlock()
	var value helpers.Export

	message := model.Message{
		ID:        value.RandString(8),
		CreatedAt: time.Now(),
		Text:      text,
		CreatedBy: username,
	}
	println("Look down")
	println(room.ID, "SSSSSss")
	instantMess := &model.InstantMessage{
		ID:      room.ID,
		Name:    room.Name,
		Message: message,
	}
	room.Messages = append(room.Messages, message)
	// this marhalled data is being non unique every time unless it's repeated
	rj, _ := json.Marshal(instantMess)
	// using sadd to make sure there's no duplicate
	// use foreign key relation to make sure things are stable
	if err := r.RedisClient.SAdd("rooms", rj).Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	r.mu.Lock()
	for _, observer := range room.Observers {
		println(observer.Username, "OBSERVER")
		println(message.CreatedBy, "CREATOR")
		// if observer.Username == "" || observer.Username == message.CreatedBy {
		observer.Message <- &message
		// }
	}
	r.mu.Unlock()
	return &message, nil
}

func (r *mutationResolver) PopAllChats(ctx context.Context) (bool, error) {
	// rms := model.Chatroom{}
	if err := r.RedisClient.Del("all_rooms").Err(); err != nil {
		log.Println(err)
		return false, err
	}
	if err := r.RedisClient.Del("rooms").Err(); err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) NewRoom(ctx context.Context, roomName string) (*model.Chatroom, error) {
	r.mu.Lock()
	room := &model.Chatroom{
		ID:   uuid.UUIDv4(),
		Name: roomName,
		Observers: map[string]struct {
			Username string
			Message  chan *model.Message
		}{},
	}
	r.Rooms[roomName] = room
	// save the room inside a new set now
	rj, _ := json.Marshal(room)

	if err := r.RedisClient.SAdd("all_rooms", rj).Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	// send notification
	// notification := "Here is some message"
	notification := model.Notifications{ID: uuid.UUIDv4(), Type: "roomNot", Text: "This is room is finna be close MF"}
	notificationRoom := r.RoomNotificationsChans["roomNotification"]
	if notificationRoom != nil {
		for _, observer := range notificationRoom.Observers {
			println(observer.Username)
			observer.Message <- &notification
		}
	}
	r.mu.Unlock()
	// end here
	return room, nil
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

func (r *queryResolver) UserData(ctx context.Context) (*model.User, error) {
	tokenStat, user, err := model.GetAuthStat(r.Coredb, ctx, "Invalid token")
	if err != nil {
		return nil, err
	}
	if tokenStat {
		return &user, nil
	}
	return nil, nil
}

func (r *queryResolver) Room(ctx context.Context, name string) (*model.Chatroom, error) {
	r.mu.Lock()
	room := r.Rooms[name]
	if room == nil {
		room = &model.Chatroom{
			Name: name,
			Observers: map[string]struct {
				Username string
				Message  chan *model.Message
			}{},
		}
		r.Rooms[name] = room
	}
	r.mu.Unlock()
	return room, nil
}

func (r *queryResolver) AllRoomsList(ctx context.Context) ([]*model.Chatroom, error) {
	cmd := r.RedisClient.SMembers("all_rooms")
	if cmd.Err() != nil {
		log.Println(cmd.Err(), "CMD ERR")
		return nil, cmd.Err()
	}
	res, err := cmd.Result()
	if err != nil {
		log.Println(err, "ERROR WITH RESULT")
		return nil, err
	}
	chatRooms := []*model.Chatroom{}
	for _, cr := range res {
		// var rm model.Chatroom
		var finRoom model.Chatroom
		// json.Unmarshal([]byte(cr), &rm)
		json.Unmarshal([]byte(cr), &finRoom)
		chatRooms = append(chatRooms, &finRoom)
	}
	return chatRooms, nil
}

func (r *queryResolver) AllRooms(ctx context.Context) ([]*model.InstatntMessage, error) {
	cmd := r.RedisClient.SMembers("rooms")
	if cmd.Err() != nil {
		log.Println(cmd.Err(), "CMD ERR")
		return nil, cmd.Err()
	}
	res, err := cmd.Result()
	if err != nil {
		log.Println(err, "ERROR WITH RESULT")
		return nil, err
	}
	chatRooms := []*model.InstatntMessage{}
	for _, cr := range res {
		// var rm model.Chatroom
		var finRoom model.InstatntMessage
		// json.Unmarshal([]byte(cr), &rm)
		json.Unmarshal([]byte(cr), &finRoom)
		chatRooms = append(chatRooms, &finRoom)
	}
	return chatRooms, nil
}

func (r *queryResolver) AllMessages(ctx context.Context) ([]*model.Chatroom, error) {
	cmd := r.RedisClient.SMembers("messages")
	if cmd.Err() != nil {
		log.Println(cmd.Err(), "CMD ERR")
		return nil, cmd.Err()
	}
	res, err := cmd.Result()
	if err != nil {
		log.Println(err, "ERROR WITH RESULT")
		return nil, err
	}
	chatRooms := []*model.Chatroom{}
	for _, cr := range res {
		// var rm model.Chatroom
		var finRoom model.Chatroom
		// json.Unmarshal([]byte(cr), &rm)
		json.Unmarshal([]byte(cr), &finRoom)
		chatRooms = append(chatRooms, &finRoom)
	}
	return chatRooms, nil
}

func (r *queryResolver) SingleRoomMessages(ctx context.Context, room string) ([]*model.Message, error) {
	cmd := r.RedisClient.SMembers("rooms")
	if cmd.Err() != nil {
		log.Println(cmd.Err(), "CMD ERR")
		return nil, cmd.Err()
	}
	res, err := cmd.Result()
	if err != nil {
		log.Println(err, "ERROR WITH RESULT")
		return nil, err
	}

	allMes := []*model.Message{}
	for _, cr := range res {
		var finRoom model.InstantMessage
		json.Unmarshal([]byte(cr), &finRoom)
		// check if the room id matches any from the database
		println(finRoom.Message.Text, "....................", room, finRoom.ID)
		if finRoom.ID == room {
			println("ROOM MESSAGES ARE HERE")
			// nested loop makes it slower
			allMes = append(allMes, &finRoom.Message)
		}
	}
	// Fixed with sort
	sort.Slice(allMes, func(i, j int) bool {
		return allMes[i].CreatedAt.Before(allMes[j].CreatedAt)
	})
	return allMes, nil
}

func (r *subscriptionResolver) UsersNotified(ctx context.Context) (<-chan *model.Notifications, error) {
	r.mu.Lock()
	roomId := "users"
	room := r.NotificationChans[roomId]
	if room == nil {
		room = &chans.NotificationChannel{
			RoomId: roomId,
			Observers: map[string]struct {
				Username string
				Message  chan *model.Notifications
			}{},
		}
		r.NotificationChans[roomId] = room
	}
	var value helpers.Export
	r.mu.Unlock()

	id := value.RandString(8)
	events := make(chan *model.Notifications, 1)

	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(room.Observers, id)
		r.mu.Unlock()
	}()

	r.mu.Lock()
	room.Observers[id] = struct {
		Username string
		Message  chan *model.Notifications
	}{Username: "hey", Message: events}
	r.mu.Unlock()
	return events, nil
}

func (r *subscriptionResolver) MessageAdded(ctx context.Context, roomName string, username string) (<-chan *model.Message, error) {
	// r.mu.Lock()
	room := r.Rooms[roomName]
	var finalRoom model.Chatroom
	isRoominDb := false
	// var finalRoom model.Chatroom
	// fmt.Println(room == nil, "NO ROOM")
	if room == nil {
		// if room in state is empty then get rooms from redis database
		cmd := r.RedisClient.SMembers("all_rooms")
		if cmd.Err() != nil {
			log.Println(cmd.Err(), "CMD ERR")
			return nil, cmd.Err()
		}
		res, err := cmd.Result()
		if err != nil {
			log.Println(err, "ERROR WITH RESULT")
			return nil, err
		}
		for _, cr := range res {
			var finRoom model.Chatroom
			// json.Unmarshal([]byte(cr), &rm)
			json.Unmarshal([]byte(cr), &finRoom)
			if finRoom.Name == roomName {
				isRoominDb = true
				finalRoom = finRoom
				// println("got room from db", room.Name, room.ID)
			}
		}
		// chcek if the outer final room is available
		room = &finalRoom
		r.Rooms[finalRoom.Name] = &finalRoom

		// if room is null create a new one
		if !isRoominDb {
			return nil, errors.New("Room not found")
		}
	}
	// r.mu.Unlock()
	var value helpers.Export
	id := value.RandString(8)
	events := make(chan *model.Message, 1)
	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(room.Observers, id)
		r.mu.Unlock()
	}()
	r.mu.Lock()
	// add a new observer
	room.Observers[id] = struct {
		Username string
		Message  chan *model.Message
	}{
		Username: username,
		Message:  events,
	}
	println(len(room.Observers), "OBSRVERS COUNT")
	r.mu.Unlock()
	return events, nil
}

func (r *subscriptionResolver) RoomNotification(ctx context.Context) (<-chan *model.Notifications, error) {
	r.mu.Lock()
	roomId := "roomNotification"
	room := r.RoomNotificationsChans[roomId]
	if room == nil {
		room = &chans.RoomNotification{
			RoomId: roomId,
			Observers: map[string]struct {
				Username string
				Message  chan *model.Notifications
			}{},
		}
		r.RoomNotificationsChans[roomId] = room
	}
	var value helpers.Export
	r.mu.Unlock()

	id := value.RandString(8)
	events := make(chan *model.Notifications, 1)

	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(room.Observers, id)
		r.mu.Unlock()
	}()

	r.mu.Lock()
	room.Observers[id] = struct {
		Username string
		Message  chan *model.Notifications
	}{Username: "hey", Message: events}
	r.mu.Unlock()
	return events, nil
}

func (r *subscriptionResolver) NewRoomAdded(ctx context.Context) (<-chan *model.Chatroom, error) {
	panic(fmt.Errorf("not implemented"))
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

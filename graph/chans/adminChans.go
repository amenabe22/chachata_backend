package chans

type CoreAdminChannel struct {
	RoomId    string
	Message   string
	Observers map[string]struct {
		Username string
		Message  chan *string
	}
}

package group

const stateSize = 24

type Group struct {
	ID          int64
	Name        string
	ScreenName  string
	IsConnected bool
}

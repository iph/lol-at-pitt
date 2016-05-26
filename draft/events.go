package draft

import "encoding/json"

const (
	TeamEventType     = "AddTeamEvent"
	PlayerEventType   = "AddPlayerEvent"
	WinEventType      = "WinEvent"
	BidEventType      = "BidEvent"
	SkipEventType     = "SkipEvent"
	ClearEventType    = "UndoWinEvent"
	PreviousEventType = "UndoSkip"
)

type DraftEvent struct {
	Type    string
	Payload json.RawMessage
}

type AddPlayerEvent struct {
	LeagueID int64
	RealName string
	Roles    []LeagueRole
}

type AddTeamEvent struct {
	TeamName        string
	DrafterLeagueID int64
	DrafterID       string
}

type WinEvent struct {
	TeamName          string
	CaptainFacebookID string
	Amount            int
	PlayerLeagueID    int64
}

type BidEvent struct {
	Amount            int
	PlayerLeagueID    int64
	DrafterFacebookID string
}

type SkipEvent struct {
	PlayerLeagueID int64
}

type TimeEvent struct {
	RemainingTime int
	TimeElapsed   int
}

type LoginEvent struct {
	DrafterFacebookID string
	TeamName          string
	Points            int
}

type DraftStore struct {
	DraftEvents []*DraftEvent
	ID          string
}

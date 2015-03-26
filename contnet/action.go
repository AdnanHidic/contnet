package contnet

type ActionType uint8

const (
	ACTION_NONE ActionType = iota
	ACTION_READ
	ACTION_UPVOTE
	ACTION_DOWNVOTE
)

type Action struct {
	IdentityID int64
	Type       ActionType
	Arguments  map[string]string
}

func NewAction(identityID int64, actionType ActionType, arguments map[string]string) *Action {
	return &Action{
		IdentityID: identityID,
		Type:       actionType,
		Arguments:  arguments,
	}
}

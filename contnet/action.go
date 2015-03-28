package contnet

type ActionType uint8

const (
	ACTION_NONE ActionType = iota
	ACTION_READ
	ACTION_UPVOTE
	ACTION_DOWNVOTE
)

type Action struct {
	ProfileID int64
	Type      ActionType
	Arguments map[string]string
}

func NewAction(profileID int64, actionType ActionType, arguments map[string]string) *Action {
	return &Action{
		ProfileID: profileID,
		Type:      actionType,
		Arguments: arguments,
	}
}

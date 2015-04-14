package contnet

import "time"

type ActionType uint8

const (
	action_none ActionType = iota
	action_read
	action_upvote
	action_downvote
)

var ActionTypes = struct {
	None     ActionType
	Read     ActionType
	Upvote   ActionType
	Downvote ActionType
}{
	None:     action_none,
	Read:     action_read,
	Upvote:   action_upvote,
	Downvote: action_downvote,
}

type Action struct {
	ProfileID ID `json:"profileID"`
	ContentID ID `json:"contentID"`
	Content   *Content
	Type      ActionType      `json:"type"`
	Arguments ActionArguments `json:"arguments"`
	Timestamp time.Time       `json:"timestamp"`
}
type ActionFactory struct{}

func (factory ActionFactory) New(profileID, contentID ID, actionType ActionType, timestamp time.Time, args ActionArguments) *Action {
	return &Action{
		ProfileID: profileID,
		ContentID: contentID,
		Type:      actionType,
		Arguments: args.Clone(),
		Timestamp: timestamp,
	}
}

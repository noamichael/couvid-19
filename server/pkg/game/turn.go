package game

import "time"

//TurnState is the current state in the multipart turn
type TurnState string

const (
	//TurnBlocked the turn was blocked
	TurnBlocked TurnState = "Blocked"
	//TurnBlocking the turn is being blocked
	TurnBlocking TurnState = "Blocking"
	//TurnBlockFailed the turn is being blocked
	TurnBlockFailed TurnState = "BlockFailed"
	//TurnCalled the turn in under call
	TurnCalled TurnState = "Called"
	//TurnCalling the turn in under call
	TurnCalling TurnState = "Calling"
	//TurnCallFailed the turn in under call
	TurnCallFailed TurnState = "CallFailed"
	//TurnAwaitingTarget - the turn is awaiting an action from the player
	TurnAwaitingTarget TurnState = "AwaitingTarget"
	//TurnPending - the turn is awaiting an action from the player
	TurnPending TurnState = "Pending"
	//TurnSubmitted - the turn was submitted and can still be blocked or called
	TurnSubmitted TurnState = "Submitted"
	//TurnComplete - the turn was completed
	TurnComplete TurnState = "Completed"
)

//Turn represents a single turn
type Turn struct {
	Player       *Player
	TargetPlayer *Player
	Action       ActionType
	State        TurnState
	Blocker      *Player
	Caller       *Player
	BlockedWith  CardType
	//Time indicates that this turn might still be blocked or called
	timer *time.Timer
	//channel for once the turn is fully resolved
	resolved chan bool
}

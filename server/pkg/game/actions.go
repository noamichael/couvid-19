package game

//ActionType are all the actions a player can take
type ActionType string

const (
	//ActionBlock is the block action
	ActionBlock ActionType = "Block"
	//ActionCall will call a bluff
	ActionCall ActionType = "Call"
	//ActionCoup will kill a player
	ActionCoup ActionType = "Coup"
	//ActionDuke will can take three
	ActionDuke ActionType = "Duke"
	//ActionSteal is the steal action
	ActionSteal ActionType = "Steal"
	//ActionForeignAid is the ability to take three coins
	ActionForeignAid ActionType = "ForeignAid"
	//ActionAssassinate is the assassinate action
	ActionAssassinate ActionType = "Assassinate"
	//ActionAmbassador is the ability to view the top two cards of the deck
	ActionAmbassador ActionType = "Ambassador"
	//ActionTakeOne is the default "take one coin" action
	ActionTakeOne ActionType = "TakeOne"
)

//Blockable returns true if this action can be blocked
func (action ActionType) Blockable() bool {
	return action != ActionCall &&
		action != ActionBlock &&
		action != ActionTakeOne &&
		action != ActionCoup
}

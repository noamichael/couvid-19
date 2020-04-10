package game

//CardType represents one of the five coup cards
type CardType string

const (
	//Ambassador is the type of the Ambassador card
	Ambassador CardType = "Ambassador"
	//Assassin is the type of the ASSASSIN card
	Assassin CardType = "Assassin"
	//Captain is the type of the CAPTAIN card
	Captain CardType = "Captain"
	//Contessa is the type of the Contessa card
	Contessa CardType = "Contessa"
	//Duke is the type of the Duke card
	Duke CardType = "Duke"
)

//CanBlock returns true if the cardType can block the given action
func (cardType CardType) CanBlock(action ActionType) bool {
	switch action {
	case ActionAmbassador:
		return Ambassador == cardType
	case ActionAssassinate:
		return Contessa == cardType
	case ActionForeignAid:
		return Duke == cardType
	case ActionSteal:
		return Captain == cardType
	}
	return false
}

//CanPerform returns true if this card type can perform the given action
func (cardType CardType) CanPerform(action ActionType) bool {
	switch action {
	case ActionAmbassador:
		return Ambassador == cardType
	case ActionAssassinate:
		return Assassin == cardType
	case ActionDuke:
		return Duke == cardType
	case ActionSteal:
		return Captain == cardType
	}
	return false
}

//CardState is whether the card is alive or dead
type CardState string

const (
	//ALIVE is for when a card is alive
	ALIVE CardState = "ALIVE"
	//DEAD is for when a card is dead
	DEAD CardState = "DEAD"
)

//Card represents a single card
type Card struct {
	CardType  CardType
	CardState CardState
}

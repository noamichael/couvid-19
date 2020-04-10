package game

import "fmt"

//Player represents a single player in the gamer
type Player struct {
	Name        string
	Coins       int
	ClientID    string
	LoseCard    bool
	TradedCards []CardType
	Cards       []*Card
	OutOfGame   bool
}

func (player *Player) getCard(cardType CardType) *Card {
	for _, card := range player.Cards {
		if card.CardType == cardType && card.CardState == ALIVE {
			return card
		}
	}
	return nil
}

func (player *Player) hasCard(cardType CardType) bool {
	return player.getCard(cardType) != nil
}

func (player *Player) canPerform(action ActionType) bool {
	return player.getCardForAction(action) != nil
}

func (player *Player) getCardForAction(action ActionType) *Card {
	for _, card := range player.Cards {
		if card.CardState == ALIVE && card.CardType.CanPerform(action) {
			return card
		}
	}
	return nil
}

func (player *Player) killCard(cardType CardType) error {
	cardToKill := player.getCard(cardType)

	if cardToKill == nil {
		return fmt.Errorf("Cannot kill card %s: player doesn't have this card", cardToKill)
	}

	cardToKill.CardState = DEAD

	return nil
}

func (player *Player) removeCard(cardType CardType) *Card {
	newCards := make([]*Card, 0, 2)
	var removedCard *Card
	for _, card := range player.Cards {
		if card.CardState == ALIVE && card.CardType == cardType {
			removedCard = card
		} else {
			newCards = append(newCards, card)
		}
	}
	player.Cards = newCards
	return removedCard
}

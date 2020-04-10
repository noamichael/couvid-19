package game

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

//State is the current status of the game
type State string

const (
	//StatePending indicates the game is not yet started
	StatePending = "Pending"
	//StatePlaying indicates the game is active
	StatePlaying = "Playing"
	//StateFinished indicates the game is finished
	StateFinished = "Finished"
)

//Game holds the state of a single Coup game
type Game struct {
	State       State
	Turns       []*Turn
	CurrentTurn *Turn
	Players     []*Player
	Coins       int
	Deck        []*Card
	events      chan Game
}

//NewGame will reinitize the game state
func NewGame() *Game {
	game := &Game{
		Players: make([]*Player, 0, 6),
		Turns:   make([]*Turn, 0),
		Coins:   50,
		State:   StatePending,
		Deck:    make([]*Card, 0, 15),
	}

	deck := game.Deck

	addToDeck := func(cardType CardType) {
		for i := 0; i < 3; i++ {
			deck = append(deck, &Card{
				CardType: cardType,
			})
		}
	}

	addToDeck(Ambassador)
	addToDeck(Assassin)
	addToDeck(Captain)
	addToDeck(Contessa)
	addToDeck(Duke)

	game.Deck = deck

	game.shuffleDeck()

	return game
}

//AddPlayer will create a player object with the given name and clientID
func (game *Game) AddPlayer(name string, clientID string) {
	player := &Player{
		Name:     name,
		ClientID: clientID,
		Coins:    2,
		Cards:    make([]*Card, 0, 2),
	}

	game.Players = append(game.Players, player)
}

//ResolveAmbassador allows the player to swap cards
func (game *Game) ResolveAmbassador(clientID string, cardToKeep, cardToReturn CardType) error {
	turn := game.CurrentTurn

	//Validate we can
	err := game.validateForTurn(clientID, ActionAmbassador, TurnAwaitingTarget)

	if err != nil {
		return err
	}

	if cardToKeep != "" && cardToKeep != cardToReturn {
		card := turn.Player.removeCard(cardToReturn)
		if game.Deck[0].CardType == cardToKeep {
			turn.Player.Cards = append(turn.Player.Cards, game.Deck[0])
			game.Deck = game.Deck[1:]
		} else if game.Deck[1].CardType == cardToKeep {
			turn.Player.Cards = append(turn.Player.Cards, game.Deck[1])
			game.Deck = append(game.Deck[2:], game.Deck[0])
		}
		game.returnToDeck(card)
	}

	return game.resolveTurn()

}

//ResolveFailedCall allows a caller of a turn to say which card they lose
func (game *Game) ResolveFailedCall(clientID string, cardTypeToLose CardType) error {
	turn := game.CurrentTurn

	if turn.State != TurnCallFailed {
		return errors.New("Cannot resolve failed call: game in invalid state")
	}

	if err := turn.Caller.killCard(cardTypeToLose); err != nil {
		return err
	}

	return game.resolveTurn()
}

//Turn will set the action for the current turn
func (game *Game) Turn(clientID string, action ActionType, targetPlayer string) error {
	turn := game.CurrentTurn

	if clientID == "" {
		return errors.New("Missing clientID")
	}

	if clientID != game.CurrentTurn.Player.ClientID {
		return errors.New("Current player does not match for turn")
	}

	if turn.State != TurnPending {
		return errors.New("Turn has already been submitted")
	}

	if targetPlayer != "" {
		turn.TargetPlayer = game.findPlayerByName(targetPlayer)
	}

	turn.Action = action
	turn.State = TurnSubmitted

	if action.Blockable() {
		//automatically resolve the turn after above timer times out
		//In other words, no one channeles it
		turn.timer = time.NewTimer(time.Second * 10)

		go func() {
			<-turn.timer.C
			game.resolveTurn()
		}()

		return nil
	}
	//Not blockable - attempt to resolve right away
	game.resolveTurn()

	return nil
}

//CallBlock allows a player to call the blocker
func (game *Game) CallBlock(callerClientID string) error {

	turn := game.CurrentTurn
	caller := game.findPlayerByClientID(callerClientID)

	if turn.State != TurnBlocking {
		return errors.New("Cannot call a block unless the turn is being contested")
	}

	if caller == turn.Blocker {
		return errors.New("Cannot call own block")
	}

	blocker := game.CurrentTurn.Blocker

	blockedWith := blocker.getCard(turn.BlockedWith)

	if blockedWith != nil && blockedWith.CardType.CanBlock(turn.Action) {
		turn.State = TurnBlocked
		return nil
	}

	turn.State = TurnBlockFailed

	return nil
}

func (game *Game) resolveTurn() error {
	turn := game.CurrentTurn
	player := game.CurrentTurn.Player
	playerName := player.Name

	switch turn.Action {
	case ActionAmbassador:
		if turn.State == TurnSubmitted {
			turn.State = TurnAwaitingTarget
			player.TradedCards = game.peakTwo()
		} else if turn.State == TurnAwaitingTarget {
			turn.State = TurnComplete
		}
		break
	case ActionAssassinate:
		if turn.State == TurnSubmitted {
			if player.Coins < 3 {
				return errors.New("Not enough coins to Assassinate")
			}
			turn.State = TurnAwaitingTarget
			turn.Player.Coins -= 3
		} else if turn.State == TurnAwaitingTarget {
			turn.State = TurnComplete
		}
		break
	case ActionCoup:
		if player.Coins < 7 {
			return errors.New("Not enough coins to Coup")
		}
		player.Coins -= 7
		turn.State = TurnAwaitingTarget
		fmt.Printf("\n%s Couped %s", playerName, turn.TargetPlayer.Name)
	case ActionForeignAid:
		turn.Player.Coins += 2
		turn.State = TurnComplete
		fmt.Printf("\n%s took three", playerName)
		break
	case ActionDuke:
		turn.Player.Coins += 3
		turn.State = TurnComplete
		fmt.Printf("\n%s took three", playerName)
		break
	case ActionTakeOne:
		turn.Player.Coins++
		turn.State = TurnComplete
		fmt.Printf("\n%s took one", playerName)
		break
	case ActionSteal:
		turn.Player.Coins += 2
		turn.TargetPlayer.Coins -= 2
		turn.State = TurnComplete
		fmt.Printf("\n%s Took 2 from %s", playerName, turn.TargetPlayer.Name)
		break
	}

	if turn.State == TurnComplete {
		turn.Player.TradedCards = nil
		//TODO: Cleanup all player statesx
		game.nextTurn()
	}

	return nil

}

//BlockTurn will allow a user to attempt to block a specific turn
func (game *Game) BlockTurn(blockerClientID string, blockWith CardType) error {
	turn := game.CurrentTurn
	blocker := game.findPlayerByClientID(blockerClientID)

	if blocker == turn.Player {
		return errors.New("User cannot block their own turn")
	}

	if !turn.Action.Blockable() {
		return fmt.Errorf("%s is not blockable", turn.Action)
	}

	if turn.State != TurnSubmitted {
		return errors.New("Turn can only be blocked when in Submitted state")
	}

	turn.timer.Stop()

	turn.State = TurnBlocking
	turn.Blocker = blocker
	turn.BlockedWith = blockWith

	return nil
}

//CallTurn call the current turn's action
func (game *Game) CallTurn(callerClientID string) error {
	caller := game.findPlayerByClientID(callerClientID)
	turn := game.CurrentTurn

	if caller == turn.Player {
		return errors.New("User cannot call their own turn")
	}

	if turn.State != TurnSubmitted {
		return errors.New("Turn cannot be called. Not in a submitted state")
	}

	if !turn.Action.Blockable() {
		return fmt.Errorf("%s is not callable", turn.Action)
	}

	turn.timer.Stop()

	turn.Caller = caller

	if card := turn.Player.getCardForAction(turn.Action); card != nil {
		turn.State = TurnCallFailed
		turn.Player.removeCard(card.CardType)
		game.returnToDeck(card)
		game.resolveTurn()

	} else {
		turn.State = TurnCalled
	}

	return nil
}

//StartGame will deal the cards and set the inital player
func (game *Game) StartGame() {
	players := game.Players
	numberOfPlayers := len(players)

	for i := 0; i < numberOfPlayers; i++ {
		player := players[i]
		for j := 0; j < 2; j++ {
			game.dealCard(player)
		}
	}

	game.CurrentTurn = &Turn{
		Player: players[0],
		State:  TurnPending,
	}

	game.State = StatePlaying
}

func (game *Game) shuffleDeck() {
	deck := game.Deck
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
}

func (game *Game) nextTurn() {
	index := -1

	for i, v := range game.Players {
		if game.CurrentTurn.Player == v {
			index = i
			break
		}
	}

	if index == -1 {
		return //TODO: Handle
	}
	index++

	if index == len(game.Players) {
		index = 0
	}

	game.Turns = append(game.Turns, game.CurrentTurn)

	game.CurrentTurn = &Turn{
		Player: game.Players[index],
		State:  TurnPending,
	}

	fmt.Printf("\nIt is %s's turns", game.CurrentTurn.Player.Name)

}

func (game *Game) validateForTurn(clientID string, action ActionType, expectedState TurnState) error {
	if clientID == "" {
		return errors.New("Missing clientID")
	}

	if clientID != game.CurrentTurn.Player.ClientID {
		return errors.New("Current player does not match for turn")
	}

	if action != game.CurrentTurn.Action {
		return errors.New("Action does not match current turn")
	}

	if expectedState != game.CurrentTurn.State {
		return errors.New("State does not match current turn")
	}

	return nil
}

func (game *Game) dealCard(player *Player) {
	card := game.Deck[0] //draw first element
	card.CardState = ALIVE
	player.Cards = append(player.Cards, card)
	game.Deck = game.Deck[1:] //remove first element
}

func (game *Game) returnToDeck(card *Card) {
	card.CardState = ALIVE
	game.Deck = append(game.Deck, card)
	game.shuffleDeck()
}

func (game *Game) peakTwo() []CardType {
	cardTypes := make([]CardType, 0, 2)
	deck := game.Deck
	numberOfCards := len(deck)

	if numberOfCards > 0 {
		cardTypes = append(cardTypes, deck[0].CardType)
	}
	if numberOfCards > 1 {
		cardTypes = append(cardTypes, deck[1].CardType)
	}

	return cardTypes
}

func (game *Game) findPlayerByClientID(clientID string) *Player {
	return game.findPlayerByPredicate(func(p *Player) bool {
		return p.ClientID == clientID
	})
}

func (game *Game) findPlayerByName(name string) *Player {
	return game.findPlayerByPredicate(func(p *Player) bool {
		return p.Name == name
	})
}

func (game *Game) findPlayerByPredicate(predicate func(*Player) bool) *Player {
	for _, p := range game.Players {
		if predicate(p) {
			return p
		}
	}
	return nil
}

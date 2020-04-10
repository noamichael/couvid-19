package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/noamichael/couvid-19/server/pkg/game"
)

func main() {
	// hub := client.NewHub()
	// go hub.Run()
	// http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	client.ServeWs(hub, w, r)
	// })

	michael, john, alex, eddy := "1", "2", "3", "4"

	coup := game.NewGame()

	coup.AddPlayer("Michael", michael)
	coup.AddPlayer("John", john)
	coup.AddPlayer("Alex", alex)
	coup.AddPlayer("Eddy", eddy)

	coup.StartGame()

	coup.Turn(michael, game.ActionAmbassador, "")
	time.Sleep(11 * time.Second)
	cardToKeep := coup.CurrentTurn.Player.TradedCards[0]
	coup.ResolveAmbassador(michael, cardToKeep, coup.CurrentTurn.Player.Cards[0].CardType)

	coup.Turn(john, game.ActionTakeOne, "")
	coup.Turn(alex, game.ActionDuke, "")
	coup.CallTurn(john)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Hit Enter to Exit")
	reader.ReadString('\n')

}

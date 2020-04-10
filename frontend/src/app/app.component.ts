import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  currentPlayer =  {
    name: 'Michael',
    coins: 7,
    cards: [
      {
        type: 'duke',
        label: 'Duke',
        state: 'alive'
      },
      {
        type: 'captin',
        label: 'Captin',
        state: 'dead'
      }
    ]
  }
  players = [
    {
      name: "John",
      coins: 3,
      cards: [ {state: 'alive'}, {state: 'alive'} ]
    },
    {
      name: "Kyle",
      coins: 8,
      cards: [ {state: 'alive'}, {type: 'duke', label: 'Duke', state: 'dead'} ]
    },
    {
      name: "Alex",
      coins: 8,
      cards: [ {type: 'captin', state: 'dead', label: 'Captin'}, {type: 'duke', label: 'Duke', state: 'dead'} ]
    }
  ]

  ngOnInit() {
    this.processPlayer(this.currentPlayer);
    this.players.forEach(p => this.processPlayer(p))
  }

  processPlayer(player) {
    player.coinsArray = new Array(player.coins || 0)
  }
}

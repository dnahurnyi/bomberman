# bomberman client
Client for the bomberman game
[![codecov.io Code Coverage](https://img.shields.io/codecov/c/github/dwyl/hapi-auth-jwt2.svg?maxAge=2592000)](https://codecov.io/github/dnahurnyi/bomberman?branch=master)

Check details at: https://dojorena.io/codenjoy-contest/resources/help/bomberman.html

To install this package run:
 > `go get github.com/dnahurnyi/bomberman`

And import github.com/dnahurnyi/bomberman to your bot

Use next code snippet to try this client:
```
package main

import (
	"log"

	"github.com/dnahurnyi/bomberman"
)

func main() {
	browserURL := "https://dojorena.io/codenjoy-contest/board/player/{player-id}?code={code}&gameName=bomberman"
	game, c := bomberman.StartGame(browserURL)

	for {
		select {
		case <-c.Done:
			log.Fatal("It's done")
		case <-c.Read:
			// Make your move
			game.Move(bomberman.ACT)
			c.Write <- struct{}{}
		}
	}
}

```

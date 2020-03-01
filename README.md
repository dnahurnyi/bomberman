# bomberman client
Client for the bomberman game

Check details at: https://dojorena.io/codenjoy-contest/resources/help/bomberman.html

Use next code snippet to try this client with keyboard:
```
package main

import (
	"log"

	"github.com/dnahurnyi/bomberman"
	"github.com/eiannone/keyboard"
)

// Example URL
// https://dojorena.io/codenjoy-contest/board/player/793wvxsasdfasdfas1z?code=53193242346800&gameName=bomberman
// https -> wss
// Scheme - wss
// Host - dojorena.io
// Path - /codenjoy-contest/ws
// Query - user=793wvxsasdfasdfas1z&code=53193242346800&gameName=bomberman

func main() {
	// Use your URL
	browserURL := "https://dojorena.io/codenjoy-contest/board/player/793wvxsasdfasdfas1z?code=53193242346800&gameName=bomberman"

	game, done := bomberman.StartGame(browserURL)

	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	for {
		select {
		case <-done:
			log.Fatal("It's done")
		default:
			char, _, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			} else if char == 'q' {
				log.Fatal("It's done")
			}
			// You have 0.5 second to make a move or game will think you 
			// make a move "STOP"
			moveAction := bomberman.Action(char)
			switch moveAction {
			case bomberman.Action('w'):
				moveAction = bomberman.UP
			case bomberman.Action('s'):
				moveAction = bomberman.DOWN
			case bomberman.Action('d'):
				moveAction = bomberman.RIGHT
			case bomberman.Action('a'):
				moveAction = bomberman.LEFT
			case bomberman.Action('c'):
				moveAction = bomberman.ACT
			}
			game.Move(moveAction)
		}
	}
}

```

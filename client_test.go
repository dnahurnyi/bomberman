package client

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	testValidBoard = "☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼ #        ####  # #&########♠ #☼☼ ☼ ☼#☼#☼ ☼ ☼ ☼#☼#☼ ☼ ☼#☼ ☼ ☼ ☼ ☼☼  # # #  ♠ ##    # ####  ☻    #☼☼☺☼#☼♥☼ ☼ ☼ ☼ ☼ ☼#☼ ☼ ☼ ☼ ☼ ☼ ☼#☼☼# ####          ♥   # #   ## ##☼☼ ☼#☼#☼#☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼#☼ ☼ ☼ ☼ ☼☼#  ###  #             ## #   ##☼☼♥☼#☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼#☼ ☼☼##                    3       #☼☼ ☼ ☼#☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼#☼ ☼ ☼ ☼ ☼☼  #  #                         ☼☼#☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼☼##                       #     ☼☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼☼                           #   ☼☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼#☼ ☼ ☼ ☼☼ #                #      #     ☼☼ ☼ ☼ ☼ ☼ ☼#☼ ☼ ☼ ☼ ☼&☼#☼ ☼#☼ ☼ ☼☼&#    #     #                  ☼☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼☼   &                     #     ☼☼#☼#☼#☼ ☼ ☼ ☼ ☼ ☼#☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼☼ #        #               &    ☼☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼#☼ ☼ ☼#☼ ☼ ☼ ☼ ☼☼     #                &        ☼☼ ☼ ☼ ☼#☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼☼  # &                    #     ☼☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼#☼ ☼ ☼#☼☼ # #         #                 ☼☼#☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼#☼ ☼&☼#☼ ☼ ☼#☼☼    &  #&          ♥        #  ☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼"
	testValidMsg = "board="+testValidBoard
	testValidStructuredBoard = `☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼
☼ #        ####  # #&########♠ #☼
☼ ☼ ☼#☼#☼ ☼ ☼ ☼#☼#☼ ☼ ☼#☼ ☼ ☼ ☼ ☼
☼  # # #  ♠ ##    # ####  ☻    #☼
☼☺☼#☼♥☼ ☼ ☼ ☼ ☼ ☼#☼ ☼ ☼ ☼ ☼ ☼ ☼#☼
☼# ####          ♥   # #   ## ##☼
☼ ☼#☼#☼#☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼#☼ ☼ ☼ ☼ ☼
☼#  ###  #             ## #   ##☼
☼♥☼#☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼#☼ ☼
☼##                    3       #☼
☼ ☼ ☼#☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼#☼ ☼ ☼ ☼ ☼
☼  #  #                         ☼
☼#☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼
☼##                       #     ☼
☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼
☼                           #   ☼
☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼#☼ ☼ ☼ ☼
☼ #                #      #     ☼
☼ ☼ ☼ ☼ ☼ ☼#☼ ☼ ☼ ☼ ☼&☼#☼ ☼#☼ ☼ ☼
☼&#    #     #                  ☼
☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼
☼   &                     #     ☼
☼#☼#☼#☼ ☼ ☼ ☼ ☼ ☼#☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼
☼ #        #               &    ☼
☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼#☼ ☼ ☼#☼ ☼ ☼ ☼ ☼
☼     #                &        ☼
☼ ☼ ☼ ☼#☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼
☼  # &                    #     ☼
☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼#☼ ☼ ☼#☼
☼ # #         #                 ☼
☼#☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼ ☼#☼ ☼&☼#☼ ☼ ☼#☼
☼    &  #&          ♥        #  ☼
☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼☼
`
)

func Test_createURL(t *testing.T) {
	type tstruct struct {
		name string
		browserUrl string
		expectedURL url.URL
		expectedError error
		setEnv func()
	}

	tests := []tstruct{
		{
			name: "Success, user input",
			browserUrl: "https://dojorena.io/codenjoy-contest/board/player/793wdxskw521spo4mn1y?code=531459153668826800&gameName=bomberman",
			expectedURL: url.URL{
				Scheme: gameProtocol,
				Host: "dojorena.io",
				Path: gamePath,
				RawQuery: fmt.Sprintf(gameQueryTemplate, "793wdxskw521spo4mn1y", "531459153668826800"),
			},
			expectedError: nil,
			setEnv: func() {
				// Don't set any variable
			},
		}, {
			name: "Success, get from env",
			browserUrl: "",
			expectedURL: url.URL{
				Scheme: gameProtocol,
				Host: "dojorena.io",
				Path: gamePath,
				RawQuery: fmt.Sprintf(gameQueryTemplate, "793wdxskw521spo4mn1y", "531459153668826800"),
			},
			expectedError: nil,
			setEnv: func() {
				os.Setenv("HOST", "dojorena.io")
				os.Setenv("PLAYER", "793wdxskw521spo4mn1y")
				os.Setenv("CODE", "531459153668826800")
			},
		}, {
			name: "Invalid host",
			browserUrl: "dojorena.iocodenjoy-contest/board/player/793wdxskw521spo4mn1y?code=531459153668826800&gameName=bomberman",
			expectedURL: url.URL{},
			expectedError: errors.New("Invalid URL, can't get host name, url: "+"dojorena.iocodenjoy-contest/board/player/793wdxskw521spo4mn1y?code=531459153668826800&gameName=bomberman"),
			setEnv: func() {},
		}, {
			name: "Invalid player ID",
			browserUrl: "https://dojorena.io/codenjoy-contest/board/player/793wdxskw521spo4mn1ycode=531459153668826800&gameName=bomberman",
			expectedURL: url.URL{},
			expectedError: errors.New("Invalid URL, can't get player ID, url: "+"https://dojorena.io/codenjoy-contest/board/player/793wdxskw521spo4mn1ycode=531459153668826800&gameName=bomberman"),
			setEnv: func() {},
		}, {
			name: "Invalid game code",
			browserUrl: "https://dojorena.io/codenjoy-contest/board/player/793wdxskw521spo4mn1y?code=531459153668826800gameName=bomberman",
			expectedURL: url.URL{},
			expectedError: errors.New("Invalid URL, can't get game code, url: "+"https://dojorena.io/codenjoy-contest/board/player/793wdxskw521spo4mn1y?code=531459153668826800gameName=bomberman"),
			setEnv: func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name,  func(t *testing.T) {
			os.Clearenv()
			tt.setEnv()
			url, err := createURL(tt.browserUrl)
			assert.Equal(t, tt.expectedURL, url)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func Test_updateBoard(t *testing.T) {
	type tstruct struct {
		name string
		msg string
		board *board
		resBoard *board
		expectedErr error
	}

	tests := []tstruct{
		{
			name: "Successful update",
			msg: testValidMsg,
			board: &board{},
			resBoard: &board{
				rawBoard: testValidBoard,
			},
			expectedErr: nil,
		}, {
			name: "Invalid message",
			msg: "Invalid message",
			board: &board{},
			resBoard: &board{},
			expectedErr: errors.New("Invalid input, board size is not valid, input msg: "+"Invalid message"),
		}, {
			name: "Empty message",
			msg: "",
			board: &board{},
			resBoard: &board{},
			expectedErr: errors.New("Invalid input, board size is not valid, input msg: "),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name,  func(t *testing.T) {
			err := updateBoard(tt.msg, tt.board)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.resBoard.rawBoard, tt.board.rawBoard)
		})
	}
}

func Test_readWriteSocket(t *testing.T) {
	type tstruct struct {
		name string
		board *board
		rawBoardRes string
		server func(response string, done chan struct{}) func (w http.ResponseWriter, r *http.Request)
		async bool
	}

	workingServer := func(response string, done chan struct{}) func (w http.ResponseWriter, r *http.Request){
		return func (w http.ResponseWriter, r *http.Request) {
			upgrader := websocket.Upgrader{}
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return
			}
			defer c.Close()
			// End on the second run
			runCount := 0
			for {
				err = c.WriteMessage(websocket.TextMessage, []byte(response))
				if err != nil {
					break
				}

				_, _, err := c.ReadMessage()
				if err != nil {
					break
				}

				if runCount == 0 {
					done <- struct{}{}
					runCount++
				}
			}
		}
	}
	brokenServer := func(response string, done chan struct{}) func (w http.ResponseWriter, r *http.Request){
		return func (w http.ResponseWriter, r *http.Request) {
			upgrader := websocket.Upgrader{}
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return
			}
			c.Close()
			// Fail on first try connect
			done <- struct{}{}
			return
		}
	}
	badDataServer := func(response string, done chan struct{}) func (w http.ResponseWriter, r *http.Request){
		return func (w http.ResponseWriter, r *http.Request) {
			upgrader := websocket.Upgrader{}
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return
			}
			defer c.Close()
			for {
				err = c.WriteMessage(websocket.TextMessage, []byte("invalid data"))
				if err != nil {
					break
				}

				_, _, err := c.ReadMessage()
				if err != nil {
					break
				}
			}
		}
	}

	tests := []tstruct{
		{
			name: "Successful board update",
			board: &board{},
			rawBoardRes: testValidBoard,
			server: workingServer,
			async: true,
		}, {
			name: "Server is down",
			board: &board{},
			rawBoardRes: "",
			server: brokenServer,
			async: false,
		},
		{
			name: "Server returns bad data",
			board: &board{},
			rawBoardRes: "",
			server: badDataServer,
			async: false,
		},
	}


	for _, tt := range tests {
		t.Run(tt.name,  func(t *testing.T) {
			// Start mock server
			readyToCheck := make(chan struct{})
			server := httptest.NewServer(http.HandlerFunc(tt.server(testValidMsg, readyToCheck)))
			defer server.Close()
			// Prepare connection
			server.URL = strings.Replace(server.URL, "http", "ws", 1)
			u, _ := url.Parse(server.URL)
			conn, _ := getConnection(*u)
			// Setup exit
			done := make(chan struct{})

			// In async case we wait unit server will allow us to check the result
			if tt.async {
				defer func() {
					done <- struct{}{}
				}()
				go readWriteSocket(tt.board, conn, done)
				<- readyToCheck
			// In sync case we run readWriteSocket synchronously and don't need to use extra mechanisms
			} else {
				readWriteSocket(tt.board, conn, done)
			}

			assert.Equal(t, tt.rawBoardRes, tt.board.rawBoard)
		})
	}
}

func Test_StartGame(t *testing.T) {
	type tstruct struct {
		name string
		browserUrl string
		boardRepresentation string
		server func(response string, done chan struct{}) func (w http.ResponseWriter, r *http.Request)
		async bool
		panicValue string
	}

	workingServer := func(response string, done chan struct{}) func (w http.ResponseWriter, r *http.Request){
		return func (w http.ResponseWriter, r *http.Request) {
			upgrader := websocket.Upgrader{}
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return
			}
			defer c.Close()
			// End on the second run
			runCount := 0
			for {
				err = c.WriteMessage(websocket.TextMessage, []byte(response))
				if err != nil {
					break
				}

				_, msg, err := c.ReadMessage()
				if len(msg) != 0 {
					switch(Action(msg)) {
					case ACT:
						response = strings.Replace(response, string(BOMBERMAN), string(BOMB_BOMBERMAN), 1)
					}
				}
				if err != nil {
					break
				}

				if runCount == 0 {
					done <- struct{}{}
					runCount++
				}
			}
		}
	}

	tests := []tstruct{
		{
			name: "Successful board update",
			browserUrl: "https://{serverHostname}/codenjoy-contest/board/player/793wdxskw521spo4mn1y?code=531459153668826800&gameName=bomberman",
			async: true,
			server: workingServer,
			boardRepresentation: testValidStructuredBoard,
		}, {
			name: "Invalid URL",
			browserUrl: "",
			async: false,
			panicValue: "Failed to create valid game url, err:  Invalid URL, can't get host name, url: \n",
		}, {
			name: "Can't create connection",
			browserUrl: "https://127.0.0.1/codenjoy-contest/board/player/793wdxskw521spo4mn1y?code=531459153668826800&gameName=bomberman",
			async: false,
			panicValue: "Failed to create connection to game, err:  dial tcp 127.0.0.1:80: connect: connection refused\n",
		},
	}


	for _, tt := range tests {
		t.Run(tt.name,  func(t *testing.T) {
			readyToCheck := make(chan struct{})
			if tt.server != nil {
				gameProtocol = "ws"
				// Start mock server
				server := httptest.NewServer(http.HandlerFunc(tt.server(testValidMsg, readyToCheck)))
				defer server.Close()
				// Prepare connection
				server.URL = strings.Replace(server.URL, "http", "ws", 1)
				serverHostname := strings.Replace(server.URL, "ws://", "", 1)
				tt.browserUrl = strings.Replace(tt.browserUrl, "{serverHostname}", serverHostname, 1)
			}

			var game Game
			var done chan struct{}
			// Setup exit
			if tt.async {
				defer func() {
					done <- struct{}{}
				}()
				game, done = StartGame(tt.browserUrl)
				<-readyToCheck
				assert.Equal(t, tt.boardRepresentation, game.Show())
				game.Move(ACT)
				time.Sleep(time.Second*1)
				assert.NotEqual(t, tt.boardRepresentation, game.Show()) // bomberman changed to bomb
			} else {
				assert.PanicsWithValue(t, tt.panicValue, func(){StartGame(tt.browserUrl)})
			}
		})
	}
}

func Test_Move(t *testing.T) {
	type tstruct struct {
		name string
		board *board
		command Action // Resulting move
		moves []Action // List of moves asked to do
		expectedErr error
	}

	tests := []tstruct{
		{
			name: "One move",
			board: &board{},
			moves: []Action{ACT},
			command: ACT,
			expectedErr: nil,
		}, {
			name: "List of moves",
			board: &board{},
			moves: []Action{UP, DOWN, ACT},
			command: ACT,
			expectedErr: nil,
		}, {
			name: "Invalid move",
			board: &board{},
			moves: []Action{"invalid"},
			command: STOP,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name,  func(t *testing.T) {
			for _, m := range tt.moves {
				tt.board.Move(m)
			}
			assert.Equal(t, tt.command, tt.board.getAction())
		})
	}
}

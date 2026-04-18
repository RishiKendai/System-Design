package gamemanager

import (
	"sync"

	"github.com/google/uuid"
	"github.com/lld/chess/board"
	"github.com/lld/chess/game"
	"github.com/lld/chess/player"
)

func NewGameManager() *GameManager {
	return &GameManager{
		Games: make(map[string]*game.Game),
		mu:    sync.RWMutex{},
	}
}

func (gm *GameManager) CreateGame(players []*player.Player) *game.Game {

	board := board.NewBoard()
	gameId := uuid.New().String()
	gm.mu.Lock()
	defer gm.mu.Unlock()
	g := game.NewGame(board, players, 0)
	g.InitializePieces()
	gm.Games[gameId] = g
	return g
}

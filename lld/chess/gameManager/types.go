package gamemanager

import (
	"sync"

	"github.com/lld/chess/game"
)

type GameManager struct {
	Games map[string]*game.Game // map[GameID]*game.Game
	mu    sync.RWMutex
}

package budgetchat

import (
	"sync"
	"sync/atomic"
	"unicode"
)

//
// === CONSTANTS === //
//

// Colours for colourising log messages
const (
	ColourReset  = "\033[0m"
	ColourRed    = "\033[31m"
	ColourGreen  = "\033[32m"
	ColourYellow = "\033[33m"
	ColourBlue   = "\033[34m"
	ColourPurple = "\033[35m"
	ColourCyan   = "\033[36m"
	ColourWhite  = "\033[37m"
)

// Prefix for server log messages
const S_PREFIX = ColourCyan + "[server]" + ColourReset + " "

//
// === STRUCTS === //
//

// IDGenerator represents a thread-safe unique ID generator
type IDGenerator struct {
	lastID int64
	mu     sync.Mutex
}

//
// === METHODS === //
//

// NextID returns a new unique ID
func (g *IDGenerator) NextID() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return atomic.AddInt64(&g.lastID, 1)
}

//
// === FUNCTIONS === //
//

// IsValidUsername checks if the username is valid:
// - must contain at least 1 character
// - must consist entirely of alphanumeric characters
func IsValidUsername(username string) bool {
	if len(username) < 1 || username == "" {
		return false
	}
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}

// NewIDGenerator creates a new instance of IDGenerator
func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

// Colourise returns a colourised string
func Colourise(txt string, colourCode string) string {
	return colourCode + txt + ColourReset
}

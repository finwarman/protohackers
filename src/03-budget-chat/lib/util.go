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
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

// Prefix for server log messages
const S_PREFIX = ColorCyan + "[server]" + ColorReset + " "

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

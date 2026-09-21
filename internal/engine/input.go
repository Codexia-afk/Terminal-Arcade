package engine

import "github.com/gdamore/tcell/v2"

// Action represents an abstract player intent independent of hardware input.
type Action int

const (
	ActionNone Action = iota
	ActionUp
	ActionDown
	ActionLeft
	ActionRight
	ActionPause
	ActionQuit
	ActionConfirm
	ActionTab
)

// Direction converts directional actions to a (dx, dy) coordinate delta.
// Returns false if the action is non-directional.
func (a Action) Direction() (Position, bool) {
	switch a {
	case ActionUp:
		return Position{X: 0, Y: -1}, true
	case ActionDown:
		return Position{X: 0, Y: 1}, true
	case ActionLeft:
		return Position{X: -1, Y: 0}, true
	case ActionRight:
		return Position{X: 1, Y: 0}, true
	default:
		return Position{}, false
	}
}

// TranslateEvent maps a tcell keyboard event to an engine Action.
// Supports standard Arrow keys, WASD movement, Space/Enter confirm, P pause, and Q/Esc quit.
func TranslateEvent(ev *tcell.EventKey) Action {
	switch ev.Key() {
	case tcell.KeyUp:
		return ActionUp
	case tcell.KeyDown:
		return ActionDown
	case tcell.KeyLeft:
		return ActionLeft
	case tcell.KeyRight:
		return ActionRight
	case tcell.KeyEscape, tcell.KeyCtrlC:
		return ActionQuit
	case tcell.KeyEnter:
		return ActionConfirm
	case tcell.KeyTab:
		return ActionTab
	}

	switch ev.Rune() {
	case 'w', 'W':
		return ActionUp
	case 's', 'S':
		return ActionDown
	case 'a', 'A':
		return ActionLeft
	case 'd', 'D':
		return ActionRight
	case 'p', 'P':
		return ActionPause
	case 'q', 'Q':
		return ActionQuit
	case ' ':
		return ActionConfirm
	}

	return ActionNone
}

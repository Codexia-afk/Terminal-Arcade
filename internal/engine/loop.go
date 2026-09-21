package engine

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

// Loop manages game execution with a fixed-interval ticker and non-blocking input handling.
type Loop struct {
	Screen   *Screen
	Input    chan Action
	interval time.Duration
}

// NewLoop constructs a Loop with the default tick interval.
func NewLoop(s *Screen, defaultInterval time.Duration) *Loop {
	if defaultInterval <= 0 {
		defaultInterval = 100 * time.Millisecond
	}
	return &Loop{
		Screen:   s,
		Input:    make(chan Action, 64),
		interval: defaultInterval,
	}
}

// Run executes the game loop until the game concludes, the user quits, or an error occurs.
func (l *Loop) Run(g Game) TickResult {
	interval := l.interval
	if tip, ok := g.(TickIntervalProvider); ok {
		if d := tip.TickInterval(); d > 0 {
			interval = d
		}
	}

	stopPoller := make(chan struct{})
	defer close(stopPoller)

	// Background input poller feeds the buffered input channel
	go func() {
		for {
			select {
			case <-stopPoller:
				return
			default:
			}

			if l.Screen == nil || l.Screen.Raw == nil {
				return
			}

			ev := l.Screen.Raw.PollEvent()
			if ev == nil {
				return
			}

			switch e := ev.(type) {
			case *tcell.EventKey:
				action := TranslateEvent(e)
				select {
				case l.Input <- action:
				default:
					// Input buffer full; drop to maintain real-time responsiveness
				}
			case *tcell.EventResize:
				l.Screen.Raw.Sync()
			}
		}
	}()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	paused := false

	// Initial render
	g.Render(l.Screen)
	l.Screen.Flush()

	for {
		select {
		case a := <-l.Input:
			if a == ActionQuit {
				return TickResult{Continue: false, Reason: "quit"}
			}
			if a == ActionPause {
				paused = !paused
				g.Render(l.Screen)
				if paused {
					l.renderPauseOverlay()
				}
				l.Screen.Flush()
				continue
			}
			if !paused {
				g.HandleInput(a)
				g.Render(l.Screen)
				l.Screen.Flush()
			}

		case <-ticker.C:
			// Non-blocking drain of any queued input before ticking
		drain:
			for {
				select {
				case a := <-l.Input:
					if a == ActionQuit {
						return TickResult{Continue: false, Reason: "quit"}
					}
					if a == ActionPause {
						paused = !paused
						continue
					}
					if !paused {
						g.HandleInput(a)
					}
				default:
					break drain
				}
			}

			if !paused {
				res := g.Tick()
				g.Render(l.Screen)
				l.Screen.Flush()
				if !res.Continue {
					return res
				}
			} else {
				g.Render(l.Screen)
				l.renderPauseOverlay()
				l.Screen.Flush()
			}
		}
	}
}

func (l *Loop) renderPauseOverlay() {
	_, h := l.Screen.Size()
	l.Screen.CenterText(h/2, " [ PAUSED ] - Press P to Resume ", tcell.ColorYellow, tcell.ColorBlack)
}

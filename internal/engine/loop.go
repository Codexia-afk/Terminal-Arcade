package engine

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

// Loop manages game execution with 60 FPS decoupled rendering and 2-step input buffering.
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

// Run executes the game loop with 60 FPS decoupled rendering and 2-step input buffering.
func (l *Loop) Run(g Game) TickResult {
	simInterval := l.interval
	if tip, ok := g.(TickIntervalProvider); ok {
		if d := tip.TickInterval(); d > 0 {
			simInterval = d
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

	// 60 FPS rendering ticker (approx 16.6ms) for silky-smooth animations and interpolation
	renderTicker := time.NewTicker(16 * time.Millisecond)
	defer renderTicker.Stop()

	// Simulation ticker for game logic steps
	simTicker := time.NewTicker(simInterval)
	defer simTicker.Stop()

	paused := false
	dirQueue := make([]Action, 0, 2)

	// Initial render
	g.Render(l.Screen)
	l.Screen.Flush()

	for {
		select {
		case a := <-l.Input:
			switch a {
			case ActionQuit:
				return TickResult{Continue: false, Reason: "quit"}
			case ActionPause:
				paused = !paused
				g.Render(l.Screen)
				if paused {
					l.renderPauseOverlay()
				}
				l.Screen.Flush()
			case ActionUp, ActionDown, ActionLeft, ActionRight:
				if !paused {
					// 2-step directional queue: buffer up to 2 directions
					if len(dirQueue) < 2 {
						dirQueue = append(dirQueue, a)
					} else {
						// Overwrite second direction with newest intent
						dirQueue[1] = a
					}
					// Immediate render on input for instantaneous feedback
					g.Render(l.Screen)
					l.Screen.Flush()
				}
			default:
				if !paused {
					g.HandleInput(a)
					g.Render(l.Screen)
					l.Screen.Flush()
				}
			}

		case <-simTicker.C:
			if !paused {
				// Process next directional action in the 2-step queue
				if len(dirQueue) > 0 {
					nextDir := dirQueue[0]
					dirQueue = dirQueue[1:]
					g.HandleInput(nextDir)
				}

				res := g.Tick()
				g.Render(l.Screen)
				l.Screen.Flush()

				// Drop-and-continue strategy: discard lagged ticks
				for len(simTicker.C) > 0 {
					<-simTicker.C
				}

				if !res.Continue {
					return res
				}
			} else {
				g.Render(l.Screen)
				l.renderPauseOverlay()
				l.Screen.Flush()
			}

		case <-renderTicker.C:
			// Continuous 60 FPS refresh for timer countdowns and visual feedback
			if !paused {
				g.Render(l.Screen)
				l.Screen.Flush()
			}
		}
	}
}

func (l *Loop) renderPauseOverlay() {
	_, h := l.Screen.Size()
	l.Screen.CenterText(h/2, " [ PAUSED ] - Press P to Resume ", tcell.ColorYellow, tcell.ColorBlack)
}

// Package engine supplies the shared terminal game runtime.
package engine

// Difficulty is the suite-wide difficulty selector.
type Difficulty string

const (
	// Easy represents the casual difficulty tier.
	Easy Difficulty = "easy"
	// Medium represents the standard arcade difficulty tier.
	Medium Difficulty = "medium"
	// Hard represents the expert difficulty tier.
	Hard Difficulty = "hard"
)

// Difficulties returns all available difficulty levels in standard order.
func Difficulties() []Difficulty {
	return []Difficulty{Easy, Medium, Hard}
}

// Valid reports whether d is a supported difficulty.
func (d Difficulty) Valid() bool {
	return d == Easy || d == Medium || d == Hard
}

// String returns the user-friendly display name.
func (d Difficulty) String() string {
	switch d {
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	case Hard:
		return "Hard"
	default:
		return "Unknown"
	}
}

// DifficultyProfile identifies a game-specific parameter set.
type DifficultyProfile interface {
	DifficultyLevel() Difficulty
}

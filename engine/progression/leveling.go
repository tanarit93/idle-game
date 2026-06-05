package progression

import (
	"math"
)

// LevelConfig holds the constants for experience calculation
const (
	BaseExp      = 100     // Exp needed for level 1 -> 2
	ExpExponent  = 1.5     // Growth factor
)

// CalculateRequiredExp returns the total exp needed to reach the next level
// Formula: BaseExp * (CurrentLevel ^ ExpExponent)
func CalculateRequiredExp(level int) int64 {
	return int64(float64(BaseExp) * math.Pow(float64(level), ExpExponent))
}

// ProcessExperience adds exp to a character and handles multiple level-ups
func (res *SimulationResult) ProcessExperience(currentLevel int, currentExp int64, expPerKill int) (newLevel int, newExp int64, levelsGained int) {
	totalGained := int64(res.TotalKills * expPerKill)
	newExp = currentExp + totalGained
	newLevel = currentLevel
	levelsGained = 0

	for {
		required := CalculateRequiredExp(newLevel)
		if newExp >= required {
			newExp -= required
			newLevel++
			levelsGained++
		} else {
			break
		}
	}

	return newLevel, newExp, levelsGained
}

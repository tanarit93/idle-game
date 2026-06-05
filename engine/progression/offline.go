package progression

import (
	"math"
)

type EntityStats struct {
	HP    float64
	ATK   float64
	DEF   float64
	ASPD  float64 // Attacks per second
}

type Loot struct {
	TemplateID string
	Tier       int // 0: Normal, 1: Rare, 2: Epic, 3: Legendary, 4: Skill Gem
}

type SimulationResult struct {
	TotalKills int
	Loot       []Loot
	GoldGained int
}

// SimulateOfflineProgress simulates combat over time when the user is away.
func SimulateOfflineProgress(timeElapsedSeconds int, player EntityStats, monster EntityStats) SimulationResult {
	// 1. Calculate Player DPS
	damagePerHit := math.Max(1, player.ATK-monster.DEF)
	playerDPS := damagePerHit * player.ASPD

	// 2. Calculate Time To Kill (TTK)
	ttk := (monster.HP / playerDPS) + 0.5

	// 3. Calculate Total Kills
	totalKills := int(math.Floor(float64(timeElapsedSeconds) / ttk))

	// 4. Generate Loot using the new Weighted Loot Table
	result := SimulationResult{
		TotalKills: totalKills,
	}

	lootTable := NewDefaultLootTable()
	for i := 0; i < totalKills; i++ {
		result.Loot = append(result.Loot, lootTable.Roll())
	}

	// 5. Calculate Gold
	result.GoldGained = CalculateGoldReward(totalKills, 5, 15)

	return result
}

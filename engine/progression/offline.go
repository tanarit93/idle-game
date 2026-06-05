package progression

import (
	"math"
	"math/rand"
	"time"
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
	// Simplified: (Atk - Def) * ASPD. In real engine, would use combat.CalculateFinalDamage for more accuracy.
	damagePerHit := math.Max(1, player.ATK-monster.DEF)
	playerDPS := damagePerHit * player.ASPD

	// 2. Calculate Time To Kill (TTK)
	// We add 0.5s for monster spawn/death animations
	ttk := (monster.HP / playerDPS) + 0.5

	// 3. Calculate Total Kills
	totalKills := int(math.Floor(float64(timeElapsedSeconds) / ttk))

	// 4. Generate Loot
	result := SimulationResult{
		TotalKills: totalKills,
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < totalKills; i++ {
		// Mock loot generation
		tierRoll := r.Intn(100)
		var tier int
		if tierRoll > 95 {
			tier = 4 // Skill Gem
		} else if tierRoll > 85 {
			tier = 3 // Legendary
		} else if tierRoll > 60 {
			tier = 1 // Rare
		} else {
			tier = 0 // Normal
		}

		result.Loot = append(result.Loot, Loot{
			TemplateID: "drop_item_id",
			Tier:       tier,
		})
	}

	return result
}

package progression

import (
	"math/rand"
	"time"
)

type Rarity int

const (
	Common Rarity = iota
	Rare
	Epic
	Legendary
	SkillGem
)

type LootEntry struct {
	TemplateID string
	Tier       Rarity
	Weight     int // Cumulative probability weight
}

type LootTable struct {
	Entries     []LootEntry
	TotalWeight int
}

// NewDefaultLootTable creates a standard drop table for regular monsters
func NewDefaultLootTable() *LootTable {
	table := &LootTable{}
	
	// Define drops with weights
	// Total weight = 10000 (allows for 0.01% precision)
	table.AddEntry("iron_scrap", Common, 7000)    // 70%
	table.AddEntry("steel_ingot", Rare, 2000)    // 20%
	table.AddEntry("magic_dust", Epic, 800)      // 8%
	table.AddEntry("hero_relic", Legendary, 150) // 1.5%
	table.AddEntry("fire_skill_gem", SkillGem, 50) // 0.5%
	
	return table
}

func (lt *LootTable) AddEntry(id string, tier Rarity, weight int) {
	lt.Entries = append(lt.Entries, LootEntry{
		TemplateID: id,
		Tier:       tier,
		Weight:     weight,
	})
	lt.TotalWeight += weight
}

// Roll executes a random drop from the table
func (lt *LootTable) Roll() Loot {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	roll := r.Intn(lt.TotalWeight)
	
	currentWeight := 0
	for _, entry := range lt.Entries {
		currentWeight += entry.Weight
		if roll < currentWeight {
			return Loot{
				TemplateID: entry.TemplateID,
				Tier:       int(entry.Tier),
			}
		}
	}
	
	// Fallback to first entry
	return Loot{TemplateID: lt.Entries[0].TemplateID, Tier: int(lt.Entries[0].Tier)}
}

// Global helper for random gold based on tier
func CalculateGoldReward(totalKills int, minGold, maxGold int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	gold := 0
	for i := 0; i < totalKills; i++ {
		gold += r.Intn(maxGold-minGold+1) + minGold
	}
	return gold
}

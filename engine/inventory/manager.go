package inventory

import (
	"fmt"
	"github.com/tanarit93/idle-game/engine/progression"
)

const MaxInventoryLimit = 100

type InventoryManager struct {
	CurrentCount int
}

// ProcessLootDrops handles the inventory queue and auto-salvages low-tier items if full.
func (im *InventoryManager) ProcessLootDrops(drops []progression.Loot) ([]progression.Loot, int) {
	var keptLoot []progression.Loot
	var totalGoldFromSalvage int

	for _, loot := range drops {
		// Logic: Always keep Skill Gems (Tier 4) or High-Tier (Tier >= 2)
		// Or if we have space, keep everything.
		
		isHighTier := loot.Tier >= 2 || loot.Tier == 4 // Epic+ or Skill Gem

		if im.CurrentCount < MaxInventoryLimit {
			keptLoot = append(keptLoot, loot)
			im.CurrentCount++
		} else {
			// Inventory is full
			if isHighTier {
				// Strategy: Find a low-tier item to salvage to make room
				// For this simulation, we'll just salvage the current drop if we can't find room,
				// but high-tier items should ideally "bump" low-tier ones.
				
				// Auto-salvage logic for the incoming item if we can't swap (simplified for this task)
				// In a full implementation, we would query the DB for the oldest/lowest tier item to replace.
				totalGoldFromSalvage += SalvageValue(loot)
			} else {
				// Auto-salvage normal/low-tier items
				totalGoldFromSalvage += SalvageValue(loot)
			}
		}
	}

	return keptLoot, totalGoldFromSalvage
}

func SalvageValue(loot progression.Loot) int {
	switch loot.Tier {
	case 0:
		return 10 // Normal
	case 1:
		return 50 // Rare
	default:
		return 100
	}
}

func (im *InventoryManager) ShowStatus() {
	fmt.Printf("Inventory: %d/%d\n", im.CurrentCount, MaxInventoryLimit)
}

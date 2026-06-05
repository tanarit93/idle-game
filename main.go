package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/tanarit93/idle-game/engine/combat"
	"github.com/tanarit93/idle-game/engine/progression"
	"github.com/tanarit93/idle-game/engine/inventory"
	"github.com/tanarit93/idle-game/engine/database"
)

type SyncRequest struct {
	CharacterId     string
	ClientTimestamp int64
}

type SyncResponse struct {
	ServerTime time.Time
	Character  *database.CharacterRecord
	Inventory  []*database.ItemRecord
}

func HandleSyncRequest(req SyncRequest) SyncResponse {
	charID, _ := uuid.Parse(req.CharacterId)
	char, err := database.GetCharacter(charID)
	if err != nil {
		fmt.Printf("[Server] Character %s not found. Creating default...\n", req.CharacterId)
		char = &database.CharacterRecord{
			Id: charID, Name: "New Hero", Level: 1, Strength: 10, Agility: 10, Vitality: 10, HP: 100,
			LastSync: time.Now().Unix() - 120, // Simulate 2 mins away
		}
		database.SaveCharacter(char)
	}

	fmt.Printf("\n[Server] Syncing state for: %s (Lv.%d)\n", char.Name, char.Level)

	// 1. Combat Math Demo
	playerDerived := combat.CalculateDerivedStats(combat.Attributes{Strength: char.Strength, Agility: char.Agility, Vitality: char.Vitality}, char.Level)
	monsterStats := combat.DerivedStats{Attack: 10, Defense: 2, MaxHP: 50, AttackSpeed: 1.0}
	
	dmg, isCrit := combat.CalculateFinalDamage(playerDerived, monsterStats, combat.ElementFire, combat.ElementWood)
	critText := ""
	if isCrit {
		critText = " (CRITICAL!)"
	}
	fmt.Printf("[Server] Combat Demo: Dealing %d%s damage (Fire vs Wood)\n", dmg, critText)

	// 2. Simulate Offline Progress
	playerStats := progression.EntityStats{
		HP: float64(char.HP), 
		ATK: playerDerived.Attack, 
		DEF: playerDerived.Defense, 
		ASPD: playerDerived.AttackSpeed,
	}
	mStats := progression.EntityStats{HP: 50, ATK: 10, DEF: 2, ASPD: 1.0}
	
	simResult := progression.SimulateOfflineProgress(int(time.Now().Unix()-char.LastSync), playerStats, mStats)
	fmt.Printf("[Server] Offline Results: %d monsters defeated!\n", simResult.TotalKills)

	// 3. Handle Experience & Level Up
	newLevel, newExp, _ := simResult.ProcessExperience(char.Level, char.Experience, 45)
	char.Level = newLevel
	char.Experience = newExp
	char.LastSync = time.Now().Unix()
	
	// 4. Handle Inventory & Loot
	invManager := &inventory.InventoryManager{CurrentCount: 0} 
	keptLoot, goldGained := invManager.ProcessLootDrops(simResult.Loot)
	char.Gold += goldGained

	// 5. Persist
	database.SaveCharacter(char)
	for _, l := range keptLoot {
		database.SaveItem(char.Id, l)
	}

	items, _ := database.GetInventory(char.Id)
	return SyncResponse{ServerTime: time.Now(), Character: char, Inventory: items}
}

func main() {
	database.InitDB()

	fmt.Println("-------------------------------------------")
	fmt.Println("⚔️  Idle RPG Server-Authoritative Engine ⚔️")
	fmt.Println("-------------------------------------------")

	// Start a listener to keep the container alive
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	fmt.Println("Status: Server Listening on :50051 (Ready for Sync)")

	// Run an initial demo sync
	fixedID := "550e8400-e29b-41d4-a716-446655440000"
	HandleSyncRequest(SyncRequest{CharacterId: fixedID})
	
	// Keep server alive
	for {
		conn, err := lis.Accept()
		if err == nil {
			conn.Close()
		}
		time.Sleep(1 * time.Second)
	}
}

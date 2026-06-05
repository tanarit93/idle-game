package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/idle-game/backend/engine/combat"
	"github.com/idle-game/backend/engine/progression"
	"github.com/idle-game/backend/engine/inventory"
	pb "github.com/idle-game/backend/proto/game" // Assuming generated proto is here
)

type gameServer struct {
	pb.UnimplementedGameServiceServer
	// In a real app, inject DB pool here
}

func (s *gameServer) SyncGameState(ctx context.Context, req *pb.SyncRequest) (*pb.SyncResponse, error) {
	fmt.Printf("Syncing state for character: %s\n", req.CharacterId)

	// 1. Fetch character from DB (Mocked for now)
	// In reality: char := db.GetCharacter(req.CharacterId)
	now := time.Now().Unix()
	timeElapsed := int(now - req.ClientTimestamp)

	// 2. Simulate Offline Progress if time passed
	playerStats := progression.EntityStats{HP: 100, ATK: 20, DEF: 5, ASPD: 1.5}
	monsterStats := progression.EntityStats{HP: 50, ATK: 10, DEF: 2, ASPD: 1.0}
	
	simResult := progression.SimulateOfflineProgress(timeElapsed, playerStats, monsterStats)

	// 3. Handle Inventory & Loot
	invManager := &inventory.InventoryManager{CurrentCount: 20}
	keptLoot, goldGained := invManager.ProcessLootDrops(simResult.Loot)

	// 4. Construct response
	resp := &pb.SyncResponse{
		ServerTime: timestamppb.Now(),
		State: &pb.GameState{
			Character: &pb.Character{
				Id:         req.CharacterId,
				Name:       "Hero",
				Level:      10,
				Experience: 5000,
				Stats: &pb.Stats{
					Strength:     15,
					Agility:      12,
					Intelligence: 10,
					Vitality:     15,
					Attack:       20,
					Defense:      5,
				},
				Resources: &pb.Resources{
					CurrentHp: 80,
					MaxHp:     100,
					CurrentMp: 40,
					MaxMp:     50,
				},
			},
			Inventory: []*pb.Item{}, // Map keptLoot to pb.Item here
		},
	}

	// Add gold gained to character in real DB update
	_ = goldGained 
	
	// Convert keptLoot to protobuf items
	for _, l := range keptLoot {
		resp.State.Inventory = append(resp.State.Inventory, &pb.Item{
			Id:         uuid.New().String(),
			TemplateId: l.TemplateID,
			Level:      1,
		})
	}

	return resp, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGameServiceServer(s, &gameServer{})

	fmt.Println("Server authoritative engine running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

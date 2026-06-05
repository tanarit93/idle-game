package database

import (
	"github.com/google/uuid"
	"github.com/tanarit93/idle-game/engine/progression"
)

type CharacterRecord struct {
	Id         uuid.UUID
	Name       string
	Level      int
	Experience int64
	Gold       int
	Strength   int
	Agility    int
	Vitality   int
	HP         int
	LastSync   int64
}

type ItemRecord struct {
	Id         uuid.UUID
	CharacterId uuid.UUID
	TemplateId string
	Tier       int
	Level      int
}

func GetCharacter(id uuid.UUID) (*CharacterRecord, error) {
	char := &CharacterRecord{}
	err := DB.QueryRow(
		"SELECT id, name, level, experience, gold, strength, agility, vitality, hp, last_sync FROM characters WHERE id = $1",
		id,
	).Scan(&char.Id, &char.Name, &char.Level, &char.Experience, &char.Gold, &char.Strength, &char.Agility, &char.Vitality, &char.HP, &char.LastSync)
	
	if err != nil {
		return nil, err
	}
	return char, nil
}

func SaveCharacter(char *CharacterRecord) error {
	_, err := DB.Exec(
		`INSERT INTO characters (id, name, level, experience, gold, strength, agility, vitality, hp, last_sync)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 ON CONFLICT (id) DO UPDATE SET
		 name = $2, level = $3, experience = $4, gold = $5, strength = $6, agility = $7, vitality = $8, hp = $9, last_sync = $10`,
		char.Id, char.Name, char.Level, char.Experience, char.Gold, char.Strength, char.Agility, char.Vitality, char.HP, char.LastSync,
	)
	return err
}

func SaveItem(charID uuid.UUID, loot progression.Loot) error {
	_, err := DB.Exec(
		"INSERT INTO items (id, character_id, template_id, tier, level) VALUES ($1, $2, $3, $4, $5)",
		uuid.New(), charID, loot.TemplateID, loot.Tier, 1,
	)
	return err
}

func GetInventory(charID uuid.UUID) ([]*ItemRecord, error) {
	rows, err := DB.Query("SELECT id, character_id, template_id, tier, level FROM items WHERE character_id = $1", charID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*ItemRecord
	for rows.Next() {
		item := &ItemRecord{}
		err := rows.Scan(&item.Id, &item.CharacterId, &item.TemplateId, &item.Tier, &item.Level)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

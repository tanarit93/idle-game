package combat

import (
	"math"
	"math/rand"
	"time"
)

type Element int

const (
	ElementNone Element = iota
	ElementFire
	ElementWood
	ElementWater
)

// Attributes represents the raw stats of an entity
type Attributes struct {
	Strength     int
	Agility      int
	Intelligence int
	Vitality     int
}

// DerivedStats are calculated from Attributes
type DerivedStats struct {
	Attack      float64
	Defense     float64
	CritRate    float64 // 0.0 to 1.0
	AttackSpeed float64
	MaxHP       int
}

// CalculateDerivedStats transforms raw attributes into combat stats
func CalculateDerivedStats(attr Attributes, level int) DerivedStats {
	return DerivedStats{
		Attack:      float64(attr.Strength*3 + level*2),
		Defense:     float64(attr.Vitality*2 + level),
		CritRate:    math.Min(0.5, float64(attr.Agility)*0.005), // Max 50% crit
		AttackSpeed: 1.0 + (float64(attr.Agility) * 0.01),      // 1.0 base + 1% per Agility
		MaxHP:       attr.Vitality*20 + level*50,
	}
}

// CalculateFinalDamage computes damage with Crit and Elements
func CalculateFinalDamage(attacker DerivedStats, target DerivedStats, attackerElem, targetElem Element) (int, bool) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// 1. Elemental Multiplier
	elemMult := GetElementalMultiplier(attackerElem, targetElem)
	
	// 2. Base Damage Calculation
	baseDmg := (attacker.Attack - target.Defense) * elemMult
	
	// 3. Critical Hit Check
	isCrit := r.Float64() < attacker.CritRate
	if isCrit {
		baseDmg *= 2.0 // 200% Crit Damage
	}
	
	finalDamage := int(math.Max(1, math.Floor(baseDmg)))
	
	return finalDamage, isCrit
}

func GetElementalMultiplier(attacker, target Element) float64 {
	if attacker == ElementNone || target == ElementNone {
		return 1.0
	}
	// Fire > Wood > Water > Fire
	if (attacker == ElementFire && target == ElementWood) ||
	   (attacker == ElementWood && target == ElementWater) ||
	   (attacker == ElementWater && target == ElementFire) {
		return 1.5
	}
	// Disadvantage
	if (attacker == ElementWood && target == ElementFire) ||
	   (attacker == ElementWater && target == ElementWood) ||
	   (attacker == ElementFire && target == ElementWater) {
		return 0.8
	}
	return 1.0
}

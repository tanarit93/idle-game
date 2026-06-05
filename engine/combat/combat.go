package combat

import (
	"math"
)

type Element int

const (
	ElementNone Element = iota
	ElementFire
	ElementWood
	ElementWater
)

// CalculateFinalDamage computes the authoritative damage dealt.
// Formula: max(1, (Base Attack - Target Defense) * Skill Multiplier * Elemental Multiplier)
func CalculateFinalDamage(attackerAtk, targetDef float64, skillMult float64, attackerElem, targetElem Element) int {
	elemMult := GetElementalMultiplier(attackerElem, targetElem)
	
	rawDamage := (attackerAtk - targetDef) * skillMult * elemMult
	
	finalDamage := int(math.Max(1, math.Floor(rawDamage)))
	
	return finalDamage
}

// GetElementalMultiplier returns the multiplier based on the advantage triangle: Fire > Wood > Water > Fire.
func GetElementalMultiplier(attacker, target Element) float64 {
	if attacker == ElementNone || target == ElementNone {
		return 1.0
	}

	// Fire > Wood
	if attacker == ElementFire && target == ElementWood {
		return 1.5
	}
	// Wood > Water
	if attacker == ElementWood && target == ElementWater {
		return 1.5
	}
	// Water > Fire
	if attacker == ElementWater && target == ElementFire {
		return 1.5
	}

	// Disadvantage cases (0.8)
	if (attacker == ElementWood && target == ElementFire) ||
		(attacker == ElementWater && target == ElementWood) ||
		(attacker == ElementFire && target == ElementWater) {
		return 0.8
	}

	return 1.0
}

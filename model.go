package main

type HealthTracker struct {
	Current int
	Max     int
}

func (s *HealthTracker) Damage(amount int) {
	s.Current -= amount
}

type FortificationType uint

const (
	Wall      = FortificationType(0)
	WallAddon = FortificationType(1)
	Garrison  = FortificationType(2)
	OuterWall = FortificationType(3)

)

type Fortification struct {
	Name            string
	Type            FortificationType
	DefenseModifier int
	AttackModifier  int
}

type Army struct {
	Name        string
	HP          *HealthTracker
	AC          int
	AttackRoll  RollSpec
	DamageRoll  RollSpec
	Location    string
	Destination string
	Allegiance  string
	Destroyed   bool
}

type Settlement struct {
	Name           string
	HP             *HealthTracker
	AttackRoll     RollSpec
	DamageRoll     RollSpec
	HasWarGuard    bool
	Fortifications []Fortification
	Allegiance     string
	Occupied       bool
	Population     uint
}

func (s *Settlement) AC() int {
	var (
		seenTypes    []FortificationType
		settlementAC = 0
	)

	for _, fortification := range s.Fortifications {
		seen := false
		for _, seenType := range seenTypes {
			if fortification.Type == seenType {
				seen = true
				break
			}
		}

		if seen {
			continue
		}

		settlementAC += fortification.DefenseModifier
	}

	return settlementAC
}

func (s *Settlement) RollAttack() (int, int, error) {
	attackModifier := 0
	for _, fortification := range s.Fortifications {
		attackModifier += fortification.AttackModifier
	}

	if attackRoll, err := s.AttackRoll.Roll(); err != nil {
		return 0, 0, err
	} else if damageRoll, err := s.DamageRoll.Roll(); err != nil {
		return 0, 0, err
	} else {
		return attackRoll + attackModifier, damageRoll, nil
	}
}

type WorldActor struct {
	Name string
}

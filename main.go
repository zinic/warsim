package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"bufio"
	"os"
	"time"
)

func parseFacesAndModifier(facesModifierStr string) (int, int, error) {
	const (
		dieFaces    = 0
		dieModifier = 1
	)

	var (
		err      error
		faces    int64 = 0
		modifier int64 = 0
	)

	// Check for additive modifiers first
	splitStr := "+"
	hasModifier := strings.Contains(facesModifierStr, splitStr)

	if !hasModifier {
		// If we didn't find an additive modifier see if there's a subtractive one
		splitStr = "-"
		hasModifier = strings.Contains(facesModifierStr, splitStr)
	}

	if hasModifier {
		if facesModifierParts := strings.Split(facesModifierStr, splitStr); len(facesModifierParts) != 2 {
			return 0, 0, fmt.Errorf("expected two values but got %v", facesModifierParts)
		} else if faces, err = strconv.ParseInt(facesModifierParts[dieFaces], 10, 32); err != nil {
			return 0, 0, err
		} else if modifier, err = strconv.ParseInt(facesModifierParts[dieModifier], 10, 32); err != nil {
			return 0, 0, err
		}

		if splitStr == "-" {
			// Make sure to make the modifier negative if this is a subtractive modifier
			modifier *= -1
		}
	} else if faces, err = strconv.ParseInt(facesModifierStr, 10, 32); err != nil {
		return 0, 0, err
	}

	return int(faces), int(modifier), nil
}

type Die string

func (s Die) String() string {
	return strings.ToLower(strings.TrimSpace(string(s)))
}

func (s Die) Roll() (int, error) {
	const (
		numDicePart          = 0
		facesAndModifierPart = 1
	)

	// Start with the default of one die
	numDice := 1

	if parts := strings.Split(s.String(), "d"); len(parts) != 2 {
		return 0, fmt.Errorf("%s is not a valid roll", s)
	} else {
		if len(parts[numDicePart]) >= 0 {
			if parsedNumDice, err := strconv.ParseInt(parts[numDicePart], 10, 32); err != nil {
				return 0, fmt.Errorf("%s is not a valid roll: %v", s, err)
			} else {
				numDice = int(parsedNumDice)
			}
		}

		if faces, modifier, err := parseFacesAndModifier(parts[facesAndModifierPart]); err != nil {
			return 0, fmt.Errorf("%s is not a valid roll: %v", s, err)
		} else {
			total := 0
			for roll := 0; roll < numDice; roll++ {
				total += rand.Int()%faces + 1 + modifier
			}

			return total, nil
		}
	}
}

type RollSpec []Die

func (s RollSpec) Roll() (int, error) {
	sum := 0
	for _, die := range s {
		if nextRoll, err := die.Roll(); err != nil {
			return 0, err
		} else {
			sum += nextRoll
		}
	}

	return sum, nil
}

type HealthTracker struct {
	Current int
	Max     int
}

func (s *HealthTracker) Damage(amount int) {
	s.Current -= amount
}

type FortificationType uint
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

type World struct {
	TurnID      uint
	Settlements map[string]*Settlement
	Armies      map[string]*Army
	Actors      map[string]*WorldActor
}

func NewWorld() *World {
	return &World{
		Settlements: make(map[string]*Settlement),
		Armies:      make(map[string]*Army),
		Actors:      make(map[string]*WorldActor),
	}
}

func (s *World) Printf(format string, args ...interface{}) {
	if strings.HasSuffix(format, "\n") {
		fmt.Printf(fmt.Sprintf("[%d] %s", s.TurnID, format), args...)
	} else {
		fmt.Printf(fmt.Sprintf("[%d] %s\n", s.TurnID, format), args...)
	}

}

func (s *World) stepArmies() bool {
	var (
		forcesReady      []*Army
		activityObserved = false
	)

	// Any army that moves may not act in the same turn. This is why we track ready armies in the
	// forcesReady array
	for _, army := range s.Armies {
		// Skip destroyed armies
		if army.Destroyed {
			continue
		}

		// If the destination of the army is not equal to the location then the army needs to move
		if army.Destination != army.Location {
			s.Printf("Army %s is travelling from %s to %s.", army.Name, army.Location, army.Destination)
			army.Location = army.Destination

			activityObserved = true
		} else {
			// If no move is required then the army is considered ready for combat
			forcesReady = append(forcesReady, army)
		}
	}

	// Allow armies to attack
	for _, army := range forcesReady {
		// Look up the settlement at the location
		if target, found := s.Settlements[army.Location]; !found {
			panic(fmt.Sprintf("Unable to find location %s that army %s reports being in.", army.Name, army.Location))
		} else if army.Allegiance == target.Allegiance {
			// If this settlement is one of ours now, let's think about what to do next

		} else {
			activityObserved = true

			if armyAttackRoll, err := army.AttackRoll.Roll(); err != nil {
				panic(fmt.Sprintf("Bad roll: %v", err))
			} else if armyAttackRoll > target.AC() {
				// If the army beats the settlement's AC value then roll the damage
				if damage, err := army.DamageRoll.Roll(); err != nil {
					panic(fmt.Sprintf("Bad roll: %v", err))
				} else {
					s.Printf("Army %s attacks settlement %s (AC: %d) with a %d attack roll and %d damage!", army.Name, target.Name, target.AC(), armyAttackRoll, damage)

					// Apply the damage and see if the settlement is overcome
					target.HP.Damage(damage)
					if target.HP.Current <= 0 {
						if target.Occupied {
							// If the settlement was occupied then we're liberating it
							s.Printf("Settlement %s has been liberated by army %s!", target.Name, army.Name)
							target.Occupied = false
							target.Allegiance = army.Allegiance
						} else {
							// If the settlement wasn't occupied then it is now
							s.Printf("Settlement %s has been occupied by army %s!", target.Name, army.Name)
							target.Occupied = true
							target.Allegiance = army.Allegiance
						}
					}
				}
			} else {
				s.Printf("Army %s misses settlement %s (AC: %d) with a(n) %d attack roll!", army.Name, target.Name, target.AC(), armyAttackRoll)
			}
		}
	}

	return activityObserved
}

func (s *World) stepSettlements() bool {
	var activityObserved = false

	for _, settlement := range s.Settlements {
		if !settlement.HasWarGuard {
			// Settlements that have no local war-trained guard may not retaliate against an
			// attacking force
			continue
		}

		if settlement.Occupied {
			// Occupied settlements get no action
			s.Printf("Settlement %s is occupied and gets no action.", settlement.Name)
			continue
		}

		for _, army := range s.Armies {
			// Settlements can only attack the first army they see currently
			if !army.Destroyed && army.Location == settlement.Name {
				activityObserved = true

				if settlementAttackRoll, err := settlement.AttackRoll.Roll(); err != nil {
					panic(fmt.Sprintf("Bad roll: %v", err))
				} else if settlementAttackRoll > army.AC {
					// If the settlement beats the army's AC then roll the damage
					if damage, err := settlement.DamageRoll.Roll(); err != nil {
						panic(fmt.Sprintf("Bad roll: %v", err))
					} else {
						s.Printf("Settlement %s attacks army %s (AC: %d) with a %d attack roll and %d damage!", settlement.Name, army.Name, army.AC, settlementAttackRoll, damage)

						// Apply the damage and see if the army falls apart
						army.HP.Damage(damage)
						if army.HP.Current <= 0 {
							// If the settlement was occupied then we're liberating it
							s.Printf("Settlement %s has destroyed army %s!", settlement.Name, army.Name)
							army.Destroyed = true
						}
					}
				} else {
					s.Printf("Settlement %s misses army %s (AC: %d) with a(n) %d attack roll!", settlement.Name, army.Name, army.AC, settlementAttackRoll)
				}

				break
			}
		}
	}

	return activityObserved
}

func (s *World) Turn() bool {
	// Move armies first
	armiesActive := s.stepArmies()

	// Allow settlements to act last
	settlementsActive := s.stepSettlements()

	// Increment this turn as having past
	s.TurnID++

	// Return whether or not any activity took place this turn
	return armiesActive || settlementsActive
}

const (
	Walls             = FortificationType(0)
	WallFortification = FortificationType(1)
)

var (
	WoodenWalls = Fortification{
		Name:            "Wooden Walls",
		Type:            Walls,
		DefenseModifier: 10,
		AttackModifier:  0,
	}

	StoneWalls = Fortification{
		Name:            "Stone Walls",
		Type:            Walls,
		DefenseModifier: 15,
		AttackModifier:  0,
	}
)

func main() {
	rand.Seed(time.Now().Unix())

	writeSeigeOfThraneBase()
}

func writeSeigeOfThraneBase() {
	var (
		AundairSettlementNames = []string{
			"Morningcrest",
			"Fort Light",
			"Rellekor",
			"Tellyn",
		}

		AundairArmies = []string{
			"First Cog",
			"Second Cog",
			"Third Cog",
			"Fourth Cog",
			"Fifth Cog",
			"First Gear",
			"Second Gear",
			"First Chain",
			"Second Chain",
			"The Cinch",
			"The Hammer",
			"The Blade",
		}

		ThraneSettlementNames = []string{
			"Daskaran",
			"Thaliost",
			"Silvercliff Castle",
			"Auxylgard",
			"Flamekeep",
			"Danthaven",
			"Athandra",
			"Traelyn",
			"Avaroth",
			"Sharavacion",
			"Shadukar",
			"Olath",
			"Angwar Keep",
			"Aelyndar",
			"Valiron",
			"Siyar",
			"Sigilstar",
			"Lessyk",
			"Nathyrr",
			"The Thornwood",
			"Arythawn Keep",
		}

		ThraneArmies = []string{
			"First Host",
			"Second Host",
			"Third Host",
			"Fourth Host",
			"Fifth Host",
			"Sixth Host",
			"First Surgeons",
			"Second Surgeons",
			"Lightbringers",
			"Demon's Bane",
			"Truthspeakers",
		}
	)
	world := NewWorld()
	world.Actors["Aundair"] = &WorldActor{
		Name: "Aundair",
	}
	world.Actors["Thrane"] = &WorldActor{
		Name: "Thrane",
	}

	for _, name := range AundairSettlementNames {
		world.Settlements[name] = &Settlement{
			Name:        name,
			Allegiance:  "Aundair",
			Occupied:    true,
			Population:  0,
			HasWarGuard: true,

			HP: &HealthTracker{
				Current: 100,
				Max:     100,
			},
			AttackRoll: RollSpec{"d20"},
			DamageRoll: RollSpec{"d20"},
			Fortifications: []Fortification{
				WoodenWalls,
			},
		}
	}

	for _, name := range ThraneSettlementNames {
		world.Settlements[name] = &Settlement{
			Name:        name,
			Allegiance:  "Thrane",
			Occupied:    false,
			Population:  0,
			HasWarGuard: true,

			HP: &HealthTracker{
				Current: 100,
				Max:     100,
			},
			AttackRoll: RollSpec{"d20"},
			DamageRoll: RollSpec{"d4"},
			Fortifications: []Fortification{
				WoodenWalls,
			},
		}
	}

	for _, name := range AundairArmies {
		world.Armies[name] = &Army{
			Name:       name,
			Allegiance: "Aundair",
			AC:         19,
			Destroyed:  false,

			HP: &HealthTracker{
				Current: 100,
				Max:     100,
			},
			AttackRoll:  RollSpec{"d20"},
			DamageRoll:  RollSpec{"d8"},
			Location:    "UNSET_LOCATION",
			Destination: "UNSET_LOCATION",
		}
	}

	for _, name := range ThraneArmies {
		world.Armies[name] = &Army{
			Name:       name,
			Allegiance: "Thrane",
			AC:         19,
			Destroyed:  false,

			HP: &HealthTracker{
				Current: 100,
				Max:     100,
			},
			AttackRoll:  RollSpec{"d20"},
			DamageRoll:  RollSpec{"d8"},
			Location:    "UNSET_LOCATION",
			Destination: "UNSET_LOCATION",
		}
	}

	stdinC := make(chan string)
	reader := bufio.NewReader(os.Stdin)

	go func() {
		for {
			if text, err := reader.ReadString('\n'); err != nil {
				os.Exit(0)
			} else {
				stdinC <- text
			}
		}
	}()

	for range stdinC {
		if !world.Turn() {
			break
		}
	}

	if err := WriteWorld("siege_of_thrane.toml", world); err != nil {
		fmt.Printf("Error writing world: %v.", err)
	}
}

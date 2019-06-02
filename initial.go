package main

import (
	"fmt"
)

var (
	WoodenWalls = Fortification{
		Name:            "Wooden Walls",
		Type:            Wall,
		DefenseModifier: 10,
		AttackModifier:  1,
	}

	StoneWalls = Fortification{
		Name:            "Stone Walls",
		Type:            Wall,
		DefenseModifier: 15,
		AttackModifier:  3,
	}
)

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

	if err := WriteWorld("siege_of_thrane.toml", world); err != nil {
		fmt.Printf("Error writing world: %v.", err)
	}
}

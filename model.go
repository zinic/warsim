package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path"
	"strings"
	"sort"
)

const (
	stateFilename = "state.toml"
)

type Simulation struct {
	Name string
	Step int

	StateDir string `toml:"-"`
	World    *World `toml:"-"`
}

func LoadSimulation(stateDir string) (*Simulation, error) {
	state := &Simulation{
		StateDir: stateDir,
	}

	if _, err := toml.DecodeFile(state.StatePath(), state); err != nil {
		return nil, err
	}

	// Load the world state file based on which step we left off at
	worldPath := state.WorldPath()
	fmt.Printf("Loading world %s\n", worldPath)

	if loadedWorld, err := LoadWorld(worldPath); err != nil {
		return nil, err
	} else {
		state.World = loadedWorld
	}

	return state, nil
}

func (s *Simulation) Write() error {
	if file, err := os.OpenFile(s.StatePath(), os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644); err != nil {
		return err
	} else if err := toml.NewEncoder(file).Encode(s); err != nil {
		return err
	}

	return nil
}

func (s *Simulation) StatePath() string {
	return path.Join(s.StateDir, stateFilename)
}

func (s *Simulation) WorldPath() string {
	sanitizedName := strings.Replace(strings.TrimSpace(strings.ToLower(s.Name)), " ", "_", -1)
	worldFilename := fmt.Sprintf("%s.%d.toml", sanitizedName, s.Step)

	return path.Join(s.StateDir, worldFilename)
}

func (s *Simulation) Turn() error {
	output := &strings.Builder{}

	// Increase our step counter then step the simulation
	s.Step++
	s.World.Turn(s.Step, output)

	// Commit a new version of this world
	if err := WriteWorld(s.WorldPath(), s.World); err != nil {
		return err
	}

	if file, err := os.OpenFile("rendered.html", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644); err != nil {
		return err
	} else if _, err := file.WriteString(fmt.Sprintf("<html><body>%s</body></html>\n", output)); err != nil {
		return err
	}

	// Commit our current state
	if err := s.Write(); err != nil {
		return err
	}

	return nil
}

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

type ArmyList []*Army

func (s ArmyList) Sorted() ArmyList {
	sorted := s
	sort.Sort(sorted)

	return sorted
}

func (s ArmyList) Len() int {
	return len(s)
}

func (s ArmyList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ArmyList) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

type Settlement struct {
	Name           string
	HP             *HealthTracker
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

func (s *Settlement) attackModifier() int {
	attackModifier := 0
	for _, fortification := range s.Fortifications {
		attackModifier += fortification.AttackModifier
	}

	return attackModifier
}

func (s *Settlement) RollAttack() (int, int, error) {
	if attackRoll, err := D20.Roll(); err != nil {
		return 0, 0, err
	} else if damageRoll, err := s.DamageRoll.Roll(); err != nil {
		return 0, 0, err
	} else {
		return attackRoll + s.attackModifier(), damageRoll, nil
	}
}

func (s *Settlement) AttackRoll() Die {
	attackMod := s.attackModifier()

	if attackMod > 0 {
		return Die(fmt.Sprintf("d20+%d", attackMod))
	} else if attackMod < 0 {
		return Die(fmt.Sprintf("d20-%d", attackMod))
	}

	return D20
}

type SettlementList []*Settlement

func SettlementListFromMap(source map[string]*Settlement) SettlementList {
	var settlements SettlementList
	for _, settlement := range source {
		settlements = append(settlements, settlement)
	}

	return settlements
}

func (s SettlementList) Sorted() SettlementList {
	sorted := s
	sort.Sort(sorted)

	return sorted
}

func (s SettlementList) Len() int {
	return len(s)
}

func (s SettlementList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SettlementList) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

type WorldActor struct {
	Name string
}

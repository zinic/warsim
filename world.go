package main

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var printer = message.NewPrinter(language.English)

func Printf(format string, args ...interface{}) {
	if strings.HasSuffix(format, "\n") {
		printer.Printf(format, args...)
	} else {
		printer.Printf(fmt.Sprintf("%s\n", format), args...)
	}
}

type World struct {
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

func (s *World) SortedActors() []*WorldActor {
	var (
		sortedNames  []string
		sortedActors []*WorldActor
	)

	for name := range s.Actors {
		sortedNames = append(sortedNames, name)
	}

	sort.Strings(sortedNames)
	for _, name := range sortedNames {
		sortedActors = append(sortedActors, s.Actors[name])
	}

	return sortedActors
}

func (s *World) SettlementsByActor() map[string]SettlementList {
	settlementMap := make(map[string]SettlementList)
	for _, settlement := range s.Settlements {
		if settlements, found := settlementMap[settlement.Allegiance]; !found {
			settlementMap[settlement.Allegiance] = SettlementList{settlement}
		} else {
			settlementMap[settlement.Allegiance] = append(settlements, settlement)
		}
	}

	return settlementMap
}

func DocumentID(name string) string {
	return strings.Replace(strings.ToLower(strings.TrimSpace(name)), " ", "_", -1)
}

func (s *World) ArmiesAt(location string) ArmyList {
	var armies ArmyList
	for _, army := range s.Armies {
		if army.Location == location {
			armies = append(armies, army)
		}
	}

	return armies
}

func (s *World) ArmiesByActor() map[string]ArmyList {
	armyMap := make(map[string]ArmyList)
	for _, army := range s.Armies {
		if armies, found := armyMap[army.Allegiance]; !found {
			armyMap[army.Allegiance] = ArmyList{army}
		} else {
			armyMap[army.Allegiance] = append(armies, army)
		}
	}

	return armyMap
}

func (s *World) Markdown(upstreamOutput *strings.Builder) {
	html := NewElement(HTRoot)
	body := html.Element(HTBody)
	rootDiv := body.Element(HTDiv)

	rootDiv.Element(HTH1).Text = "Animus Warsim"

	sortedActors := s.SortedActors()
	settlementsByActor := s.SettlementsByActor()

	settlementsDiv := rootDiv.Element(HTDiv)
	for _, actor := range sortedActors {
		settlementsDiv.Element(HTH2).Text = fmt.Sprintf("%s Occupied Settlements", actor.Name)

		settlementList := settlementsDiv.Element(HTUL)
		for _, settlement := range settlementsByActor[actor.Name].Sorted() {
			settlementLink := settlementList.Element(HTLI).Element(HTA)
			settlementLink.Attributes["href"] = fmt.Sprintf("#%s", DocumentID(settlement.Name))
			settlementLink.Text = settlement.Name
		}
	}

	armiesDiv := rootDiv.Element(HTDiv)
	armiesByActor := s.ArmiesByActor()
	for _, actor := range sortedActors {
		armiesDiv.Element(HTH2).Text = fmt.Sprintf("%s Armies", actor.Name)

		armyList := armiesDiv.Element(HTUL)
		for _, army := range armiesByActor[actor.Name].Sorted() {
			armyLink := armyList.Element(HTLI).Element(HTA)
			armyLink.Attributes["href"] = fmt.Sprintf("#%s", DocumentID(army.Name))
			armyLink.Text = army.Name
		}
	}

	for i := 0; i < 10; i++ {
		rootDiv.Element(HTBR)
	}

	settlementDetailsDiv := rootDiv.Element(HTDiv)
	settlementDetailsDiv.Element(HTH1).Text = "Settlement Details"
	for _, settlement := range SettlementListFromMap(s.Settlements).Sorted() {
		settlementDiv := settlementDetailsDiv.Element(HTDiv)
		settlementDiv.Element(SPAN).Attributes["id"] = DocumentID(settlement.Name)
		settlementDiv.Element(HTH2).Text = settlement.Name

		table := settlementDiv.Element(TABLE)
		table.Element(TR).Do(func(row *Element) {
			fieldName := row.Element(TD).Element(SPAN)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "Population"

			row.Element(TD).Element(SPAN).Text = printer.Sprint(settlement.Population)
		})

		table.Element(TR).Do(func(row *Element) {
			fieldName := row.Element(TD).Element(SPAN)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "HP"

			row.Element(TD).Element(SPAN).Text = printer.Sprintf("%d / %d", settlement.HP.Current, settlement.HP.Max)
		})

		table.Element(TR).Do(func(row *Element) {
			fieldName := row.Element(TD).Element(SPAN)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "AC"

			row.Element(TD).Element(SPAN).Text = printer.Sprint(settlement.AC())
		})

		table.Element(TR).Do(func(row *Element) {
			fieldName := row.Element(TD).Element(SPAN)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "Attack Roll"

			row.Element(TD).Element(SPAN).Text = printer.Sprint(settlement.AttackRoll())
		})

		table.Element(TR).Do(func(row *Element) {
			fieldName := row.Element(TD).Element(SPAN)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "Attack Damage"

			row.Element(TD).Element(SPAN).Text = printer.Sprint(settlement.DamageRoll)
		})

		settlementDiv.Element(HTH3).Text = "Settlement Fortifications"
		fortificationList := settlementDiv.Element(HTUL)

		for _, fortification := range settlement.Fortifications {
			listItem := fortificationList.Element(HTLI)
			listItem.Attributes["style"] = "padding-top: 10px;"
			listItem.Element(SPAN).Text = fortification.Name

			table = listItem.Element(TABLE)
			table.Attributes["style"] = "margin-left: 20px;"

			table.Element(TR).Do(func(row *Element) {
				fieldName := row.Element(TD).Element(SPAN)
				fieldName.Attributes["style"] = "font-weight: bold;"
				fieldName.Text = "Defense Modifier"

				row.Element(TD).Element(SPAN).Text = printer.Sprint(fortification.DefenseModifier)
			})

			table.Element(TR).Do(func(row *Element) {
				fieldName := row.Element(TD).Element(SPAN)
				fieldName.Attributes["style"] = "font-weight: bold;"
				fieldName.Text = "Attack Modifier"

				row.Element(TD).Element(SPAN).Text = printer.Sprint(fortification.AttackModifier)
			})
		}
	}

	armyDetailsDiv := rootDiv.Element(HTDiv)
	armyDetailsDiv.Element(HTH1).Text = "Army Details"

	for _, army := range s.Armies {
		armyDiv := settlementDetailsDiv.Element(HTDiv)
		armyDiv.Element(SPAN).Attributes["id"] = DocumentID(army.Name)
		armyDiv.Element(HTH2).Text = army.Name

		table := armyDiv.Element(TABLE)
		table.Element(TR).Do(func(row *Element) {
			fieldName := row.Element(TD).Element(SPAN)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "Location"

			row.Element(TD).Element(SPAN).Text = printer.Sprint(army.Location)
		})

		table.Element(TR).Do(func(row *Element) {
			fieldName := row.Element(TD).Element(SPAN)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "HP"

			row.Element(TD).Element(SPAN).Text = printer.Sprintf("%d / %d", army.HP.Current, army.HP.Max)
		})

		table.Element(TR).Do(func(row *Element) {
			fieldName := row.Element(TD).Element(SPAN)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "AC"

			row.Element(TD).Element(SPAN).Text = printer.Sprint(army.AC)
		})

		table.Element(TR).Do(func(row *Element) {
			fieldName := row.Element(TD).Element(SPAN)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "Attack Roll"

			row.Element(TD).Element(SPAN).Text = army.AttackRoll.String()
		})

		table.Element(TR).Do(func(row *Element) {
			fieldName := row.Element(TD).Element(SPAN)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "Attack Damage"

			row.Element(TD).Element(SPAN).Text = printer.Sprint(army.DamageRoll)
		})
	}

	html.Output(upstreamOutput)
}

func (s *World) stepArmies(output *strings.Builder) bool {
	var (
		forcesReady      ArmyList
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
			output.WriteString(fmt.Sprintf("Army %s is travelling from %s to %s.", army.Name, army.Location, army.Destination))
			army.Location = army.Destination

			activityObserved = true
		} else {
			// If no move is required then the army is considered ready for combat
			forcesReady = append(forcesReady, army)
		}
	}

	// Allow armies to attack
	for _, army := range forcesReady.Sorted() {
		attackedOtherArmy := false

		// Check to see if there are any armies in our location first
		for actorName, otherArmies := range s.ArmiesByActor() {
			if army.Allegiance != actorName {
				// Pick the first army in our location
				for _, target := range otherArmies {
					if army.Location != target.Location {
						continue
					}

					activityObserved = true
					attackedOtherArmy = true

					if armyAttackRoll, err := army.AttackRoll.Roll(); err != nil {
						panic(fmt.Sprintf("Bad roll: %v", err))
					} else if armyAttackRoll > target.AC {
						// If the army beats the settlement's AC value then roll the damage
						if damage, err := army.DamageRoll.Roll(); err != nil {
							panic(fmt.Sprintf("Bad roll: %v", err))
						} else {
							output.WriteString(fmt.Sprintf("Army <span style=\"font-weight: bold;\">%s</span> attacks army <span style=\"font-weight: bold;\">%s</span> (AC: %d) with a %d attack roll and %d damage!<br />\n", army.Name, target.Name, target.AC, armyAttackRoll, damage))

							// Apply the damage and see if the settlement is overcome
							target.HP.Damage(damage)
							if target.HP.Current <= 0 {
								output.WriteString(fmt.Sprintf("Army <span style=\"font-weight: bold;\">%s</span> has destroyed army <span style=\"font-weight: bold;\">%s</span>!<br />\n", army.Name, target.Name))
								army.Destroyed = true
							}
						}
					} else {
						output.WriteString(fmt.Sprintf("Army <span style=\"font-weight: bold;\">%s</span> misses army <span style=\"font-weight: bold;\">%s</span> (AC: %d) with a(n) %d attack roll!<br />\n", army.Name, target.Name, target.AC, armyAttackRoll))
					}

					break
				}
			}
		}

		if attackedOtherArmy {
			continue
		}

		// Look up the settlement at the location
		if target, found := s.Settlements[army.Location]; !found {
			panic(fmt.Sprintf("Unable to find location %s that army %s reports being in.", army.Location, army.Name))
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
					output.WriteString(fmt.Sprintf("Army <span style=\"font-weight: bold;\">%s</span> attacks settlement <span style=\"font-weight: bold;\">%s</span> (AC: %d) with a %d attack roll and %d damage!<br />\n", army.Name, target.Name, target.AC(), armyAttackRoll, damage))

					// Apply the damage and see if the settlement is overcome
					target.HP.Damage(damage)
					if target.HP.Current <= 0 {
						if target.Occupied {
							// If the settlement was occupied then we're liberating it
							output.WriteString(fmt.Sprintf("Settlement <span style=\"font-weight: bold;\">%s</span> has been liberated by army <span style=\"font-weight: bold;\">%s</span>!<br />\n", target.Name, army.Name))
							target.Occupied = false
							target.Allegiance = army.Allegiance
						} else {
							// If the settlement wasn't occupied then it is now
							output.WriteString(fmt.Sprintf("Settlement <span style=\"font-weight: bold;\">%s</span> has been occupied by army <span style=\"font-weight: bold;\">%s</span>!<br />\n", target.Name, army.Name))
							target.Occupied = true
							target.Allegiance = army.Allegiance
						}
					}
				}
			} else {
				output.WriteString(fmt.Sprintf("Army <span style=\"font-weight: bold;\">%s</span> misses settlement <span style=\"font-weight: bold;\">%s</span> (AC: %d) with a(n) %d attack roll!<br />\n", army.Name, target.Name, target.AC(), armyAttackRoll))
			}
		}
	}

	return activityObserved
}

func (s *World) stepSettlements(output *strings.Builder) bool {
	var activityObserved = false

	for _, settlement := range SettlementListFromMap(s.Settlements).Sorted() {
		if !settlement.HasWarGuard {
			// Settlements that have no local war-trained guard may not retaliate against an
			// attacking force
			continue
		}

		if settlement.Occupied {
			// Occupied Settlements get no action
			output.WriteString(fmt.Sprintf("Settlement <span style=\"font-weight: bold;\">%s is occupied</span> and gets no action.<br />\n", settlement.Name))
			continue
		}

		for _, army := range s.Armies {
			// Settlements can only attack the first army they see currently
			if !army.Destroyed && army.Location == settlement.Name && army.Allegiance != settlement.Allegiance {
				activityObserved = true

				if settlementAttackRoll, damage, err := settlement.RollAttack(); err != nil {
					panic(fmt.Sprintf("Bad roll: %v", err))
				} else if settlementAttackRoll > army.AC {
					// If the settlement beats the army's AC then roll the damage
					output.WriteString(fmt.Sprintf("Settlement <span style=\"font-weight: bold;\">%s</span> attacks army <span style=\"font-weight: bold;\">%s</span> (AC: %d) with a %d attack roll and %d damage!<br />\n", settlement.Name, army.Name, army.AC, settlementAttackRoll, damage))

					// Apply the damage and see if the army falls apart
					army.HP.Damage(damage)
					if army.HP.Current <= 0 {
						output.WriteString(fmt.Sprintf("Settlement <span style=\"font-weight: bold;\">%s</span> has destroyed army <span style=\"font-weight: bold;\">%s</span>!<br />\n", settlement.Name, army.Name))
						army.Destroyed = true
					}
				} else {
					output.WriteString(fmt.Sprintf("Settlement <span style=\"font-weight: bold;\">%s</span> misses army <span style=\"font-weight: bold;\">%s</span> (AC: %d) with a(n) %d attack roll!<br />\n", settlement.Name, army.Name, army.AC, settlementAttackRoll))
				}

				break
			}
		}
	}

	return activityObserved
}

func (s *World) Turn(turnID int, output *strings.Builder) bool {
	output.WriteString("<div style=\"float: right; background: #E0E0E0;\">\n")
	output.WriteString(fmt.Sprintf("<h1>Combat Log for Turn %d</h1>\n", turnID))

	// Move armies first
	armiesActive := s.stepArmies(output)

	output.WriteString("<br />\n")

	// Allow Settlements to act last
	settlementsActive := s.stepSettlements(output)

	output.WriteString("</p></div>\n")

	// Dump the world state
	s.Markdown(output)

	// Return whether or not any activity took place this turn
	return armiesActive || settlementsActive
}

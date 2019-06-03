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

func (s *World) WriteWorld(body *DocumentElement) {
	rootDiv := body.Element(Division)
	rootDiv.Element(H1).Text = "Animus Warsim"

	sortedActors := s.SortedActors()
	settlementsByActor := s.SettlementsByActor()

	settlementsDiv := rootDiv.Element(Division)
	for _, actor := range sortedActors {
		settlementsDiv.Element(H2).Text = fmt.Sprintf("%s Occupied Settlements", actor.Name)

		settlementList := settlementsDiv.Element(UnorderedList)
		for _, settlement := range settlementsByActor[actor.Name].Sorted() {
			settlementLink := settlementList.Element(ListItem).Element(Anchor)
			settlementLink.Attributes["href"] = fmt.Sprintf("#%s", DocumentID(settlement.Name))
			settlementLink.Text = settlement.Name
		}
	}

	armiesDiv := rootDiv.Element(Division)
	armiesByActor := s.ArmiesByActor()
	for _, actor := range sortedActors {
		armiesDiv.Element(H2).Text = fmt.Sprintf("%s Armies", actor.Name)

		armyList := armiesDiv.Element(UnorderedList)
		for _, army := range armiesByActor[actor.Name].Sorted() {
			armyLink := armyList.Element(ListItem).Element(Anchor)
			armyLink.Attributes["href"] = fmt.Sprintf("#%s", DocumentID(army.Name))
			armyLink.Text = army.Name
		}
	}

	for i := 0; i < 4; i++ {
		rootDiv.Element(BR)
	}

	settlementDetailsDiv := rootDiv.Element(Division)
	settlementDetailsDiv.Element(H1).Text = "Settlement Details"

	for _, settlement := range SettlementListFromMap(s.Settlements).Sorted() {
		settlementDiv := settlementDetailsDiv.Element(Division)
		settlementDiv.Attributes["style"] = "display: inline-block; background: #EFEFEF; margin-bottom: 30px; padding-left: 15px; padding-right: 15px; padding-top: 1px; border: solid 1px black;"
		settlementDiv.Element(Span).Attributes["id"] = DocumentID(settlement.Name)
		settlementDiv.Element(H2).Text = settlement.Name

		settlementDetailsTable := settlementDiv.Element(Table)
		statsTableHeadersRow := settlementDetailsTable.Element(TableHeaders).Element(TableRow)
		statsTableHeadersRow.Element(TableCell).Text = NewElement(Span, "Statistics", ElementAttributes{
			"style": "font-weight: bold; font-size: 14pt;",
		}).String()

		statsTableHeadersRow.Element(TableCell).Do(func(cell *DocumentElement) {
			cell.Text = NewElement(Span, "Fortifications", ElementAttributes{
				"style": "font-weight: bold; font-size: 14pt;",
			}).String()

			cell.Attributes["style"] = "padding-left: 15px;"
		})

		statsRow := settlementDetailsTable.Element(TableRow)
		detailsCell := statsRow.Element(TableCell)
		detailsCell.Attributes["style"] = "vertical-align: top;"

		statsTable := detailsCell.Element(Table)
		statsTable.Element(TableRow).Do(func(row *DocumentElement) {
			statsCell := row.Element(TableCell)
			statsCell.Attributes["style"] = "padding-right: 30px;"

			fieldName := statsCell.Element(Span)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "Population"

			row.Element(TableCell).Element(Span).Text = printer.Sprint(settlement.Population)
		})

		statsTable.Element(TableRow).Do(func(row *DocumentElement) {
			statsCell := row.Element(TableCell)
			statsCell.Attributes["style"] = "padding-right: 30px;"

			fieldName := statsCell.Element(Span)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "Allegiance"

			row.Element(TableCell).Element(Span).Text = printer.Sprint(settlement.Allegiance)
		})

		statsTable.Element(TableRow).Do(func(row *DocumentElement) {
			statsCell := row.Element(TableCell)
			statsCell.Attributes["style"] = "padding-right: 30px;"

			fieldName := statsCell.Element(Span)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "HP"

			row.Element(TableCell).Element(Span).Text = printer.Sprintf("%d / %d", settlement.HP.Current, settlement.HP.Max)
		})

		statsTable.Element(TableRow).Do(func(row *DocumentElement) {
			statsCell := row.Element(TableCell)
			statsCell.Attributes["style"] = "padding-right: 30px;"

			fieldName := statsCell.Element(Span)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "AC"

			row.Element(TableCell).Element(Span).Text = printer.Sprint(settlement.AC())
		})

		statsTable.Element(TableRow).Do(func(row *DocumentElement) {
			statsCell := row.Element(TableCell)
			statsCell.Attributes["style"] = "padding-right: 30px;"

			fieldName := statsCell.Element(Span)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "Attack Roll"

			row.Element(TableCell).Element(Span).Text = printer.Sprint(settlement.AttackRoll())
		})

		statsTable.Element(TableRow).Do(func(row *DocumentElement) {
			statsCell := row.Element(TableCell)
			statsCell.Attributes["style"] = "padding-right: 30px;"

			fieldName := statsCell.Element(Span)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "Attack Damage"

			row.Element(TableCell).Element(Span).Text = printer.Sprint(settlement.DamageRoll)
		})

		detailsCell = statsRow.Element(TableCell)
		fortificationList := detailsCell.Element(UnorderedList)
		fortificationList.Attributes["style"] = "list-style-type: none;"

		for _, fortification := range settlement.Fortifications {
			listItem := fortificationList.Element(ListItem)
			listItem.Attributes["style"] = "padding-top: 10px;"
			listItem.Element(Span).Text = fortification.Name

			statsTable = listItem.Element(Table)
			statsTable.Attributes["style"] = "margin-left: 20px;"

			statsTable.Element(TableRow).Do(func(row *DocumentElement) {
				fieldName := row.Element(TableCell).Element(Span)
				fieldName.Attributes["style"] = "font-weight: bold;"
				fieldName.Text = "Defense Modifier"

				row.Element(TableCell).Element(Span).Text = printer.Sprint(fortification.DefenseModifier)
			})

			statsTable.Element(TableRow).Do(func(row *DocumentElement) {
				fieldName := row.Element(TableCell).Element(Span)
				fieldName.Attributes["style"] = "font-weight: bold;"
				fieldName.Text = "Attack Modifier"

				row.Element(TableCell).Element(Span).Text = printer.Sprint(fortification.AttackModifier)
			})
		}

		settlementDetailsDiv.Element(Division).Attributes["style"] = "display: block;"
	}

	armyDetailsDiv := rootDiv.Element(Division)
	armyDetailsDiv.Element(H1).Text = "Army Details"

	for _, army := range s.Armies {
		armyDiv := armyDetailsDiv.Element(Division)
		armyDiv.Element(Span).Attributes["id"] = DocumentID(army.Name)
		armyDiv.Element(H2).Text = army.Name

		table := armyDiv.Element(Table)
		table.Element(TableRow).Do(func(row *DocumentElement) {
			fieldName := row.Element(TableCell).Element(Span)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "Location"

			row.Element(TableCell).Element(Span).Text = printer.Sprint(army.Location)
		})

		table.Element(TableRow).Do(func(row *DocumentElement) {
			fieldName := row.Element(TableCell).Element(Span)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "Allegiance"

			row.Element(TableCell).Element(Span).Text = army.Allegiance
		})

		table.Element(TableRow).Do(func(row *DocumentElement) {
			fieldName := row.Element(TableCell).Element(Span)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "HP"

			row.Element(TableCell).Element(Span).Text = printer.Sprintf("%d / %d", army.HP.Current, army.HP.Max)
		})

		table.Element(TableRow).Do(func(row *DocumentElement) {
			fieldName := row.Element(TableCell).Element(Span)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "AC"

			row.Element(TableCell).Element(Span).Text = printer.Sprint(army.AC)
		})

		table.Element(TableRow).Do(func(row *DocumentElement) {
			fieldName := row.Element(TableCell).Element(Span)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "Attack Roll"

			row.Element(TableCell).Element(Span).Text = army.AttackRoll.String()
		})

		table.Element(TableRow).Do(func(row *DocumentElement) {
			fieldName := row.Element(TableCell).Element(Span)
			fieldName.Attributes["style"] = "font-weight: bold;"
			fieldName.Text = "Attack Damage"

			row.Element(TableCell).Element(Span).Text = printer.Sprint(army.DamageRoll)
		})
	}
}

func (s *World) stepArmies(log *DocumentElement) bool {
	var (
		forcesReady      ArmyList
		activityObserved = false
	)

	actionList := log.Element(UnorderedList)
	actionList.Attributes["style"] = "list-style-type: none;"

	// Any army that moves may not act in the same turn. This is why we track ready armies in the
	// forcesReady array
	for _, army := range s.Armies {
		// Skip destroyed armies
		if army.Destroyed {
			continue
		}

		// If the destination of the army is not equal to the location then the army needs to move
		if army.Destination != army.Location {
			actionList.Element(ListItem).Text = fmt.Sprintf("Army %s is travelling from %s to %s.", NameLink(army.Name), army.Location, army.Destination)

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
					if target.Destroyed || army.Location != target.Location {
						continue
					}

					activityObserved = true
					attackedOtherArmy = true

					if armyAttackRoll, err := army.AttackRoll.Roll(); err != nil {
						panic(fmt.Sprintf("Bad roll: %v", err))
					} else if armyAttackRoll >= target.AC {
						// If the army beats the settlement's AC value then roll the damage
						if damage, err := army.DamageRoll.Roll(); err != nil {
							panic(fmt.Sprintf("Bad roll: %v", err))
						} else {
							actionList.Element(ListItem).Text = fmt.Sprintf("Army %s attacks army %s (AC: %d) rolling a %d for attack and %d for damage!",
								NameLink(army.Name), NameLink(target.Name), target.AC, armyAttackRoll, damage)

							// Apply the damage and see if the settlement is overcome
							if target.Damage(damage); target.Destroyed {
								actionList.Element(ListItem).Text = fmt.Sprintf("Army %s has destroyed army %s!", NameLink(army.Name), NameLink(target.Name))
							}
						}
					} else {
						actionList.Element(ListItem).Text = fmt.Sprintf("Army %s misses army %s (AC: %d) rolling a %d for attack.",
							NameLink(army.Name), NameLink(target.Name), target.AC, armyAttackRoll)
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
			} else if armyAttackRoll >= target.AC() {
				// If the army beats the settlement's AC value then roll the damage
				if damage, err := army.DamageRoll.Roll(); err != nil {
					panic(fmt.Sprintf("Bad roll: %v", err))
				} else {
					actionList.Element(ListItem).Text = fmt.Sprintf("Army %s attacks settlement %s (AC: %d) rolling a %d for attack and %d for damage!",
						NameLink(army.Name), NameLink(target.Name), target.AC(), armyAttackRoll, damage)

					// Apply the damage and see if the settlement is overcome
					target.HP.Damage(damage)
					if target.HP.Current <= 0 {
						if target.Occupied {
							// If the settlement was occupied then we're liberating it
							actionList.Element(ListItem).Text = fmt.Sprintf("Settlement %s has been liberated by army %s!",
								NameLink(target.Name), NameLink(army.Name))

							target.Occupied = false
							target.Allegiance = army.Allegiance
						} else {
							// If the settlement wasn't occupied then it is now
							actionList.Element(ListItem).Text = fmt.Sprintf("Settlement %s has been occupied by army %s!",
								NameLink(target.Name), NameLink(army.Name))

							target.Occupied = true
							target.Allegiance = army.Allegiance
						}
					}
				}
			} else {
				actionList.Element(ListItem).Text = fmt.Sprintf("Army %s misses settlement %s (AC: %d) rolling a %d for attack.",
					NameLink(army.Name), NameLink(target.Name), target.AC(), armyAttackRoll)
			}
		}
	}

	return activityObserved
}

func NameLink(name string) *DocumentElement {
	anchor := Element(Anchor)
	anchor.Attributes["href"] = fmt.Sprintf("#%s", DocumentID(name))
	anchor.Text = name

	return anchor
}

func (s *World) stepSettlements(log *DocumentElement) bool {
	var activityObserved = false

	actionList := log.Element(UnorderedList)
	actionList.Attributes["style"] = "list-style-type: none;"

	for _, settlement := range SettlementListFromMap(s.Settlements).Sorted() {
		if !settlement.HasWarGuard {
			// Settlements that have no local war-trained guard may not retaliate against an
			// attacking force
			continue
		}

		if settlement.Occupied {
			// Occupied Settlements get no action
			actionList.Element(ListItem).Text = fmt.Sprintf("Settlement %s is occupied and gets no action.", NameLink(settlement.Name))
			continue
		}

		for _, army := range s.Armies {
			// Settlements can only attack the first army they see currently
			if !army.Destroyed && army.Location == settlement.Name && army.Allegiance != settlement.Allegiance {
				activityObserved = true

				if settlementAttackRoll, damage, err := settlement.RollAttack(); err != nil {
					panic(fmt.Sprintf("Bad roll: %v", err))
				} else if settlementAttackRoll >= army.AC {
					// If the settlement beats the army's AC then roll the damage
					actionList.Element(ListItem).Text = fmt.Sprintf("Settlement %s attacks army %s (AC: %d) rolling a %d for attack and %d for damage!",
						NameLink(settlement.Name), NameLink(army.Name), army.AC, settlementAttackRoll, damage)

					// Apply the damage and see if the army falls apart
					army.Damage(damage)
					if army.Destroyed {
						actionList.Element(ListItem).Text = fmt.Sprintf("Settlement %s has destroyed army %s!",
							NameLink(settlement.Name), NameLink(army.Name))
					}
				} else {
					actionList.Element(ListItem).Text = fmt.Sprintf("Settlement %s misses army %s (AC: %d) rolling a %d for attack.",
						NameLink(settlement.Name), NameLink(army.Name), army.AC, settlementAttackRoll)
				}

				break
			}
		}
	}

	return activityObserved
}

func (s *World) Turn(turnID int, output *strings.Builder) bool {
	html := Element("html")
	body := html.Element(HTBody)

	combatLogDiv := body.Element(Division)
	combatLogDiv.Attributes["style"] = "float: right; background: #EFEFEF; padding-left: 10px; padding-right: 10px; border: solid 1px black; width: 50%;"
	combatLogDiv.Element(H1).Text = "Combat Log"
	combatLogDiv.Element(H3).Text = fmt.Sprintf("Sim Turn: %d", turnID)

	// Move armies first
	armiesActive := s.stepArmies(combatLogDiv)

	// Allow Settlements to act last
	settlementsActive := s.stepSettlements(combatLogDiv)

	// Dump the world state
	s.WriteWorld(body)
	html.Output(output)

	// Return whether or not any activity took place this turn
	return armiesActive || settlementsActive
}

package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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
		if len(parts[numDicePart]) > 0 {
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

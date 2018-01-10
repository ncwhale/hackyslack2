package roll

import (
	"testing"
)

var ParseTests = []RollTest{
	{"2d6", []Dice{{Operator: Add, Number: 2, Sides: 6}}},
	{"", []Dice{{Operator: Add, Number: 1, Sides: 100}}},
	{"blah", []Dice{{Operator: Add, Number: 1, Sides: 100}}},
	{"0d0", []Dice{{Operator: Add, Number: 1, Sides: 1}}},
	{"d", []Dice{{Operator: Add, Number: 1, Sides: 100}}},
	{"d%", []Dice{{Operator: Add, Number: 1, Sides: 100}}},
	{"2d%", []Dice{{Operator: Add, Number: 2, Sides: 100}}},
	{"0d1", []Dice{{Operator: Add, Number: 1, Sides: 1}}},
	{"1d6!", []Dice{{Operator: Add, Number: 1, Sides: 6, Explode: true}}},
	{"4f", []Dice{{Operator: Add, Number: 4, Sides: 3, Fudge: true}}},
	{"4df", []Dice{{Operator: Add, Number: 4, Sides: 3, Fudge: true}}},
	{"999d9999", []Dice{{Operator: Add, Number: 100, Sides: 1000}}},
	{"2d6×5", []Dice{
		{Operator: Add, Number: 2, Sides: 6},
		{Operator: Multiply, Number: 5, Sides: 1},
	}},
	{"Skill=50", []Dice{
		{Operator: Add, Number: 50, Sides: 1},
		{Operator: Subtract, Number: 1, Sides: 100},
	}},
	{"1d20-1", []Dice{
		{Operator: Add, Number: 1, Sides: 20},
		{Operator: Subtract, Number: 1, Sides: 1},
	}},
	{"2d20+12345", []Dice{
		{Operator: Add, Number: 2, Sides: 20},
		{Operator: Add, Number: 12345, Sides: 1},
	}},
	{"2d2+1 1d6", []Dice{
		{Operator: Add, Number: 2, Sides: 2},
		{Operator: Add, Number: 1, Sides: 1},
		{Operator: Add, Number: 1, Sides: 6},
	}},
	{"1d20, 2d6-10", []Dice{
		{Operator: Add, Number: 1, Sides: 20},
		{Operator: Add, Number: 2, Sides: 6},
		{Operator: Subtract, Number: 10, Sides: 1},
	}},
	{"1d1+1 2d2-2 3d3+3", []Dice{
		{Operator: Add, Number: 1, Sides: 1},
		{Operator: Add, Number: 1, Sides: 1},
		{Operator: Add, Number: 2, Sides: 2},
		{Operator: Subtract, Number: 2, Sides: 1},
		{Operator: Add, Number: 3, Sides: 3},
		{Operator: Add, Number: 3, Sides: 1},
	}},
	{"2d6>5", []Dice{{Operator: Add, Number: 2, Sides: 6, Minimum: 5}}},
	{"2d6<2", []Dice{{Operator: Add, Number: 2, Sides: 6, Maximum: 2}}},
	{"2d6>6", []Dice{{Operator: Add, Number: 2, Sides: 6, Minimum: 5}}},
	{"2d6<1", []Dice{{Operator: Add, Number: 2, Sides: 6, Maximum: 2}}},
	{"6d6k5", []Dice{{Operator: Add, Number: 6, Sides: 6, Keep: 5}}},
	{"2d6k5", []Dice{{Operator: Add, Number: 2, Sides: 6, Keep: 2}}},
	{"6d6k-4", []Dice{{Operator: Add, Number: 6, Sides: 6, Keep: -4}}},
	{"2d6k-5", []Dice{{Operator: Add, Number: 2, Sides: 6, Keep: -2}}},
	{"2d6 + 1", []Dice{
		{Operator: Add, Number: 2, Sides: 6},
		{Operator: Add, Number: 1, Sides: 1},
	}},
	{"3d3 - 5", []Dice{
		{Operator: Add, Number: 3, Sides: 3},
		{Operator: Subtract, Number: 5, Sides: 1},
	}},
}

type RollTest struct {
	Text  string
	Rolls []Dice
}

func compareRolls(a Dice, b Dice) bool {
	return a.Number == b.Number && a.Sides == b.Sides &&
		a.Operator == b.Operator && a.Keep == b.Keep &&
		a.Minimum == b.Minimum && a.Maximum == b.Maximum &&
		a.Explode == b.Explode
}

func TestRoll(t *testing.T) {
	for _, test := range ParseTests {
		rolls := Parse(test.Text)
		for i, result := range rolls {
			if i >= len(test.Rolls) {
				t.Error("Failed", test, "got extra", *result)
				break
			}
			if !compareRolls(*result, test.Rolls[i]) {
				t.Error("Failed", test, "got", *result)
			}
			for i := 0; i < 10; i++ {
				result.Roll()
				if !result.Fudge && result.Keep == 0 && result.Total < result.Number {
					t.Error(test.Text, "Rolled too low", *result)
				}
				if !result.Explode && result.Keep == 0 && result.Total > result.Number*result.Sides {
					t.Error(test.Text, "Rolled too high", *result)
				}
			}
		}
	}
}

func TestDiceGenerate(t *testing.T) {
	if DiceGenerate(1) < 1 {
		t.Error("This roll should be 1")
	}

	if DiceGenerate(1000) > 1000 {
		t.Error("This roll should less than 1000")
	}

	if DiceGenerate(0) != 0 {
		t.Error("This roll should be 0")
	}

	if DiceGenerate(-1) != 0 {
		t.Error("This roll should be 0")
	}
}

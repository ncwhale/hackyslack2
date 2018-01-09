package dicebot

import (
	"fmt"
	"testing"

	"github.com/ncwhale/hackyslack2"
	"github.com/ncwhale/hackyslack2/dicebot/roll"
)

var FormatTests = []string{
	"",
	"d6",
	"10d10-5",
	"123456789d123456789+123456789",
	"1d1 2d2 3d3",
	"10d10>5",
	"10d10<5",
	"10d10k5",
	"2d6 / 2d6",
	"mini 2d6 / 2d6",
	"-1d% Skill=50",
	"-1d% * 5",
}

func TestRollFormat(t *testing.T) {
	username := "TestUser"
	for _, test := range FormatTests {
		rolls := roll.Parse(test)
		verify := fmt.Sprint("@", username, " rolled ")
		sum := 0
		for i, result := range rolls {
			result.Roll()
			// Use the last roll as the final.
			total := result.Total
			switch result.Operator {
			case "+":
				sum += result.Total
			case "-":
				sum -= result.Total
			case "*":
				sum *= result.Total
			case "/":
				sum /= result.Total
			}
			if i == 0 {
				if result.Operator == "-" {
					verify += fmt.Sprint("*-", total, "*")
				} else {
					verify += fmt.Sprint("*", total, "*")
				}
			} else {
				if result.Operator == "*" {
					verify += fmt.Sprint(" × *", total, "*")
				} else {
					verify += fmt.Sprint(" ", result.Operator, " *", total, "*")
				}
			}
		}
		if len(rolls) > 1 {
			verify += fmt.Sprint(" = *", sum, "*")
		}
		resp := formatRoll(username, rolls)
		if resp["response_type"] != "in_channel" {
			t.Error("Incorrect response type", resp["response_type"])
		}
		attach := resp["attachments"].([]hackyslack.D)[0]
		if attach["color"] == "" {
			t.Error("Missing color")
		}
		if attach["text"] != verify {
			t.Error("Incorrect response text", attach["text"], "instead of", verify)
		}
	}
}

// Test for panics.
func TestCommand(t *testing.T) {
	for _, test := range FormatTests {
		resp := command(hackyslack.Args{
			Text:     test,
			UserName: "CommandTest",
		})
		if resp["response_type"] != "in_channel" {
			t.Error("Incorrect response type", resp["response_type"])
		}
	}
}

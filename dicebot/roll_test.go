package dicebot

import (
	"encoding/json"
	"testing"

	"github.com/ncwhale/hackyslack2"
	"github.com/ncwhale/hackyslack2/dicebot/roll"
)

var CommandTests = []CommandTest{
	{
		"",
		[]*Block{
			{
				false,
				[]*roll.Dice{{Operator: roll.Add, Number: 1, Sides: 100}}},
		},
	},
	{
		"mini Skill=50, -1d% + 20 for fire, 100f for fire, 1d%>50, 1d%<50",
		[]*Block{
			{
				false,
				[]*roll.Dice{{Operator: roll.Add, Number: 1, Sides: 100}}},
		},
	},
	{
		"-1d% + 2d6 - 1d6 * 1d6 / 2, 10d6k4, 10d6k-4",
		[]*Block{
			{
				false,
				[]*roll.Dice{{Operator: roll.Add, Number: 1, Sides: 100}}},
		},
	},
}

type CommandTest struct {
	Text   string
	Result []*Block
}

// Test for panics.
func TestCommand(t *testing.T) {
	for _, test := range CommandTests {
		args := hackyslack.Args{
			Text:     test.Text,
			UserId:   "000L",
			UserName: "Test",
		}

		blocks := ParseCommand(args)
		// TODO: check output blocks.
		resp := formatBlocks(args, blocks)
		// TODO: check format output.
		_, err := json.Marshal(resp)

		if err != nil {
			t.Error("JSON stringify error:", err)
		}

		if resp["response_type"] != "in_channel" {
			t.Error("Incorrect response type", resp["response_type"])
		}
	}
}

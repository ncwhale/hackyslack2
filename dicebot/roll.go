package dicebot

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/ncwhale/hackyslack2"
	"github.com/ncwhale/hackyslack2/dicebot/roll"
)

var (
	blockRegexp = regexp.MustCompile(`(?i)(?P<mini>mini |m )?(?P<block>[^,;]+)[\s,;]*`)
)

// Block for every command block result.
type Block struct {
	mini  bool
	dices []*roll.Dice
}

func init() {
	hackyslack.Register("r", command)
	hackyslack.Register("roll", command)
}

func formatBlocks(args hackyslack.Args, results []*Block) hackyslack.D {
	var (
		attachments []hackyslack.D
	)

	// Get every block calculated.
	for _, block := range results {
		var (
			color    string
			fields   []hackyslack.D
			final    int
			text     string
			fallback string
		)
		for i, result := range block.dices {
			result.Roll()

			if i == 0 {
				if result.Operator == "-" {
					final -= result.Total
				} else {
					final = result.Total
				}
				text = fmt.Sprint("*", final, "*")
				fallback = fmt.Sprint(final)
				if result.For != "" {
					text += fmt.Sprint(" for *", result.For, "*")
					fallback += fmt.Sprint(" for ", result.For)
				}
			} else {
				op := result.Operator
				switch result.Operator {
				case roll.Add:
					final += result.Total
				case roll.Subtract:
					final -= result.Total
				case roll.Multiply:
					final *= result.Total
					op = "×"
				case roll.Divide:
					final /= result.Total
				case roll.Max:
					if result.Total > final {
						final = result.Total
					}
				case roll.Min:
					if result.Total < final {
						final = result.Total
					}
				}
				text += fmt.Sprint(" ", op, " *", result.Total, "*")
				fallback += fmt.Sprint(" ", result.Operator, " ", result.Total)
				if result.For != "" {
					text += fmt.Sprint(" for *", result.For, "*")
					fallback += fmt.Sprint(" for ", result.For)
				}
				if i == len(results)-1 {
					text += fmt.Sprint(" = *", final, "*")
					fallback += fmt.Sprint(" = ", final)
				}
			}

			if result.Sides == 1 {
				continue
			}
			rollText := fmt.Sprint(result.Rolls)
			if result.Fudge {
				if result.Total > 0 {
					color = "good"
				} else if result.Total == 0 {
					color = "warning"
				} else {
					color = "danger"
				}
			} else {
				single := result.Number * result.Sides / 3.0
				if result.Total > single*2 {
					color = "good"
				} else if result.Total > single+result.Number-1 {
					color = "warning"
				} else {
					color = "danger"
				}
			}
			dice := fmt.Sprint(result.Number, "d", result.Sides)
			if result.Fudge {
				dice = fmt.Sprint(result.Number, "df")
			}
			fields = append(fields, hackyslack.D{
				"title": "Dice",
				"value": dice,
				"short": true,
			}, hackyslack.D{
				"title": "Rolls",
				"value": rollText[1 : len(rollText)-1],
				"short": true,
			})
			if result.Minimum != 0 {
				count := 0
				for _, r := range result.Rolls {
					if r > result.Minimum {
						count++
					}
				}
				fields = append(fields, hackyslack.D{
					"title": "Minimum",
					"value": strconv.Itoa(result.Minimum),
					"short": true,
				}, hackyslack.D{
					"title": "Over",
					"value": strconv.Itoa(count),
					"short": true,
				})
			}
			if result.Maximum != 0 {
				count := 0
				for _, r := range result.Rolls {
					if r < result.Maximum {
						count++
					}
				}
				fields = append(fields, hackyslack.D{
					"title": "Maximum",
					"value": strconv.Itoa(result.Maximum),
					"short": true,
				}, hackyslack.D{
					"title": "Under",
					"value": strconv.Itoa(count),
					"short": true,
				})
			}
			if result.Keep != 0 {
				removed := fmt.Sprint(result.Removed)
				fields = append(fields, hackyslack.D{
					"title": "Keep",
					"value": strconv.Itoa(result.Keep),
					"short": true,
				}, hackyslack.D{
					"title": "Removed",
					"value": removed[1 : len(removed)-1],
					"short": true,
				})
			}
		}

		if block.mini {
			fields = []hackyslack.D{}
		}

		attachments = append(attachments, hackyslack.D{
			// <@UserID>  for slack shortcut.
			"fallback": fmt.Sprint("@", args.UserName, " rolled ", fallback),
			"pretext":  fmt.Sprint("<@", args.UserId, "> rolled ", text),
			// TODO: Color just uses the last color chosen.
			"color":     color,
			"fields":    fields,
			"mrkdwn_in": []string{"pretext"},
		})
	}

	return hackyslack.D{
		"response_type": "in_channel",
		"attachments":   attachments,
	}
}

// ParseCommand split slack command text to blocks.
func ParseCommand(args hackyslack.Args) []*Block {
	var result []*Block
	for _, m := range blockRegexp.FindAllStringSubmatch(args.Text, 5) {
		block := &Block{}
		for i, name := range blockRegexp.SubexpNames() {
			switch name {
			case "mini":
				block.mini = m[i] != ""
			case "block":
				block.dices = roll.Parse(m[i])
			}
		}

		result = append(result, block)
	}

	if len(result) == 0 {
		result = append(result, &Block{
			mini:  false,
			dices: roll.Parse(""),
		})
	}

	return result
}

func command(args hackyslack.Args) hackyslack.D {
	result := ParseCommand(args)
	return formatBlocks(args, result)
}

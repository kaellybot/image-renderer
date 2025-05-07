package constants

import (
	"image"
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
)

type Command struct {
	Name      string
	Duration  time.Duration
	Bounds    image.Rectangle
	Arguments map[amqp.Language][]string
}

func GetCommands() map[string]Command {
	return map[string]Command{
		"about": {
			Name:     "about",
			Duration: 10 * time.Second,
			Bounds:   image.Rect(350, 200, 1000, 1010),
		},
		"job_get": {
			Name:     "job_get",
			Duration: 12 * time.Second,
			Bounds:   image.Rect(350, 200, 1000, 1010),
			Arguments: map[amqp.Language][]string{
				amqp.Language_FR: {"Alchimiste"},
				amqp.Language_EN: {"Alchemist"},
				amqp.Language_ES: {"Alquimista"},
				amqp.Language_DE: {"Alchemist"},
				amqp.Language_PT: {"Alquimista"},
			},
		},
	}
}

func TODO() []string {
	return []string{
		"align_get", "align_set",
		"almanax_day", "almanax_effects", "almanax_resources",
		"config_get", "config_almanax", "config_rss", "config_server", "config_twitter",
		"help",
		"item",
		"job_set",
		"map",
		"pos",
		"set",
	}
}

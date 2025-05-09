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
			Duration: 8 * time.Second,
			Bounds:   image.Rect(350, 200, 1000, 1010),
		},
		"align_get": {
			Name:     "align_get",
			Duration: 8 * time.Second,
			Bounds:   image.Rect(350, 400, 1000, 1010),
		},
		"align_set": {
			Name:     "align_set",
			Duration: 12 * time.Second,
			Bounds:   image.Rect(350, 400, 1000, 1010),
			Arguments: map[amqp.Language][]string{
				amqp.Language_FR: {"Bonta", "Esprit"},
				amqp.Language_EN: {"Bonta", "Spirit"},
				amqp.Language_ES: {"Bonta", "Espíritu"},
				amqp.Language_DE: {"Bonta", "Geist"},
				amqp.Language_PT: {"Bonta", "Espírito"},
			},
		},
		"job_get": {
			Name:     "job_get",
			Duration: 10 * time.Second,
			Bounds:   image.Rect(350, 400, 1000, 1010),
			Arguments: map[amqp.Language][]string{
				amqp.Language_FR: {"Alchimiste"},
				amqp.Language_EN: {"Alchemist"},
				amqp.Language_ES: {"Alquimista"},
				amqp.Language_DE: {"Alchemist"},
				amqp.Language_PT: {"Alquimista"},
			},
		},
		"job_set": {
			Name:     "job_set",
			Duration: 10 * time.Second,
			Bounds:   image.Rect(350, 400, 1000, 1010),
			Arguments: map[amqp.Language][]string{
				amqp.Language_FR: {"Alchimiste"},
				amqp.Language_EN: {"Alchemist"},
				amqp.Language_ES: {"Alquimista"},
				amqp.Language_DE: {"Alchemist"},
				amqp.Language_PT: {"Alquimista"},
			},
		},
		"map": {
			Name:     "map",
			Duration: 8 * time.Second,
			Bounds:   image.Rect(350, 400, 1000, 1010),
		},
	}
}

func TODO() []string {
	return []string{
		"almanax_day", "almanax_effects", "almanax_resources",
		"config_get", "config_almanax", "config_rss", "config_server", "config_twitter",
		"help",
		"item",
		"pos",
		"set",
	}
}

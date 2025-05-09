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
		"almanax_day": {
			Name:     "almanax_day",
			Duration: 8 * time.Second,
			Bounds:   image.Rect(350, 400, 1000, 1010),
		},
		"almanax_effects": {
			Name:     "almanax_effects",
			Duration: 10 * time.Second,
			Bounds:   image.Rect(350, 300, 1000, 1010),
			Arguments: map[amqp.Language][]string{
				amqp.Language_FR: {"Économie"},
				amqp.Language_EN: {"Saved"},
				amqp.Language_ES: {"Economía"},
				amqp.Language_DE: {"Zutateneinsparung"},
				amqp.Language_PT: {"Economia"},
			},
		},
		"almanax_resources": {
			Name:     "almanax_resources",
			Duration: 10 * time.Second,
			Bounds:   image.Rect(350, 400, 1000, 1010),
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
		"item": {
			Name:     "item",
			Duration: 10 * time.Second,
			Bounds:   image.Rect(350, 270, 1000, 1010),
			Arguments: map[amqp.Language][]string{
				amqp.Language_FR: {"Amulette Bouftou"},
				amqp.Language_EN: {"Gobball Amulet"},
				amqp.Language_ES: {"Amuleto Jalató"},
				amqp.Language_DE: {"Amulett Fresssacks"},
				amqp.Language_PT: {"Amuleto Papatudo"},
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
		"pos": {
			Name:     "pos",
			Duration: 8 * time.Second,
			Bounds:   image.Rect(350, 150, 1000, 1010),
		},
		"set": {
			Name:     "set",
			Duration: 10 * time.Second,
			Bounds:   image.Rect(350, 270, 1000, 1010),
			Arguments: map[amqp.Language][]string{
				amqp.Language_FR: {"Bouftou"},
				amqp.Language_EN: {"Gobball"},
				amqp.Language_ES: {"Jalató"},
				amqp.Language_DE: {"Fresssack"},
				amqp.Language_PT: {"Papatudo"},
			},
		},
	}
}

func TODO() []string {
	return []string{
		"config_get",
		"config_almanax",
		"config_rss",
		"config_server",
		"config_twitter",
	}
}

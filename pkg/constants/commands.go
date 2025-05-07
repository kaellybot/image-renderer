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
		"about": Command{
			Name:     "about",
			Duration: 10 * time.Second,
			Bounds:   image.Rect(350, 200, 1000, 1010),
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
		"job_get", "job_set",
		"map",
		"pos",
		"set",
	}
}

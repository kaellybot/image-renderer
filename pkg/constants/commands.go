package constants

import (
	"image"
	"time"
)

type Command struct {
	Name     string
	Duration time.Duration
	Bounds   image.Rectangle
}

func GetCommandNames() []string {
	return []string{
		"about",
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

package constants

import amqp "github.com/kaellybot/kaelly-amqp"

type Command struct {
	Name        string
	Automations map[amqp.Language]func()
}

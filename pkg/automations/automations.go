package automations

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-vgo/robotgo"
)

func SetupDiscordTutorial() error {
	pids, err := robotgo.FindIds("Discord")
	if err != nil {
		return err
	}
	if len(pids) == 0 {
		return fmt.Errorf("discord not found")
	}

	// Focus the first matching Discord window
	robotgo.ActivePid(pids[0])

	// Click server "tutorial"
	robotgo.MoveSmooth(40, 250)
	robotgo.Click("left")
	time.Sleep(500 * time.Millisecond)

	// Click channel inside that server
	robotgo.MoveSmooth(150, 220)
	robotgo.Click("left")
	return nil
}

func TypeTextSlowly(text string, minDelay, maxDelay int) {
	for _, char := range text {
		robotgo.TypeStr(string(char))
		delay := time.Duration(rand.Intn(maxDelay-minDelay+1)+minDelay) * time.Millisecond
		time.Sleep(delay)
	}
}

package automations

import (
	"embed"
	"fmt"
	"kaellybot/image-renderer/pkg/constants"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	amqp "github.com/kaellybot/kaelly-amqp"
)

const autoHotkeyExt = "ahk"

//go:embed *.ahk
var Folder embed.FS

func SetupDiscordTutorial() error {
	setup := constants.Command{
		Name: "_tuto_setup",
	}
	return runAHKScript(setup, amqp.Language_ANY)
}

func RunCommandTutorial(command constants.Command, locale amqp.Language) error {
	return runAHKScript(command, locale)
}

func runAHKScript(command constants.Command, locale amqp.Language) error {
	rawScript, err := Folder.ReadFile(fmt.Sprintf("cmd_%v.%v", command.Name, autoHotkeyExt))
	if err != nil {
		return err
	}

	tmpDir := os.TempDir()
	scriptPath := filepath.Join(tmpDir, fmt.Sprintf(command.Name, ".", autoHotkeyExt))

	if len(command.Arguments) > 0 {
		script := string(rawScript)
		args, found := command.Arguments[locale]
		if !found {
			return fmt.Errorf("missing locale %v for command '%v'", locale, command.Name)
		}

		for _, arg := range args {
			script = strings.Replace(script, "{{ . }}", arg, 1)
		}

		rawScript = []byte(script)
	}

	if err := os.WriteFile(scriptPath, rawScript, 0644); err != nil {
		return err
	}
	defer os.Remove(scriptPath)

	cmd := exec.Command("AutoHotkeyUX.exe", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

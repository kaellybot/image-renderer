package automations

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	amqp "github.com/kaellybot/kaelly-amqp"
)

//go:embed *.ahk
var Folder embed.FS

func SetupDiscordTutorial() error {
	return runAHKScript("_discord_tutorial_setup.ahk")
}

func RunCommandTutorial(commandName string, locale amqp.Language) error {
	return runAHKScript(fmt.Sprintf("%v_%v.ahk", commandName, locale))
}

func runAHKScript(name string) error {
	script, err := Folder.ReadFile(name)
	if err != nil {
		return err
	}

	tmpDir := os.TempDir()
	scriptPath := filepath.Join(tmpDir, name)

	if err := os.WriteFile(scriptPath, script, 0644); err != nil {
		return err
	}
	defer os.Remove(scriptPath)

	cmd := exec.Command("AutoHotkeyUX.exe", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

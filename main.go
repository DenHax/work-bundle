package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func initConfigs() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("не удалось получить домашнюю директорию: %v", err)
	}

	configs := map[string]string{
		"nvim": filepath.Join(homeDir, ".config", "nvim", "init.lua"),
		"vim":  filepath.Join(homeDir, ".vimrc"),
		"tmux": filepath.Join(homeDir, ".tmux.conf"),
		"zsh":  filepath.Join(homeDir, ".zshrc"),
	}

	for tool, destPath := range configs {
		srcPath := fmt.Sprintf("./configs/%s", tool)
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			return fmt.Errorf("конфигурация для %s не найдена", tool)
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("не удалось создать директорию для %s: %v", tool, err)
		}

		if err := copyFile(srcPath, destPath); err != nil {
			return fmt.Errorf("не удалось скопировать конфигурацию для %s: %v", tool, err)
		}

		fmt.Printf("Конфигурация для %s успешно установлена\n", tool)
	}

	return nil
}

func copyFile(src, dest string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dest, input, 0644)
}

func installZsh() error {
	if err := runCommand("sudo", "apt", "install", "-y", "zsh"); err != nil {
		return fmt.Errorf("ошибка при установке Zsh: %v", err)
	}

	if err := runCommand("sh", "-c", "$(curl -fsSL curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"); err != nil {
		return fmt.Errorf("ошибка при установке Oh My Zsh: %v", err)
	}

	zshrcPath := filepath.Join(os.Getenv("HOME"), ".zshrc")
	localZshrc, err := os.ReadFile("./configs/zsh/.zshrc")
	if err != nil {
		return fmt.Errorf("не удалось прочитать локальный .zshrc: %v", err)
	}

	file, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("не удалось открыть .zshrc: %v", err)
	}
	defer file.Close()

	if _, err := file.Write(localZshrc); err != nil {
		return fmt.Errorf("не удалось записать в .zshrc: %v", err)
	}

	fmt.Println("Zsh и Oh My Zsh успешно установлены, .zshrc обновлен")
	return nil
}

func runCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: ./setup-tools <команда>")
		fmt.Println("Команды: init, install-zsh")
		return
	}

	switch os.Args[1] {
	case "init":
		if err := initConfigs(); err != nil {
			fmt.Println("Ошибка:", err)
			os.Exit(1)
		}
	case "install-zsh":
		if err := installZsh(); err != nil {
			fmt.Println("Ошибка:", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Неизвестная команда")
		os.Exit(1)
	}
}

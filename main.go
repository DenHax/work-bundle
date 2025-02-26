package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type PackageManager struct {
	Name         string
	InstallCmd   string
	UpdateCmd    string
	PrivilegeCmd string
}

func detectPackageManager() *PackageManager {
	switch runtime.GOOS {
	case "linux":

		if _, err := exec.LookPath("apt"); err == nil {
			return &PackageManager{
				Name:         "apt",
				InstallCmd:   "install -y",
				UpdateCmd:    "update",
				PrivilegeCmd: "sudo",
			}
		}

		if _, err := exec.LookPath("dnf"); err == nil {
			return &PackageManager{
				Name:         "dnf",
				InstallCmd:   "install -y",
				UpdateCmd:    "makecache",
				PrivilegeCmd: "sudo",
			}
		}

		if _, err := exec.LookPath("pacman"); err == nil {
			return &PackageManager{
				Name:         "pacman",
				InstallCmd:   "-S --noconfirm",
				UpdateCmd:    "-Sy",
				PrivilegeCmd: "sudo",
			}
		}
	case "darwin":

		if _, err := exec.LookPath("brew"); err == nil {
			return &PackageManager{
				Name:         "brew",
				InstallCmd:   "install",
				UpdateCmd:    "update",
				PrivilegeCmd: "",
			}
		}
	}
	return nil
}

func readPackagesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", filename, err)
	}
	defer file.Close()

	var packages []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		packages = append(packages, strings.TrimSpace(scanner.Text()))
	}
	return packages, scanner.Err()
}

func installPackages(pm *PackageManager, packages []string) error {
	if len(packages) == 0 {
		return fmt.Errorf("no packages to install")
	}

	if pm.UpdateCmd != "" {
		updateCmd := exec.Command(pm.PrivilegeCmd, pm.Name, pm.UpdateCmd)
		updateCmd.Stdout = os.Stdout
		updateCmd.Stderr = os.Stderr
		fmt.Printf("Updating package manager (%s)...\n", pm.Name)
		if err := updateCmd.Run(); err != nil {
			return fmt.Errorf("failed to update package manager: %v", err)
		}
	}

	installArgs := strings.Fields(pm.InstallCmd)
	installArgs = append(installArgs, packages...)
	installCmd := exec.Command(pm.PrivilegeCmd, pm.Name)
	installCmd.Args = append(installCmd.Args, installArgs...)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	fmt.Printf("Installing packages: %v\n", packages)
	return installCmd.Run()
}

func backupConfig(path string) error {
	if _, err := os.Stat(path); err == nil {
		backupPath := path + ".bak"
		fmt.Printf("Backing up %s to %s\n", path, backupPath)
		return os.Rename(path, backupPath)
	}
	return nil
}

func initConfigs() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	configs := map[string]string{
		"nvim": filepath.Join(homeDir, ".config", "nvim", "init.lua"),
		"vim":  filepath.Join(homeDir, ".vimrc"),
		"tmux": filepath.Join(homeDir, ".tmux.conf"),
		"zsh":  filepath.Join(homeDir, ".zshrc"),
	}

	for tool, destPath := range configs {
		srcPath := filepath.Join("./configs", tool)

		if err := backupConfig(destPath); err != nil {
			return fmt.Errorf("failed to backup %s: %v", tool, err)
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %v", tool, err)
		}

		if err := copyFile(srcPath, destPath); err != nil {
			return fmt.Errorf("failed to copy configuration for %s: %v", tool, err)
		}

		fmt.Printf("Configuration for %s successfully installed\n", tool)
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
		return fmt.Errorf("error installing Zsh: %v", err)
	}

	if err := runCommand("sh", "-c", "$(curl -fsSL curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"); err != nil {
		return fmt.Errorf("error installing Oh My Zsh: %v", err)
	}

	zshrcPath := filepath.Join(os.Getenv("HOME"), ".zshrc")
	localZshrc, err := os.ReadFile("./configs/zsh/.zshrc")
	if err != nil {
		return fmt.Errorf("failed to read local .zshrc: %v", err)
	}

	file, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open .zshrc: %v", err)
	}
	defer file.Close()

	if _, err := file.Write(localZshrc); err != nil {
		return fmt.Errorf("failed to write to .zshrc: %v", err)
	}

	fmt.Println("Zsh and Oh My Zsh successfully installed, .zshrc updated")
	return nil
}

func runCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {

	pm := detectPackageManager()
	if pm == nil {
		fmt.Println("Unsupported system or package manager not found")
		return
	}
	fmt.Printf("Detected package manager: %s\n", pm.Name)

	packages, err := readPackagesFromFile("packages.txt")
	if err != nil {
		fmt.Printf("Error reading packages.txt: %v\n", err)
		return
	}

	if err := installPackages(pm, packages); err != nil {
		fmt.Printf("Error installing packages: %v\n", err)
		return
	}

	if err := initConfigs(); err != nil {
		fmt.Printf("Error initializing configurations: %v\n", err)
		return
	}

	if err := installZsh(); err != nil {
		fmt.Printf("Error installing Zsh and Oh My Zsh: %v\n", err)
		return
	}

	fmt.Println("Setup completed successfully!")
}

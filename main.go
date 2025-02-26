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

// PackageManager represents a package manager and its commands
type PackageManager struct {
	Name         string
	InstallCmd   string
	UpdateCmd    string
	PrivilegeCmd string
}

// Detect the package manager for the current system
func detectPackageManager() *PackageManager {
	switch runtime.GOOS {
	case "linux":
		// Check for APT (Debian/Ubuntu)
		if _, err := exec.LookPath("apt"); err == nil {
			return &PackageManager{
				Name:         "apt",
				InstallCmd:   "install -y",
				UpdateCmd:    "update",
				PrivilegeCmd: "sudo",
			}
		}
		// Check for DNF (Fedora)
		if _, err := exec.LookPath("dnf"); err == nil {
			return &PackageManager{
				Name:         "dnf",
				InstallCmd:   "install -y",
				UpdateCmd:    "makecache",
				PrivilegeCmd: "sudo",
			}
		}
		// Check for Pacman (Arch)
		if _, err := exec.LookPath("pacman"); err == nil {
			return &PackageManager{
				Name:         "pacman",
				InstallCmd:   "-S --noconfirm",
				UpdateCmd:    "-Sy",
				PrivilegeCmd: "sudo",
			}
		}
	case "darwin":
		// Check for Homebrew (macOS)
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

// Read packages from packages.txt
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

// Install packages using the detected package manager
func installPackages(pm *PackageManager, packages []string) error {
	if len(packages) == 0 {
		return fmt.Errorf("no packages to install")
	}

	// Update package manager
	if pm.UpdateCmd != "" {
		updateCmd := exec.Command(pm.PrivilegeCmd, pm.Name, pm.UpdateCmd)
		updateCmd.Stdout = os.Stdout
		updateCmd.Stderr = os.Stderr
		fmt.Printf("Updating package manager (%s)...\n", pm.Name)
		if err := updateCmd.Run(); err != nil {
			return fmt.Errorf("failed to update package manager: %v", err)
		}
	}

	// Install packages
	installArgs := strings.Fields(pm.InstallCmd)
	installArgs = append(installArgs, packages...)
	installCmd := exec.Command(pm.PrivilegeCmd, pm.Name)
	installCmd.Args = append(installCmd.Args, installArgs...)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	fmt.Printf("Installing packages: %v\n", packages)
	return installCmd.Run()
}

// Backup existing config files
func backupConfig(path string) error {
	if _, err := os.Stat(path); err == nil {
		backupPath := path + ".bak"
		fmt.Printf("Backing up %s to %s\n", path, backupPath)
		return os.Rename(path, backupPath)
	}
	return nil
}

// Copy a file or directory recursively
func copyConfig(src, dest string) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %v", src, err)
	}

	if info.IsDir() {
		// Create the destination directory
		if err := os.MkdirAll(dest, info.Mode()); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dest, err)
		}

		// Read the source directory
		entries, err := os.ReadDir(src)
		if err != nil {
			return fmt.Errorf("failed to read directory %s: %v", src, err)
		}

		// Copy each entry
		for _, entry := range entries {
			srcPath := filepath.Join(src, entry.Name())
			destPath := filepath.Join(dest, entry.Name())

			if err := copyConfig(srcPath, destPath); err != nil {
				return err
			}
		}
	} else {
		// Copy the file
		input, err := os.ReadFile(src)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %v", src, err)
		}
		if err := os.WriteFile(dest, input, info.Mode()); err != nil {
			return fmt.Errorf("failed to write file %s: %v", dest, err)
		}
	}

	return nil
}

// Initialize configurations from the configs folder
func initConfigs() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	// Define configurations and their destination paths
	configs := map[string]string{
		"nvim": filepath.Join(homeDir, ".config", "nvim"),
		"vim":  filepath.Join(homeDir, ".vimrc"),
		"tmux": filepath.Join(homeDir, ".tmux.conf"),
		"zsh":  filepath.Join(homeDir, ".zshrc"),
	}

	for tool, destPath := range configs {
		srcPath := filepath.Join("./configs", tool)

		// Backup existing config
		if err := backupConfig(destPath); err != nil {
			return fmt.Errorf("failed to backup %s: %v", tool, err)
		}

		// Copy the configuration file or directory
		if err := copyConfig(srcPath, destPath); err != nil {
			return fmt.Errorf("failed to copy configuration for %s: %v", tool, err)
		}

		fmt.Printf("Configuration for %s successfully installed\n", tool)
	}

	return nil
}

// Install Zsh and Oh My Zsh
func installZsh() error {
	// Install Zsh
	if err := runCommand("sudo", "apt", "install", "-y", "zsh"); err != nil {
		return fmt.Errorf("error installing Zsh: %v", err)
	}

	// Install Oh My Zsh
	if err := runCommand("sh", "-c", "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"); err != nil {
		return fmt.Errorf("error installing Oh My Zsh: %v", err)
	}

	// Append to .zshrc
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

// Run a shell command
func runCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	// Detect package manager
	pm := detectPackageManager()
	if pm == nil {
		fmt.Println("Unsupported system or package manager not found")
		return
	}
	fmt.Printf("Detected package manager: %s\n", pm.Name)

	// Read packages from packages.txt
	packages, err := readPackagesFromFile("packages.txt")
	if err != nil {
		fmt.Printf("Error reading packages.txt: %v\n", err)
		return
	}

	// Install packages
	if err := installPackages(pm, packages); err != nil {
		fmt.Printf("Error installing packages: %v\n", err)
		return
	}

	// Initialize configurations
	if err := initConfigs(); err != nil {
		fmt.Printf("Error initializing configurations: %v\n", err)
		return
	}

	// Install Zsh and Oh My Zsh
	if err := installZsh(); err != nil {
		fmt.Printf("Error installing Zsh and Oh My Zsh: %v\n", err)
		return
	}

	fmt.Println("Setup completed successfully!")
}

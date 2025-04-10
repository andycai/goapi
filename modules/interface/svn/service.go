package svn

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var svnCmd string

func initService() {
	svnCmd = getSvnCommand()
}

// getSvnCommand returns the appropriate Svn command based on the OS
func getSvnCommand() string {
	if runtime.GOOS == "windows" {
		return "svn.exe"
	}
	return "svn"
}

// isValidSvnPath checks if the path is safe for Svn operations
func isValidSvnPath(path string) bool {
	// Clean and normalize the path
	cleanPath := filepath.Clean(path)

	// Check for suspicious patterns
	suspicious := []string{
		"..",
		"~",
		"$",
		"|",
		">",
		"<",
		"&",
		"`",
		"*",
		"?",
		"[",
		"]",
	}

	for _, pattern := range suspicious {
		if strings.Contains(cleanPath, pattern) {
			return false
		}
	}

	return true
}

// execSvnCommand executes an Svn command and returns its output
func execSvnCommand(args ...string) (string, error) {
	cmd := exec.Command(svnCmd, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Svn command failed: %v\nError: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// Checkout performs Svn checkout operation
func Checkout(url, path, username, password string) error {
	if !isValidSvnPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	args := []string{
		"checkout",
		url,
		path,
		"--username", username,
		"--password", password,
		"--non-interactive",
		"--trust-server-cert",
	}

	_, err := execSvnCommand(args...)
	return err
}

// Update performs Svn update operation
func Update(path string) error {
	if !isValidSvnPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	_, err := execSvnCommand("update", path, "--non-interactive")
	return err
}

// Commit performs Svn commit operation
func Commit(path, message string) error {
	if !isValidSvnPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	_, err := execSvnCommand("commit", path, "-m", message, "--non-interactive")
	return err
}

// Status gets Svn working copy status
func Status(path string) (string, error) {
	if !isValidSvnPath(path) {
		return "", fmt.Errorf("invalid path: %s", path)
	}

	return execSvnCommand("status", path)
}

// Info gets Svn repository information
func Info(path string) (string, error) {
	if !isValidSvnPath(path) {
		return "", fmt.Errorf("invalid path: %s", path)
	}

	return execSvnCommand("info", path)
}

// Log gets Svn commit history
func Log(path string, limit int) (string, error) {
	if !isValidSvnPath(path) {
		return "", fmt.Errorf("invalid path: %s", path)
	}

	return execSvnCommand("log", path, "-l", fmt.Sprintf("%d", limit))
}

// Revert reverts Svn working copy changes
func Revert(path string) error {
	if !isValidSvnPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	_, err := execSvnCommand("revert", path, "-R")
	return err
}

// Add adds files to Svn version control
func Add(path string) error {
	if !isValidSvnPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	_, err := execSvnCommand("add", path, "--non-interactive")
	return err
}

// Delete removes files from Svn version control
func Delete(path string) error {
	if !isValidSvnPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	_, err := execSvnCommand("delete", path, "--non-interactive")
	return err
}

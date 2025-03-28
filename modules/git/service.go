package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var gitService *GitService

type GitService struct {
	gitCmd string
}

func InitService() {
	gitService = &GitService{
		gitCmd: getGitCommand(),
	}
}

func GetService() *GitService {
	return gitService
}

// getGitCommand returns the appropriate Git command based on the OS
func getGitCommand() string {
	if runtime.GOOS == "windows" {
		return "git.exe"
	}
	return "git"
}

// isValidGitPath checks if the path is safe for Git operations
func (s *GitService) isValidGitPath(path string) bool {
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

// ExecGitCommand executes a Git command and returns its output
func (s *GitService) ExecGitCommand(dir string, args ...string) (string, error) {
	cmd := exec.Command(s.gitCmd, args...)
	if dir != "" {
		cmd.Dir = dir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Git command failed: %v\nError: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// Clone performs Git clone operation
func (s *GitService) Clone(url, path, branch string, username, password string) error {
	if !s.isValidGitPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	args := []string{"clone"}
	if branch != "" {
		args = append(args, "-b", branch)
	}

	if username != "" && password != "" {
		// Insert credentials into URL
		url = strings.Replace(url, "https://", fmt.Sprintf("https://%s:%s@", username, password), 1)
	}

	args = append(args, url, path)
	_, err := s.ExecGitCommand("", args...)
	return err
}

// Pull performs Git pull operation
func (s *GitService) Pull(path string) error {
	if !s.isValidGitPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	_, err := s.ExecGitCommand(path, "pull")
	return err
}

// Push performs Git push operation
func (s *GitService) Push(path string, branch string) error {
	if !s.isValidGitPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	args := []string{"push"}
	if branch != "" {
		args = append(args, "origin", branch)
	}

	_, err := s.ExecGitCommand(path, args...)
	return err
}

// Status gets Git repository status
func (s *GitService) Status(path string) (string, error) {
	if !s.isValidGitPath(path) {
		return "", fmt.Errorf("invalid path: %s", path)
	}

	return s.ExecGitCommand(path, "status")
}

// Log gets Git commit history
func (s *GitService) Log(path string, limit int) (string, error) {
	if !s.isValidGitPath(path) {
		return "", fmt.Errorf("invalid path: %s", path)
	}

	args := []string{"log", "--oneline"}
	if limit > 0 {
		args = append(args, fmt.Sprintf("-%d", limit))
	}

	return s.ExecGitCommand(path, args...)
}

// Commit performs Git commit operation
func (s *GitService) Commit(path, message string) error {
	if !s.isValidGitPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	_, err := s.ExecGitCommand(path, "commit", "-m", message)
	return err
}

// Checkout performs Git checkout operation
func (s *GitService) Checkout(path, branch string, create bool) error {
	if !s.isValidGitPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	args := []string{"checkout"}
	if create {
		args = append(args, "-b")
	}
	args = append(args, branch)

	_, err := s.ExecGitCommand(path, args...)
	return err
}

// Branch performs Git branch operations
func (s *GitService) Branch(path string, create bool, name string) (string, error) {
	if !s.isValidGitPath(path) {
		return "", fmt.Errorf("invalid path: %s", path)
	}

	args := []string{"branch"}
	if create && name != "" {
		args = append(args, name)
		_, err := s.ExecGitCommand(path, args...)
		return "", err
	}

	// List branches if not creating
	return s.ExecGitCommand(path, args...)
}

// Merge performs Git merge operation
func (s *GitService) Merge(path, branch string) error {
	if !s.isValidGitPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	_, err := s.ExecGitCommand(path, "merge", branch)
	return err
}

// Reset performs Git reset operation
func (s *GitService) Reset(path string, hard bool) error {
	if !s.isValidGitPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	args := []string{"reset"}
	if hard {
		args = append(args, "--hard")
	}

	_, err := s.ExecGitCommand(path, args...)
	return err
}

// Stash performs Git stash operations
func (s *GitService) Stash(path string, pop bool) error {
	if !s.isValidGitPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	args := []string{"stash"}
	if pop {
		args = append(args, "pop")
	}

	_, err := s.ExecGitCommand(path, args...)
	return err
}

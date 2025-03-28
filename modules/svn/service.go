package svn

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var srv *SVNService

type SVNService struct {
	svnCmd string
}

func initService() *SVNService {
	srv = &SVNService{
		svnCmd: getSVNCommand(),
	}

	return srv
}

// getSVNCommand returns the appropriate SVN command based on the OS
func getSVNCommand() string {
	if runtime.GOOS == "windows" {
		return "svn.exe"
	}
	return "svn"
}

// isValidSVNPath checks if the path is safe for SVN operations
func (s *SVNService) isValidSVNPath(path string) bool {
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

// execSVNCommand executes an SVN command and returns its output
func (s *SVNService) execSVNCommand(args ...string) (string, error) {
	cmd := exec.Command(s.svnCmd, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("SVN command failed: %v\nError: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// Checkout performs SVN checkout operation
func (s *SVNService) Checkout(url, path, username, password string) error {
	if !s.isValidSVNPath(path) {
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

	_, err := s.execSVNCommand(args...)
	return err
}

// Update performs SVN update operation
func (s *SVNService) Update(path string) error {
	if !s.isValidSVNPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	_, err := s.execSVNCommand("update", path, "--non-interactive")
	return err
}

// Commit performs SVN commit operation
func (s *SVNService) Commit(path, message string) error {
	if !s.isValidSVNPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	_, err := s.execSVNCommand("commit", path, "-m", message, "--non-interactive")
	return err
}

// Status gets SVN working copy status
func (s *SVNService) Status(path string) (string, error) {
	if !s.isValidSVNPath(path) {
		return "", fmt.Errorf("invalid path: %s", path)
	}

	return s.execSVNCommand("status", path)
}

// Info gets SVN repository information
func (s *SVNService) Info(path string) (string, error) {
	if !s.isValidSVNPath(path) {
		return "", fmt.Errorf("invalid path: %s", path)
	}

	return s.execSVNCommand("info", path)
}

// Log gets SVN commit history
func (s *SVNService) Log(path string, limit int) (string, error) {
	if !s.isValidSVNPath(path) {
		return "", fmt.Errorf("invalid path: %s", path)
	}

	return s.execSVNCommand("log", path, "-l", fmt.Sprintf("%d", limit))
}

// Revert reverts SVN working copy changes
func (s *SVNService) Revert(path string) error {
	if !s.isValidSVNPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	_, err := s.execSVNCommand("revert", path, "-R")
	return err
}

// Add adds files to SVN version control
func (s *SVNService) Add(path string) error {
	if !s.isValidSVNPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	_, err := s.execSVNCommand("add", path, "--non-interactive")
	return err
}

// Delete removes files from SVN version control
func (s *SVNService) Delete(path string) error {
	if !s.isValidSVNPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	_, err := s.execSVNCommand("delete", path, "--non-interactive")
	return err
}

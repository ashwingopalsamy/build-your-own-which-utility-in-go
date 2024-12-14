package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: which <command>")
		os.Exit(1)
	}

	cmd := os.Args[1]

	path := os.Getenv("PATH")
	if path == "" {
		fmt.Fprintln(os.Stderr, "PATH environment variable not set.")
		os.Exit(1)
	}

	for _, dir := range filepath.SplitList(path) {
		if fullPath := findExecutable(dir, cmd); fullPath != "" {
			fmt.Println(fullPath)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "%s not found in PATH\n", cmd)
	os.Exit(1)
}

func findExecutable(dir, cmd string) string {
	// Handle explicit paths directly
	if strings.Contains(cmd, string(os.PathSeparator)) {
		if isExecutable(cmd) {
			return cmd
		}
		return "" // Not executable, bail out
	}

	fullPath := filepath.Join(dir, cmd)
	if isExecutable(fullPath) {
		return fullPath
	}

	// Windows extension handling
	if runtime.GOOS == "windows" && !strings.Contains(cmd, ".") {
		for _, ext := range getExecutableExtensions() {
			candidate := filepath.Join(dir, cmd+ext)
			if isExecutable(candidate) {
				return candidate
			}
		}
	}

	return "" // Not found in this directory
}

func isExecutable(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false // File doesn't exist
	}

	m := fi.Mode()
	if m.IsDir() || !m.IsRegular() {
		return false // Not a regular file
	}

	if runtime.GOOS == "windows" {
		return true // Existence and regular file are enough on Windows
	}

	return m.Perm()&0111 != 0 // Check execute permissions on Unix-like
}

func getExecutableExtensions() []string {
	if runtime.GOOS == "windows" {
		pathExt := os.Getenv("PATHEXT")
		if pathExt != "" {
			exts := strings.Split(strings.ToLower(pathExt), ";")
			// Sanitize extensions (remove leading dots if present)
			for i, ext := range exts {
				if strings.HasPrefix(ext, ".") {
					exts[i] = ext[1:]
				}
			}
			return exts
		}
		return []string{"com", "exe", "bat", "cmd"} // Sensible defaults
	}
	return []string{""} // No extensions on Unix-like
}

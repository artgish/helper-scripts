package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	Latest   string
	Previous string
	Current  string
)

func parseSemver(version string) (major int, minor int, patch int, hasVPrefix bool, err error) {
	if version == "" {
		return 0, 0, 0, false, fmt.Errorf("empty version")
	}
	hasVPrefix = len(version) > 0 && (version[0] == 'v' || version[0] == 'V')
	if hasVPrefix {
		version = version[1:]
	}
	// Strip pre-release and build metadata if present
	if idx := strings.IndexByte(version, '-'); idx >= 0 {
		version = version[:idx]
	}
	if idx := strings.IndexByte(version, '+'); idx >= 0 {
		version = version[:idx]
	}
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return 0, 0, 0, hasVPrefix, fmt.Errorf("invalid semver: %q", version)
	}
	major, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, hasVPrefix, fmt.Errorf("invalid major: %w", err)
	}
	minor, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, hasVPrefix, fmt.Errorf("invalid minor: %w", err)
	}
	patch, err = strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, 0, hasVPrefix, fmt.Errorf("invalid patch: %w", err)
	}
	return major, minor, patch, hasVPrefix, nil
}

func formatSemver(major int, minor int, patch int, withV bool) string {
	if withV {
		return fmt.Sprintf("v%d.%d.%d", major, minor, patch)
	}
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}

func computeApplied(previous string, latest string, current string) (string, error) {
	prevMajor, prevMinor, prevPatch, _, err := parseSemver(previous)
	if err != nil {
		return "", err
	}
	latestMajor, latestMinor, latestPatch, _, err := parseSemver(latest)
	if err != nil {
		return "", err
	}
	currMajor, currMinor, currPatch, hasV, err := parseSemver(current)
	if err != nil {
		return "", err
	}

	// Determine highest-order change and apply same delta to current
	if latestMajor != prevMajor {
		delta := latestMajor - prevMajor
		newMajor := currMajor + delta
		if newMajor < 0 {
			return "", fmt.Errorf("resulting major would be negative")
		}
		return formatSemver(newMajor, 0, 0, hasV), nil
	}
	if latestMinor != prevMinor {
		delta := latestMinor - prevMinor
		newMinor := currMinor + delta
		if newMinor < 0 {
			return "", fmt.Errorf("resulting minor would be negative")
		}
		return formatSemver(currMajor, newMinor, 0, hasV), nil
	}
	// Patch-level change (or no change)
	delta := latestPatch - prevPatch
	newPatch := currPatch + delta
	if newPatch < 0 {
		return "", fmt.Errorf("resulting patch would be negative")
	}
	return formatSemver(currMajor, currMinor, newPatch, hasV), nil
}

func loadArgs() (string, string, string) {
	flag.StringVar(&Latest, "latest", "", "latest semver (second positional arg)")
	flag.StringVar(&Previous, "previous", "", "previous semver (first positional arg)")
	flag.StringVar(&Current, "current", "", "current semver to apply change to (third positional arg)")
	flag.Parse()

	if Latest != "" && Previous != "" && Current != "" {
		return Previous, Latest, Current
	}
	args := flag.Args()
	if len(args) >= 3 {
		// Positional: previous, latest, current
		return args[0], args[1], args[2]
	}
	fmt.Fprintln(os.Stderr, "Usage: get-forth-semver [previous latest current] or with flags -previous -latest -current")
	os.Exit(2)
	return "", "", ""
}

func main() {
	previous, latest, current := loadArgs()
	out, err := computeApplied(previous, latest, current)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
}

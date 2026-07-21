// Package buildinfo provides one source of truth for CLI and API version data.
package buildinfo

import (
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

// Info describes the source revision used to build Chronicle.
type Info struct {
	Version       string `json:"version"`
	ModuleVersion string `json:"module_version,omitempty"`
	Commit        string `json:"commit,omitempty"`
	CommitDate    string `json:"commit_date,omitempty"`
	BuildTime     string `json:"build_time,omitempty"`
	Modified      bool   `json:"modified,omitempty"`
}

var (
	injectedMu sync.RWMutex
	injected   Info
)

// Set records values injected by a local build or GoReleaser.
func Set(version, commit, commitDate, buildTime string) {
	injectedMu.Lock()
	injected = Info{
		Version:    version,
		Commit:     shortenCommit(commit),
		CommitDate: shortDate(commitDate),
		BuildTime:  buildTime,
	}
	injectedMu.Unlock()
}

// Current returns build information injected with ldflags, supplemented by Go's
// module and VCS build metadata. A Go pseudo-version is presented as "dev" while
// its timestamp and revision are exposed as human-readable fields.
func Current() Info {
	injectedMu.RLock()
	result := injected
	injectedMu.RUnlock()

	moduleVersion := ""
	settings := map[string]string{}
	if raw, ok := debug.ReadBuildInfo(); ok {
		moduleVersion = raw.Main.Version
		for _, setting := range raw.Settings {
			settings[setting.Key] = setting.Value
		}
	}

	return resolve(result, moduleVersion, settings)
}

func resolve(result Info, moduleVersion string, settings map[string]string) Info {
	if moduleVersion != "" && moduleVersion != "(devel)" {
		result.ModuleVersion = moduleVersion
	}

	if result.Version == "" {
		result.Version = moduleVersion
	}

	if commit := settings["vcs.revision"]; result.Commit == "" && commit != "" {
		result.Commit = shortenCommit(commit)
	}
	if commitTime := settings["vcs.time"]; result.CommitDate == "" && commitTime != "" {
		result.CommitDate = shortDate(commitTime)
	}
	result.Modified = settings["vcs.modified"] == "true" || strings.HasSuffix(moduleVersion, "+dirty")

	pseudoVersion := strings.TrimSuffix(moduleVersion, "+dirty")
	if commit, commitDate, ok := parsePseudoVersion(pseudoVersion); ok {
		if result.Commit == "" {
			result.Commit = commit
		}
		if result.CommitDate == "" {
			result.CommitDate = commitDate
		}
		if result.Version == moduleVersion || result.Version == pseudoVersion {
			result.Version = "dev"
		}
	}

	result.Version = normalizeVersion(result.Version)
	return result
}

// String formats the concise, human-facing version line.
func (info Info) String() string {
	details := make([]string, 0, 3)
	if info.Commit != "" {
		details = append(details, "commit "+info.Commit)
	}
	if info.CommitDate != "" {
		details = append(details, info.CommitDate)
	}
	if info.Modified {
		details = append(details, "modified")
	}

	line := "chronicle " + normalizeVersion(info.Version)
	if len(details) > 0 {
		line += " (" + strings.Join(details, ", ") + ")"
	}
	return line
}

func normalizeVersion(version string) string {
	version = strings.TrimSpace(version)
	if version == "" || version == "(devel)" || version == "(unknown)" {
		return "dev"
	}
	gitVersion := strings.TrimSuffix(version, "-dirty")
	if (len(gitVersion) >= 7 && isHex(gitVersion)) || isGitDescribeVersion(gitVersion) {
		return "dev"
	}
	if version[0] >= '0' && version[0] <= '9' {
		return "v" + version
	}
	return version
}

func isGitDescribeVersion(version string) bool {
	marker := strings.LastIndex(version, "-g")
	return marker > 0 && len(version[marker+2:]) >= 7 && isHex(version[marker+2:])
}

func shortenCommit(commit string) string {
	commit = strings.TrimSpace(commit)
	if len(commit) > 12 {
		return commit[:12]
	}
	return commit
}

func shortDate(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= len("2006-01-02") {
		return value[:len("2006-01-02")]
	}
	return value
}

func parsePseudoVersion(version string) (commit, commitDate string, ok bool) {
	parts := strings.Split(version, "-")
	if len(parts) < 3 {
		return "", "", false
	}
	timestamp := parts[len(parts)-2]
	commit = parts[len(parts)-1]
	if len(timestamp) != len("20060102150405") || len(commit) < 7 {
		return "", "", false
	}
	parsed, err := time.Parse("20060102150405", timestamp)
	if err != nil || !isHex(commit) {
		return "", "", false
	}
	return shortenCommit(commit), parsed.Format("2006-01-02"), true
}

func isHex(value string) bool {
	for _, char := range value {
		if !strings.ContainsRune("0123456789abcdefABCDEF", char) {
			return false
		}
	}
	return true
}

// API returns the stable field names exposed by /api/v1/version.
func (info Info) API() map[string]string {
	return map[string]string{
		"version":        info.Version,
		"module_version": info.ModuleVersion,
		"git_commit":     info.Commit,
		"git_date":       info.CommitDate,
		"build_time":     info.BuildTime,
		"modified":       fmt.Sprintf("%t", info.Modified),
	}
}

package qdrant

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const fullVersionParts = 3
const reducedVersionParts = 2
const unknownVersion = "Unknown"

type Version struct {
	Major int
	Minor int
}

func getServerVersion(clientConn *GrpcClient) string {
	logger := slog.Default()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	healthCheckResult, err := clientConn.qdrant.HealthCheck(ctx, &HealthCheckRequest{})
	if err != nil {
		logger.Warn("Unable to get server version, use default", "err", err, "default", unknownVersion)
		return unknownVersion
	}
	serverVersion := healthCheckResult.GetVersion()

	return serverVersion
}

func removeLeadingNonNumeric(versionStr string) string {
	return strings.TrimLeftFunc(versionStr, func(r rune) bool {
		return !unicode.IsDigit(r)
	})
}

// ParseVersion converts a version string "x.y[.z]" into a Version struct.
func ParseVersion(versionStr string) (*Version, error) {
	cleanedVersionStr := removeLeadingNonNumeric(versionStr)
	parts := strings.SplitN(cleanedVersionStr, ".", fullVersionParts)
	if len(parts) < reducedVersionParts {
		return nil, fmt.Errorf("unable to parse version, expected format: x.y[.z], found: %s", cleanedVersionStr)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse major version: %w", err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse minor version: %w", err)
	}

	return &Version{
		Major: major,
		Minor: minor,
	}, nil
}

func IsCompatible(clientVersion, serverVersion string) bool {
	if clientVersion == serverVersion {
		return true
	}
	logger := slog.Default()
	client, err := ParseVersion(clientVersion)
	if err != nil {
		logger.Warn("Unable to compare versions", "err", err)
		return false
	}

	server, err := ParseVersion(serverVersion)
	if err != nil {
		logger.Warn("Unable to compare versions", "err", err)
		return false
	}

	if client.Major != server.Major {
		return false
	}

	diff := client.Minor - server.Minor
	return diff <= 1 && diff >= -1
}

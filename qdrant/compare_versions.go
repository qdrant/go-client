package qdrant

import (
	"context"
	"fmt"
	"log"
	"math"
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
	Rest  string
}

func getServerVersion(clientConn *GrpcClient) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	healthCheckResult, err := clientConn.qdrant.HealthCheck(ctx, &HealthCheckRequest{})
	if err != nil {
		log.Printf("Unable to get server version: %v, server version defaults to `%s`", err, unknownVersion)
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

// ParseVersion converts a version string "x.y.z" into a Version struct.
func ParseVersion(versionStr string) (*Version, error) {
	cleanedVersionStr := removeLeadingNonNumeric(versionStr)
	parts := strings.SplitN(cleanedVersionStr, ".", fullVersionParts)
	if len(parts) < reducedVersionParts {
		return nil, fmt.Errorf("unable to parse version, expected format: x.y.z, found: %s", cleanedVersionStr)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse major version: %w", err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse minor version: %w", err)
	}

	rest := ""
	if len(parts) == fullVersionParts {
		rest = parts[2]
	}

	return &Version{
		Major: major,
		Minor: minor,
		Rest:  rest,
	}, nil
}

func IsCompatible(clientVersion, serverVersion *string) bool {
	if *clientVersion == *serverVersion {
		return true
	}

	parsedClientVersion, err := ParseVersion(*clientVersion)
	if err != nil {
		log.Printf("Unable to compare versions: %v", err)
		return false
	}
	parsedServerVersion, err := ParseVersion(*serverVersion)
	if err != nil {
		log.Printf("Unable to compare versions: %v", err)
		return false
	}
	majorDiff := int(math.Abs(float64(parsedClientVersion.Major - parsedServerVersion.Major)))
	if majorDiff >= 1 {
		return false
	}
	return int(math.Abs(float64(parsedClientVersion.Minor-parsedServerVersion.Minor))) <= 1
}

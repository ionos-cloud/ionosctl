package doc

import (
	"strings"
)

// FlagDefaultHandler is a Strategy for handling pflag default values while generating docs
type FlagDefaultHandler func(flagDescription, defaultValue string) string

// MaintenanceHandler is a concrete strategy for --maintenance-day
func MaintenanceHandler(flagDescription, defaultValue string) string {
	if strings.Contains(flagDescription, "Defaults to a random day") {
		return "Random (Mon-Fri 10:00-16:00)"
	}
	return defaultValue
}

func DataplatformUsesLatestVersion(flagDescription, defaultValue string) string {
	if strings.Contains(flagDescription, "dataplatform") {
		return "same as 'dataplatform version active'"
	}
	return defaultValue
}

func getStrategyForFlag(flagName string) FlagDefaultHandler {
	switch flagName {
	case "maintenance-day", "maintenance-time":
		return MaintenanceHandler
	case "version":
		return DataplatformUsesLatestVersion
	default:
		return nil
	}
}

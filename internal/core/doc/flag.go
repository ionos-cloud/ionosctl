package doc

import (
	"strings"
)

// FlagDefaultDocumentationHandler is a Strategy for handling pflag default values while generating docs
type FlagDefaultDocumentationHandler func(flagDescription, defaultValue string) string

// RandomDayDescriptionHandler is a concrete strategy which changes the default value
// to "Random (Mon-Fri 10:00-16:00)" instead of an actual random day i.e. Thursday
func RandomDayDescriptionHandler(flagDescription, defaultValue string) string {
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

func CustomNamesHandler(flagDescription, defaultValue string) string {
	if strings.Contains(flagDescription, "custom names for each spec") {
		return "<Overriden with sdk-go-bundle product names: [authentication=auth, certificatemanager=cert, cloud=compute, object‑storage=objectstorage, object‑storage‑management=objectstoragemanagement, mongodb=mongo, postgresql=psql]>"
	}
	return defaultValue
}

func getStrategyForFlag(flagName string) FlagDefaultDocumentationHandler {
	switch flagName {
	case "maintenance-day", "maintenance-time":
		return RandomDayDescriptionHandler
	case "garbage-collection-schedule-days", "garbage-collection-schedule-time":
		return RandomDayDescriptionHandler
	case "version":
		return DataplatformUsesLatestVersion
	case "custom-names":
		return CustomNamesHandler
	default:
		return nil
	}
}

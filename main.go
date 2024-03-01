// Package main is a standalone utility that allows a user to parse a CSV file that contains Mattermost User IDs,
// and convert them to either usernames or full names.
package main

import (
	"flag"
	"log"
	"os"
)

var debugMode bool = false
var fullnameMode bool = false

// LogLevel is used to refer to the type of message that will be written using the logging code.
type LogLevel string

const (
	debugLevel   LogLevel = "DEBUG"
	infoLevel    LogLevel = "INFO"
	warningLevel LogLevel = "WARNING"
	errorLevel   LogLevel = "ERROR"
)

const (
	defaultPort   = "8065"
	defaultScheme = "http"
)

// Logging functions

// LogMessage logs a formatted message to stdout or stderr
func LogMessage(level LogLevel, message string) {
	if level == errorLevel {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(os.Stdout)
	}
	log.SetFlags(log.Ldate | log.Ltime)
	log.Printf("[%s] %s\n", level, message)
}

// DebugPrint allows us to add debug messages into our code, which are only printed if we're running in debug more.
// Note that the command line parameter '-debug' can be used to enable this at runtime.
func DebugPrint(message string) {
	if debugMode {
		LogMessage(debugLevel, message)
	}
}

// getEnvWithDefaults allows us to retrieve Environment variables, and to return either the current value or a supplied default
func getEnvWithDefault(key string, defaultValue interface{}) interface{} {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// Main section

func main() {

	// Parse Command Line
	DebugPrint("Parsing command line")

	var MattermostURL string
	var MattermostPort string
	var MattermostScheme string
	var MattermostToken string
	var InputCSVFilename string
	var OutputCSVFilename string
	var UserIDColumnName string
	var FullnameFlag bool
	var DebugFlag bool

	flag.StringVar(&MattermostURL, "url", "", "The URL of the Mattermost instance (without the HTTP scheme)")
	flag.StringVar(&MattermostPort, "port", "", "The TCP port used by Mattermost. [Default: "+defaultPort+"]")
	flag.StringVar(&MattermostScheme, "scheme", "", "The HTTP scheme to be used (http/https). [Default: "+defaultScheme+"]")
	flag.StringVar(&MattermostToken, "token", "", "The auth token used to connect to Mattermost")
	flag.StringVar(&InputCSVFilename, "infile", "", "*Required* The name of the CSV file to be processed")
	flag.StringVar(&OutputCSVFilename, "outfile", "", "The name of the output file that the CSV should be written to.  Note that if this parameter is omitted, the output will be written to stdout.")
	flag.StringVar(&UserIDColumnName, "column", "", "*Required* The name of the column within the CSV file that contains the user ID")
	flag.BoolVar(&FullnameFlag, "fullname", false, "Return the full name of the Mattermost user, instead of the username")
	flag.BoolVar(&DebugFlag, "debug", false, "Enable debug output")

	flag.Parse()

	// If parameters have not been passed on the command line, check for the presence of environment variables or defaults.
	if MattermostURL == "" {
		MattermostURL = getEnvWithDefault("MM_URL", "").(string)
	}
	if MattermostPort == "" {
		MattermostPort = getEnvWithDefault("MM_PORT", defaultPort).(string)
	}
	if MattermostScheme == "" {
		MattermostScheme = getEnvWithDefault("MM_SCHEME", defaultScheme).(string)
	}
	if MattermostToken == "" {
		MattermostToken = getEnvWithDefault("MM_TOKEN", "").(string)
	}
	if !DebugFlag {
		DebugFlag = getEnvWithDefault("MM_DEBUG", debugMode).(bool)
	}

	DebugPrint("Parameters: MattermostURL=" + MattermostURL + " MattermostPort=" + MattermostPort + " MattermostScheme=" + MattermostScheme + " MattermostToken=" + MattermostToken + " InputCSVFilename=" + InputCSVFilename + " OutputCSVFilename='" + OutputCSVFilename + "' UserIDColumnName='" + UserIDColumnName + "'")
	if FullnameFlag {
		DebugPrint("Fullname flag is set")
	}

	// Validate required parameters
	DebugPrint("Validating parameters")
	var cliErrors bool = false
	if MattermostURL == "" {
		LogMessage(errorLevel, "The Mattermost URL must be supplied either on the command line of vie the MM_URL environment variable")
		cliErrors = true
	}
	if MattermostScheme == "" {
		LogMessage(errorLevel, "The Mattermost HTTP scheme must be supplied either on the command line of vie the MM_SCHEME environment variable")
		cliErrors = true
	}
	if MattermostToken == "" {
		LogMessage(errorLevel, "The Mattermost auth token must be supplied either on the command line of vie the MM_TOKEN environment variable")
		cliErrors = true
	}
	if InputCSVFilename == "" {
		LogMessage(errorLevel, "The CSV input file must be supplied as a command line parameter")
		cliErrors = true
	}
	if UserIDColumnName == "" {
		LogMessage(errorLevel, "The user ID column name from the CSV must be supplied as a command line parameter")
		cliErrors = true
	}
	if cliErrors {
		flag.Usage()
		os.Exit(1)
	}

	debugMode = DebugFlag
	fullnameMode = FullnameFlag

}
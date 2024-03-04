// Package main is a standalone utility that allows a user to parse a CSV file that contains Mattermost User IDs,
// and convert them to either usernames or full names.
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
)

var debugMode bool = false
var fullnameMode bool = false

// LogLevel is used to refer to the type of message that will be written using the logging code.
type LogLevel string

type mmConnection struct {
	mmURL    string
	mmPort   string
	mmScheme string
	mmToken  string
}

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

// findStringInSlice searches for a string in a slice and returns its index.
// If the string is not found, it returns -1.
func findStringInSlice(slice []string, value string) int {
	for i, item := range slice {
		if item == value {
			return i
		}
	}
	return -1 // Not found
}

func getUserFromMattermost(mattermostCon mmConnection, userID string, fullnameFlag bool) (string, bool) {
	DebugPrint("Retrieving user data from Mattermost for user ID: " + userID)

	userData := ""

	url := fmt.Sprintf("%s://%s:%s/api/v4/users/%s", mattermostCon.mmScheme, mattermostCon.mmURL, mattermostCon.mmPort, userID)
	DebugPrint("URL to call: " + url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		LogMessage(errorLevel, "Error preparing GET")
		log.Fatal(err)
	}
	// Add the bearer token as a header
	req.Header.Add("Authorization", "Bearer "+mattermostCon.mmToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		LogMessage(errorLevel, "Failed to query Mattermost")
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Extract the body of the message
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		LogMessage(errorLevel, "Unable to extract body data from Mqattermost response")
		log.Fatal(err)
	}

	// Parse the response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		LogMessage(errorLevel, "Failed to convert body data")
		log.Fatal(err)
	}

	// Convert the data to a string to return to the calling function
	mmUserData, err := json.Marshal(result)
	if err != nil {
		LogMessage(errorLevel, "Unable to convert user data to string")
		log.Fatal(err)
	}

	username, err := jsonparser.GetString([]byte(mmUserData), "username")
	if err != nil {
		LogMessage(warningLevel, "Error processing JSON response data for user ID: "+userID)
		return "", false
	}
	userEmail, err := jsonparser.GetString([]byte(mmUserData), "email")
	if err != nil {
		LogMessage(warningLevel, "Error processing JSON response data for user ID: "+userID)
		return "", false
	}
	userFirstName, err := jsonparser.GetString([]byte(mmUserData), "first_name")
	if err != nil {
		LogMessage(warningLevel, "Error processing JSON response data for user ID: "+userID)
		return "", false
	}
	userLastName, err := jsonparser.GetString([]byte(mmUserData), "last_name")
	if err != nil {
		LogMessage(warningLevel, "Error processing JSON response data for user ID: "+userID)
		return "", false
	}
	userFullName := fmt.Sprintf("%s %s", userFirstName, userLastName)
	DebugPrint("Username: " + username + " Email: " + userEmail + " Full Name: " + userFullName)

	if fullnameFlag {
		if userFullName == " " {
			userData = username
		} else {
			userData = userFullName
		}
	} else {
		userData = username
	}

	return userData, true
}

func processCSVFile(mattermostCon mmConnection, csvInputFile string, csvOuputFIle string, userIDColumn string, fullnameFlag bool) bool {
	DebugPrint("Starting to process CSV file")

	LogMessage(infoLevel, "Processing data from file: "+csvInputFile)
	LogMessage(infoLevel, "Writing output to file:    "+csvOuputFIle)

	file, err := os.Open(csvInputFile)
	if err != nil {
		log.Fatal("Error reading inpur file", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// We need to read the header row before starting to process the rest of the file, in order to
	// identify which entry contains the user ID
	header, err := reader.Read()
	if err != nil {
		LogMessage(errorLevel, "Unable to read header record from CSV file: "+err.Error())
		return false
	}
	DebugPrint("CSV Header: " + strings.Join(header, ", "))
	index := findStringInSlice(header, userIDColumn)
	if index < 0 {
		LogMessage(errorLevel, "Unable to find column '"+userIDColumn+"' in CSV header")
		return false
	}
	DebugPrint("Selected column is at index: " + strconv.Itoa(index) + " (zero-based)")

	outfile, err := os.Create(csvOuputFIle)
	if err != nil {
		LogMessage(warningLevel, "Unable to create output file - writing to stdout")
		outfile = os.Stdout
	}

	// Initialise CSV writer
	writer := csv.NewWriter(outfile)

	defer writer.Flush()

	// Write out the header row
	writer.Write(header)

	// At this point, we've read the first line of the CSV file (the header) and we know at which
	// position the user ID column is located.  We can now process the rest of the file.

	recordsProcessed := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			LogMessage(errorLevel, "Unable to process CSV record: "+err.Error())
			return false
		}
		DebugPrint("Current record: [ " + strings.Join(record, ", ") + " ]")
		currentUserID := record[index]
		DebugPrint("User ID: " + currentUserID)
		userData, success := getUserFromMattermost(mattermostCon, currentUserID, fullnameFlag)
		if !success {
			LogMessage(warningLevel, "Error looking up User ID - skipping record!")
			continue
		}
		DebugPrint("User data from Mattermost: " + userData)

		// Now that we have the updated record, we can simply replace the relevant entry in the array
		record[index] = userData
		writer.Write(record)
		recordsProcessed += 1
	}

	if err := writer.Error(); err != nil {
		LogMessage(errorLevel, "Error writing CSV file!")
		log.Fatal(err)
	}

	processedRecordsMessage := fmt.Sprintf("Records processed: %d", recordsProcessed)
	LogMessage(infoLevel, processedRecordsMessage)

	return true
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
	flag.StringVar(&OutputCSVFilename, "outfile", "", "*Required* The name of the output file that the CSV should be written to.")
	flag.StringVar(&UserIDColumnName, "column", "", "*Required* The name of the column within the CSV file that contains the user ID")
	flag.BoolVar(&FullnameFlag, "fullname", false, "Return the full name of the Mattermost user, instead of the username (if a full name is available)")
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
	if OutputCSVFilename == "" {
		LogMessage(errorLevel, "The CSV output file must be supplied as a command line parameter")
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

	mattermostConenction := mmConnection{
		mmURL:    MattermostURL,
		mmPort:   MattermostPort,
		mmScheme: MattermostScheme,
		mmToken:  MattermostToken,
	}

	processCSVFile(mattermostConenction, InputCSVFilename, OutputCSVFilename, UserIDColumnName, fullnameMode)

	LogMessage(infoLevel, "CSV processing complete!")

}

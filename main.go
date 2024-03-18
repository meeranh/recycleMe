package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"github.com/fatih/color"
	"bufio"
)

// Program constants
const BaseURL = "https://api.zerogpt.com/api/detect/detectText"
const MaxAllowedCharacters = 10000

// To process the response from GptZero
type Response struct {
	Data struct {
		PhrasesToHumanize []string	`json:"h"`
		PercentageOfAI		float64		`json:"fakePercentage"`
	} `json:"data"`
}

// Error handling
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Getting the contents of the file
func readFile(path string) string {
	data, err := os.ReadFile(path)
	check(err)
	return string(data)
}

// File path generation
func getFilePath() string {
	if len(os.Args) < 2 {
		fmt.Println("Usage: recycleMe <file_path>")
		os.Exit(1)
	}

	// Just a welcome message :)
	fmt.Println("Loading :D")

	// Crafting the absolute file path
	absolutePath, err := filepath.Abs(os.Args[1])
	check(err)

	return absolutePath
}

// This is were the HTTP request is crafted
func makeRequest(fileContent string) []byte {

	// Headers
	const UserAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0"
	const Accept = "application/json, text/plain, */*"
	const AcceptLanguage = "en-US,en;q=0.5"
	const AcceptEncoding = "gzip, deflate, br"
	const ContentType = "application/json"
	const Origin = "https://www.zerogpt.com"
	const Connection = "keep-alive"
	const Referer = "https://www.zerogpt.com/"
	const SecFetchDest = "empty"
	const SecFetchMode = "cors"
	const SecFetchSite = "same-site"

	// Prepare the post body
	postBody, err := json.Marshal(map[string]string {
		"input_text": fileContent,
		})
	check(err)

	// Create a new request
	req, err := http.NewRequest("POST", BaseURL, bytes.NewBuffer(postBody))
	check(err)

	// Append the headers
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("Accept", Accept)
	req.Header.Add("Accept-Language", AcceptLanguage)
	req.Header.Add("Accept-Encoding", AcceptEncoding)
	req.Header.Add("Content-Type", ContentType)
	req.Header.Add("Origin", Origin)
	req.Header.Add("Connection", Connection)
	req.Header.Add("Referer", Referer)
	req.Header.Add("Sec-Fetch-Dest", SecFetchDest)
	req.Header.Add("Sec-Fetch-Mode", SecFetchMode)
	req.Header.Add("Sec-Fetch-Site", SecFetchSite)

	resp, err := http.DefaultClient.Do(req)
	check(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	check(err)
	return body
}

// Extracting the content details from the response
func extractContentDetails(rawJson []byte) ([]string, float64) {
	var parsedJson Response
	json.Unmarshal(rawJson, &parsedJson)
	return parsedJson.Data.PhrasesToHumanize, parsedJson.Data.PercentageOfAI
}

// Creating a duplicate file for editing
func createDuplicateFile(originalFilePath string, fileContent string) string {

	// Extract file data
	originalFileName := filepath.Base(originalFilePath)
	fileExtension := filepath.Ext(originalFilePath)
	fileName := strings.TrimSuffix(originalFileName, fileExtension)

	// Create a new duplicate file in the current working directory
	newFileName := fileName + "_humanized" + fileExtension
	currentWorkingDirectory, err := os.Getwd()
	check(err)
	newFilePath := filepath.Join(currentWorkingDirectory, newFileName)
	newFile, err := os.Create(newFilePath)
	check(err)

	// Write to the file and close it
	_, err = newFile.WriteString(fileContent)
	check(err)
	newFile.Close()

	return newFilePath
}

// Replacing AI content in the duplicate file
func replaceString(replaceString string, matchString string, filePath string) {
	fileContent := readFile(filePath)
	newFileContent := strings.ReplaceAll(fileContent, matchString, replaceString)
	err := os.WriteFile(filePath, []byte(newFileContent), 0644)
	check(err)
}

// The iteration that starts the humanization process
func startHumanization(arrToHumanize []string, newFilePath string) {

	// Colors
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)

	fmt.Printf("Rewrite these sentences in your own words\n")

	scanner := bufio.NewScanner(os.Stdin)

	for i, v := range arrToHumanize {

		yellow.Printf("\n[%d/%d] ", i+1, len(arrToHumanize))
		blue.Printf("%s\n", v)
		green.Printf("Humanized: ")

		// Taking in user input from stdin
		var userInput string
		scanner.Scan()
		userInput = scanner.Text()

		// Editing the duplicated file
		replaceString(userInput, v, newFilePath)
	}
}

// Checks if the text needs to be segmented
func isTooBig(fileContent string) bool {
	return len(fileContent) > MaxAllowedCharacters
}

// Segments the text into allowable chunks
func segmentText(fileContent string) []string {
	var segments []string
	segment := ""
	for _, char := range fileContent {
		segment += string(char)
		if len(segment) > MaxAllowedCharacters {
			segments = append(segments, segment)
			segment = ""
		}
	}
	segments = append(segments, segment)
	return segments
}

// Send the segmented text to ZeroGPT
func sendSegments(segments []string) ([]string, float64) {
	toHumanizeSlice := make([]string, 0)
	avgAiPercentage := 0.0
	for i, segment := range segments {
		fmt.Printf("[%d/%d] Sending current segment\n", i+1, len(segments))
		rawJson := makeRequest(segment)
		stringsToHumanize, aiPercentage := extractContentDetails(rawJson)
		toHumanizeSlice = append(toHumanizeSlice, stringsToHumanize...)
		avgAiPercentage += aiPercentage / float64(len(segments))
	}
	return toHumanizeSlice, avgAiPercentage
}

func main() {
	// Preparing a duplicated file for editing
	path := getFilePath()
	fileContent := readFile(path)
	newFilePath := createDuplicateFile(path, fileContent)

	// Initialize slice & float for the humanization process
	stringsToHumanize := make([]string, 0)
	aiPercentage := 0.0

	// Check if file needs segmentation
	if isTooBig(fileContent) {
		fmt.Printf("Uh oh! The file is too big. We need to segment it.\n\n")
		segments := segmentText(fileContent)
		stringsToHumanize, aiPercentage = sendSegments(segments)
	} else {
		rawJson := makeRequest(fileContent)
		stringsToHumanize, aiPercentage = extractContentDetails(rawJson)
	}

	// Printing out the flagged AI percentage
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("\nAI Percentage (%%): %s\n", red(aiPercentage))

	// Starting the editing iteration
	startHumanization(stringsToHumanize, newFilePath)
}

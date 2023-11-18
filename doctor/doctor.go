package doctor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const cloudflareAPIEndpoint = "https://api.cloudflare.com/client/v4/accounts/0843a42dc7915eb7a5ca1e3bb05cfce2/ai/run/@cf/meta/llama-2-7b-chat-int8"
const accessToken = "wwJH-Ts41Wquh49IXlT50krftRa4E3w7SbAq7gL3"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}


func sendToCloudflare(mdFilename string) error {
    // Read Markdown file content
    content, err := ioutil.ReadFile(mdFilename)
    if err != nil {
        return err
    }

    // Create messages slice
   messages := []Message{
    {"system", "You specialize in generating Mermaid Markdown syntax for Go projects."},
    {"user", fmt.Sprintf("Generate a Mermaid Markdown overview for the following Go file:\n```go\n%s\n```", content)},
}

    // Prepare input JSON
    inputs := map[string]interface{}{"messages": messages}
    inputJSON, err := json.Marshal(inputs)
    if err != nil {
        return err
    }

    // Create HTTP request
    req, err := http.NewRequest("POST", cloudflareAPIEndpoint, bytes.NewBuffer(inputJSON))
    if err != nil {
        return err
    }

    // Set Authorization and Content-Type headers
    req.Header.Set("Authorization", "Bearer "+accessToken)
    req.Header.Set("Content-Type", "application/json")

    // Make the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatalf("Error making request to Cloudflare API: %v", err)
        return err
    }
    defer resp.Body.Close()
// Unmarshal JSON response into the result struct
var result struct {
    Result struct {
        Response string `json:"response"`
    } `json:"result"`
}

  responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Error reading response body: %v", err)
        return err
    }

	log.Println(string(responseBody))


if err := json.Unmarshal(responseBody, &result); err != nil {
    log.Fatalf("Error unmarshaling JSON response: %v", err)
    return err
}

log.Println(result.Result.Response)

// Append the response to the Markdown file
responseMarkdown := fmt.Sprintf("```mermaid\n\n%s\n\n```\n", result.Result.Response)
newContent := string(content) + "\n" + responseMarkdown

    // Update the Markdown file
    err = os.WriteFile(mdFilename, []byte(newContent), 0644)
    if err != nil {
        return err
    }

    return nil
}



 func Run(filePath string)  {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run script.go <filepath>")
		return
	}

	filepathParam := os.Args[1]

	// Get filename without extension
	filename := strings.TrimSuffix(filepath.Base(filepathParam), filepath.Ext(filepathParam))

	// Create JSON filename
	jsonFilename := filename + ".json"

	// Create MD filename
	mdFilename := filename + ".md"

	// Create .doctor.md filename
	doctorFilename := filename + ".doctor.md"

	// Require and process JSON file
	processJSONFile(jsonFilename)

	// Require and log MD file for each key-value pair in JSON
	logMDFile(mdFilename)

	// Write logs to .doctor.md file
	writeLogsToDoctorFile(doctorFilename)
}

func processJSONFile(filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(file, &data); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Processing %s:\n", filename)
	for key, value := range data {
		// Perform any specific processing/logic here based on key-value pairs
		fmt.Printf("Key: %s, Value: %v\n", key, value)
	}
}

func logMDFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Log to MD file for each key-value pair in JSON
	log.SetOutput(file)

	// Example log message, modify as needed
	log.Println("Log entry for MD file.")
}

func writeLogsToDoctorFile(filename string) {
	// Open or create .doctor.md file
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read MD log file
	mdLogContent, err := ioutil.ReadFile(filename + ".md")
	if err != nil {
		log.Fatal(err)
	}

	// Write MD log content to .doctor.md file
	_, err = file.Write(mdLogContent)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Logs written to %s\n", filename)
}

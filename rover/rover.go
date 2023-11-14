package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"log"
	"net/http"
)

const cloudflareAPIEndpoint = "https://api.cloudflare.com/client/v4/accounts/0843a42dc7915eb7a5ca1e3bb05cfce2/ai/run/@cf/meta/llama-2-7b-chat-int8"
const accessToken = "wwJH-Ts41Wquh49IXlT50krftRa4E3w7SbAq7gL3"

func convertToMarkdown(filename string) string {
	// Your logic to convert file content to markdown goes here
	// For simplicity, let's just add a markdown extension to the original filename
	return strings.TrimSuffix(filename, filepath.Ext(filename)) + ".md"
}

func copyFileWithCodeBlock(src, dest, fileExtension string) error {
	if fileExtension == "md" {
		// Skip processing .md files
		return nil
	}

	content, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	// Wrap the content with the specified markdown code block without the dot in the file extension
	wrappedContent := fmt.Sprintf("```%s\n\n%s\n\n```\n", fileExtension, string(content))

	return ioutil.WriteFile(dest, []byte(wrappedContent), 0644)
}


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

  responseBody, err := ioutil.ReadAll(resp.Body)
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

func findNearestGitignore(dir string) (string, error) {
	for dir != "/" {
		gitignorePath := filepath.Join(dir, ".gitignore")
		if _, err := os.Stat(gitignorePath); err == nil {
			return gitignorePath, nil
		}
		dir = filepath.Dir(dir)
	}
	return "", nil
}

func isPathIgnored(path, gitignorePath string) (bool, error) {
	gitignoreData, err := ioutil.ReadFile(gitignorePath)
	if err != nil {
		return false, err
	}

	gitignoreLines := strings.Split(string(gitignoreData), "\n")
	for _, pattern := range gitignoreLines {
		// Skip comments and empty lines
		if strings.HasPrefix(pattern, "#") || strings.TrimSpace(pattern) == "" {
			continue
		}

		matched, err := filepath.Match(pattern, path)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}

	return false, nil
}

func processFile(path string, info os.FileInfo, err error) error {
		if err != nil {
		return err
	}

	// Find the nearest .gitignore file
	gitignorePath, err := findNearestGitignore(filepath.Dir(path))
	if err != nil {
		return err
	}

	// Check if the file should be ignored based on .gitignore rules
	if gitignorePath != "" {
		matches, err := isPathIgnored(path, gitignorePath)
		if err != nil {
			return err
		}
		if matches {
			fmt.Printf("Skipping %s based on .gitignore rules\n", path)
			return nil
		}
	}

	// Check if the file is an .md file
	if filepath.Ext(path) == ".md" {
		fmt.Printf("Skipping %s (Markdown file)\n", path)
		return nil
	}

	// Check if the file is an .mod file
	if filepath.Ext(path) == ".mod" {
		fmt.Printf("Skipping %s (Mod file)\n", path)
		return nil
	}

		// Check if the file is an .sum file
	if filepath.Ext(path) == ".sum" {
		fmt.Printf("Skipping %s (Sum file)\n", path)
		return nil
	}


	if info.IsDir() {
		return nil
	}
	mdFilename := convertToMarkdown(path)
	fmt.Printf("Copying %s to %s\n", path, mdFilename)

	fileExtension := strings.TrimPrefix(filepath.Ext(path), ".")
	if err := copyFileWithCodeBlock(path, mdFilename, fileExtension); err != nil {
		return err
	}

	if err := sendToCloudflare(mdFilename); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// func processFile(path string, info os.FileInfo, err error) error {
// 	if err != nil {
// 		return err
// 	}

// 	if info.IsDir() {
// 		return nil
// 	}

// 	mdFilename := convertToMarkdown(path)
// 	fmt.Printf("Copying %s to %s\n", path, mdFilename)

// 	fileExtension := strings.TrimPrefix(filepath.Ext(path), ".")
// 	if err := copyFileWithCodeBlock(path, mdFilename, fileExtension); err != nil {
// 		return err
// 	}

// 	if err := sendToCloudflare(mdFilename); err != nil {
//        log.Fatal(err)
// 		return err
// 	}

// 	return nil
// }

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <repository_path>")
		return
	}

	repoPath := os.Args[1]
	err := filepath.Walk(repoPath, processFile)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

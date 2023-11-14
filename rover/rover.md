```go

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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


func processFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	// Skip processing .md files
	if filepath.Ext(path) == ".md" {
		fmt.Printf("Skipping %s\n", path)
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


```

```mermaid

To generate the Mermaid Markdown overview for the Go file, you can use the following steps:
1. Read the Go file content using `ioutil.ReadFile()` function.
2. Create a `Message` struct to hold the Markdown content and the original file content.
3. Use the `json` package to marshal the `Message` struct into a JSON object.
4. Create an HTTP request to the Cloudflare API with the JSON object as the request body.
5. Set the `Authorization` and `Content-Type` headers with the appropriate values.
6. Make the HTTP request using the `http.Client` type.
7. Unmarshal the JSON response into the `result` struct using the `json` package.
8. Print the `response` field of the `result` struct to the console.
9. Append the response to the Markdown file using the `fmt` package.
10. Update the Markdown file content using the `os` package.

Here is an example of how the Mermaid Markdown overview for the Go file might look like:
```mermaid
graph LR
    A[

```

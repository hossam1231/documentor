package geek

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// Create a struct to hold analysis results
type AnalysisResult struct {
	PackageName     string   `json:"packageName"`
	Imports         []string `json:"imports"`
	Structs         []string `json:"structs"`
	Variables       []string `json:"variables"`
	Constants       []string `json:"constants"`
	Comments        []string `json:"comments"`
	Interfaces      []string `json:"interfaces"`
	Methods         []string `json:"methods"`
	Channels        []string `json:"channels"`
	ErrorHandling   []string `json:"errorHandling"`
	TypeAssertions  []string `json:"typeAssertions"`
	ControlFlow     []string `json:"controlFlow"`
	DeferStatements []string `json:"deferStatements"`
	PanicRecover    []string `json:"panicRecover"`
	FunctionCalls   []string `json:"functionCalls"`
}

func removeDuplicates[T comparable](input []T) []T {
    encountered := map[T]bool{}
    result := []T{}

    for _, value := range input {
        if !encountered[value] {
            encountered[value] = true
            result = append(result, value)
        }
    }

    return result
}

func Run(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file:", err)
		return
	}
	defer file.Close()

	// Specify the target file path in the same directory with a different filename
	targetFileName := filePath+ ".json"
	targetFilePath := filepath.Join(filepath.Dir(filePath), targetFileName)

	importRegex := regexp.MustCompile(`^\s*import\s+\(\s*"([^"]+)"\s*\)`)
	var imports []string
	commentRegex := regexp.MustCompile(`^\s*\/\/(.*)|^\s*\/\*([^*]*\*\/)`)
	structRegex := regexp.MustCompile(`^\s*type\s+([a-zA-Z_]\w*)\s+struct\s*{`)
	varRegex := regexp.MustCompile(`var\s+([a-zA-Z_]\w*)\s+([a-zA-Z_]\w*)`)
	constRegex := regexp.MustCompile(`const\s+([a-zA-Z_]\w*)\s+([a-zA-Z_]\w*)`)
	interfaceRegex := regexp.MustCompile(`type\s+([a-zA-Z_]\w*)\s+interface\s*{`)
	methodRegex := regexp.MustCompile(`func\s+\(([a-zA-Z_]\w*)\s*\*?([a-zA-Z_]\w*)\)\s*([a-zA-Z_]\w*)\(`)
	channelRegex := regexp.MustCompile(`(make\()?(chan<-[^\s]*)`)
	packageRegex := regexp.MustCompile(`^package\s+([a-zA-Z_]\w*)`)
	errorHandlingRegex := regexp.MustCompile(`(\w+),\s*(\w+)\s*:=\s*(\w+)\(.*\)`)
	typeAssertionRegex := regexp.MustCompile(`([a-zA-Z_]\w*)\s*:=\s*\(([^)]+)\)`)
	controlFlowRegex := regexp.MustCompile(`(if|else|switch|case|default|select)\s*{`)
	deferRegex := regexp.MustCompile(`defer\s+([a-zA-Z_]\w*)\(`)
	panicRecoverRegex := regexp.MustCompile(`(panic|recover)\(.*\)`)
	functionCallRegex := regexp.MustCompile(`([a-zA-Z_]\w*)\(`)

	var (
		structs         []string
		variables       []string
		constants       []string
		comments        []string
		interfaces      []string
		methods         []string
		channels        []string
		packageName     string
		errorHandling   []string
		typeAssertions  []string
		controlFlow     []string
		deferStatements []string
		panicRecover    []string
		functionCalls   []string
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if match := importRegex.FindStringSubmatch(line); len(match) > 1 {
			imports = append(imports, match[1])
		}

		if match := commentRegex.FindStringSubmatch(line); len(match) > 0 {
			comments = append(comments, match[1])
		}

		// Check for imports
		if match := importRegex.FindStringSubmatch(line); len(match) > 1 {
			imports = append(imports, match[1])
		}

		// Check for structs
		if match := structRegex.FindStringSubmatch(line); len(match) > 1 {
			structs = append(structs, match[1])
		}

		// Check for variables
		if match := varRegex.FindStringSubmatch(line); len(match) > 2 {
			variables = append(variables, match[2])
		}

		// Check for constants
		if match := constRegex.FindStringSubmatch(line); len(match) > 2 {
			constants = append(constants, match[2])
		}

		// Check for comments
		if match := commentRegex.FindStringSubmatch(line); len(match) > 0 {
			comments = append(comments, match[1])
		}

		// Check for interfaces
		if match := interfaceRegex.FindStringSubmatch(line); len(match) > 1 {
			interfaces = append(interfaces, match[1])
		}

		// Check for channels
		if match := channelRegex.FindStringSubmatch(line); len(match) > 1 {
			channels = append(channels, match[2])
		}

		// Check for package name
		if match := packageRegex.FindStringSubmatch(line); len(match) > 1 {
			packageName = match[1]
		}

		// Check for error handling
		if match := errorHandlingRegex.FindStringSubmatch(line); len(match) > 2 {
			errorHandling = append(errorHandling, match[1], match[2], match[3])
		}

		// Check for type assertions
		if match := typeAssertionRegex.FindStringSubmatch(line); len(match) > 1 {
			typeAssertions = append(typeAssertions, match[1])
		}

		// Check for control flow statements
		if match := controlFlowRegex.FindStringSubmatch(line); len(match) > 0 {
			controlFlow = append(controlFlow, match[1])
		}

		// Check for defer statements
		if match := deferRegex.FindStringSubmatch(line); len(match) > 1 {
			deferStatements = append(deferStatements, match[1])
		}

		// Check for panic/recover statements
		if match := panicRecoverRegex.FindStringSubmatch(line); len(match) > 0 {
			panicRecover = append(panicRecover, match[1])
		}

		// Check for methods
if match := methodRegex.FindStringSubmatch(line); len(match) > 2 {
    methods = append(methods, match[1]+"."+match[3]) // Assuming you want the full method name with the struct/interface name
}

// Check for function calls
if match := functionCallRegex.FindStringSubmatch(line); len(match) > 1 {
    functionCalls = append(functionCalls, match[1])
}

	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading file:", err)
		return
	}

// Usage
result := AnalysisResult{
    PackageName:     packageName,
    Imports:         removeDuplicates(imports),
    Structs:         removeDuplicates(structs),
    Variables:       removeDuplicates(variables),
    Constants:       removeDuplicates(constants),
    Comments:        removeDuplicates(comments),
    Interfaces:      removeDuplicates(interfaces),
    Methods:         removeDuplicates(methods),
    Channels:        removeDuplicates(channels),
    ErrorHandling:   removeDuplicates(errorHandling),
    TypeAssertions:  removeDuplicates(typeAssertions),
    ControlFlow:     removeDuplicates(controlFlow),
    DeferStatements: removeDuplicates(deferStatements),
    PanicRecover:    removeDuplicates(panicRecover),
    FunctionCalls:   removeDuplicates(functionCalls),
}
	// Create or open the target file
	targetFile, err := os.Create(targetFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating file:", err)
		return
	}
	defer targetFile.Close()



	
// Convert the result to JSON
resultJSON, err := json.MarshalIndent(result, "", "  ")
if err != nil {
    fmt.Fprintln(os.Stderr, "Error marshalling result to JSON:", err)
    return
}

// Write the JSON to the target file
_, err = targetFile.Write(resultJSON)
if err != nil {
    fmt.Fprintln(os.Stderr, "Error writing JSON to file:", err)
    return
}

// Close the file and print a message
fmt.Println("Analysis result written to", targetFilePath)

	fmt.Println(result)
}

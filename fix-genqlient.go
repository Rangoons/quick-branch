package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run scripts/fix-genqlient.go <path-to-generated-file>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	// Read the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	file.Close()

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Pattern to match slice fields with json tags that don't have omitempty
	// Example: FieldName []string `json:"fieldName"`
	pattern := regexp.MustCompile(`^(\s+\w+\s+\[\][\w\*]+\s+` + "`" + `json:"[^"]+)"` + "`" + `$`)

	modified := false
	for i, line := range lines {
		if match := pattern.FindStringSubmatch(line); match != nil {
			// Check if omitempty is already present
			if !strings.Contains(line, "omitempty") {
				// Add omitempty before the closing backtick
				lines[i] = match[1] + ",omitempty\"`"
				modified = true
			}
		}
	}

	if !modified {
		fmt.Println("No changes needed.")
		return
	}

	// Write back to file
	output, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}
	defer output.Close()

	writer := bufio.NewWriter(output)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()

	fmt.Println("Successfully added omitempty to slice fields.")
}

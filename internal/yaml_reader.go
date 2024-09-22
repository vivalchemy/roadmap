package internal

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Subject represents a subject with topics.
type Subject struct {
	Name   string  `yaml:"name"`
	Topics []Topic `yaml:"topics"`
}

// Topic represents a topic with subtopics.
type Topic struct {
	Name      string   `yaml:"name"`
	Subtopics []string `yaml:"subtopics"`
}

/*
ReadYAMLFile reads a YAML file and returns a Subject struct.

Example usage:

	subject, err := ReadYAMLFile("subject.yaml")
	if err != nil {
	    log.Fatalf("Failed to read YAML file: %v", err)
	}
	fmt.Printf("Subject Name: %s\n", subject.Name)
*/
func ReadYAMLFile(fileName string) (Subject, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return Subject{}, err
	}
	if len(file) == 0 {
		return Subject{}, fmt.Errorf("error: file is empty")
	}

	var subject Subject
	err = yaml.Unmarshal(file, &subject)
	if err != nil {
		return Subject{}, err
	}

	return subject, nil
}

/*
ValidateFileName checks if the provided filename has a valid YAML extension.

Example usage:

	filename := "example.yaml"
	isValid := ValidateFileName(filename)
	if isValid {
	    fmt.Println("Valid YAML file.")
	} else {
	    fmt.Println("Invalid file. Please provide a .yaml or .yml file.")
	}
*/
func ValidateFileName(fileName string) bool {
	return len(fileName) >= 5 && (strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml"))
}

package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/vivalchemy/roadmap/internal"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Get the file name from the command line
	if len(os.Args) < 2 {
		log.Fatal("Error: No file name provided.\nUsage: roadmap <filename>")
		return
	}

	fileName := os.Args[1]
	if !internal.ValidateFileName(fileName) {
		log.Fatal("Error: File is not a yaml file.\nUsage: roadmap <filename>")
		return
	}

	subject, err := internal.ReadYAMLFile(fileName)
	if err != nil {
		log.Fatal(err)
		return
	}

	subjectDir := internal.CreateDirectory(subject.Name)
	internal.GenerateSubjectMarkdown(subject, subjectDir)

	// internal.ExampleGenerativeModel_GenerateContent_textOnly("Create a summary for the go programming language covering all the major features")
	// internal.Youtube_Search("k means clustering")
}

package main

import (
	"github.com/vivalchemy/roadmap/internal"
	"log"
	"os"
)

func main() {
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
}

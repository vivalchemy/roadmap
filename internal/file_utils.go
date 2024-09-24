package internal

import (
	"fmt"
	"log"
	"os"
	"strings"
)

/*
CreateDirectory creates a directory for the subject based on its name.

Example usage:

subjectDir := CreateDirectory("new-directory")
fmt.Printf("Directory created at: %s\n", subjectDir)
*/
func CreateDirectory(subjectName string) string {
	baseDir := os.Getenv("ROADMAP_HOME")
	if baseDir == "" {
		baseDir = "/home/shadow/Public/roadmap/tmp"
	}
	if baseDir[len(baseDir)-1:] == "/" {
		baseDir = baseDir[:len(baseDir)-1]
	}
	subjectDir := fmt.Sprintf("%s/%s", baseDir, strings.ReplaceAll(subjectName, " ", "-"))
	err := os.MkdirAll(subjectDir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	return subjectDir
}

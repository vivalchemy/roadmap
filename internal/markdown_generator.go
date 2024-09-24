// internal/markdown.go

package internal

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// GenerateSubjectMarkdown generates markdown files for the subject and its topics.
/*
Example usage:

	subject := internal.Subject{
		Name: "Golang Basics",
		Topics: []internal.Topic{
			{
				Name:      "Introduction to Go",
				Subtopics: []string{"Setting Up Go", "First Go Program"},
			},
			{
				Name:      "Go Data Types",
				Subtopics: []string{"Variables", "Constants"},
			},
		},
	}

	// Define the base directory where the markdown files will be created
	baseDir := "./roadmap"
  GenerateSubjectMarkdown(subject, subjectDir)
  fmt.Println("Markdown files generated successfully.")
*/
func GenerateSubjectMarkdown(subject Subject, subjectDir string) {
	// Create a markdown file for the subject
	subjectREADME, err := os.Create(fmt.Sprintf("%s/README.md", subjectDir))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer subjectREADME.Close()

	writeSubjectHeader(subjectREADME, subject.Name)

	for _, topic := range subject.Topics {
		topicDir := createTopicDirectory(subjectDir, topic.Name)
		writeTopicHeader(subjectREADME, topic.Name)

		for _, subtopic := range topic.Subtopics {
			fmt.Fprintf(
				subjectREADME, "- [ ] [[%s/%s/%s | %s]]\n",
				strings.ReplaceAll(subject.Name, " ", "-"),
				strings.ReplaceAll(topic.Name, " ", "-"),
				strings.ReplaceAll(subtopic, " ", "_"),
				subtopic)
			createSubtopicMarkdown(topicDir, subtopic)
		}
		fmt.Fprintln(subjectREADME)
	}
}

func writeSubjectHeader(file *os.File, subjectName string) {
	fmt.Fprintln(file, "# ", subjectName)
	fmt.Fprintln(file)
}

func createTopicDirectory(subjectDir string, topicName string) string {
	topicDir := strings.ReplaceAll(topicName, " ", "-")
	createTopicDir := fmt.Sprintf("%s/%s", subjectDir, topicDir)
	os.MkdirAll(createTopicDir, 0755)
	return createTopicDir
}

func writeTopicHeader(file *os.File, topicName string) {
	fmt.Fprintln(file, "##", topicName, "\n---\n\n")
}

func createSubtopicMarkdown(topicDir string, subtopic string) {
	subtopicFile := strings.ReplaceAll(subtopic, " ", "_")
	createSubTopicFile := fmt.Sprintf("%s/%s.md", topicDir, subtopicFile)
	subtopicREADME, err := os.Create(createSubTopicFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer subtopicREADME.Close()

	fmt.Fprintf(subtopicREADME, "# %s\n---\n\n", subtopic)
	GenerateVideoMarkdown(subtopicREADME, subtopic)
}

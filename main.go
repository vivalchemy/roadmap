package main

import (
	"log"

	"os"

	"github.com/joho/godotenv"
	"github.com/vivalchemy/roadmap/internal"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Get the file name from the command line
	if len(os.Args) < 2 {
		log.Fatal("Error: No file name provided.\nUsage: roadmap <filename.yaml>")
		return
	}

	fileName := os.Args[1]
	if !internal.ValidateFileName(fileName) {
		log.Fatal("Error: File is not a yaml file.\nUsage: roadmap <filename.yaml>")
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
	// var YOUTUBE_API_KEY string
	// var exist bool
	//
	// if YOUTUBE_API_KEY, exist = os.LookupEnv("YOUTUBE_API_KEY"); !exist {
	// 	log.Fatal("YOUTUBE_API_KEY doesn't exists")
	// }
	//
	// yt := internal.NewYouTube(YOUTUBE_API_KEY)
	// const YOUTUBE_VIDEO_URL = "https://www.youtube.com/watch?v=d0PyfYpD4bw"
	//
	// transcript, err := yt.GrabTranscriptForUrl(YOUTUBE_VIDEO_URL)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(transcript)
}

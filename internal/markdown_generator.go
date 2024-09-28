package internal

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

func GenerateVideoMarkdown(ctx context.Context, file *os.File, rawQuery ...string) {
	// File name is generated previously
	searchQuery := strings.Join(rawQuery, " ")

	var YOUTUBE_API_KEY string
	var exist bool

	if YOUTUBE_API_KEY, exist = os.LookupEnv("YOUTUBE_API_KEY"); !exist {
		log.Println("YOUTUBE_API_KEY doesn't exists")
	}

	yt := NewYouTube(YOUTUBE_API_KEY)

	var curTopic *CurrectTopic = ctx.Value(CurTopicKey).(*CurrectTopic)

	searchResults, err := yt.Search(curTopic.Topic, searchQuery)
	if err != nil {
		log.Fatalf("SearchVideos failed: %v", err)
	}

	firstResult := searchResults[0]

	log.Printf("Adding data to markdown file: %s\n", firstResult.VideoTitle)

	if err := addVideoHeader(file, firstResult); err != nil {
		log.Fatalf("Failed to add video header: %v", err)
	}

	if err := addExplanationSection(ctx, file, yt, firstResult.VideoId, firstResult.VideoTitle); err != nil {
		log.Printf("Failed to add explanation section: %v", err)
	}

	if err := addNotesSection(file); err != nil {
		log.Fatalf("Failed to add notes section: %v", err)
	}

	if err := addRelatedVideosSection(file, searchResults[1:]); err != nil {
		log.Fatalf("Failed to add related videos section: %v", err)
	}

	log.Printf("Finished adding content for video: %s\n", firstResult.VideoTitle)
}

// addVideoHeader writes the video title, channel, and thumbnail to the markdown file.
func addVideoHeader(file *os.File, video YoutubeSearchResult) error {
	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.VideoId)
	markdownHeader := fmt.Sprintf("### %s by %s\n![%s by %s](%s)\n\n---\n\n",
		video.VideoTitle, video.ChannelTitle,
		video.VideoTitle, video.ChannelTitle,
		videoURL)

	_, err := file.WriteString(markdownHeader)
	if err != nil {
		return fmt.Errorf("failed to write video header: %w", err)
	}
	return nil
}

// addExplanationSection writes the explanation section with AI-generated insights.
func addExplanationSection(ctx context.Context, file *os.File, yt *YouTube, videoID string, videoTitle string) error {
	transcript, err := yt.GrabTranscript(videoID)
	if err != nil {
		log.Printf("Failed to grab transcript: %v", err)
		return nil // Proceed without explanation
	}

	promptText, err := loadPrompt("/home/shadow/Public/roadmap/patterns/prompt.md")
	if err != nil {
		return fmt.Errorf("failed to load prompt: %w", err)
	}

	aiResponse, err := generateAIInsights(ctx, yt, videoTitle, transcript, promptText)
	if err != nil {
		return fmt.Errorf("failed to generate AI insights: %w", err)
	}

	if _, err := file.WriteString(aiResponse); err != nil {
		return fmt.Errorf("failed to write AI response: %w", err)
	}

	log.Println(file.Stat())
	log.Println(file.Name())
	// log.Printf("AI response added: %s\n", aiResponse)
	return nil
}

// loadPrompt reads the prompt text from the specified file.
func loadPrompt(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to open prompt file: %w", err)
	}
	defer file.Close()

	var promptBuilder strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		promptBuilder.WriteString(scanner.Text())
		promptBuilder.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading prompt file: %w", err)
	}

	return promptBuilder.String(), nil
}

// generateAIInsights interacts with the generative model to create insights.
func generateAIInsights(ctx context.Context, yt *YouTube, videoTitle, transcript, prompt string) (string, error) {
	var curTopic *CurrectTopic = ctx.Value(CurTopicKey).(*CurrectTopic)
	// Assuming ExampleGenerativeModel_GenerateContent_textOnly is a predefined function.
	// You might need to replace this with the actual implementation or dependency injection.
	return ExampleGenerativeModel_GenerateContent_textOnly("Subject: ", curTopic.Subject, "\n\n", "Topic: ", curTopic.Topic, "\n\n", "# Title: ", videoTitle, "\n\n", transcript, prompt), nil
}

// addNotesSection writes the notes section to the markdown file.
func addNotesSection(file *os.File) error {
	if _, err := file.WriteString("## Notes\n---\n\n"); err != nil {
		return fmt.Errorf("failed to write notes section: %w", err)
	}
	return nil
}

// addRelatedVideosSection appends related videos to the markdown file.
func addRelatedVideosSection(file *os.File, related []YoutubeSearchResult) error {
	if len(related) == 0 {
		return nil
	}

	if _, err := file.WriteString("---\n## Related Videos\n"); err != nil {
		return fmt.Errorf("failed to write related videos header: %w", err)
	}

	for _, result := range related {
		if err := writeRelatedVideo(file, result); err != nil {
			return err
		}
	}

	return nil
}

// writeRelatedVideo writes a single related video entry to the markdown file.
func writeRelatedVideo(file *os.File, video YoutubeSearchResult) error {
	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.VideoId)
	relatedVideo := fmt.Sprintf("### %s by %s\n![%s by %s](%s)\n\n",
		video.VideoTitle, video.ChannelTitle,
		video.VideoTitle, video.ChannelTitle,
		videoURL)

	if _, err := file.WriteString(relatedVideo); err != nil {
		return fmt.Errorf("failed to write related video: %w", err)
	}
	return nil
}

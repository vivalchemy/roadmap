package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeSearchResult struct {
	VideoId      string `json:"videoId,omitempty"`
	VideoTitle   string `json:"title,omitempty"`
	ChannelTitle string `json:"channelTitle,omitempty"`
}

func Youtube_Search(rawQuery ...string) []YoutubeSearchResult {
	searchQuery := strings.Join(rawQuery, " ")

	developerKey, doesExist := os.LookupEnv("YOUTUBE_API_KEY")
	if !doesExist {
		log.Fatal("YouTube API key not found")
	}

	service, err := youtube.NewService(context.Background(), option.WithAPIKey(developerKey))
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	// Make the API call to YouTube.
	call := service.Search.List([]string{"snippet"}).
		Q(searchQuery). // Replace with your actual search query
		MaxResults(6)   // Only retrieve 5 results

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}

	searchResult := []YoutubeSearchResult{}
	// Check if the response contains items
	if len(response.Items) == 0 {
		log.Println("No results found")
		return searchResult
	}

	for _, item := range response.Items {
		if item.Id.Kind == "youtube#video" {
			videoID := item.Id.VideoId
			videoTitle := item.Snippet.Title
			channel := item.Snippet.ChannelTitle
			searchResult = append(searchResult, YoutubeSearchResult{videoID, videoTitle, channel})
			// videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
		}
	}
	return searchResult
}

func GenerateVideoMarkdown(file *os.File, rawQuery ...string) {
	// File name is generated previously
	searchQuery := strings.Join(rawQuery, " ")
	searchResult := Youtube_Search(searchQuery)

	firstResult := searchResult[0]
	videoUrl := fmt.Sprintf("https://www.youtube.com/watch?v=%s", firstResult.VideoId)
	fmt.Fprintf(file, "### %s by %s\n![%s by %s](%s)\n\n---\n\n", firstResult.VideoTitle, firstResult.ChannelTitle, firstResult.VideoTitle, firstResult.ChannelTitle, videoUrl)

	fmt.Fprintln(file, "## Explanation\n---\n\n")

	fmt.Fprintln(file, "## Notes\n---\n\n")

	fmt.Fprintln(file, "---\n## Related videos")
	for i := 1; i < len(searchResult); i++ {
		result := searchResult[i]
		videoUrl := fmt.Sprintf("https://www.youtube.com/watch?v=%s", result.VideoId)
		fmt.Fprintf(file, "### %s by %s\n![%s by %s](%s)\n\n", result.VideoTitle, result.ChannelTitle, result.VideoTitle, result.ChannelTitle, videoUrl)
	}

}

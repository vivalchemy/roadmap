package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/anaskhan96/soup"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Setting struct {
	EnvVariable string
	Value       string
	Required    bool
}

type YouTube struct {
	service *youtube.Service
	ApiKey  Setting
}

func NewYouTube(apiKey string) *YouTube {
	return &YouTube{
		ApiKey: Setting{
			EnvVariable: "YOUTUBE_API_KEY",
			Value:       apiKey,
			Required:    true,
		},
	}
}

func (o *YouTube) initService() (err error) {
	if o.service == nil {
		ctx := context.Background()
		o.service, err = youtube.NewService(ctx, option.WithAPIKey(o.ApiKey.Value))
	}
	return
}

func (o *YouTube) GetVideoId(url string) (ret string, err error) {
	if err = o.initService(); err != nil {
		return
	}

	pattern := `(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/(?:[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(url)
	fmt.Println(match)
	if len(match) > 1 {
		ret = match[1]
	} else {
		err = fmt.Errorf("invalid YouTube URL, can't get video ID")
	}
	return
}

func (o *YouTube) GrabTranscriptForUrl(url string) (ret string, err error) {
	var videoId string
	if videoId, err = o.GetVideoId(url); err != nil {
		return
	}
	return o.GrabTranscript(videoId)
}

func (o *YouTube) GrabTranscript(videoId string) (ret string, err error) {
	var transcript string
	if transcript, err = o.GrabTranscriptBase(videoId); err != nil {
		err = fmt.Errorf("transcript not available. (%v)", err)
		return
	}

	// Parse the XML transcript
	doc := soup.HTMLParse(transcript)
	// Extract the text content from the <text> tags
	textTags := doc.FindAll("text")
	var textBuilder strings.Builder
	for _, textTag := range textTags {
		textBuilder.WriteString(textTag.Text())
		textBuilder.WriteString(" ")
		ret = textBuilder.String()
	}
	return
}

func (o *YouTube) GrabTranscriptBase(videoId string) (ret string, err error) {
	if err = o.initService(); err != nil {
		return
	}

	url := "https://www.youtube.com/watch?v=" + videoId
	var resp string
	if resp, err = soup.Get(url); err != nil {
		return
	}

	doc := soup.HTMLParse(resp)
	scriptTags := doc.FindAll("script")
	for _, scriptTag := range scriptTags {
		if strings.Contains(scriptTag.Text(), "captionTracks") {
			regex := regexp.MustCompile(`"captionTracks":(\[.*?\])`)
			match := regex.FindStringSubmatch(scriptTag.Text())
			if len(match) > 1 {
				var captionTracks []struct {
					BaseURL string `json:"baseUrl"`
				}

				if err = json.Unmarshal([]byte(match[1]), &captionTracks); err != nil {
					return
				}

				if len(captionTracks) > 0 {
					transcriptURL := captionTracks[0].BaseURL
					ret, err = soup.Get(transcriptURL)
					return
				}
			}
		}
	}
	err = fmt.Errorf("transcript not found")
	return
}

type YoutubeSearchResult struct {
	VideoId      string `json:"videoId,omitempty"`
	VideoTitle   string `json:"title,omitempty"`
	ChannelTitle string `json:"channelTitle,omitempty"`
}

func (o *YouTube) Search(rawQuery ...string) ([]YoutubeSearchResult, error) {
	searchQuery := strings.Join(rawQuery, " ")

	if err := o.initService(); err != nil {
		return nil, err
	}

	// Make the API call to YouTube.
	call := o.service.Search.List([]string{"snippet"}).
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
		return searchResult, errors.New("No results found in youtube search")
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
	return searchResult, nil
}

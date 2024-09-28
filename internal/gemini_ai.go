package internal

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func ExampleGenerativeModel_GenerateContent_textOnly(rawPrompt ...string) string {
	prompt := strings.Join(rawPrompt, "\n")
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// [START text_gen_text_only_prompt]
	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}

	return printResponse(resp)
	// [END text_gen_text_only_prompt]
}

func printResponse(resp *genai.GenerateContentResponse) (res string) {
	// hijack the stdout
	ogStdOut := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout = w

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Fprintln(w, part)
			}
		}
	}
	w.Close()
	os.Stdout = ogStdOut

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		res += scanner.Text()
	}
	return res
}

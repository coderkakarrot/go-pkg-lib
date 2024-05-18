package main

import (
	"flag"
	"fmt"
	"github.com/eritikass/githubmarkdownconvertergo"
	"log"
)

func main() {
	releaseNotes := flag.String("release-notes", "", "release notes of the github release")
	flag.Parse()
	if releaseNotes == nil {
		log.Fatalf("release notes is required")
	}

	//Check for the length of the releaseNotes (Should be greater than 0)
	if len(*releaseNotes) == 0 {
		log.Fatalf("release notes is required")
	}

	// Convert markdown to Slack format only if there's content
	var slackText string
	if len(*releaseNotes) > 0 {
		slackText = githubmarkdownconvertergo.Slack(*releaseNotes, githubmarkdownconvertergo.SlackConvertOptions{
			Headlines: true,
		})
	} else {
		slackText = "No release notes provided." // or handle empty input differently
	}

	//test := fmt.Sprintf("{\"text\": \"%s\"}", slackText)

	fmt.Println(slackText)

}

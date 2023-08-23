package cmd

import (
	"bytes"
	"strings"

	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"

	"github.com/Delta456/box-cli-maker/v2"
	markdown "github.com/MichaelMure/go-term-markdown"
)

func renderAlertTimelineBox(b box.Box, timeline api.AlertTimeline) {
	if timeline.Message.Value == "" {
		return
	}
	value := []byte(timeline.Message.Value)
	if strings.HasPrefix(timeline.Message.Format, api.AlertCommentFormatMarkdown.String()) {
		// replace CR LF \r\n (windows) with LF \n (unix)
		value = bytes.Replace(value, []byte{13, 10}, []byte{10}, -1)
		// replace CF \r (mac) with LF \n (unix)
		value = bytes.Replace(value, []byte{13}, []byte{10}, -1)
		value = markdown.Render(timeline.Message.Value, 80, 0)
	}
	title := timeline.User.Name
	// we need to avoid a box panic when the comment is shorter than the title; specifically
	// panic: Title must be shorter than the Top & Bottom Bars
	// to do this, we will set horizontal padding to 1 instead of 2...and
	// add our own space-based padding as a prefix
	customPadding := []byte(strings.Repeat(" ", len(title)) + "\n")
	value = append(customPadding, value...)
	// we will also make sure that the value trails with a newline...
	// such that we have consistent horizontal bottom padding
	if !strings.HasSuffix(string(value), "\n") {
		value = append(value, []byte("\n")...)
	}
	b.Println(title, string(value))
}

func renderAlertTimelineBoxes(timelines []api.AlertTimeline) {
	timelineBox := box.New(box.Config{Px: 2, Py: 2, Type: "Single", Color: "Cyan", TitlePos: "Top"})

	for _, t := range timelines {
		renderAlertTimelineBox(timelineBox, t)
	}
}

func showAlertTimeline(id int) error {
	cli.StartProgress(" Fetching alert timeline...")
	timelineResponse, err := cli.LwApi.V2.Alerts.GetTimeline(id)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to show alert")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(timelineResponse.Data)
	}

	if len(timelineResponse.Data) == 0 {
		cli.OutputHuman("There are no timeline entries for the specified alert.\n")
		return nil
	}

	renderAlertTimelineBoxes(timelineResponse.Data)
	return nil
}

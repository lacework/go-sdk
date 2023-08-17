package cmd

import (
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
		value = markdown.Render(timeline.Message.Value, 80, 0)
	}
	b.Println(timeline.User.Name, string(value))
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

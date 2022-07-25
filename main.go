package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
)

type Status struct {
	Pass  int `json:"pass,omitempty"`
	Fail  int `json:"fail,omitempty"`
	Total int `json:"total,omitempty"`
	Skip  int `json:"skip,omitempty"`
}
type cmdFlags struct {
	inputFile string
	resultURL string
}

func main() {
	rootCmd := initRootCommand()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

}

func initRootCommand() *cobra.Command {
	flags := &cmdFlags{}
	rootCmd := &cobra.Command{
		Use:  "send-slack",
		Long: "send slack message with data retrived from status file",
		Run: func(cmd *cobra.Command, args []string) {
			sendSlack(flags)
		},
	}
	rootCmd.PersistentFlags().StringVarP(&flags.inputFile,
		"input",
		"i",
		"",
		"the json input file")
	rootCmd.MarkPersistentFlagRequired("input")
	rootCmd.PersistentFlags().StringVarP(&flags.resultURL,
		"url",
		"u",
		"",
		"result url")
	return rootCmd
}

func sendSlack(flags *cmdFlags) {
	status := readStatus(flags.inputFile)
	token := os.Getenv("SLACK_AUTH_TOKEN")
	channelID := os.Getenv("SLACK_CHANNEL_ID")

	client := slack.New(token, slack.OptionDebug(true))

	header := slack.Attachment{
		Pretext: fmt.Sprintf(":fire:  Agni E2E Report : <%s> ", flags.resultURL),
		Text:    time.Now().Format(time.ANSIC),
		Color:   "#b380ff",
	}

	total := slack.Attachment{
		Color: "#b380ff",
		Text:  fmt.Sprintf("TOTAL : %d", status.Total),
	}
	pass := slack.Attachment{
		Color: "#33cc33",
		Text:  fmt.Sprintf("PASS : %d", status.Pass),
		// Pretext: ":large_green_square: 10",
	}
	fail := slack.Attachment{
		Color: "#cc0000",
		Text:  fmt.Sprintf("FAIL : %d", status.Fail),
		// Pretext: ":large_red_square: 10",
	}
	skip := slack.Attachment{
		Color: "#cccccc",
		Text:  fmt.Sprintf("SKIP : %d", status.Skip),
		// Pretext: ":large_yellow_square: 10",
	}

	_, _, err := client.PostMessage(
		channelID,
		slack.MsgOptionAttachments(header, total, pass, skip, fail),
	)

	if err != nil {
		panic(err)
	}
}

func readStatus(statusFile string) *Status {

	status := Status{}
	jsonFile, err := os.Open(statusFile)
	if err != nil {
		return &status
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &status)
	return &status
}

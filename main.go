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

	token := os.Getenv("SLACK_AUTH_TOKEN")
	channelID := os.Getenv("SLACK_CHANNEL_ID")

	client := slack.New(token, slack.OptionDebug(true))

	attachment, status := buildSlackMsg(flags)
	_, _, _, err := client.SendMessage(
		channelID,
		slack.MsgOptionAttachments(*attachment),
	)

	if err != nil {
		panic(err)
	}
	os.Exit(status)
}

func buildSlackMsg(flags *cmdFlags) (*slack.Attachment, int) {
	status := readStatus(flags.inputFile)

	line1 := fmt.Sprintf("")
	line2 := fmt.Sprintf("\nTOTAL: %v", status.Total)
	line3 := fmt.Sprintf("\n:white_check_mark:PASS: %d		:x:FAIL: %d		:leftwards_arrow_with_hook:SKIP: %d", status.Pass, status.Fail, status.Skip)
	message := line1 + line2 + line3
	attachment := slack.Attachment{
		Title:     fmt.Sprintf("Agni E2E Report"),
		TitleLink: flags.resultURL,
		Text:      message,
		Color:     "#b380ff",
		Fields: []slack.AttachmentField{
			{}},
		Footer: fmt.Sprintf("%v", time.Now().Format("02 Jan 06 15:04 MST")),
	}
	return &attachment, status.Fail
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

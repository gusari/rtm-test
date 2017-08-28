package main

import (
	"log"
	"os"

	"fmt"
	"github.com/nlopes/slack"
	"strings"

	"net/http"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func checkSession(r *http.Request) {
    // Get a session. Get() always returns a session, even if empty.
    session, err := store.Get(r, "session-name")
    if err != nil {
        log.Print(err.Error())
        return
    }
}

func updateSession(r *http.Request) {
    // Get a session. Get() always returns a session, even if empty.
    session, err := store.Get(r, "session-name")
    if err != nil {
        log.Print(err.Error())
        return
    }

    // Set some session values.
    session.Values["foo"] = "bar"
    session.Values[42] = 43
    // Save it before we write to the response/return from the handler.
    session.Save(r, w)
}

func run(api *slack.Client) int {
	log.Print("go run start!")
	rtm := api.NewRTM()
	go rtm.ManageConnection()
	_, conn, _ := rtm.connect(0, rtm.useRTMStart)
	r := conn.Request()
	updateSession(r)
	checkSession(r)

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				log.Print("Hello Event")

			case *slack.MessageEvent:
				log.Printf("Message: %v\n", ev.Text)

			case *slack.InvalidAuthEvent:
				log.Print("Invalid credentials")
				return 1

			}
		}
	}
}

func main() {
	api := slack.New("xoxb-232574333637-8CPgYpujD530qwXl5IoYcasJ")

	os.Exit(run(api))
}

func makeThreadWeb(api *slack.Client, threadTimestamp, channel string) {
	attachment := slack.Attachment{
		Text:       "Which group do you want to join? :smile:",
		Color:      "#f9a41b",
		CallbackID: "hoge",
		Actions: []slack.AttachmentAction{
			{
				Name: "actionSelect",
				Type: "select",
				Options: []slack.AttachmentActionOption{
					{
						Text:  "GitHub",
						Value: "github",
					},
					{
						Text:  "まかれる",
						Value: "m",
					},
					{
						Text:  "Slack",
						Value: "slack",
					},
				},
			},
			{
				Name:  "actionCancel",
				Text:  "Cancel",
				Type:  "button",
				Style: "danger",
			},
		},
	}

	params := slack.PostMessageParameters{
		ThreadTimestamp: threadTimestamp,
		Attachments: []slack.Attachment{
			attachment,
		},
	}
	if _, _, err := api.PostMessage(channel, "てきすと", params); err != nil {
		fmt.Errorf("failed to post message: %s", err)
		return
	}
	return
}

/*
func makeThreadRtm (rtm *RTM,timestamp, channel string)  {
	//rtm:=api.NewRTM()
	x:=rtm.NewOutgoingMessage("どこに招待します :simple_smile: ", channel)
	x.ThreadTimestamp = timestamp
	rtm.SendMessage(x)
	return
}*/

func detectTread(rtm *slack.RTM) {

}

func filterMessage(message string) bool {
	callBotText := "招待"
	flag := strings.Contains(message, callBotText)
	if flag {
		return true
	} else {
		return false
	}
}

//親のtsがthread_tsになる
func isThreadExist(thread_ts, ts string) string {
	if thread_ts == "" { //親メッセージ
		return ts
	} else {
		return thread_ts
	}
}

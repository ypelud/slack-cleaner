package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"flag"

	"github.com/nlopes/slack"
)

func ctv(ts string) time.Time {
	i, err := strconv.ParseFloat(ts, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(int64(i), 0)
}

func usage() {
	fmt.Println("slack-cleaner is a tool for removing messages into a slack channel")
	fmt.Println("Usage:")
	fmt.Println("\tslack-cleaner <options> channel")
	fmt.Println("\tchannel\t:The name of the channel")
	fmt.Println("options:")

	fmt.Printf("\t-date\t:Remove every message before that date 'YYYYMMDD'\n")
	fmt.Printf("\t-user\t:Message user\n\n")
	fmt.Printf("\t-bot\t:bot id\n\n")
}

func main() {
	//convert into options
	help := flag.Bool("h", false, "help") 
	selectedDate := flag.String("date", "", "delete messages until this date") 
	selectedUser := flag.String("user", "", "only for this user") 
	selectedBotID := flag.String("bot", "", "only for this bot id") 

	flag.Parse()
	myArgs :=  flag.Args()

	fmt.Println(myArgs)

	if len(myArgs) != 1 || *help {
		usage()
		os.Exit(1)
	}

	selectedChannel := myArgs[0]

	fmt.Println("selectedChannel:", selectedChannel)
	fmt.Println("date:", *selectedDate)
	fmt.Println("user:", *selectedUser)
	
	dateSelectedString := ""
	if *selectedDate != "" {
		const shortForm = "20060102"
		t, _ := time.Parse(shortForm, *selectedDate)
		dateSelectedString = strconv.FormatFloat(float64(t.Unix()), 'f', 6, 64)
		fmt.Printf("Remove every message on <%s> before <%s>\n", selectedChannel, t.Format(shortForm))	
	}

	api := slack.New("xoxp-11888837425-11886289616-13691440419-33c723bbbb")

	params := &slack.HistoryParameters{
		Latest:    dateSelectedString,
		Oldest:    "0",
		Count:     4000,
		Inclusive: false,
		Unreads:   false,
	}

	channels, err := api.GetChannels(false)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	selectedChannelID := ""		

	for _, channel := range channels {
		if selectedChannel == channel.Name {
			selectedChannelID = channel.ID
		}
	}

	if selectedChannelID == "" {
		fmt.Printf("Channel not found <%s>\n", selectedChannel)
		os.Exit(1)
	}

	selectedUserID := ""

	if *selectedUser != "" {
		users, err := api.GetUsers()
		if err != nil {
			fmt.Printf("%s\n", err)
			return		
		}

		for _, user := range users {
			fmt.Printf("user: %s %s\n", user.ID, user.Profile.RealName)
			if *selectedUser == user.Profile.RealName {
				selectedUserID = user.ID
			}
	 
		}
	
	}

	history, err := api.GetChannelHistory(selectedChannelID, *params)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}



	for _, message := range history.Messages {
		var skip = false

		if selectedUserID != "" && selectedUserID != message.User {
			skip = true
		}		

		if *selectedBotID != "" && *selectedBotID != message.BotID {
			skip = true
		}		

		if !skip {
			respChannel, respTimestamp, err := api.DeleteMessage(selectedChannelID, message.Timestamp)
			fmt.Printf("res: %s %s %s \n", respChannel, ctv(respTimestamp), err)
		}
		
		time.Sleep(2 * time.Second)
	}

}

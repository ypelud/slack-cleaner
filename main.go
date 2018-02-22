package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"flag"
	"log"
	"github.com/BurntSushi/toml"
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
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("\tslack-cleaner <options> channel")
	fmt.Println("\tchannel\t:The name of the channel")
	fmt.Println("")
	fmt.Println("options:")

	fmt.Println("\t-date\t:Remove every message before that date 'YYYYMMDD'")
	fmt.Println("\t-user\t:Message user")
	fmt.Println("\t-bot\t:bot id")
	fmt.Println("\t-n\t:dry run")
	fmt.Println("")
}

// Config ...
type Config struct {
	SlackID   string
}

// ReadConfig ...
func ReadConfig() Config {
	var configfile = "config.toml"
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}

	return config
}

func main() {
	//convert into options
	help := flag.Bool("h", false, "help") 
	selectedDate := flag.String("date", "", "delete messages until this date") 
	selectedUser := flag.String("user", "", "only for this user") 
	selectedBotID := flag.String("bot", "", "only for this bot id") 
	dryRun := flag.Bool("n", false, "dry run") 

	flag.Parse()
	myArgs :=  flag.Args()

	if len(myArgs) != 1 || *help {
		usage()
		os.Exit(1)
	}

	selectedChannel := myArgs[0]

	log.Println("selectedChannel:", selectedChannel)
	log.Println("date:", *selectedDate)
	log.Println("user:", *selectedUser)
	log.Println("dry:", *dryRun)
	
	dateSelectedString := ""
	if *selectedDate != "" {
		const shortForm = "20060102"
		t, _ := time.Parse(shortForm, *selectedDate)
		dateSelectedString = strconv.FormatFloat(float64(t.Unix()), 'f', 6, 64)
		log.Printf("Remove every message on <%s> before <%s>\n", selectedChannel, t.Format(shortForm))	
	}

	var config = ReadConfig()
	api := slack.New(config.SlackID)

	params := &slack.HistoryParameters{
		Latest:    dateSelectedString,
		Oldest:    "0",
		Count:     4000,
		Inclusive: false,
		Unreads:   false,
	}

	channels, err := api.GetChannels(false)
	if err != nil {
		log.Fatal("Error :", err)
		return
	}

	selectedChannelID := ""		

	for _, channel := range channels {
		if selectedChannel == channel.Name {
			selectedChannelID = channel.ID
		}
	}

	if selectedChannelID == "" {
		log.Fatalf("Channel not found <%s>", selectedChannel)
		os.Exit(1)
	}

	selectedUserID := ""

	if *selectedUser != "" {
		users, err := api.GetUsers()
		if err != nil {
			log.Fatal("Error :", err)
			return		
		}

		for _, user := range users {
			log.Printf("user: %s %s\n", user.ID, user.Profile.RealName)
			if *selectedUser == user.Profile.RealName {
				selectedUserID = user.ID
			}
	 
		}
	
	}

	history, err := api.GetChannelHistory(selectedChannelID, *params)
	if err != nil {
		log.Fatal("Error :", err)
		return
	}



	for _, message := range history.Messages {
		var skip = false;

		if selectedUserID != "" && selectedUserID != message.User {
			skip = true
		}		

		if *selectedBotID != "" && *selectedBotID != message.BotID {
			skip = true
		}		

		if !skip {

			if (*dryRun) {
				log.Printf("(dry) Delete : %s \n", ctv(message.Timestamp))				
			} else {
				log.Printf("Delete : %s %s", ctv(message.Timestamp), message.Text)
				respChannel, respTimestamp, err := api.DeleteMessage(selectedChannelID, message.Timestamp)
				if err != nil {
					log.Fatal("Error :", err)
				} else {
					log.Printf("res: %s %s", respChannel, ctv(respTimestamp))
				}	
				// slack rates...
				time.Sleep(2 * time.Second)
			}
		
		}
		
	}

}

package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"ovh/commands"
)

func main() {
	discord, err := discordgo.New("Bot " + "")
	if err != nil {
		log.Fatalln(err)
	}
	discord.Identify.Presence.Game.Name = "test"
	err = discord.Open()
	if err != nil {
		log.Fatalln(err)
	}
	commands.RegisterCommands(discord)
	go commands.StartCronScheduler(discord)
	defer func() {
		err := discord.Close()
		if err != nil {
			log.Fatal("Error while closing the session ", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Connected as ", discord.State.User.Username)
	log.Println("----------------------")
	log.Println("Starting logs")
	log.Println("Press CTRL+C to exit")
	<-stop
}

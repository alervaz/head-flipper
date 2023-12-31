package main

import (
	"log"
	"os"

	"github.com/alervaz/head-flipper/model"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	if err := model.PingDB(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer model.DB.Close()

	dg, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	dg.AddHandler(Flip)

	dg.Identify.Intents = discordgo.IntentsAll

	if err := dg.Open(); err != nil {
		log.Fatal(err)
	}

	log.Println("Bot is running")
	ch := make(chan struct{})
	<-ch

	dg.Close()
}

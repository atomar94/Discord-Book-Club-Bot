package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"bookclubbot.com/main/views"
)

func main() {
	if os.Getenv("DISCORD_BOT_TOKEN") == "" {
		log.Fatalf("DISCORD_BOT_TOKEN environment variable not set")
	}
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Critical error loading data: %v", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMessageReactions

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	dg.AddHandler(views.HandleVotingReactions)

	// Open the connection
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	defer dg.Close()

	views.RegisterInteractionCreateHandler(dg)

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	// 4. Graceful Shutdown (to optionally delete commands on exit)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Optional: Delete command on exit to keep test environment clean
	// dg.ApplicationCommandDelete(dg.State.User.ID, "", registeredCmd.ID)

	fmt.Println("Gracefully shutting down...")
}

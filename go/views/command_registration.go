package views

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type SlashCommand struct {
	discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate) error
}

type ModalHandler struct {
	CustomIdPrefix string
	Handler        func(s *discordgo.Session, i *discordgo.InteractionCreate) error
}

func getSlashCommands() []SlashCommand {

	commands := []SlashCommand{
		{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "recommend-a-book",
				Description: "Submit a book recommendation",
			},
			Handler: HandleRecommendABook,
		},
		{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "view-schedule",
				Description: "View the book club schedule",
			},
			Handler: DiscordResponseWrapper(HandleSchedule),
		},
		{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "recommend-a-cafe",
				Description: "Submit a cafe recommendation",
			},
			Handler: HandleRecommendACafe,
		},
		{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "make-announcement",
				Description: "View the list of recommended cafes",
			},
			Handler: DiscordResponseWrapper(HandleSchedule),
		},
	}

	return commands
}

func getModalHandlers() []ModalHandler {
	handlers := []ModalHandler{
		{
			CustomIdPrefix: "book_recommendation",
			Handler:        HandleRecommendABookModalResponse,
		},
		{
			CustomIdPrefix: "cafe_recommendation",
			Handler:        HandleRecommendACafeModalResponse,
		},
	}
	return handlers
}

func makeSlashCommandHandler(cmds []SlashCommand) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmdMap := make(map[string]SlashCommand, len(cmds))
	for _, cmd := range cmds {
		cmdMap[cmd.Name] = cmd
	}

	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		data := i.ApplicationCommandData()
		if handlerCmd, ok := cmdMap[data.Name]; ok {
			fmt.Println("Handling slash command:", data.Name)
			err := handlerCmd.Handler(s, i)
			if err != nil {
				// Log the error
				// In a real application, you might want to send an error response to the user
				// Here we just log it for simplicity
				fmt.Println("Error handling command ", data.Name, ":", err)
			}
		} else {
			log.Printf("No handler found for command: %s", data.Name)
		}
	}
}

func makeModalHandler(handlers []ModalHandler) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	handlerMap := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) error)
	for _, handler := range handlers {
		handlerMap[handler.CustomIdPrefix] = handler.Handler
	}

	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		d := i.ModalSubmitData()
		for CustomIdPrefix, handler := range handlerMap {
			if strings.HasPrefix(d.CustomID, CustomIdPrefix) {
				fmt.Println("Handling modal:", CustomIdPrefix)
				err := handler(s, i)
				if err != nil {
					fmt.Println("Error handling modal ", CustomIdPrefix, ":", err)
				}
				return
			}
		}
	}
}

func makeInteractionCreateHandler(commands []SlashCommand, modalHandlers []ModalHandler) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	slashCommandHandler := makeSlashCommandHandler(commands)
	modalHandler := makeModalHandler(modalHandlers)
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			slashCommandHandler(s, i)
		case discordgo.InteractionModalSubmit:
			modalHandler(s, i)
		}
	}
}
func RegisterInteractionCreateHandler(s *discordgo.Session) {
	commands := getSlashCommands()
	modalHandlers := getModalHandlers()
	s.AddHandler(makeInteractionCreateHandler(commands, modalHandlers))
	// Register the command with Discord
	// NOTE: Passing "" as the GuildID makes it a Global Command (can take 1 hour to appear).
	// For testing, replace "" with your specific Guild ID (Server ID) for instant updates.
	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", &cmd.ApplicationCommand)
		if err != nil {
			log.Panicf("Cannot create slash command %s: %v", cmd.Name, err)
		}
	}
}

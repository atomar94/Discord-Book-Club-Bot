package views

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"

	"bookclubbot.com/main/controllers"
	"bookclubbot.com/main/models"
)

func HandleSchedule(t models.ClubTable) (string, error) {
	controllers.SanitizeClubTable(t)

	err := controllers.AssignDatesToSchedule(t.Schedule)
	if err != nil {
		return "", fmt.Errorf("Unable to assign dates: %v", err)
	}

	err = controllers.AssignBooksToSchedule(t.BookPool, t.Schedule)
	if err != nil {
		return "", fmt.Errorf("Unable to assign books: %v", err)
	}

	err = controllers.AssignCafesToSchedule(t.CafePool, t.Schedule)
	if err != nil {
		return "", fmt.Errorf("Unable to assign cafes: %v", err)
	}

	response, err := t.RenderSchedule()
	if err != nil {
		return "", fmt.Errorf("Unable to render schedule: %v", err)
	}

	return response, nil
}

func DiscordResponseWrapper(handler func(models.ClubTable) (string, error)) func(*discordgo.Session, *discordgo.InteractionCreate) error {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		t := loadClubTable()
		response, err := handler(t)
		if err != nil {
			return err
		}

		_, err = s.ChannelMessageSend(i.ChannelID, response)
		if err != nil {
			return fmt.Errorf("Unable to send schedule message: %v", err)
		}

		saveClubTable(t)
		return nil
	}
}

func HandleRecommendABook(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "book_recommendation_" + i.Interaction.Member.User.ID,
			Title:    "Book Recommendation",
			Flags:    discordgo.MessageFlagsIsComponentsV2,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "book_title",
							Label:     "Book Title",
							Style:     discordgo.TextInputShort,
							Required:  true,
							MaxLength: 200,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "book_author",
							Label:     "Book Author",
							Style:     discordgo.TextInputShort,
							Required:  true,
							MaxLength: 100,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "book_goodreads_link",
							Label:     "Goodreads Link (optional)",
							Style:     discordgo.TextInputShort,
							Required:  false,
							MaxLength: 200,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "book_description",
							Label:     "Description (optional)",
							Style:     discordgo.TextInputParagraph,
							Required:  false,
							MaxLength: 500,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("Unable to send book recommendation modal: %v", err)
	}
	return nil
}

func HandleRecommendABookModalResponse(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	t := loadClubTable()

	d := i.ModalSubmitData()
	title := d.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value         // book title
	author := d.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value        // book author
	goodreadsLink := d.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value // goodreads link (optional)
	description := d.Components[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value   // book description (optional)

	fmt.Println("Received book recommendation:", title, author, goodreadsLink, description)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Thank you for your recommendation! 📚",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return fmt.Errorf("Unable to send book recommendation confirmation: %v", err)
	}

	err = controllers.AddBook(t.BookPool, title, author, goodreadsLink, description)
	if err != nil {
		return fmt.Errorf("Unable to add book to the pool: %v", err)
	}
	saveClubTable(t)

	embed := discordgo.MessageEmbed{
		Title: "New Book Recommendation Received! 📚",
		Description: fmt.Sprintf("%s recommended a new book! ", i.Interaction.Member.User.DisplayName()) +
			"If you want to read this book for book club please leave a ❤️ reaction below!",
		Color: 0x00ff00, // Green color
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Title",
				Value:  title,
				Inline: true,
			},
			{
				Name:   "Author",
				Value:  author,
				Inline: true,
			},
			{
				Name:  "Goodreads Link",
				Value: goodreadsLink,
			},
			{
				Name:  "Description",
				Value: description,
			},
		},
	}

	message, err := s.ChannelMessageSendEmbed(i.ChannelID, &embed)
	if err != nil {
		return fmt.Errorf("Unable to send book recommendation embed: %v", err)
	}
	err = s.MessageReactionAdd(i.ChannelID, message.ID, "❤️")
	if err != nil {
		return fmt.Errorf("Unable to add reaction to book recommendation embed: %v", err)
	}

	return nil
}

func HandleRecommendACafe(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "cafe_recommendation_" + i.Interaction.Member.User.ID,
			Title:    "Cafe Recommendation",
			Flags:    discordgo.MessageFlagsIsComponentsV2,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "cafe_name",
							Label:     "Cafe Name",
							Style:     discordgo.TextInputShort,
							Required:  true,
							MaxLength: 200,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "google_maps_link",
							Label:     "Google Maps Link",
							Style:     discordgo.TextInputShort,
							Required:  true,
							MaxLength: 200,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("Unable to send cafe recommendation modal: %v", err)
	}
	return nil
}

func HandleRecommendACafeModalResponse(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	t := loadClubTable()

	d := i.ModalSubmitData()
	cafeName := d.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value       // cafe name
	googleMapsLink := d.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value // google maps link

	fmt.Println("Received cafe recommendation:", cafeName, googleMapsLink)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Thank you for your recommendation! 📚",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return fmt.Errorf("Unable to send cafe recommendation confirmation: %v", err)
	}

	err = controllers.AddCafe(t.CafePool, cafeName, googleMapsLink)
	if err != nil {
		return fmt.Errorf("Unable to add cafe to the pool: %v", err)
	}
	saveClubTable(t)
	return nil
}

func HandleVotingReactions(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	log.Println("HandleVotingReactions")
	// A. Ignore the bot's own reactions
	if r.UserID == s.State.User.ID {
		return
	}
	// Skip non-heart reactions
	if r.Emoji.Name != "❤️" {
		return
	}

	// B. Check if the message belongs to the Bot ("My Post")
	// The event (r) only gives us the MessageID, not the Message author.
	// We must fetch the message to see who wrote it.
	// Try the local cache first (fastest)
	msg, err := s.State.Message(r.ChannelID, r.MessageID)
	if err != nil {
		// If not in cache (common for older messages), fetch from API
		msg, err = s.ChannelMessage(r.ChannelID, r.MessageID)
		if err != nil {
			log.Println("Could not fetch message:", err)
			return
		}
		// Add to cache for next time
		s.State.MessageAdd(msg)
	}

	// Skip reaction on messages that aren't from the bot
	if msg.Author.ID != s.State.User.ID {
		return
	}

	t := loadClubTable()

	vote_count := 0
	for _, react := range msg.Reactions {
		if react.Emoji.Name == "❤️" {
			vote_count += react.Count
		}
	}

	// Get the book title from the embed
	if len(msg.Embeds) == 0 {
		log.Println("No embeds found in the message")
		return
	}
	for _, embed := range msg.Embeds {
		for _, field := range embed.Fields {
			if field.Name == "Title" {
				bookName := field.Value
				log.Println("Updating votes for book:", bookName, "to", vote_count)
				err := controllers.UpdateVotes(t.BookPool, bookName, vote_count)
				if err != nil {
					log.Println("Error updating book vote count:", err)
				}
				saveClubTable(t)
				return
			}
		}
	}
}

func loadClubTable() models.ClubTable {
	file, _ := os.ReadFile("club_table.json")
	var data models.ClubTable
	json.Unmarshal(file, &data)
	return data
}

func saveClubTable(t models.ClubTable) {
	data, _ := json.MarshalIndent(t, "", "  ")
	os.WriteFile("club_table.json", data, 0644)
}

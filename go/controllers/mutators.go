package controllers

import (
	"fmt"
	"math/rand"
	"time"

	"bookclubbot.com/main/models"
)

const TIME_FORMAT string = "January 2, 2006"

func SanitizeClubTable(t models.ClubTable) {
	for i := range t.Schedule {
		if t.Schedule[i].Id == "" {
			t.Schedule[i].Id = models.GenerateId()
		}
	}
	for i := range t.BookPool {
		if t.BookPool[i].Id == "" {
			t.BookPool[i].Id = models.GenerateId()
		}
	}
	for i := range t.CafePool {
		if t.CafePool[i].Id == "" {
			t.CafePool[i].Id = models.GenerateId()
		}
	}
}

func AssignDatesToSchedule(schedules []models.ScheduleEntry) error {
	if len(schedules) < 1 {
		return fmt.Errorf("Received schedules list was empty.")
	}
	var initialDate time.Time = time.Now()

	if schedules[0].Date == "" {
		now := time.Now()
		// Find days until next Saturday (Saturday is 6 in Go)
		daysUntilSat := (int(time.Saturday) - int(now.Weekday()) + 7) % 7
		if daysUntilSat == 0 && now.Hour() >= 14 { // Skip to next week if past 2PM Sat
			daysUntilSat = 7
		}

		initialDate = now.AddDate(0, 0, daysUntilSat)
	} else {
		parsedDate, err := time.Parse(TIME_FORMAT, schedules[0].Date)
		if err != nil {
			return fmt.Errorf("Unable to parse date when assigning schedules. Found '%s'. %w", schedules[0].Date, err)
		}
		initialDate = parsedDate
	}

	for i := range schedules {
		schedules[i].Date = initialDate.Format(TIME_FORMAT)
		initialDate = initialDate.AddDate(0, 0, 7)
	}

	return nil
}

func AssignBooksToSchedule(books []models.BookEntry, schedules []models.ScheduleEntry) error {
	if len(schedules) < 1 {
		return nil
	}

	for i, s := range schedules {
		// don't assign a book if we have one for at least the next 2 weeks
		if i > 2 {
			break
		}
		if s.BookId != "" {
			continue
		}

		next_book, err := selectNextBook(books)
		if err != nil {
			return fmt.Errorf("Unable to select next book for our schedule: %w", err)
		}
		found_book := false
		for i, book := range books {
			if book.Id == next_book.Id {
				books[i].Read = true
				found_book = true
				break
			}
		}
		if !found_book {
			return fmt.Errorf("A book was selected but it isn't in our database. BookId=%s not found.", next_book.Id)
		}

		book_week_duration := 4
		if len(schedules) < i+book_week_duration {
			return fmt.Errorf("We don't have enough weeks in the schedule to schedule this book.")
		}
		for nth_week_of_book := range schedules[i : i+book_week_duration] {
			schedules[i+nth_week_of_book].BookId = next_book.Id
		}
		return nil
	}

	return nil
}

func AssignCafesToSchedule(cafes []models.CafeEntry, schedules []models.ScheduleEntry) error {
	if len(schedules) < 1 {
		return nil
	}
	if len(cafes) < 1 {
		return fmt.Errorf("No cafes available to assign to schedule.")
	}

	for i := range schedules {
		random_index := rand.Intn(len(cafes))
		schedules[i].CafeId = cafes[random_index].Id
	}
	return nil
}

func AddBook(books []models.BookEntry, title string, author string, goodreadsLink string, description string) error {
	new_book := models.BookEntry{
		Id:          models.GenerateId(),
		Name:        title,
		Author:      author,
		Link:        goodreadsLink,
		Description: description,
		Votes:       0,
		Read:        false,
	}
	books = append(books, new_book)
	return nil
}

func AddCafe(cafes []models.CafeEntry, name string, googleMapsLink string) error {
	new_cafe := models.CafeEntry{
		Id:   models.GenerateId(),
		Name: name,
		Link: googleMapsLink,
	}
	cafes = append(cafes, new_cafe)
	return nil
}

func UpdateVotes(books []models.BookEntry, bookName string, vote_count int) error {
	for i, book := range books {
		if book.Name == bookName {
			books[i].Votes = vote_count
			return nil
		}
	}
	return fmt.Errorf("Book with name '%s' not found", bookName)
}

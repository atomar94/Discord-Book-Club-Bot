package models

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func (t *ClubTable) RenderSchedule() (string, error) {

	readFileContents, err := os.ReadFile("models/templates/schedule.template.md")
	if err != nil {
		return "", fmt.Errorf("Error opening schedule.template.md: %v", err)
	}

	tmpl, err := template.New("schedule_template").Parse(string(readFileContents))

	var buf strings.Builder

	var rendered_schedule_data = struct {
		CurrentBook       string
		CurrentAuthor     string
		NextBook          string
		NextAuthor        string
		NextBookStartDate string
		Schedule          []struct {
			Date     string
			Link     string
			CafeName string
			BookName string
		}
	}{}

	current_book, err := t.GetBookById(t.Schedule[0].BookId)
	if err != nil {
		return "", fmt.Errorf("Error getting current book: %v", err)
	}
	rendered_schedule_data.CurrentBook = current_book.Name
	rendered_schedule_data.CurrentAuthor = current_book.Author

	for _, schedule_entry := range t.Schedule {
		if schedule_entry.BookId != t.Schedule[0].BookId && schedule_entry.BookId != "" {
			next_book, err := t.GetBookById(schedule_entry.BookId)
			if err != nil {
				return "", fmt.Errorf("Error getting next book: %v", err)
			}
			rendered_schedule_data.NextBook = next_book.Name
			rendered_schedule_data.NextAuthor = next_book.Author
			rendered_schedule_data.NextBookStartDate = schedule_entry.Date
			break
		}
	}

	max_schedule_entries := 4
	for i, schedule_entry := range t.Schedule {
		if i >= max_schedule_entries {
			break
		}
		cafe, err := t.GetCafeById(schedule_entry.CafeId)
		if err != nil {
			return "", fmt.Errorf("Error getting cafe: %v", err)
		}
		var book_name string
		if schedule_entry.BookId == "" {
			book_name = "TBD"
		} else {
			book, err := t.GetBookById(schedule_entry.BookId)
			if err != nil {
				return "", fmt.Errorf("Error getting book: %v", err)
			}
			book_name = book.Name
		}
		rendered_schedule_data.Schedule = append(rendered_schedule_data.Schedule, struct {
			Date     string
			Link     string
			CafeName string
			BookName string
		}{
			Date:     schedule_entry.Date,
			Link:     cafe.Link,
			CafeName: cafe.Name,
			BookName: book_name,
		})
	}

	err = tmpl.Execute(&buf, rendered_schedule_data)
	if err != nil {
		return "", fmt.Errorf("Error executing template: %v", err)
	}

	return buf.String(), nil
}

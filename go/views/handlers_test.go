package views

import (
	"fmt"
	"testing"

	"bookclubbot.com/main/models"
)

func TestHandleSchedule(t *testing.T) {
	dated_schedule := models.ScheduleEntry{
		Id:     "123",
		Date:   "December 20, 2025",
		BookId: "book-99",
		CafeId: "cafe-42",
	}
	undated_schedule := models.ScheduleEntry{
		Id:     "124",
		Date:   "",
		BookId: "book-99",
		CafeId: "cafe-42",
	}

	table := models.ClubTable{}
	table.Schedule = append(table.Schedule, dated_schedule)
	table.Schedule = append(table.Schedule, undated_schedule)

	response, err := HandleSchedule(table)
	if err != nil {
		t.Error("Internal error %w", err)
	}
	if response == "" {
		t.Error("Empty response returned.")
	}
	// run with -v to see
	fmt.Print(response)
}

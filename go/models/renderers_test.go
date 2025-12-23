package models

import (
	"fmt"
	"testing"
)

func TestRenderSchedule(t *testing.T) {
	dated_schedule := ScheduleEntry{
		Id:     "123",
		Date:   "December 20, 2025",
		BookId: "book-1",
		CafeId: "cafe-1",
	}
	undated_schedule := ScheduleEntry{
		Id:     "124",
		Date:   "December 27, 2025",
		BookId: "book-2",
		CafeId: "cafe-2",
	}

	table := ClubTable{}
	table.Schedule = append(table.Schedule, dated_schedule)
	table.Schedule = append(table.Schedule, undated_schedule)

	table.BookPool = append(table.BookPool, BookEntry{
		Id:   "book-1",
		Name: "Example Book",
		Link: "https://booklink.com",
	})
	table.BookPool = append(table.BookPool, BookEntry{
		Id:   "book-2",
		Name: "Example Book 2",
		Link: "https://booklink.com",
	})

	table.CafePool = append(table.CafePool, CafeEntry{
		Id:   "cafe-1",
		Name: "Example Cafe",
		Link: "https://cafelink.com",
	})
	table.CafePool = append(table.CafePool, CafeEntry{
		Id:   "cafe-2",
		Name: "Example Cafe 2",
		Link: "https://cafelink.com",
	})

	response, err := table.RenderSchedule()
	if err != nil {
		t.Error("Internal error %w", err)
	}
	if response == "" {
		t.Error("Empty response returned.")
	}
	// run with -v to see
	fmt.Print(response)
}

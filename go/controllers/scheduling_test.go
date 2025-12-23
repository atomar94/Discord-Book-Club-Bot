package controllers

import (
	"slices"
	"testing"

	"bookclubbot.com/main/models"
)

func TestAssignDatesToSchedule_InitializeEmptySchedule(t *testing.T) {

	undated_schedule := models.ScheduleEntry{
		Id:     "123",
		Date:   "",
		BookId: "book-99",
		CafeId: "cafe-42",
	}
	uninitialized_schedule := []models.ScheduleEntry{undated_schedule, undated_schedule}

	err := AssignDatesToSchedule(uninitialized_schedule)
	if err != nil {
		t.Error(err)
	}
	for _, schedule := range uninitialized_schedule {
		if schedule.Date == "" {
			t.Errorf("Date was not initialized.")
		}
	}
}

func TestAssignDatesToSchedule_PicksConsecutiveSaturdays(t *testing.T) {
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
	partial_schedule := []models.ScheduleEntry{dated_schedule, undated_schedule}

	err := AssignDatesToSchedule(partial_schedule)
	if err != nil {
		t.Error(err)
	}

	if partial_schedule[1].Date == "" {
		t.Errorf("Schedule was not extended.")
	}
	if partial_schedule[0].Date == partial_schedule[1].Date {
		t.Errorf("Date was not incremented.")
	}
}

func TestSelectNextBook_TestCases(t *testing.T) {
	tests := []struct {
		name             string
		inputVotes       []int // the vote counts for books going into the algorithm
		validResultVotes []int // The pool of acceptable books, by their vote count
		expectError      bool
	}{
		{
			name:             "Selects from all when less than 5 books",
			inputVotes:       []int{10, 3, 7},
			validResultVotes: []int{10, 3, 7}, // Any of these are valid winners
			expectError:      false,
		},
		{
			name:       "Selects only from top 5 when many books exist",
			inputVotes: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			// Logic: Sort desc -> 10, 9, 8, 7, 6, ...
			// Only the top 5 are eligible. 1-5 should never be returned.
			validResultVotes: []int{10, 9, 8, 7, 6},
			expectError:      false,
		},
		{
			name:       "Handles ties in top 5",
			inputVotes: []int{50, 50, 50, 50, 50, 10},
			// Only the 50s get in. The 10 is the 6th element.
			validResultVotes: []int{50},
			expectError:      false,
		},
		{
			name:             "Error on empty list",
			inputVotes:       []int{},
			validResultVotes: []int{},
			expectError:      true,
		},
	}

	for range 5 {
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// 1. Helper to build the struct inputs from integers
				inputBooks := make([]models.BookEntry, len(tc.inputVotes))
				for i, v := range tc.inputVotes {
					inputBooks[i] = models.BookEntry{
						Votes: v,
						Read:  false, // For this test all books should be considered candidates, so we keep them read.
					}
				}

				// 2. Call the function
				got, err := selectNextBook(inputBooks)

				// 3. Check Error State
				if tc.expectError {
					if err == nil {
						t.Errorf("Expected an error, but got none")
					}
					return // Stop here if we expected an error
				}
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				// 4. Verify the result is in the allowed list
				isValid := slices.Contains(tc.validResultVotes, got.Votes)

				if !isValid {
					t.Errorf("Result vote %d was not in the list of valid outcomes %v", got.Votes, tc.validResultVotes)
				}
			})
		}
	}
}

func TestAssignBooksToSchedule_NoBookSelected(t *testing.T) {
	books := []models.BookEntry{
		{
			Id:   "1",
			Read: false,
		},
		{
			Id:   "2",
			Read: false,
		},
		{
			Id:   "3",
			Read: false,
		},
		{
			Id:   "4",
			Read: false,
		},
	}
	schedules := []models.ScheduleEntry{
		{
			Id: "s1",
		},
		{
			Id: "s2",
		},
		{
			Id: "s3",
		},
		{
			Id: "s4",
		},
		{
			Id: "s5",
		},
		{
			Id: "s6",
		},
		{
			Id: "s7",
		},
		{
			Id: "s8",
		},
	}
	err := AssignBooksToSchedule(books, schedules)
	if err != nil {
		t.Errorf("Internal Error %v", err)
	}
	// Assert book is set to read.
	found_read_book := false
	for _, book := range books {
		if book.Read {
			found_read_book = true
			break
		}
	}
	if !found_read_book {
		t.Errorf("No book was marked as Read.")
	}

	// Assert schedule entries were given a book ID.
	for i, schedule := range schedules {
		if i >= 4 {
			break
		}
		if schedule.BookId == "" {
			t.Errorf("A book was not set for the first 4 schedule entries")
		}
	}
}

func TestAssignBooksToSchedule_FullSchedule(t *testing.T) {
	books := []models.BookEntry{
		{
			Id:   "1",
			Read: true,
		},
		{
			Id:   "2",
			Read: false,
		},
		{
			Id:   "3",
			Read: false,
		},
		{
			Id:   "4",
			Read: false,
		},
	}
	schedules := []models.ScheduleEntry{
		{
			Id:     "s1",
			BookId: "1",
		},
		{
			Id:     "s2",
			BookId: "1",
		},
		{
			Id:     "s3",
			BookId: "1",
		},
		{
			Id:     "s4",
			BookId: "1",
		},
		{
			Id: "s5",
		},
		{
			Id: "s6",
		},
		{
			Id: "s7",
		},
		{
			Id: "s8",
		},
	}
	err := AssignBooksToSchedule(books, schedules)
	if err != nil {
		t.Errorf("Internal Error %v", err)
	}

	// Assert earlier schedules have a book and later schedules do not
	for i, schedule := range schedules {
		if i < 4 {
			if schedule.BookId != "1" {
				t.Errorf("A book was changed at the beginning of our schedule")
			}
		} else {
			if schedule.BookId != "" {
				t.Errorf("A book was set too far in the future")
			}
		}
	}
}

func TestAssignBooksToSchedule_PartialSchedule(t *testing.T) {
	books := []models.BookEntry{
		{
			Id:   "1",
			Read: true,
		},
		{
			Id:   "2",
			Read: false,
		},
		{
			Id:   "3",
			Read: false,
		},
		{
			Id:   "4",
			Read: false,
		},
	}
	schedules := []models.ScheduleEntry{
		{
			Id:     "s1",
			BookId: "1",
		},
		{
			Id:     "s2",
			BookId: "1",
		},
		{
			Id: "s3",
		},
		{
			Id: "s4",
		},
		{
			Id: "s5",
		},
		{
			Id: "s6",
		},
		{
			Id: "s7",
		},
		{
			Id: "s8",
		},
	}
	err := AssignBooksToSchedule(books, schedules)
	if err != nil {
		t.Errorf("Internal Error %v", err)
	}

	// Assert earlier schedules have a book and later schedules do not
	for i, schedule := range schedules {
		if i < 2 {
			if schedule.BookId != "1" {
				t.Errorf("A book was overridden in our first 2 schedule entries")
			}
		}
		if i >= 2 && i < 6 {
			if schedule.BookId == "" {
				t.Errorf("A new book was not added to the partial schedule")
			}
		}
		if i >= 6 {
			if schedule.BookId != "" {
				t.Errorf("A book was set too far in the future")
			}
		}
	}
}

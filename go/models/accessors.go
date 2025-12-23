package models

import "fmt"

func (t *ClubTable) GetCafeById(id string) (CafeEntry, error) {
	for _, c := range t.CafePool {
		if c.Id == id {
			return c, nil
		}
	}
	return CafeEntry{}, fmt.Errorf("No cafe with ID %s", id)
}

func (t *ClubTable) GetScheduleById(id string) (ScheduleEntry, error) {
	for _, s := range t.Schedule {
		if s.Id == id {
			return s, nil
		}
	}
	return ScheduleEntry{}, fmt.Errorf("No schedule with ID %s", id)
}

func (t *ClubTable) GetBookById(id string) (BookEntry, error) {
	for _, b := range t.BookPool {
		if b.Id == id {
			return b, nil
		}
	}
	return BookEntry{}, fmt.Errorf("No book with ID %s", id)
}

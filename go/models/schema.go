package models

type CafeEntry struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type ScheduleEntry struct {
	Id     string `json:"id"`
	Date   string `json:"date"`
	BookId string `json:"book_id"`
	CafeId string `json:"cafe_id"`
}

type BookEntry struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Votes       int    `json:"votes"`
	Read        bool   `json:"read"`
}

type ClubTable struct {
	CafePool []CafeEntry     `json:"cafe_pool"`
	Schedule []ScheduleEntry `json:"schedule"`
	BookPool []BookEntry     `json:"book_pool"`
}

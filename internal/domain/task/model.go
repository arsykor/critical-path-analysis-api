package task

import "time"

type Task struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Predecessors []int     `json:"predecessors"`
	Duration     int       `json:"-"`
	IsConsidered bool      `json:"-"`
	Rank         int       `json:"-"`
}

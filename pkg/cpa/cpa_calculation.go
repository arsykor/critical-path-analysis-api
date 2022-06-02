package cpa

import (
	"critical-path-analysis-api/internal/domain/task"
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"time"
)

var zeroTime time.Time

func Arrange(tasks *[]task.Task) error {
	if len(*tasks) < 2 {
		return nil
	}

	var iterations int
	var rank int
	var tasksConsideredLocally []int

	zeroTime.IsZero()
	projectStartDate := zeroTime

	for i, t := range *tasks {
		(*tasks)[i].Duration = durationExcludingWeekends(t.StartDate, t.EndDate)

		if len(t.Predecessors) == 0 {
			(*tasks)[i].IsConsidered = true
			if projectStartDate == zeroTime || projectStartDate.Before(t.StartDate) {
				projectStartDate = t.StartDate
			}
			iterations++
		}
	}
	fmt.Println("*** for 1")

	//Dividing tasks by ranks:
Loop:
	for true {
		rank++

		for i, t := range *tasks {
			allPredConsidered, err := allPreviousTasksConsidered(t, tasks)
			if err != nil {
				return err
			}

			if !t.IsConsidered && allPredConsidered {
				(*tasks)[i].Rank = rank
				tasksConsideredLocally = append(tasksConsideredLocally, t.Id)
				iterations++
			}
		}

		for _, t := range tasksConsideredLocally {
			getTaskById(t, tasks).IsConsidered = true
		}
		tasksConsideredLocally = nil

		if iterations == len(*tasks) {
			break Loop
		}
	}
	fmt.Println("*** Loop")
	//Enter the data rank-by-rank:
	for i := 1; i <= rank; i++ {
		tasksOfIRang := filterTasks(func(task task.Task) bool { return task.Rank == i }, tasks)

		for _, t := range *tasksOfIRang {
			currentTask := getTaskById(t.Id, tasks)
			latestPredecessor := getLatestPredecessor(t, tasks)

			if latestPredecessor.After(t.StartDate) {
				(*currentTask).StartDate = latestPredecessor.AddDate(0, 0, 1)
				(*currentTask).EndDate = sumIncludingWeekends(currentTask.StartDate, currentTask.Duration)
			}
		}
	}
	fmt.Println("*** rank-by-rank")

	return nil
}

func durationExcludingWeekends(from time.Time, to time.Time) int {
	n := 0
	if to == from {
		return 0
	}
	nextDate := from
	for nextDate.Before(to) {
		if nextDate.Weekday() != 6 && nextDate.Weekday() != 0 {
			n++
		}
		nextDate = nextDate.AddDate(0, 0, 1)
	}
	return n + 1
}

func sumIncludingWeekends(date time.Time, duration int) time.Time {
	if duration == 0 {
		return date
	}

	switch date.Weekday() {
	case time.Weekday(6):
		date = date.AddDate(0, 0, 2)
		duration--
	case time.Weekday(0):
		date = date.AddDate(0, 0, 1)
		duration--
	}

	date = date.AddDate(0, 0, duration/5*7)
	extraDays := duration % 5

	if int(date.Weekday())+extraDays > 5 {
		extraDays += 2
	}

	return date.AddDate(0, 0, extraDays-1)
}

func getLatestPredecessor(task task.Task, tasks *[]task.Task) time.Time {
	max := zeroTime
	for _, p := range task.Predecessors {
		pEndTime := getTaskById(p, tasks).EndDate
		if pEndTime.After(max) {
			max = pEndTime
		}
	}
	return max
}

func filterTasks(suitable func(task.Task) bool, tasks *[]task.Task) *[]task.Task {
	var filteredTasks []task.Task
	for _, t := range *tasks {
		if suitable(t) {
			filteredTasks = append(filteredTasks, t)
		}
	}
	return &filteredTasks
}

func getTaskById(id int, tasks *[]task.Task) *task.Task {
	index := slices.IndexFunc(*tasks, func(t task.Task) bool { return t.Id == id })
	return &(*tasks)[index]
}

func allPreviousTasksConsidered(t task.Task, tasks *[]task.Task) (bool, error) {
	for _, p := range t.Predecessors {
		if p == t.Id {
			resp := fmt.Sprintf("the task cannot include itself in the array of predecessors, task id: %d", t.Id)
			return false, errors.New(resp)
		}

		if !getTaskById(p, tasks).IsConsidered {
			return false, nil
		}
	}
	return true, nil
}

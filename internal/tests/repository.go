package tests

import (
	"critical-path-analysis-api/internal/domain/task"
	"critical-path-analysis-api/pkg/cpa"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"io/ioutil"
	"os"
)

var TestRepository testRepository

type testRepository struct {
	Tasks []task.Task
}

func (tr *testRepository) GetOne(id int) (*task.Task, error) {
	index := slices.IndexFunc(tr.Tasks, func(t task.Task) bool { return t.Id == id })
	if index == -1 {
		return nil, errors.New("there is no task with this id in the array")
	}
	return &tr.Tasks[index], nil
}

func (tr *testRepository) GetAll() *[]task.Task {
	return &tr.Tasks
}

func (tr *testRepository) Create(tasks *[]task.Task) *[]task.Task {
	for _, t := range *tasks {
		tr.Tasks = append(tr.Tasks, t)
	}

	cpa.Arrange(&tr.Tasks)
	return &tr.Tasks
}

func InitTestRepository() {
	var testRep testRepository

	jsonFile, err := os.Open("internal/tests/repository.json")
	if err != nil {
		fmt.Println(err)
	}

	jsonData, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(jsonData, &testRep.Tasks)
	if err != nil {
		fmt.Println(err)
	}

	TestRepository = testRep
}

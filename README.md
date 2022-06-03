# REST API service for Task Management

This is the REST API service for managing tasks. 

The solution provides the following functionality:

* storing tasks in a database (PostgreSQL client or in-memory)
* Critical Path Method technique used for planning and scheduling tasks
 
Using CPA provides the possibility of considering the interconnectedness of the tasks, which leads to the fact that if the tasks lying on the critical path are delayed, trimed or removed the whole project completion will be changed.

## Getting Started
1. Clone the project from GitHub:

```shell
# Clone from GitHub
git clone https://github.com/arsykor/critical-path-analysis-api
```

2. Build and run the project:

```shell
# Build and run
cd critical-path-analysis-api
go build ./critical-path-analysis-api
```

3. Set up a PostgreSQL Database:

All necessary scripts are located in the `database.sql` file.

## API endpoints

When a RESTful API server is ready, it provides the following endpoints (by default it runs at `http://127.0.0.1:8080`):

* `POST /task/create` - creare tasks*
* `GET /task/get/:id` - get task by id
* `GET /task/get` - get all tasks
* `POST /task/update` - update task
* `GET /task/delete/:id` - delete task by id

*You can send all the tasks in the request, regardless of whether it already exists in the database or not. The system will determine which tasks are new, the parameters of which tasks have been changed. The tasks with changed deadlines will be added to the directory of changes. In the future, it will be possible to check the historicity of changes to each task.

You can check the correctness of the application build via `POST /task/create` request with test data.

```shell
# create tasks via: POST /task/create

curl -X POST -H "Content-Type: application/json" -d '
[
  {
    "id": 1,
    "start_date": "2022-01-01T00:00:00Z",
    "end_date": "2022-08-08T00:00:00Z",
    "predecessors": [],
    "name": "Task A"
  },
  {
    "id": 2,
    "start_date": "2022-03-04T00:00:00Z",
    "end_date": "2022-04-17T00:00:00Z",
    "predecessors": [1],
    "name": "Task B"
  }
]
' http://localhost:8080/task/create
```
You can grab some test tasks from repository.json - `critical-path-analysis-api/blob/master/internal/tests/repository.json`

Specify the tasks-predecessors' IDs in the "predecessors" variable `[1,2,3,4]` for tasks, witch are supposed to be started after predecessors ones. For independent tasks or for first tasks of the project use an empty array - `[]`.

## Critical Path Analysis

Interconnection of tasks is requered for tracking changes in task deadlines. If the tasks lying on the critical path are delayed, trimed or removed,  the whole project completion will be changed.

By changing the duration of the task, all the durations of subsequent related tasks are automatically arranged depending on whether the task is on a critical path or not. The shift takes into account the days off of the present calendar.

On the Gantt Chart below Task 4 depends on the completion of Task 2 and Task 3, it can be seen that Task 3 does not lie on the critical path, thus start date for Task 4 is calculated based on the end date of Task 2 (start date is laid out at the earliest):

``` 
+
|                                                                                        
|                                                                                        
| +------------------+                                                                   
| |      Task 1      |-------------+                                                     
| +------------------+             |                                                     
|                                  v                                                     
|                    +---------------------------+                                       
|                    |          Task 2           |-----------+                           
|                    +---------------------------+           |                           
|                                                            |                           
|                    +----------+                            |                           
|                    |  Task 3  |----------------------------|                           
|                    +----------+                            |                           
|                                                            v                           
|                                                 +---------------------+                
|                                                 |       Task 4        |                
|                                                 +---------------------+                
+--------------------------------------------------------------------------------+                                                      
```

## Project Layout
 
```
.
├── cmd                  main application of the project
├── internal             private application and library code
│   ├── adapters         convenient for external agencys middlewares
│   ├── composites       composed into tree structures objects
│   ├── config           configuration
│   ├── domain           enterprise logic and types
│   └── tests            repository for in-memory approach and test data script
├── pkg                  public library code
│   ├── client           database client
│   └── cpa              CPA calculation
├── database.sql         scripts for database creation
└── config.yml           configuration file
```

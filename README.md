# Orchestrator

A lightweight orchestrator that provides automation for deploying, scaling, and otherwise managing docker containers.

## Requirements
- Linux Operating System
- Docker
- Go (v1.20+)

## ⚙️ How to run

1. *Clone the repo:*

```bash
git clone https://github.com/sulaimonshittu/orchestrator.git
```

2. *Install dependencies:*

```bash
go mod tidy
```

3. *Run the application*

```bash
$ WORKER_HOST=localhost WORKER_PORT=8080 MANAGER_HOST=localhost MANAGER_PORT=4000 go run main.go
```

4. *Call the manager endpoints*

- GET /tasks Gets a list of all tasks
- POST /tasks Creates a task JSON-encoded task.TaskEvent body
- DELETE /tasks/{taskID} Stop the task identified by taskID


## Components

- *TASK:* The task is basically a Docker container running a process. it will provide CPU cycles, memory, and networking according to
  the needs of the task.
- *WORKER:* The Worker represents a physical or virtual machine and the worker component of the orchestration system that runs on that
  machine. It is responsible for running/stopping the tasks assigned to it by the manager and keeping track of them.
- *MANAGER:* The manager is the main entry point for users. The manager collects tasks from the users and schedules them to a worker where the task can run. The manager
  also periodically collects metrics from each of its workers, which are used in the scheduling process

## 🔥 Future Improvements
- Implement the enhanced parallel virtual machine (E-PVM) scheduler
- Implement persistent storage of the tasks by adding an embedded Key-Value Store
- Implement a CLI to interact with the Manager API
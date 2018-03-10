package main

import (
  "context"
  "fmt"
  tug "github.com/buchanae/tugboat"
  "github.com/buchanae/tugboat/docker"
)

func main() {
  ctx := context.Background()
  log := tug.EmptyLogger{}
  store := &tug.EmptyStorage{}

	stage, err := tug.NewStage("foo", 0755)
  if err != nil {
    panic(err)
  }

  exec := &docker.Docker{
    Logger: log,
  }

  task := &tug.Task{
    ID: "test1",
    ContainerImage: "alpine",
    Command: []string{"echo", "hello world!"},
  }

  err = tug.Run(ctx, task, stage, log, store, exec)
  if err != nil {
    fmt.Println("RESULT", err)
  } else {
    fmt.Println("Success")
  }
}

package main

import (
  "context"
  "fmt"
  tug "github.com/buchanae/tugboat"
  "github.com/buchanae/tugboat/docker"
  "github.com/buchanae/tugboat/storage/local"
)

func main() {
  ctx := context.Background()
  log := tug.EmptyLogger{}
  store := &local.Local{}

	stage, err := tug.NewStage("tug-workdir", 0755)
  if err != nil {
    panic(err)
  }
  //stage.LeaveDir = true
  defer stage.RemoveAll()

  exec := &docker.Docker{
    Logger: log,
  }

  task := &tug.Task{
    ID: "test1",
    ContainerImage: "alpine",
    Command: []string{"md5sum", "/inputs/infile.txt"},
    Stdout: "out.txt",
    Inputs: []tug.File{
      {
        URL: "inputs/in.txt",
        Path: "/inputs/infile.txt",
      },
    },
    Outputs: []tug.File{
      {
        URL: "output",
        Path: "out.txt",
      },
    },
  }

  err = tug.Run(ctx, task, stage, log, store, exec)
  if err != nil {
    fmt.Println("RESULT", err)
  } else {
    fmt.Println("Success")
  }
}

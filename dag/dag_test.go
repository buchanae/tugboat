package dag

import (
  "fmt"
  "testing"
  "reflect"
  "github.com/kr/pretty"
)


func TestIdle(t *testing.T) {
  steps := []Step{
    tstep{"01", false, false, nil},
    tstep{"02", true, false, nil},
    tstep{"03", false, true, nil},
    tstep{"04", false, false, nil},
    tstep{"05", false, false, nil},
  }
  expected := []Step{steps[0], steps[3], steps[4]}
  idle := Idle(steps)

  if !reflect.DeepEqual(idle, expected) {
    t.Error("unexpected idle")
    pretty.Ldiff(t, idle, expected)
  }
}

func TestRunning(t *testing.T) {
  steps := []Step{
    tstep{"01", false, false, nil},
    tstep{"02", true, false, nil},
    tstep{"03", false, true, nil},
    tstep{"04", false, false, nil},
    tstep{"05", false, false, nil},
  }
  expected := []Step{steps[2]}
  running := Running(steps)

  if !reflect.DeepEqual(running, expected) {
    t.Error("unexpected running")
    pretty.Ldiff(t, running, expected)
  }
}

func TestAllDone(t *testing.T) {
  steps := []Step{
    tstep{"01", false, false, nil},
    tstep{"02", true, false, nil},
    tstep{"03", false, true, nil},
    tstep{"04", false, false, nil},
    tstep{"05", false, false, nil},
  }
  if AllDone(steps) {
    t.Error("steps should not be all done")
  }

  done := []Step{
    tstep{"done-01", true, false, nil},
    tstep{"done-02", true, false, nil},
    tstep{"done-03", true, false, nil},
  }
  if !AllDone(done) {
    t.Error("steps should be all done")
  }
}

func TestDone(t *testing.T) {
  steps := []Step{
    tstep{"01", false, false, nil},
    tstep{"02", true, false, nil},
    tstep{"03", false, true, nil},
    tstep{"04", false, false, nil},
    tstep{"05", true, false, nil},
  }
  expected := []Step{steps[1], steps[4]}
  done := Done(steps)
  if !reflect.DeepEqual(done, expected) {
    t.Errorf("unexpected done")
    pretty.Ldiff(t, done, expected)
  }
}

func TestFailed(t *testing.T) {
  steps := []Step{
    tstep{"01", false, false, nil},
    tstep{"02", true, false, nil},
    tstep{"03", false, false, fmt.Errorf("err")},
  }
  expected := []Step{steps[2]}
  failed := Failed(steps)

  if !reflect.DeepEqual(failed, expected) {
    t.Errorf("unexpected failed")
    pretty.Ldiff(t, failed, expected)
  }
}

func TestErrors(t *testing.T) {
  err1 := fmt.Errorf("err 1")
  err2 := fmt.Errorf("err 2")
  steps := []Step{
    tstep{"01", false, false, nil},
    tstep{"02", true, false, nil},
    tstep{"03", false, false, err1},
    tstep{"04", false, false, err2},
    tstep{"05", false, false, err1},
  }
  expected := []error{err1, err2, err1}
  errors := Errors(steps)

  if !reflect.DeepEqual(errors, expected) {
    pretty.Ldiff(t, errors, expected)
  }
}

func TestLinks(t *testing.T) {
  d, _ := dag1()

  // Test upstream links for step "07"
  up07 := d.Upstream[d.Steps["07"]]
  ex07 := d.GetSteps("05", "06")
  if !reflect.DeepEqual(up07, ex07) {
    t.Error("unexpected upstream")
    pretty.Ldiff(t, up07, ex07)
  }

  // Test downstream links for "03"
  down03 := d.Downstream[d.Steps["03"]]
  ex03 := d.GetSteps("06", "08")
  if !reflect.DeepEqual(down03, ex03) {
    t.Error("unexpected downstream")
    pretty.Ldiff(t, down03, ex03)
  }
}

func TestAllUpstream(t *testing.T) {
  d, _ := dag1()
  up := AllUpstream(d, d.Steps["07"])
  ex := d.GetSteps("05", "01", "06", "02", "03")
  if !reflect.DeepEqual(up, ex) {
    t.Error("unexpected all upstream")
    pretty.Ldiff(t, up, ex)
  }
}

func TestAllDownstream(t *testing.T) {
  d, _ := dag1()
  dn := AllDownstream(d, d.Steps["03"])
  ex := d.GetSteps("06", "07", "08")

  if !reflect.DeepEqual(dn, ex) {
    t.Error("unexpected all downstream")
    pretty.Ldiff(t, dn, ex)
  }
}

func TestReady(t *testing.T) {
  d, steps := dag1()
  ready := Ready(d, steps)
  ex := d.GetSteps("03", "04")

  if !reflect.DeepEqual(ready, ex) {
    t.Error("unexpected ready")
    pretty.Ldiff(t, ready, ex)
  }
}

func TestBlocked(t *testing.T) {
  d, steps := dag1()
  blocked := Blocked(d, steps)
  ex := d.GetSteps("07")

  if !reflect.DeepEqual(blocked, ex) {
    t.Error("unexpected blocked")
    pretty.Ldiff(t, blocked, ex)
  }
}

func TestTerminals(t *testing.T) {
  d, steps := dag1()
  term := Terminals(d, steps)
  ex := d.GetSteps("04", "07", "08")

  if !reflect.DeepEqual(term, ex) {
    t.Error("unexpected term")
    pretty.Ldiff(t, term, ex)
  }
}

func TestNext(t *testing.T) {
  d, steps := dag1()
  next := Next(d, steps)
  ex := d.GetSteps("03", "04")

  if !reflect.DeepEqual(next, ex) {
    t.Error("unexpected next")
    pretty.Ldiff(t, next, ex)
  }
}

func TestCounts(t *testing.T) {
  d, steps := dag1()
  counts := Count(d, steps)

  if counts.Total != 8 {
    t.Errorf("expected total to be 8, but got %d", counts.Total)
  }
  if counts.Idle != 4 {
    t.Errorf("expected idle to be 4, but got %d", counts.Idle)
  }
  if counts.Ready != 2 {
    t.Errorf("expected ready to be 2, but got %d", counts.Ready)
  }
  if counts.Running != 1 {
    t.Errorf("expected running to be 1, but got %d", counts.Running)
  }
  if counts.Done != 3 {
    t.Errorf("expected done to be 3, but got %d", counts.Done)
  }
  if counts.Blocked != 1 {
    t.Errorf("expected blocked to be 1, but got %d", counts.Blocked)
  }
  if counts.Failed != 1 {
    t.Errorf("expected failed to be 1, but got %d", counts.Failed)
  }
}

func dag1() (*DAG, []Step) {
  /*
    01 (D) ---05 (E) ---07
                       /
    02 (D) ---06 (R) --
              /
    03 (I) --
        \
         08 (I)

    04 (I)

    I = Idle
    D = Done
    R = Running
    E = Error
  */
  err1 := fmt.Errorf("err 1")
  steps := []Step{
    tstep{"01", true, false, nil},
    tstep{"02", true, false, nil},
    tstep{"03", false, false, nil},
    tstep{"04", false, false, nil},
    tstep{"05", true, false, err1},
    tstep{"06", false, true, nil},
    tstep{"07", false, false, nil},
    tstep{"08", false, false, nil},
  }
  dag := NewDAG()

  for _, step := range steps {
    t := step.(tstep)
    dag.AddStep(t.id, step)
  }
  must(dag.AddDep("05", "01"))
  must(dag.AddDep("07", "05"))
  must(dag.AddDep("07", "06"))
  must(dag.AddDep("06", "02"))
  must(dag.AddDep("06", "03"))
  must(dag.AddDep("08", "03"))
  return dag, steps
}

func must(err error) {
  if err != nil {
    panic(err)
  }
}

type tstep struct {
  id string
  done, running bool
  err error
}
func (ts tstep) Done() bool {
  return ts.done
}
func (ts tstep) Running() bool {
  return ts.running
}
func (ts tstep) Error() error {
  return ts.err
}
func (ts tstep) ID() string {
  return ts.id
}

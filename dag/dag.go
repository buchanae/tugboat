package dag

import (
	"fmt"
)

// TODO serializable DAG

//type Version interface{}

type Step interface {
	//Version() Version
	Done() bool
	Running() bool
	Error() error

	//Start()
	//Stop()
	//Reset()
}

type DAG struct {
  Steps map[string]Step
	Upstream   map[Step][]Step
	Downstream map[Step][]Step
}

func NewDAG() *DAG {
  return &DAG{
    Steps: map[string]Step{},
    Upstream: map[Step][]Step{},
    Downstream: map[Step][]Step{},
  }
}

func (l *DAG) AddStep(id string, step Step) {
  l.Steps[id] = step
}

func (l *DAG) GetSteps(ids ...string) []Step {
  var steps []Step
  for _, id := range ids {
    steps = append(steps, l.Steps[id])
  }
  return steps
}

func (l *DAG) AddDep(stepID, depID string) error {
  dep, ok := l.Steps[depID]
  if !ok {
    return fmt.Errorf(`missing dependency "%s"`, depID)
  }

  step, ok := l.Steps[stepID]
  if !ok {
    return fmt.Errorf(`missing step "%s"`, stepID)
  }

  l.Upstream[step] = append(l.Upstream[step], dep)
  l.Downstream[dep] = append(l.Downstream[dep], step)
  return nil
}

/* TODO
- create new version of DAG by modifying a step
- manually invalidate a step
- when an intermediate, finished step is invalidated by a new version,
  how are the downstream steps invalidated?
  - how are running, invalidated, downstream steps stopped?

- change links between steps in dag?
  - or just create a new dag at that point? but lose caching?

misc:
- how are task retries handled?
- timeouts
- if a step change version, but that version ends up creating the same outputs as
  the previous version, it's possible to optimize and sort of re-cache this new version.
  how would this work? Does a step need to include its output hashes in its verison hash?
*/

func Idle(steps []Step) []Step {
	var idle []Step
	for _, step := range steps {
		if !step.Done() && !step.Running() {
			idle = append(idle, step)
		}
	}
	return idle
}

func Running(steps []Step) []Step {
	var running []Step
	for _, step := range steps {
		if step.Running() {
			running = append(running, step)
		}
	}
	return running
}

func AllDone(steps []Step) bool {
	for _, step := range steps {
		if !step.Done() {
			return false
		}
	}
	return true
}

func Done(steps []Step) []Step {
	var done []Step
	for _, step := range steps {
		if step.Done() {
			done = append(done, step)
		}
	}
	return done
}

func Failed(steps []Step) []Step {
	var failed []Step
	for _, step := range steps {
		if step.Error() != nil {
			failed = append(failed, step)
		}
	}
	return failed
}

func Errors(steps []Step) []error {
	var errors []error
	for _, step := range steps {
		if err := step.Error(); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func AllUpstream(dag *DAG, step Step) []Step {
	var upstream []Step
	for _, up := range dag.Upstream[step] {
		upstream = append(upstream, up)
		upstream = append(upstream, AllUpstream(dag, up)...)
	}
	return upstream
}

func AllDownstream(dag *DAG, step Step) []Step {
	var downstream []Step
	for _, down := range dag.Downstream[step] {
		downstream = append(downstream, down)
		downstream = append(downstream, AllDownstream(dag, down)...)
	}
	return downstream
}

func Ready(dag *DAG, steps []Step) []Step {
	var ready []Step
	for _, step := range steps {
		if IsReady(dag, step) {
			ready = append(ready, step)
		}
	}
	return ready
}

func Blocked(dag *DAG, steps []Step) []Step {
	var blocked []Step
	for _, step := range steps {
		if IsBlocked(dag, step) {
			blocked = append(blocked, step)
		}
	}
	return blocked
}

func IsBlocked(dag *DAG, step Step) bool {
	for _, upstream := range AllUpstream(dag, step) {
		if upstream.Error() != nil {
			return true
		}
	}
	return false
}

func IsReady(dag *DAG, step Step) bool {
	if step.Done() || step.Running() {
		return false
	}
	for _, dep := range dag.Upstream[step] {
		if !dep.Done() || dep.Error() != nil {
			return false
		}
	}
	return true
}

func Terminals(dag *DAG, steps []Step) []Step {
	var terminals []Step
	for _, step := range steps {
		if len(dag.Downstream[step]) == 0 {
			terminals = append(terminals, step)
		}
	}
	return terminals
}

func Next(dag *DAG, steps []Step) []Step {
	var next []Step
	for _, step := range steps {
		if IsReady(dag, step) {
			next = append(next, step)
		}
	}
	return next
}

type Categories struct {
	Idle,
	Ready,
	Running,
	Done,
	Blocked,
	Failed []Step
}

func Categorize(dag *DAG, steps []Step) Categories {
	return Categories{
		Idle:    Idle(steps),
		Ready:   Ready(dag, steps),
		Running: Running(steps),
		Done:    Done(steps),
		Blocked: Blocked(dag, steps),
		Failed:  Failed(steps),
	}
}

type Counts struct {
	Total,
	Idle,
	Ready,
	Running,
	Done,
	Blocked,
	Failed int
}

func Count(dag *DAG, steps []Step) Counts {
	c := Categorize(dag, steps)
	return Counts{
		Total:   len(steps),
		Idle:    len(c.Idle),
		Ready:   len(c.Ready),
		Running: len(c.Running),
		Done:    len(c.Done),
		Blocked: len(c.Blocked),
		Failed:  len(c.Failed),
	}
}

func FailFast(dag *DAG, steps []Step) ([]Step, error) {
	if AllDone(steps) {
		return nil, nil
	}
	errs := Errors(steps)
	if errs != nil {
		return nil, &ErrorList{errs}
	}
	return Next(dag, steps), nil
}

var ErrNoRunnableSteps = fmt.Errorf("no runnable steps")

func BestEffort(dag *DAG, steps []Step) ([]Step, error) {
	if AllDone(steps) {
		return nil, nil
	}
	next := Next(dag, steps)
	if next == nil {
		return nil, ErrNoRunnableSteps
	}
	return next, nil
}

type ErrorList struct {
	Errors []error
}

func (e *ErrorList) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}

	s := "Errors:"
	for _, err := range e.Errors {
		s += "- " + err.Error()
	}
	return s
}

package simsettlers

import (
	"log"

	"github.com/Flokey82/go_gens/vectors"
)

// TaskStatus is the status of a task.
type TaskStatus byte

const (
	TaskStatusNotStarted TaskStatus = iota
	TaskStatusInProgress
	TaskStatusCompleted
	TaskStatusFailed
)

// Tree is a chain of tasks that can be executed by a person.
type Tree struct {
	Root Task

	OnSuccess func()
	OnFailure func()
}

func NewTree(root Task, onSuccess, onFailure func()) Tree {
	return Tree{
		Root:      root,
		OnSuccess: onSuccess,
		OnFailure: onFailure,
	}
}

// Do executes the tree using a variable time step.
// The tree will be executed from start to finish.
// DO NOT USE THIS FUNCTION, it is a left-over from the previous implementation.
func (t *Tree) Do(elapsed float64) bool {
	task := t.Root
	for task != nil {
		status := task.Do(elapsed)
		if status == TaskStatusCompleted {
			task = task.Next()
		} else if status == TaskStatusFailed {
			t.OnFailure()
			return false
		}
	}
	t.OnSuccess()
	return true
}

// Step executes the next step in the tree using a variable time step.
// The tree will be evaluated from the root, so if a task step fails or
// is invalid (e.g. the person is not at the right location), the tree
// will be restarted from the root.
func (t *Tree) Step(elapsed float64) bool {
	// Do until we either get InProgress or Failed.
	task := t.Root
	for task != nil {
		status := task.Do(elapsed)
		if status == TaskStatusCompleted {
			task = task.Next()
		} else if status == TaskStatusFailed {
			t.OnFailure()
			return false
		} else {
			return true
		}
	}
	t.OnSuccess()
	return true
}

// Task is a task that can be executed by a person.
// A task can be chained together with other tasks.
type Task interface {
	Do(elapsed float64) TaskStatus // Executes the task and returns the status of the task.
	Next() Task                    // Returns the next task to do in the chain of tasks.
	Then(Task) Task                // Returns a new task that will be executed after this task.
}

// TaskThen is a component that can be embedded in a task to allow chaining of tasks.
type TaskThen struct {
	next Task
}

// Then returns a new task that will be executed after this task.
func (t *TaskThen) Then(next Task) Task {
	t.next = next
	return next
}

// Next returns the next task in the chain of tasks.
func (t *TaskThen) Next() Task {
	return t.next
}

// TaskMoveToXY is a task that moves a person to a specific location.
type TaskMoveToXY struct {
	Person *Person
	X, Y   float64
	TaskThen
}

func NewTaskMoveToXY(p *Person, x, y int) *TaskMoveToXY {
	return &TaskMoveToXY{
		Person: p,
		X:      float64(x),
		Y:      float64(y),
	}
}

// arrivalDistance is the distance to the target location that is considered "arrived".
const arrivalDistance = 0.1

// Do executes the task and returns the status of the task.
func (t *TaskMoveToXY) Do(elapsed float64) TaskStatus {
	// Check if we are close enough to the target location.
	if t.Person.distanceTo(t.X, t.Y) < arrivalDistance {
		log.Printf("%s arrived at %d,%d", t.Person.String(), int(t.X), int(t.Y))
		t.Person.Speed = vectors.NewVec2(0, 0)
		return TaskStatusCompleted
	}

	// Move towards the target location.
	// TODO: We should only set the direction if we are not already moving in the right direction.
	log.Printf("%s moving to %d,%d from %d, %d", t.Person.String(), int(t.X), int(t.Y), int(t.Person.X), int(t.Person.Y))
	t.Person.SetDirection(t.X, t.Y)
	t.Person.Move(elapsed)
	return TaskStatusInProgress
}

// TaskMoveToLocation is a task that moves a person to a specific location.
// This task uses a function to get the location, so that the location can be
// updated dynamically.
type TaskMoveToLocation struct {
	Person      *Person
	GetLocation func() vectors.Vec2
	TaskThen
}

func NewTaskMoveToLocation(p *Person, getLocation func() vectors.Vec2) *TaskMoveToLocation {
	return &TaskMoveToLocation{
		Person:      p,
		GetLocation: getLocation,
	}
}

// Do executes the task and returns the status of the task.
func (t *TaskMoveToLocation) Do(elapsed float64) TaskStatus {
	// Get the location to move to.
	location := t.GetLocation()

	// Check if we are close enough to the target location.
	if t.Person.distanceTo(location.X, location.Y) < arrivalDistance {
		log.Printf("%s arrived at %d,%d", t.Person.String(), int(location.X), int(location.Y))
		t.Person.Speed = vectors.NewVec2(0, 0)
		return TaskStatusCompleted
	}

	// Move towards the target location.
	// TODO: We should only set the direction if we are not already moving in the right direction.
	log.Printf("%s moving to %d,%d from %d, %d", t.Person.String(), int(location.X), int(location.Y), int(t.Person.X), int(t.Person.Y))
	t.Person.SetDirection(location.X, location.Y)
	t.Person.Move(elapsed)
	return TaskStatusInProgress
}

// TaskGeneric is a generic task that can be used to implement custom tasks.
type TaskGeneric struct {
	Person *Person
	Name   string
	TaskThen
	do func(elapsed float64) TaskStatus
}

func NewTaskGeneric(p *Person, name string, do func(elapsed float64) TaskStatus) *TaskGeneric {
	return &TaskGeneric{
		Person: p,
		Name:   name,
		do:     do,
	}
}

// Do executes the task and returns the status of the task.
func (t *TaskGeneric) Do(elapsed float64) TaskStatus {
	return t.do(elapsed)
}

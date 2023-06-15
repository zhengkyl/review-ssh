package enums

import (
	"encoding/json"
	"fmt"
)

type Category uint8

const (
	Film Category = iota
	Show
)

type Status uint8

const (
	PlanToWatch Status = iota
	Watching
	Completed
	Dropped
)

// Should not be used to marshal json
func (c Category) String() string {
	switch c {
	case Film:
		return "Film"
	case Show:
		return "Show"
	}
	return "Invalid Category"
}

// These are display strings, should not be used to marshal json
func (s Status) String() string {
	switch s {
	case PlanToWatch:
		return "Plan To Watch"
	case Watching:
		return "Watching"
	case Completed:
		return "Completed"
	case Dropped:
		return "Dropped"
	}
	return "Invalid Status"
}

func (c *Category) UnmarshalJSON(data []byte) (err error) {
	var category string
	if err := json.Unmarshal(data, &category); err != nil {
		return err
	}
	switch category {
	case "Film":
		*c = Film
	case "Show":
		*c = Show
	default:
		return fmt.Errorf("%q is not a valid Category", category)
	}
	return nil
}

func (s *Status) UnmarshalJSON(data []byte) (err error) {
	var status string
	if err := json.Unmarshal(data, &status); err != nil {
		return err
	}
	switch status {
	case "PlanToWatch":
		*s = PlanToWatch
	case "Watching":
		*s = Watching
	case "Completed":
		*s = Completed
	case "Dropped":
		*s = Dropped
	default:
		return fmt.Errorf("%q is not a valid Status", status)
	}
	return nil
}

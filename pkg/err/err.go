package err

import (
	"fmt"

)

type SimulationError struct {
	JobNumber int
	Err       error
}

func (e *SimulationError) SimError() string {
	return fmt.Sprintf("Error in simulation %d: %s", e.JobNumber, e.Err)
}

func (e *SimulationError) Unwrap() error {
	return e.Err
}

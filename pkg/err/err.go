package err

import (
	"fmt"

)

/*
   error codes:
   1: failed to read input csv
   2: meshing error
   3: results saving error
   4: simulation error during iteration probably diverged
   5: other error (probably sim error from mesh)
*/

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

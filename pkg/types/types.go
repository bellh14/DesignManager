package types

type SystemArgs struct {
	InputFile string
}

type DesignParameter struct {
	Name   string
	Min    float64
	Max    float64
	Step   float64
	Mean   float64
	StdDev float64
}

type SystemResourcesType struct {
	Partition string
	Nodes     int
	Ntasks    int
}

type StarCCM struct {
	StarPath  string
	PodKey    string
	JavaMacro string
	SimFile   string
}

type DesignManagerInputParameters struct {
	NumSims               int
	NtasksPerSim          int
	StudyType             string
	OptimizationAlgorithm string
	DesignParameters      []DesignParameter
	DesignObjectives      []DesignObjective
}

type DesignObjective struct {
	Name   string
	Goal   string
	Weight float64
}

type ConfigFile struct {
	// System Resources
	SystemResources SystemResourcesType

	// Project
	WorkingDir string

	// STAR-CCM+
	StarCCM StarCCM

	// Design Manager Parameters
	DesignManagerInputParameters DesignManagerInputParameters
}

type JobSubmissionType struct {
	WorkingDir       string
	Ntasks           int
	StarPath         string
	PodKey           string
	JavaMacro        string
	SimFile          string
	DesignParameters []DesignParameter
}

type ParameterSamples struct {
	Samples []float64
}

type FailedSimulation struct {
	JobNumber int
	Cause     string
}

type DesignObjectiveResult struct {
	DesignObjectiveName   string
	DesignObjectiveResult float64
}

type SimInput struct {
	Name  string
	Value float64
}

type SimulationResult struct {
	JobNumber              int
	InputParameters        []SimInput
	DesignObjectiveResults []float64
}

type GenerationResults struct {
	Number        int
	FailedSims    []FailedSimulation
	SucceededSims []int
}

package types

type SystemArgs struct {
	InputFile string
}

type DesignParameter struct {
	Name string
	Min  float64
	Max  float64
	Step float64
}

type SystemResourcesType struct {
	Partition string
	Nodes     int
	Ntasks    int
}

type StarCCM struct {
	Path string
	PodKey	string
	JavaMacro string
	SimFile string
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
	WorkingDir string
	Ntasks int
	Path string
	PodKey	string
	JavaMacro string
	SimFile string
	JobNumber int
}

type ParameterSamples struct {
	Samples []float64
}


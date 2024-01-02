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
	// System Architecture
	Partition string
	Nodes     int
	Ntasks    int

	// Project
	WorkDir string

	// STAR-CCM+
	PodKey    string
	JavaMacro string // name of the java macro file
	SimFile   string // name of the sim file

	// Design Manager Parameters
	DesignManagerParameters DesignManagerInputParameters
}

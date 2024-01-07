#!/bin/bash

WorkingDir=../../../scripts/run_simulation.sh
Ntasks=4
Path=/opt/Siemens/17.04.008-R8/STAR-CCM+17.04.008-R8/star/bin/
PodKey=1234-5678-9012-3456
JavaMacro=macro.java
SimFile=sim.sim
JobNumber=1
mkdir -p $WORKING_DIR/$JOB_NUMBER

STARCCM_PATH/starccm+ -power -licpath 1999@flex.cd-adapco.com -podkey $PODKEY -batch $WORKING_DIR/$JOB_NUMBER/$JAVA_MACRO $WORKING_DIR/$JOB_NUMBER/$SIM_FILE -np $NCPU -bs slurm -time -batch-report
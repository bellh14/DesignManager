#!/bin/bash

WORKING_DIR=../../../scripts/run_simulation.sh
NCPU=4
PODKEY=1234-5678-9012-3456
JAVA_MACRO=macro.java
SIM_FILE=sim.sim
JOB_NUMBER=1

mkdir -p $WORKING_DIR/$JOB_NUMBER

module load starccm/17.04.007
starccm+ -power -licpath 1999@flex.cd-adapco.com -podkey $PODKEY -batch $WORKING_DIR/$JOB_NUMBER/$JAVA_MACRO $WORKING_DIR/$JOB_NUMBER/$SIM_FILE -np $NCPU -bs slurm -time -batch-report
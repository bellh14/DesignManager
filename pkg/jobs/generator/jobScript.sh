#!/bin/bash

WORKING_DIR=../../../scripts/run_simulation.sh
NCPU=4
PODKEY=1234-5678-9012-3456
JAVA_MACRO=macro.java
SIM_FILE=sim.sim
JOB_NUMBER=1

STARCCM_PATH=/opt/Siemens/17.04.008-R8/STAR-CCM+17.04.008-R8/star/bin/

mkdir -p $WORKING_DIR/$JOB_NUMBER

STARCCM_PATH/starccm+ -power -licpath 1999@flex.cd-adapco.com -podkey $PODKEY -batch $WORKING_DIR/$JOB_NUMBER/$JAVA_MACRO $WORKING_DIR/$JOB_NUMBER/$SIM_FILE -np $NCPU -bs slurm -time -batch-report
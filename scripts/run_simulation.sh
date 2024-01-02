#!/bin/bash

WORKING_DIR="/scratch/ganymede/<user>/NewMapTest/"
NCPU=96
PODKEY="something"
JAVA_MACRO="HeaveSweep.java"
SIM_FILE="2023UTASpec_Rideheight_Final.sim"
JOB_NUMBER=1

module load starccm/17.04.007
starccm+ -power -licpath 1999@flex.cd-adapco.com -podkey $PODKEY -batch "/scratch/ganymede/<user>/NewMapTest/batch_8/HeaveSweep.java" -np $NCPU "/scratch/ganymede/<user>/NewMapTest/batch_8/2023UTASpec_Rideheight_Final.sim" -bs slurm -time -batch-report

#!/bin/bash

#SBATCH --partition=normal
#SBATCH --nodes=6
#SBATCH --ntasks=96
#SBATCH --time=96:00:00
#SBATCH --mail-type=ALL
#SBATCH --output=/scratch/ganymede/<user>/NewMapTest/batch_8/batch_8.out

NCPU=96
PODKEY=""

echo -e "This job allocated $NCPU cores\nJob is allocated on node(s): $SLURM_JOB_NODELIST" >/scratch/ganymede/<user>/NewMapTest/batch_8/batch_8_log.out

module load starccm/17.04.007

starccm+ -power -licpath 1999@flex.cd-adapco.com -podkey $PODKEY -batch "/scratch/ganymede/<user>/NewMapTest/batch_8/HeaveSweep.java" -np $NCPU "/scratch/ganymede/<user>/NewMapTest/batch_8/2023UTASpec_Rideheight_Final.sim" -bs slurm -time -batch-report

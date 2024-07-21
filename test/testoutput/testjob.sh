#!/bin/bash

#SBATCH		-J "testjob"
#SBATCH		-p icx
#SBATCH		-N 1
#SBATCH		-n 4
#SBATCH		-t 24:00:00
#SBATCH		-mail-user=test@gmail.com
#SBATCH		-mail-type=all
#SBATCH		-o "output.txt"
#SBATCH		-e "error.txt"

./DesignManager -config configfile
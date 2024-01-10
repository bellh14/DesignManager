#!/bin/bash

WorkingDir=.
Ntasks=4
Path=/opt/Siemens/17.04.008-R8/STAR-CCM+17.04.008-R8/star/bin/
PodKey=1234-5678-9012-3456
JavaMacro=macro.java
SimFile=sim.sim
JobNumber=1
mkdir $WorkingDir/$JobNumber

$Path/starccm+ -power -licpath 1999@flex.cd-adapco.com -podkey $PodKey -batch $WorkingDir/$JobNumber/$JavaMacro $WorkingDir/$JobNumber/$SimFile -np $Ntasks -time -batch-report
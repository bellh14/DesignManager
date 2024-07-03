#!/bin/bash

WorkingDir=../../../test/testoutput/
Ntasks=4
Path=/opt/Siemens/17.04.008-R8/STAR-CCM+17.04.008-R8/star/bin/
PodKey=1234-5678-9012-3456
JavaMacro=DMPareto.java
SimFile=sim.sim
DesignParameters=[]
module load starccm/17.04.007
starccm+ -power -licpath 1999@flex.cd-adapco.com -podkey $PodKey -batch $WorkingDir/$JavaMacro $WorkingDir/$SimFile -np $Ntasks -time -batch-report

exit_code=$?
if [ $exit_code -ne 0 ]; then
    echo "Error: StarCCM+ exited with non-zero exit code: $exit_code" >&2
    exit $exit_code
fi


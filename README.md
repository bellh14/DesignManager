# Custom Design Manager for STAR-CCM+ FSAE Style CFD Simulations

## Justifications

- StarCCM's built in Design mangager is difficult to setup and unreliable
- Automatic resource handling to parallelize running simulations on any number of nodes
- This allows us to better suit our avaliable HPC resources
- Allows for much better error handling and logging
- Can use whatever DOE method works best for us

### MOO Study
![image](https://github.com/user-attachments/assets/088ef39d-e901-42cf-b9cb-02fa3bd0445d)

### Example Config file for MOO study
```json
{
  "UseDM": true,
  "OutputDir": "./",
  "SlurmConfig": {
    "HostName": "",
    "JobName": "2024MOOBaseline",
    "Partition": "icx",
    "Nodes": 16,
    "Ntasks": 1280,
    "WallTime": "12:00:00",
    "Email": "email",
    "MailType": "all",
    "OutputFile": "output.txt",
    "ErrorFile": "error.txt",
    "WorkingDir": "."
  },
  "WorkingDir": ".",
  "StarCCM": {
    "StarPath": "/18.06.007/STAR-CCM+18.06.007/star/bin",
    "PodKey": "",
    "JavaMacro": "HandleSim.java",
    "SimFile": "2024Car.sim",
    "WorkingDir": "2024BaselineMOO"
  },
  "DesignStudyConfig": {
    "StudyType": "Pareto",
    "MOOConfig": {
      "NumGenerations": 8,
      "NumSimsPerGeneration": 16,
      "OptimizationAlgorithm": "Genetic",
      "MutationRate": 0.125
    },
    "NumSims": 128,
    "NtasksPerSim": 80,
    "OptimizationAlgorithm": "None",
    "DesignParameters": [
      {
        "Name": "Biplane1 AOA",
        "Units": "deg",
        "Min": 0,
        "Max": 10,
        "NumSims": 128
      },
      {
        "Name": "Biplane2 AOA",
        "Units": "deg",
        "Min": -10,
        "Max": 10,
        "NumSims": 128
      },
      {
        "Name": "BiplaneGapSize",
        "Units": "mm",
        "Min": -10,
        "Max": 5,
        "NumSims": 128
      },
      {
        "Name": "FW 4th Element AOA",
        "Units": "deg",
        "Min": -7,
        "Max": 5,
        "NumSims": 128
      },
      {
        "Name": "BiplanePosition",
        "Units": "mm",
        "Min": -380,
        "Max": 0,
        "NumSims": 128
      }
    ],
    "DesignObjectives": [
      {
        "Name": "Aero Efficiency",
        "Weight": 1.0,
        "Goal": "Maximize"
      },
      {
        "Name": "Rear Axle Downforce",
        "Weight": 0.15,
        "Goal": "Maximize"
      }
    ]
  }
}

```


### Work in Progress future GUI
![image](https://github.com/user-attachments/assets/89bf1c42-9a4c-4b00-b028-e9b0cb353aa5)


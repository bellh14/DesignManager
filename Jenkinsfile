pipeline {
    agent any
    triggers {
        githubPush()
        pollSCM 'H/5 * * * *'
    }
    stages {
        stage("Cleaning"){
          steps {
              echo "Cleaning bulid dir"
              sh 'make clean'
          }
        }
        stage("Format"){
            steps {
                echo "Formating using Gofumpt"
                sh 'make format'
            }
        }
       stage("Build"){
            steps {
                echo "Buliding started"
                sh 'make build'
            }
        }
        stage("Test"){
            steps {
                echo "Unit Testing Started"
                sh 'make test'
            }
        }
    }
}

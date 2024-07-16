pipeline {
    agent any
    triggers {
        githubPush()
        pollSCM 'H/5 * * * *'
    }
    stages {
        stage("Format"){
            steps {
                echo "Formating using Gofumpt"
                sh 'make format'
            }
        }
        stage("Lint"){
            steps {
                echo "Linting using Golangci-lint"
                sh 'make lint'
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

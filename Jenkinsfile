pipeline {
    agent any
    triggers {
        githubPush()
        pollSCM 'H/5 * * * *'
    }
    stages {
        stage("Format"){
            steps {
                echo "installing Gofumpt"
                sh 'go install mvdan.cc/gofumpt@latest'
                echo "Formating using Gofumpt"
                sh 'make format'
            }
        }
        stage("Lint"){
            steps {
                echo "Installing linter"
                sh "go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1"
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

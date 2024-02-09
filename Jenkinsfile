pipeline {
    agent any

    stages {
        stage('Build and Deploy Dev') {
            steps {
                script {
                    sh '''
                        docker rm -f gin-container-dev || true
                        docker image prune -af
                        docker build -t sj2go-gin-dev:latest .                       
                        docker run -d --name gin-container-dev -p 8081:8081 --network prod sj2go-gin-dev:latest
                    '''
                }
            }
        }
    }
}

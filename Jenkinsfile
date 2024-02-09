pipeline {
    agent any

    stages {
       stage('Remove dev') {
            steps {
                script {
                    sh '''
                        docker rm -f gin-container-dev || true
                        docker image prune -af
                    '''
                }
            }
        },
         stage('Build dev') {
            steps {
                script {
                    sh '''
                        docker build \
                        --build-arg ENV=${ENV} \
                        -t sj2go-gin-dev:latest .                       
                    '''
                }
            }
        }
        stage('Deploy Dev') {
            steps {
                script {
                    sh '''                     
                        docker run -d --name gin-container-dev -p 8081:8081 --network dev sj2go-gin-dev:latest
                    '''
                }
            }
        }
    }
}

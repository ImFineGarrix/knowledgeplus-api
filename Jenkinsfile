pipeline {
    agent any

    stages {
       stage('Remove') {
            steps {
                script {
                    sh '''
                        docker rm -f gin-container-${ENV} || true
                        docker image prune -af
                    '''
                }
            }
        }
         stage('Build') {
            steps {
                script {
                    sh '''
                        docker build \
                        --build-arg ENV=${ENV} \
                        -t sj2go-gin-${ENV}:latest .                       
                    '''
                }
            }
        }
        stage('Deploy') {
            steps {
                script {
                    sh '''                     
                        docker run -d --name gin-container-${ENV} -p :8081 --network ${ENV} sj2go-gin-${ENV}:latest
                    '''
                }
            }
        }
    }
}

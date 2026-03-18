pipeline {
    agent any

    options {
        disableConcurrentBuilds()
        timestamps()
    }

    stages {
        stage('Probar conexion SSH') {
            steps {
                sh '''
                    ssh -o StrictHostKeyChecking=no root@187.124.147.66 "
                        echo Conexion_OK &&
                        whoami &&
                        hostname &&
                        pwd
                    "
                '''
            }
        }

        stage('Deploy PHG API') {
            steps {
                sh '''
                    ssh -o StrictHostKeyChecking=no root@187.124.147.66 "
                        set -e
                        /root/opt/scripts/deploy-phg-api.sh
                    "
                '''
            }
        }
    }
}
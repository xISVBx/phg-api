pipeline {
    agent any

    options {
        disableConcurrentBuilds()
    }

    stages {
        stage('Probar conexion SSH') {
            steps {
                sh '''
                    ssh root@187.124.147.66 "
                        echo Conexion_OK &&
                        whoami &&
                        hostname &&
                        pwd
                    "
                '''
            }
        }
    }
}
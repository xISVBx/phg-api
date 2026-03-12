pipeline {
    agent any

    options {
        disableConcurrentBuilds()
    }

    stages {
        stage('Probar conexion SSH') {
            steps {
                sh '''
                    ssh hostinger "
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
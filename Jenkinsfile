pipeline {
    agent none

    options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
        timeout(time: 1, unit: 'HOURS')
    }

    environment {
        PROJECT_NAME = 'jokes'
        DOCKER_IMAGE = 'ghcr.io/apimgr/jokes'
        GO_VERSION = '1.21'
    }

    stages {
        stage('Build') {
            parallel {
                stage('Build AMD64') {
                    agent {
                        label 'amd64'
                    }
                    steps {
                        script {
                            sh '''
                                echo "🔨 Building for AMD64..."
                                make clean
                                make build-host
                                ./binaries/jokes --version
                            '''
                        }
                    }
                }

                stage('Build ARM64') {
                    agent {
                        label 'arm64'
                    }
                    steps {
                        script {
                            sh '''
                                echo "🔨 Building for ARM64..."
                                make clean
                                make build-host
                                ./binaries/jokes --version
                            '''
                        }
                    }
                }
            }
        }

        stage('Test') {
            agent {
                label 'amd64'
            }
            steps {
                script {
                    sh '''
                        echo "🧪 Running tests..."
                        make test
                    '''
                }
            }
        }

        stage('Docker Build') {
            when {
                branch 'main'
            }
            agent {
                label 'amd64'
            }
            steps {
                script {
                    sh '''
                        echo "🐳 Building Docker images..."
                        make docker
                    '''
                }
            }
        }

        stage('Release') {
            when {
                tag pattern: "v\\d+\\.\\d+\\.\\d+", comparator: "REGEXP"
            }
            agent {
                label 'amd64'
            }
            steps {
                script {
                    sh '''
                        echo "🚀 Creating release..."
                        make release
                    '''
                }
            }
        }
    }

    post {
        success {
            echo '✅ Pipeline completed successfully!'
        }
        failure {
            echo '❌ Pipeline failed!'
        }
        cleanup {
            cleanWs()
        }
    }
}

# Balance challenge
This is a Go application that processes a file from a mounted directory and sends summary information to a user via email.

## Prerequisites
- Go (version 1.22.0)
- Docker
- An AWS account

## Getting Started
The application can be run in Docker containers or deployed to AWS Lambda.

### Clone the Repository
```sh
git clone https://github.com/yourusername/your-repo.git
cd balance
```

### Run the app on Containers
The application has two Dockerfiles: one for the MySQL database and another for the Go application. 
Follow these steps to build and run the application:
- Create a Docoker bridge network.
```sh 
docker network create local_network
```
- Build and run mysql image.
```sh 
cd database
docker run -it --rm -p 3306:3306 --name database mysql:v1
docker network connect local_network database
```
-Build and run Go application image.
```sh 
cd ..
docker build -t golang:v1 .
docker run -it --rm --network local_network --name balance golang:v1
```

### Deploy the app to AWS Lambda.
Follow these steps to deploy the application to AWS Lambda:
- Modify main.go
    Uncomment the dependency: github.com/aws/aws-lambda-go/lambda (line 17)
    Comment out the lines 34-39
    Uncomment line 42
- Install the AWS Lambda Go SDK
```sh 
go get "github.com/aws/aws-lambda-go v1.47.0" 
```
- Build the application.
```sh 
GOOS=linux go build -o balance main.go
```
- Create deployment package.
```sh 
zip challenge.zip balance
```
- Upload the Deployment Package to AWS
    In your AWS account, upload the zipped deployment package to an S3 bucket.
    Create an AWS Lambda function with the runtime: Amazon Linux 2023.
    Upload the code from the S3 bucket where the deployment package is stored.

## Important Details
Configure the application using environment variables.
The application uses a .env file to set various configuration parameters, including database connection details and SMTP settings for sending emails.




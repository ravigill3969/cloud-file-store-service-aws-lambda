# cloud-file-store-service-aws-lambda

go mod init go-lambda

go get github.com/aws/aws-lambda-go/lambda

GOOS=linux CG0_ENABLED=0 GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go

zip go-lambda.zip bootstrap
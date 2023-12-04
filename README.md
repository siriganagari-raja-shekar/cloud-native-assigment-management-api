# webapp
Rest API for CSYE6225 Fall 2023

## Steps to install and start web server

1. Install Go dependencies
   
   ```
   go get ./...
   ```

2. Run tests
  
   ```
   go test ./...
   ```

3. Build the application and run

   ```
   go build cmd/app/main.go
   ./main.exe
   ```
   
## Step to import SSL certificate into AWS
   
   ```
   aws acm import-certificate --certificate fileb://your/dir/certificate-file --private-key fileb://your/dir/private-key-file --certificate-chain fileb://your/dir/certificate-chain-file --profile your_profile --region your_aws_region
   ```

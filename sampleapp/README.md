## Getting Started

1. Install Go
2. Clone this version of the [opentelemetry-go-contrib repo](https://github.com/open-o11y/opentelemetry-go-contrib.git) into your GOPATH
3. Go into the newly cloned repo and make a new branch off master
4. Merge `wilguo-idgenerator` branch into your new branch
5. Merge `wilguo-eksdetector` branch into your new branch
6. Clone this version of [opentelemetry-go](https://github.com/Aneurysm9/opentelemetry-go.git) into your GOPATH and checkout the `OTG-1351` branch
7. Clone this version of [opentelemetry-go-contrib](https://github.com/Tenaria/opentelemetry-go-contrib/tree/aws-sdk-s3) 
8. Clone the sample app from the [aws-otel-go repo here](https://github.com/aws-observability/aws-otel-go)
9. In the `aws-otel-go` repo, switch to the `wilguo-sample-app` branch (You should now see a folder called `sampleapp`)
10. Open the `go.mod` file change lines 6-8 in the `replace()` to point the imports to point to your local copy of those repositories
11. Start the OTEL collector
12. Go into the `sampleapp` folder and run the following commands
    ```
    go build
    ./sampleapp
    ```
11. Visit the following endpoints
    ```
    localhost:8080/aws-sdk-call
    localhost:8080/outgoing-http-call
    ```
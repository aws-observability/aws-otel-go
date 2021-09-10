## Getting Started

1. Install Go 1.15 or higher
2. Go into the `sampleapp` directory and run the following commands
    ```
    go build -o sampleapp (build binry)
    ./sampleapp (execute binary)
    ```
3. Visit the following endpoints
    ```
    localhost:8080/aws-sdk-call or localhost:8080/outgoing-http-call
    ```
4. Make sure to set `OTEL_EXPORTER_OTLP_ENDPOINT` (default value set by sample app is `0.0.0.0:4317`)

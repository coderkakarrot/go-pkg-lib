# Pubsub example

This is a simple example of how to publish a message to a Pub/Sub topic using pubsub library.

## Requirements

- Go 1.22 or higher
- A Google Cloud account
- A Google Cloud Pub/Sub topic

## Usage

1. Clone the repository:

```bash
git clone https://github.com/coderkakarrots/go-pkg-lib.git
```

2. Navigate to the project directory:

`cd examples/pubsub/`

3. Update the following config:

```go
const (
	PROJECTID = "<gcp-project>>"
	TOPICNAME = "<topic-name>"
)
```

4. Run the application

```bash
go mod download
go run main.go
```
It should print the message something like below

`Message ID:  11267691719266129`

5. Verify the message is published

Navigate to the default subscription on the topic and click `Pull`, you should see the `Hello World!` message

https://console.cloud.google.com/cloudpubsub/subscription/detail/<subscription-name>?project=<project>&tab=messages

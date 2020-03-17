# datastore-exporter

Go CLI tool to export Google Cloud Datastore. Exports all kinds and namespaces and saves to a GCS bucket. As the Datastore export operation is asynchronous the process will poll the state of the operation reporting any errors.

To run locally a service account is required with the correct permissions to export from Datastore: 
1. Create a service account with the roles of Cloud Datastore Import Export Admin
2. Download the key as JSON to your local machine
3. Set the GOOGLE_APPLICATION_CREDENTIALS environment variable to the path of the account key e.g. `export GOOGLE_APPLICATION_CREDENTIALS=/path/to/sa.json

## Development

 - Go 1.11+
 - Dependencies managed with `go mod`

### Setup

These steps will describe how to setup this project for active development. Adjust paths to your desire.

1. Clone the repository: `git clone github.com/thisissoon/datastore-exporter datastore-exporter`
2. Build: `make build`
3. 🍻

### Dependencies

Dependencies are managed using `go mod` (introduced in 1.11), their versions
are tracked in `go.mod`.

To add a dependency:
```
go get url/to/origin
```

### Configuration

Configuration can be provided through a toml file, these are loaded
in order from:

- `/etc/datastore-exporter/datastore-exporter.toml`
- `$HOME/.config/datastore-exporter.toml`

Alternatively a config file path can be provided through the
-c/--config CLI flag.

#### Example datastore-exporter.toml
```toml
timeout = "1h" # duration the process will run for (default is 1 hour)

[log]
console = true
level = "debug"  # [debug|info|error]

[gcs]
projectID = "project-id" # ID of the Google Cloud Project
bucketName = "bucket-name" # name of the bucket for exports to be saved in
```

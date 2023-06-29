# Datastore Exporter

Helm Chart to run the Datastore Exporter application as a Kubernetes CronJob.

## Installation

```
helm upgrade --install \
	datastore-exporter \
	./charts/datastore-exporter
```

## Configuration

| Parameter                                | Description                                 | Default                              |
|------------------------------------------|---------------------------------------------|--------------------------------------|
| `image.repository`                       | Container image to deploy                   | `soon/datastore-exporter`            |
| `image.pullPolicy`                       | Image pull policy                           | `IfNotPresent`                       |
| `schedule`                               | Cron schedule to trigger job                | `0 0 * * *`                          |
| `gcloud.project`                         | ID of the Google Cloud project              | ``                                   |
| `gcloud.bucket`                          | GCS Bucket name you wish to export into     | ``                                   |
| `gcloud.serviceAccount.secretName`       | Secret housing the Google Cloud credentials | `datastore-exporter-service-account` |
| `gcloud.serviceAccount.key`              | Key of the credentials file                 | `credentials.json`                   |
| `gcloud.serviceAccount.workloadIdentity` | Name of Workload Identity Service Account   | `` 									|
| `imagePullSecrets`                       |                                             | `{}`                                 |
| `nameOverride`                           |                                             | ``                                   |
| `fullnameOverride`                       |                                             | ``                                   |
| `resources`                              | Resource allocation (YAML)                  | `{}`                                 |
| `nodeSelector`                           |                                             | `{}`                                 |
| `tolerations`                            |                                             | `[]`                                 |
| `affinity`                               |                                             | `{}`                                 |

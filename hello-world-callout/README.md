# Hello world

This repo shows a sample of how to use the External Callout policy with Apigee X/hybrid.

## Components

### Hello world service

A hello-world sample external callout service is implemented [here](./cmd/server). You can install the service on Cloud Run (GCE and GKE also work perfectly fine). A Cloud Build install script is included [here](./cloudbuild.yaml). You can build & deploy with

```sh
gcloud builds submit --substitutions=_$PROJECT_ID="<PROJECT_NAME_HERE>"
```

### Setup Target Server

Configure a target server to call the External Service.

```sh
auth="Authorization: Bearer $(gcloud auth print-access-token)"

curl "https://apigee.googleapis.com/v1/organizations/$ORG/environments/$ENV/targetservers" -X POST -H $auth -H "Content-Type: application/json" --data-raw '{
  "name": "grpcserver",
  "host": "PUT_HOSTNAME_HERE",
  "port": 443,
  "isEnabled": true,
  "protocol": "GRPC",
  "sSLInfo": {
      "enabled": true,
      "ignoreValidationErrors": true
  }
}'
```

### API Proxy

Deploy the API Proxy included [here](./apiproxy).

___

## Support

This is not an officially supported Google product
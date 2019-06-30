# go-functions

The first time, execute the follows:

```sh
gcloud beta functions deploy lunch --set-env-vars PROJECT_NAME=xxx,SLACK_TOKEN=xxx --project xxx --runtime go111 --trigger-http --entry-point Lunch
```

From the second time, execute the follows:

```sh
gcloud beta functions deploy lunch --project xxx
```

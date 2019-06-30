!#/bin/bash

gcloud beta functions deploy lunch --runtime go111 --entry-point Lunch --trigger-http --project my-project-go-funcntion-test
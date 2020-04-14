#!/bin/sh
gcloud app deploy app.yaml
gcloud app deploy queue.yaml
gcloud app deploy cron.yaml
gcloud datastore indexes create index.yaml

#gcloud app deploy app.yaml --project=kaffeeshare
#gcloud app deploy queue.yaml --project=kaffeeshare
#gcloud app deploy cron.yaml --project=kaffeeshare
#gcloud datastore indexes create index.yaml --project=kaffeeshare

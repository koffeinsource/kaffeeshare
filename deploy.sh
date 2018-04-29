#!/bin/sh
gcloud app deploy app.yaml
gcloud app deploy queue.yaml
gcloud app deploy cron.yaml

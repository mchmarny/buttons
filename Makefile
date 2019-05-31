PROJECT=s9-demo
REGION=us-central1
TOPIC=button-actions
FUNC=button-action-handler
SCHEMA=button
TABLE=actions

all: url

topic:
	gcloud pubsub topics create ${TOPIC}

deploy:
	gcloud alpha functions deploy $(FUNC) \
		--entry-point GitHubEventHandler \
		--set-env-vars SEC=$(HOOK_SECRET),TOP=$(PUBSUB_TOPIC),PRJ=$(GCP_PROJECT) \
		--memory 128MB \
		--region $(REGION) \
		--runtime go112 \
		--trigger-http

policy:
	gcloud alpha functions add-iam-policy-binding $(FUNC) \
		--region $(REGION) \
		--member allUsers \
		--role roles/cloudfunctions.invoker

url:
	gcloud functions describe github-event-handler \
		--region $(REGION) \
		--format='value(httpsTrigger.url)'

table:
	bq mk $(SCHEMA)
	bq mk --schema id:string,repo:string,type:string,actor:string,event_time:timestamp,countable:boolean -t $(SCHEMA).$(TABLE)

test:
	go test ./... -v

deps:
	go mod tidy


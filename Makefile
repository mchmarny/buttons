TOPIC=button-actions

topic:
	gcloud pubsub topics create ${TOPIC}

topicless:
	gcloud pubsub topics delete ${TOPIC}

mod:
	go mod tidy
	go mod vendor

image: mod
	gcloud builds submit \
		--project ${GCP_PROJECT} \
		--tag gcr.io/${GCP_PROJECT}/buttons:0.1.1

service:
	gcloud beta run deploy buttons \
		--region=us-central1 \
		--concurrency=80 \
		--memory=256Mi \
		--allow-unauthenticated \
		--image=gcr.io/${GCP_PROJECT}/buttons:0.1.1 \
		--update-env-vars="secret=${HOOK_SECRET},project=${GCP_PROJECT},topic=${TOPIC}"

test:
	go test ./... -v

post:
	curl -H "content-type: application/json" -H "token: ${HOOK_SECRET}" \
		-d '{ "version": "v0.1.0", "type": "button", "color": "white", "action": "single-click" }' \
		-X POST https://buttons-4afw4gizxa-uc.a.run.app
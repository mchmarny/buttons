
topic:
	gcloud pubsub topics create clicks

topicless:
	gcloud pubsub topics delete clicks

mod:
	go mod tidy
	go mod vendor

image: mod
	gcloud builds submit \
		--project knative-samples \
		--tag gcr.io/knative-samples/buttons:0.1.2

service:
	gcloud beta run deploy buttons \
		--region=us-central1 \
		--concurrency=80 \
		--memory=256Mi \
		--allow-unauthenticated \
		--image=gcr.io/knative-samples/buttons:0.1.2 \
		--update-env-vars="secret=${HOOK_SECRET}"

serviceless:
	gcloud beta run services delete buttons

test:
	go test ./... -v

post:
	curl -H "content-type: application/json" -H "token: ${HOOK_SECRET}" \
		-d '{ "version": "v0.1.0", "type": "button", "color": "white", "click": 2 }' \
		-X POST https://buttons-4afw4gizxa-uc.a.run.app
# Using Flick buttons with Cloud Run on GCP

<img align="right" src="image/flic.png" alt="Flic button">

Co-worker recently told me about this simple IoT button from [flic.io](https://flic.io/) that allows you to wire-up all kinds of custom actions to a single, or double click, or even press-and-hold type of action.

Having done this short demo I can already think about few really interesting applications, but, to start with though, I wanted to create a simple service that would allow me to push the data sent from the button to Cloud PubSub which then connects me to the entire potfolio of actuation options on through GCP servcies.

The recently launched [Cloud Run](https://cloud.google.com/run/) seamed like a perfect platform to take the newelly ordered [3-pack](https://flic.io/shop/flic-4pack) for a spin.

In this repo I will show you how to:

* Deploy generic Cloud Run service that will persist sent data to Cloud PubSub
* Configure Flic buttons to sent data to Cloud Run service


## Deploy Cloud Run Service

First, create a pub/sub topic called `clicks`

```shell
gcloud pubsub topics create clicks
```

That should result with

```shell
Created topic [projects/YOUR_PROJECT_ID/topics/clicks].
```

Next deploy the generic Cloud Run service called `buttons`. The code for that service is in this repo for you to review. There is already public image available (see below) but if you want to, you can build your own copy with this command:

```shell
gcloud builds submit --tag gcr.io/YOUR_PROJECT_ID/buttons:0.1.1
```

> Quickstart on building images using Cloud Build is available [here](https://cloud.google.com/run/docs/quickstarts/build-and-deploy)

If you don't want to build images yourself however you can use the one I already published (`gcr.io/knative-samples/buttons:0.1.1`).

Before we deploy the Cloud Run service we have to create a `secret` which we will use to ensure that only your buttons data will be accepted. To do that, run the following command to after your replace the `your-long-and-super-secret-string` string with something more secure ;)

```shell
export SECRET="your-long-and-super-secret-string"
```

> For more secure way to defining secrets on GCP you can use [berglas](https://github.com/GoogleCloudPlatform/berglas)

Now that we have the `SECRET` defined, we can deploy the Cloud Run service. A couple of flags worth highlighting in the bellow command:

* `concurrency` - the button service is thread safe and doesn't store any internal state so we can turn the concurrency to maximum. More on concurrency [here](https://cloud.google.com/run/docs/about-concurrency)
* `allow-unauthenticated` - By default Cloud Run creates private services which can't be access by anonymous users. Since our buttons don't support more complex authentication, we will expose the Cloud Run service to the public and validate each request using `token` string in request header. More on allowing public access [here](https://cloud.google.com/run/docs/authenticating/public)


```shell
gcloud beta run deploy buttons \
    --region=us-central1 \
    --concurrency=80 \
    --allow-unauthenticated \
    --image=gcr.io/knative-samples/buttons:0.1.2 \
    --update-env-vars="SECRET=${SECRET}"
```

The response from the above command should look something like this

```shell
Deploying container to Cloud Run service [buttons] in project [YOUR_PROJECT_ID] region [us-central1]
✓ Deploying... Done.
  ✓ Creating Revision...
  ✓ Routing traffic...
  ✓ Setting IAM Policy...
Done.
Service [buttons] revision [buttons-00001] has been deployed and is serving traffic at https://buttons-*******-uc.a.run.app
```

You can test the depliyed service using `curl`. Just make sure you replace the `*` part of the URL with the actual `URL` returned by the above command.

```shell
curl -H "content-type: application/json" -H "token: ${SECRET}" \
    -d '{ "version": "v0.1.0", "type": "button", "color": "white", "click": 1 }' \
    -X POST https://buttons-*******-uc.a.run.app
```

## Configuring Flic buttons

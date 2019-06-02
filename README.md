# Using Flick buttons with Cloud Run on GCP

<img align="right" src="image/flic.png" alt="Flic button">

A co-worker recently told me about [flic.io](https://flic.io/) buttons that can easily wire up single or double click, or even press-and-hold type of action to all kinds of actions.

After quick review I could quickly think of a few really interesting applications. To start with though, I wanted to create a simple service that would allow me to push the data sent from the button to Cloud PubSub which then connects me to the entire portfolio of actuation options through GCP APIs and services.

The recently launched [Cloud Run](https://cloud.google.com/run/) seamed like a perfect platform to develop such a service on so I went ahead and ordered [3-pack](https://flic.io/shop/flic-4pack).

In this demo I will show you how to:

* Deploy Cloud Run service that will persist sent data to Cloud PubSub
* and, how to configure Flic buttons to sent data to Cloud Run service


## Prerequisites

* gcloud command-line tool ([how-to](https://cloud.google.com/sdk/gcloud/))
* Enabled GCP API
  * `gcloud services enable pubsub.googleapis.com`
  * `gcloud services enable run.googleapis.com`


## Deployment

### Cloud PubSub Topic

To store the data sent from each button action, first, we need to create a Cloud PubSub topic named `clicks`

```shell
gcloud pubsub topics create clicks
```

That should result with

```shell
Created topic [projects/YOUR_PROJECT_ID/topics/clicks].
```

### Cloud Run Service

Next deploy the generic Cloud Run service called `buttons`. The code for that service is in this repository for you to review. There is already a public image available (see below), but if you want to, you can build your own copy with this command:

```shell
gcloud builds submit --tag gcr.io/YOUR_PROJECT_ID/buttons:0.1.1
```

> For more information on how to build images using Cloud Build see [here](https://cloud.google.com/run/docs/quickstarts/build-and-deploy). If however you want to skip building your image, you can use the already pre-build public image (`gcr.io/knative-samples/buttons:0.1.1`).

Before we deploy the Cloud Run service we have to create a `secret` which we will use to ensure that only data from your buttons will be accepted. To do that, replace the `your-long-and-super-secret-string` string below with something more secure and define an environment variable that will hold your secret using this command:

```shell
export SECRET="your-long-and-super-secret-string"
```

> For more secure way to defining secrets on GCP you can use [berglas](https://github.com/GoogleCloudPlatform/berglas)

Now that we have the `SECRET` defined, we can deploy the Cloud Run service. A couple of flags worth highlighting in the bellow command:

* `concurrency` - the button service is thread safe and doesn't store any internal state so we can turn the concurrency to maximum. More on concurrency [here](https://cloud.google.com/run/docs/about-concurrency)
* `allow-unauthenticated` - By default Cloud Run creates private services which can't be accessed by anonymous users. Since our buttons don't support more complex authentication scheme, we will expose the Cloud Run service to the public and validate each request using `token` string in request header. More on allowing public access [here](https://cloud.google.com/run/docs/authenticating/public)


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

You should be able to see that service in Cloud Run UI now

<img src="image/cr-ui.png" alt="Cloud Run UI">

You can also test the deployed service using `curl`. Just make sure you replace the `***` part of the URL with the actual `URL` returned by the above command.

```shell
curl -H "content-type: application/json" -H "token: ${SECRET}" \
    -d '{ "version": "v0.1.0", "type": "button", "color": "white", "click": 1 }' \
    -X POST https://buttons-*******-uc.a.run.app
```

## Configuring Flic Button

To setup Flic buttons on your device () follow the [start instructions](https://start.flic.io/). The Flic buttons also come with all kinds of [preprogrammed actions](https://flic.io/all-functions).

<img src="image/flic-req.png" alt="Flic Internet Request Action">

To execute above configured Cloud Run service in this demo, we will use `Internet Request` action. To do that you will need to select one of your buttons, expand the `Advanced` category, and configure an action for either `single`, `double`, or `hold` click.

<img src="image/flic-nav.png" alt="Flic Internet Request Action">

That will get you to the `HTTP Internet Request` action configuration screen

<img src="image/flic-conf.png" alt="Flic Internet Request Action Configuration">

Few fields to configure here:

* **URL** - The full URL of the Cloud Run Service
* **Method** - `POST`
* **Content Type** - `application/json`
* **HTTP Header** - Kye: `token`, Value: the value of the previously defined `SECRET`. Make sure you click the `Insert` button to "save" the header parameter before clicking `Done` to save the entire action.

Now your Flic button is configured for use with Cloud Run

## Demo

Assuming all the above deployment steps were completed successfully, you should be able to click the button and see the following:

1. Entries in the Cloud Run service log

<img src="image/cr-log.png" alt="Cloud Run Log">

2. Entries in the Cloud Run metrics chart

<img src="image/cr-metric.png" alt="Cloud Run Metric">

3. Stackdriver PubSub topic (`clikc`) metric chart

<img src="image/sd-metric.png" alt="Stackdriver Metric">

## Summary

Hopefully this demo gave you an idea on how to connect Cloud Run services to IOT devices. With the basic implementation in place we can start working on more creative solutions.







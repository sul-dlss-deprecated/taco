# taco 🌮🌮🌮
The next generation repository system for DLSS
![taco](https://user-images.githubusercontent.com/92044/34897877-016a4e36-f7b6-11e7-80e3-4edecfb2f89d.gif)

## Developing

1. Install serverless: `npm install -g serverless`

## Deploying Lambdas

1. Change into the directory of the service you care about
1. Set up your credentials. https://serverless.com/framework/docs/providers/aws/guide/credentials/
1. `AWS_PROFILE=my-taco-profile sls deploy`


## Swagger API

This configuration is for AWS API Gateway.  It was retrieved by going to the API, selecting the "prod" under "Stages" and then doing "Export" and selecting "Export as Swagger + API Gateway Extensions"


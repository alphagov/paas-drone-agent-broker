paas-drone-agent-broker
=======================

[![Build Status](https://cloud.drone.io/api/badges/richardTowers/paas-drone-agent-broker/status.svg)](https://cloud.drone.io/richardTowers/paas-drone-agent-broker)

A service broker for managing [drone.io](https://drone.io) agents on AWS EC2.

Use case
--------

Modern development teams use continuous deployment to deliver software faster,
with lower risk. To be able to deploy changes teams need a system which can
build, test and deploy their code in a repeatable way.

Tools like [Jenkins](https://jenkins.io/),
[Travis](https://docs.travis-ci.com/), [CircleCI](https://circleci.com/) and
[drone.io](https://drone.io) enable this.

Small teams need a simple, low cost way of using one of these systems.
Unfortunately, most solutions tend to be one or more of:

* Expensive (particularly for closed-source products)
* Slow due to limits on concurrent jobs
* Difficult to host

This project allows users of platforms like GOV.UK PaaS to use a self-hosted
version of the community edition of [drone.io](https://drone.io) which is
low-cost, secure, and is only limited by the cost of infrastructure.

```
# Desired user experience on cloud foundry:

## Create the DB
$ cf create-service postgres tiny-unencrypted-9.5 drone-db

## Create the drone server (by deploying its docker image to CF)
$ cf push drone-server ...

## Create the drone agent using the service broker
$ cf create-service drone-agent tiny my-drone-agent -c '{server: "https://my-drone.cloudapps.digital"}'

## Bind the drone agent to the server
$ cf bind-service ...

## Restage the drone server to pick up the bound service
$ cf restage drone-server
```

Building the code
-----------------

```
export GO111MODULES=on
go build -mod=vendor
```


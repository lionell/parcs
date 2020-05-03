# PARCS

This is a state-of-the-art implementation of the PARCS system described in [the paper][paper]. It's heterogeneous in the nature implying that it's language agnostic.
Heavily relies on the [Docker Swarm][swarm] and Docker's networking features. Under the hood it's a set of a libraries that allow one to operate a swarm in a PARCS-specific fashion.

Features:
* Works on any Docker Swarm cluster
* Language agnostic design
* Web UI powered by [Swarmpit][swarmpit]
* Full isolation (thanks to Docker)
* Flexible dependency management
* Centralized module repository

## Cluster management

One of the biggest advantages of this PARCS implementation is that it doesn't require any special setup. Any Docker Swarm cluster will do!
There are many nice tutorials that teach how to bring up a Docker Swarm cluster: [1][cluster-1], [2][cluster-2] and [3][cluster-3].

Some of you might be wondering why Swarm and not Kubernetes. The reason here is that Swarm has all the features that I need to implement
PARCS on top of it and now that it ships with a default Docker installation it's a no brainer.

Now let me repeat these tutorials and show how easy it is to deploy Docker Swarm on the [Google Cloud Platform][gcp].

### Docker Swarm on GCP

I'm gonna be using a Google Cloud CLI tool called [gcloud][gcloud] to orchestrate the cluster from the command line. Also we'll use
[Google Compute Engine][gce] as a basement for the cluster.

1. First of all you need to make sure that `gcloud` is linked to your accont. If it's the first time you use the CLI just fire `gcloud init`
and follow the instructions. I'm also gonna set up the sensible defaults for my region via `gcloud config set compute/zone us-west1-b`.

2. Now let's start a couple of VMs that are will form a cluster later. Here I'm creating a cluster of 4 nodes that will be managed by a leader.

```
$ gcloud compute instances create leader worker-1 worker-2 worker-3

Created [https://www.googleapis.com/compute/v1/projects/genuine-ember-275604/zones/us-west1-b/instances/leader].
Created [https://www.googleapis.com/compute/v1/projects/genuine-ember-275604/zones/us-west1-b/instances/worker-1].
Created [https://www.googleapis.com/compute/v1/projects/genuine-ember-275604/zones/us-west1-b/instances/worker-2].
Created [https://www.googleapis.com/compute/v1/projects/genuine-ember-275604/zones/us-west1-b/instances/worker-3].

NAME      ZONE        MACHINE_TYPE   PREEMPTIBLE  INTERNAL_IP  EXTERNAL_IP     STATUS
leader    us-west1-b  n1-standard-1               10.138.0.6   35.247.55.235   RUNNING
worker-1  us-west1-b  n1-standard-1               10.138.0.8   35.233.219.127  RUNNING
worker-2  us-west1-b  n1-standard-1               10.138.0.7   34.83.142.137   RUNNING
worker-3  us-west1-b  n1-standard-1               10.138.0.5   35.247.116.107  RUNNING
```

[paper]: https://www.scirp.org/journal/paperinformation.aspx?paperid=78011 
[swarm]: https://docs.docker.com/engine/swarm
[swarmpit]: https://swarmpit.io
[cluster-1]: https://docs.docker.com/engine/swarm/swarm-tutorial/create-swarm
[cluster-2]: https://training.play-with-docker.com/swarm-service-discovery
[cluster-3]: https://rominirani.com/docker-swarm-tutorial-b67470cf8872
[gcp]: http://cloud.google.com
[gcloud]: https://cloud.google.com/sdk/gcloud
[gce]: https://cloud.google.com/compute

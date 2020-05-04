# PARCS

This is a state-of-the-art implementation of the PARCS system described in [the paper][paper]. It's heterogeneous in the nature implying that
it's language agnostic. Heavily relies on the [Docker Swarm][swarm] and [Docker's][docker] networking features. Under the hood it's a set of a libraries
that allow one to operate a swarm in a PARCS-specific fashion.

Features:
* Works on any Docker Swarm cluster
* Language agnostic design
* Web UI powered by [Swarmpit][swarmpit]
* Full isolation (containerization)
* Flexible dependency management
* Centralized module repository

## Cluster management

One of the biggest advantages of this PARCS implementation is that it doesn't require any special setup. Any Docker Swarm cluster will do!
There are many nice tutorials that teach how to bring up a Docker Swarm cluster: [[1]][cluster-1], [[2]][cluster-2] and [[3]][cluster-3].

### Why Swarm and not Kubernetes?

Some of you might be wondering why I decided to base project on Swarm and not on [Kubernetes][kubernetes]. There are multiple reasons for
this decision:

* Kubernetes cluster is harder to set up and maintain.
* PARCS only need a small subset of the features provided by platform.
* Docker Swarm is now shipped as part of the standard Docker distribution.
* Sticking to the [KISS principle][kiss].

As a matter of fact I'm only using load balancing(distributed service deployment) and automatic service discovery(overlay network) features
provided by both platforms. Everything else is a plain old Docker and it's friends.

There are more detailed articles on the web that compare two platforms on the multiple axes: [[1]][swarm-vs-kubernetes-1], [[2]][swarm-vs-kubernetes-2].
Now let's take a look at how easy it is to bring up a Docker Swarm cluster on the [Google Cloud Platform][gcp].

### Swarm cluster on GCP

I'm gonna be using a Google Cloud CLI tool called [gcloud][gcloud] to orchestrate the cluster from the command line. Also we'll use
[Google Compute Engine][gce] as a basement for the cluster.

1. First of all you need to make sure that `gcloud` is linked to your accont. If it's the first time you use the CLI just fire `gcloud init`
and follow the instructions. I'm also gonna set up the sensible defaults for my region via `gcloud config set compute/zone us-west1-b`.

2. Now let's start a couple of VMs that are will form a cluster later. Here I'm creating a cluster of 4 nodes that will be managed by a leader.

    ```console
    me@laptop~:$ gcloud compute instances create leader worker-1 worker-2 worker-3

    Created [https://www.googleapis.com/compute/v1/projects/ember-27/zones/us-west1-b/instances/leader].
    Created [https://www.googleapis.com/compute/v1/projects/ember-27/zones/us-west1-b/instances/worker-1].
    Created [https://www.googleapis.com/compute/v1/projects/ember-27/zones/us-west1-b/instances/worker-2].
    Created [https://www.googleapis.com/compute/v1/projects/ember-27/zones/us-west1-b/instances/worker-3].
    
    NAME      ZONE        MACHINE_TYPE   PREEMPTIBLE  INTERNAL_IP  EXTERNAL_IP     STATUS
    leader    us-west1-b  n1-standard-1               10.138.0.6   35.247.55.235   RUNNING
    worker-1  us-west1-b  n1-standard-1               10.138.0.8   35.233.219.127  RUNNING
    worker-2  us-west1-b  n1-standard-1               10.138.0.7   34.83.142.137   RUNNING
    worker-3  us-west1-b  n1-standard-1               10.138.0.5   35.247.116.107  RUNNING
    ```

3. Unfortunatelly the default Debian image doesn't ship Docker by default, but we can use this [convenience script][convenience-script] to install
the engine as follows

    ```console
    me@laptop~:$ gcloud compute ssh leader
    me@leader~:$ curl -fsSL https://get.docker.com -o get-docker.sh
    me@leader~:$ sudo sh get-docker.sh
    ```

    Make sure that you do this step for every node in the cluster replacing `leader` with a corresponding name.

4. It's time to initialize a swarm. We can do this by `ssh`-ing into a `leader` and running commands:

    ```console
    me@laptop~:$ gcloud compute ssh leader
    me@leader~:$ sudo docker swarm init
    
    Swarm initialized: current node (p7ywd9wbh6th1hy6t5hlsqv0w) is now a manager.
    
    To add a worker to this swarm, run the following command:
    
        docker swarm join --token \
          SWMTKN-1-4cj55yg229l3updnigyz86p63x9bb599htytlmtbhulo4m633d-4kcfduodzvitw4y52flh19g32 \
          10.138.0.6:2377
    ```

5. Having a `join-token` from the previous step we can connect `worker` nodes to a `leader` like follows:

    ```console
    me@laptop~:$ gcloud compute ssh worker-1
    me@worker-1~:$ sudo docker swarm join --token \
             SWMTKN-1-4cj55yg229l3updnigyz86p63x9bb599htytlmtbhulo4m633d-4kcfduodzvitw4y52flh19g32 \
             10.138.0.6:2377
    
    This node joined a swarm as a worker.
    ```

    Don't forget to do this step **for each one** of the `worker` nodes you created.

6. **IMPORTANT!** PARCS needs `leader`-s Docker Engine to listen on the port `4321`.

    This is the only extra step that users have to take to be able to run PARCS on a
    barebones Docker Swarm cluster. Here are commands that do exactly that

    ```console
    me@laptop~:$ gcloud compute ssh leader
    me@leader~:$ sudo sed -i '/ExecStart/ s/$/ -H tcp:\/\/0.0.0.0:4321/' \
                        /lib/systemd/system/docker.service
    me@leader~:$ sudo systemctl daemon-reload
    me@leader~:$ sudo systemctl restart docker
    ```

Now we have a fully configured Docker Swarm cluster ready to run PARCS services.

### Running PARCS module

All the PARCS modules (aka services) should be accessible from some Docker registry. We're going to use a
default [Docker Hub][docker-hub] registry here as an example. All the example code can be found in this repo
under `go/` and `py/` subdirs.

Let's take a look at the simple PARCS service written in Python that given a number `N` and a range `[A; B)`
just iterates in a range looking for divisors of `N`. This service can be implemented like:

```python
import logging
from parcs.server import Service, serve

class Factor(Service):
    def run(self):
        n, a, b = self.recv(), self.recv(), self.recv()
        facts = []
        for i in range(a, b):
            if n % i == 0:
                facts.append(i)
        self.send(facts)

serve(Factor())
```

Now assuming this code lives in the file `main.py` we can build a Docker image for this program by running:

```console
me@laptop~:$ wget https://raw.githubusercontent.com/lionell/parcs/master/py/Dockerfile
me@laptop~:$ cat Dockerfile
FROM lionell/parcs-py

COPY main.py .
CMD [ "python", "main.py" ]

me@laptop~:$ docker build -t lionell/factor .
Sending build context to Docker daemon  3.072kB
Step 1/3 : FROM lionell/parcs-py
 ---> ef17f28e7f39
Step 2/3 : COPY main.py .
 ---> Using cache
 ---> 28f7e7b055d6
Step 3/3 : CMD [ "python", "main.py" ]
 ---> Using cache
 ---> ea69c5b3c156
Successfully built ea69c5b3c156
Successfully tagged lionell/factor:latest
```

### Cleaning up

Don't forget to remove all created VMs. If you don't do it GCP can charge you!

```console
me@laptop~:$ gcloud compute instances delete leader worker-1 worker-2 worker-3
```

[paper]: https://www.scirp.org/journal/paperinformation.aspx?paperid=78011 
[docker]: https://www.docker.com
[swarm]: https://docs.docker.com/engine/swarm
[swarmpit]: https://swarmpit.io
[cluster-1]: https://docs.docker.com/engine/swarm/swarm-tutorial/create-swarm
[cluster-2]: https://training.play-with-docker.com/swarm-service-discovery
[cluster-3]: https://rominirani.com/docker-swarm-tutorial-b67470cf8872
[kubernetes]: https://kubernetes.io
[kiss]: https://en.wikipedia.org/wiki/KISS_principle
[swarm-vs-kubernetes-1]: https://vexxhost.com/blog/kubernetes-vs-docker-swarm-containerization-platforms
[swarm-vs-kubernetes-2]: https://thenewstack.io/kubernetes-vs-docker-swarm-whats-the-difference
[gcp]: http://cloud.google.com
[gcloud]: https://cloud.google.com/sdk/gcloud
[gce]: https://cloud.google.com/compute
[convenience-script]: https://docs.docker.com/engine/install/debian/#install-using-the-convenience-script
[docker-hub]: https://hub.docker.com

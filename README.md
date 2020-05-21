![PARCS logo](/logo.png)

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

### Why not Kubernetes?

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

#### TL; DR One can use `gce-parcs.sh` script to bootstrap cluster in under 5 mins!

```console
me@laptop~:$ sh gce-parcs.sh
Number of workers: 3
GCE instances for leader and 3 workers created
Docker installed on leader
Docker installed on worker-1
Docker installed on worker-2
Docker installed on worker-3
Docker Swarm initialized
PARCS port (4321) is open on leader
Overlay network created for PARCS
Swarmpit installed
Firewall rule for Swarmpit created
---------------------------------------
LEADER_URL=tcp://10.138.0.25:4321
Dashboard URL: http://34.83.234.248:888
Login: admin
Password: password
```

#### Manual setup

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
    me@leader~:$ curl -fsSL https://get.docker.com | sudo sh
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
    
7. **IMPORTANT!** PARCS is also utilizing a custom overlay network that one can create by typing:

    ```console
    me@laptop~:$ gcloud compute ssh leader
    me@leader~:$ sudo docker network create -d overlay parcs
    ```

Now we have a fully configured Docker Swarm cluster ready to run PARCS services.

### Writing PARCS modules

All the PARCS modules (aka services) should be accessible from some Docker registry. We're going to use a
default [Docker Hub][docker-hub] registry here as an example. All the example code can be found in this repo
under `go/` and `py/` subdirs.

#### Example service

Let's take a look at the simple PARCS service written in Python that given a number `N` and a range `[A; B)`
just iterates in a range looking for divisors of `N`. This service can be implemented like:

```python
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

me@laptop~:$ docker build -t lionell/factor-py .
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
Successfully tagged lionell/factor-py:latest

me@laptop~:$ docker push lionell/factor-py:latest
```

PARCS provides base Docker images for all supported languages: [lionell/parcs-py][parcs-py], [lionell/parcs-go][parcs-go]

#### Example runner

PARCS needs a special type of jobs that will kick off the computation. These are **Runners** and they can be
implemented in a very similar way to a servises:

```go
package main

import (
	"github.com/lionell/parcs/go/parcs"
	"log"
	"os"
	"strconv"
)

type Program struct {
	*parcs.Runner
}

func (h *Program) Run() {
	n, _ := strconv.Atoi(os.Getenv("N"))
	t, _ := h.Start("lionell/factor-py")
	t.SendAll(n, 1, n+1)
	var facts []int
	t.Recv(&facts)
	log.Printf("Factors %v", facts)
	t.Shutdown()
}

func main() {
	parcs.Exec(&Program{parcs.DefaultRunner()})
}
```

Again, assuming this code lives in the file `main.go` we can build a Docker image for this program by running:

```console
me@laptop~:$ wget https://raw.githubusercontent.com/lionell/parcs/master/go/Dockerfile
me@laptop~:$ cat Dockerfile
FROM lionell/parcs-go

COPY main.go .
CMD [ "go", "run ", "main.go" ]

me@laptop~:$ docker build -t lionell/runner-go .
Sending build context to Docker daemon  3.584kB
Step 1/3 : FROM lionell/parcs-go
 ---> 4d6be0e795ec
Step 2/3 : COPY main.go .
 ---> f1fe810151ba
Step 3/3 : CMD [ "go", "run", "main.go" ]
 ---> Running in bd9c0480b072
Removing intermediate container bd9c0480b072
 ---> 63a8a590eefc
Successfully built 63a8a590eefc
Successfully tagged lionell/runner-go:latest

me@laptop~:$ docker push lionell/runner-go:latest
```

### Running PARCS modules

In order to run a PARCS runner on a cluster you need to know **internal IP of the leader**. It can be obtained from
the Google Compute Engine UI or by firing this command:

```console
me@laptop~:$ gcloud compute instances list | grep leader | awk '{print "tcp://" $4 ":4321"}'
tcp://10.138.0.6:4321
```

Now to start a service just do this:

```console
me@laptop~:$ gcloud compute ssh leader
me@leader~:$ sudo docker service create \
                    --network parcs \
                    --restart-condition none \
                    --env LEADER_URL=tcp://<LEADER INTERNAL IP>:4321 \
                    --name runner \
                    --env N=123456789 \
                    lionell/runner-go

bjchstu57756oq5ppa6lgg1c3
overall progress: 1 out of 1 tasks 
1/1: running   [==================================================>] 
verify: Service converged 

me@leader~:$ sudo docker service logs -f runner
runner.1.luy@worker-2 05:59:31 Welcome to PARCS-Go!
runner.1.luy@worker-2 05:59:31 Running your program...
runner.1.luy@worker-2 06:00:06 Connection to silly_shtern established
runner.1.luy@worker-2 06:00:17 Factors [1 3 9 3607 3803 10821 11409 32463 34227 13717421 41152263 123456789]
runner.1.luy@worker-2 06:00:17 Connection to silly_shtern closed
runner.1.luy@worker-2 06:00:17 Bye!

me@leader~:$ sudo docker service rm runner
```

Last 3 parameters are ones that change between invocations:
* `--name` is a way to give an invocation a name that can later be used to obtain the results of a specific invocation.
* `--env N=123456789` is specific for this particular runner and tells that we're interested in divisors of that number.
* `lionell/runner-go` is a Docker image that contains the runner itself.

### PARCS Web UI

#### Installing Swarmpit

Comprehensive guide for the installation can be found on the [Github page][swarmpit-install]. Here I'll show the easiest
way to do it

```console
me@laptop~:$ gcloud compute ssh leader
me@leader~:$ sudo docker run -it --rm \
                    --name swarmpit-installer \
                    --volume /var/run/docker.sock:/var/run/docker.sock \
                    swarmpit/install:1.9
...
Summary
Username: admin
Password: password
Swarmpit is running on port :888

Enjoy :)
```

It will ask you to set up an admin account for the control panel. Make sure that you remember the credentials as you'll
need them later to access the dashboard.

#### Setting up a firewall

The last step is to make a firewall aware of the Swarmpit. We want to expose a default port `888` to the outside world.

```console
me@laptop~:$ gcloud compute firewall-rules create swarmpit --allow tcp:888
Creating firewall... Done.

NAME     NETWORK  DIRECTION  PRIORITY  ALLOW    DENY  DISABLED
swarmpit default  INGRESS    1000      tcp:888        False
```

After this is done you can navigate to the external IP address of the `leader` and use a beautiful web UI to manage the
cluster. Here's one way to obtain the URL

```console
me@laptop~:$ gcloud compute instances list | grep leader | awk '{print "http://" $5 ":888"}'
http://35.247.55.235:888
```

![Swarmpit dashboard image](/dashboard.png)

### Cleaning up

Don't forget to remove all created VMs. If you don't do it GCP can charge you!

```console
me@laptop~:$ gcloud compute instances delete leader worker-1 worker-2 worker-3
me@laptop~:$ gcloud compute firewall-rules delete swarmpit
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
[parcs-py]: https://hub.docker.com/repository/docker/lionell/parcs-py
[parcs-go]: https://hub.docker.com/repository/docker/lionell/parcs-go
[swarmpit-install]: https://github.com/swarmpit/swarmpit

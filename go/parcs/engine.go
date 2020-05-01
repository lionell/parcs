package parcs

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/docker/docker/client"
	"log"
	"os"
)

const PARCSNetwork = "parcs"

type Engine struct {
	leaderUrl string
	client    *docker.Client
}

func NewEnvEngine() *Engine {
	return NewEngine(os.Getenv("LEADER_URL"))
}

func NewEngine(leaderUrl string) *Engine {
	client, err := docker.NewClient(leaderUrl, "", nil, nil)
	if err != nil {
		log.Fatalf("Error while connecting to the Swarm Leader: %v", err)
	}
	_, err = client.Ping(context.Background())
	if err != nil {
		log.Fatalf("Error while pinging Swarm Leader: %v", err)
	}
	return &Engine{
		leaderUrl: leaderUrl,
		client:    client,
	}
}

func (e *Engine) Start(image string) (*Task, error) {
	return NewTask(image, e)
}

func (e *Engine) createService(image string) (id string, err error) {
	spec := swarm.ServiceSpec{
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image: image,
				Env:   []string{"LEADER_URL=" + e.leaderUrl},
			},
			RestartPolicy: &swarm.RestartPolicy{
				Condition: swarm.RestartPolicyConditionNone,
			},
			Networks: []swarm.NetworkAttachmentConfig{
				swarm.NetworkAttachmentConfig{
					Target: PARCSNetwork,
				},
			},
		},
	}
	s, err := e.client.ServiceCreate(
		context.Background(),
		spec,
		types.ServiceCreateOptions{},
	)
	if err != nil {
		return
	}
	return s.ID, nil
}

func (e *Engine) queryServiceName(id string) (name string, err error) {
	s, _, err := e.client.ServiceInspectWithRaw(
		context.Background(),
		id,
		types.ServiceInspectOptions{}
	)
	if err != nil {
		return
	}
	return s.Spec.Name, nil
}

func (e *Engine) removeService(id string) error {
	return e.client.ServiceRemove(context.Background(), id)
}

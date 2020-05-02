import docker
from parcs.task import Task

PARCS_OVERLAY_NETWORK = 'parcs'

class Engine:
    def __init__(self, leader_url=None):
        if not leader_url:
            import os
            leader_url = os.environ.get('LEADER_URL')
        self.leader_url = leader_url

        self.client = docker.DockerClient(base_url=leader_url)
        assert(self.client.ping())

    def run(self, image):
        return Task(
            self.client.services.create(
              image=image,
              env=[f'LEADER_URL={self.leader_url}'],
              networks=[PARCS_OVERLAY_NETWORK],
              restart_policy=docker.types.RestartPolicy()
            )
        )


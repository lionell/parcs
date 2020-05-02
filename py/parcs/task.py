import socket
import logging
from parcs.network import *

class Task:
    def __init__(self, service):
        self.service = service
        self.logger = logging.getLogger(service.name)
        self.client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

        ip = dns_lookup(service.name)
        connect(self.client, ip)
        handshake(self.client, side='client')
        self.logger.info(f'Connection to {service.name} established')

    def name(self):
        return self.service.name

    def shutdown(self):
        self.client.close()
        self.service.remove()
        self.logger.info(f'Connection to {service.name} closed')

    def send(self, obj):
        send(self.client, obj)

    def send_all(self, *objs):
        for obj in objs:
            self.send(obj)

    def recv(self):
        return recv(self.client)


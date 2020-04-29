import socket
import logging
from parcs.network import (
    dns_lookup,
    connect,
    handshake,
    send,
    recv,
)

class Task:
    def __init__(self, service):
        self.service = service
        self.logger = logging.getLogger(service.name)
        self.client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

        ip = dns_lookup(service.name)
        self.logger.info(f'Found {service.name} at {ip}')
        connect(self.client, ip)
        self.logger.info('Connected successfully')
        handshake(self.client, side='client')
        self.logger.info('Connection established!')

    def name(self):
        return self.service.name

    def shutdown(self):
        self.client.close()
        self.service.remove()

    def send(self, data):
        self.logger.info(f'Sending {data} over the wire')
        send(self.client, data)

    def recv(self):
        data = recv(self.client)
        self.logger.info(f'Received {data} from the wire')
        return data


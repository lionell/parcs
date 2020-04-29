import socket
import logging
from parcs.engine import Engine
from parcs.network import PORT, send, recv, handshake

class Service:
    def __init__(self):
        self.logger = logging.getLogger('Service')
        self.engine = Engine()

        self.server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.server.bind(('', PORT))
        self.server.listen()

    def start(self):
        self.client, (ip, unused_port) = self.server.accept()
        handshake(self.client, side='server')
        self.logger.info(f'Client from {ip} connected')

    def shutdown(self):
        self.client.close()
        self.server.close()

    def send(self, data):
        self.logger.info(f'Sending {data} over the wire')
        send(self.client, data)

    def recv(self):
        data = recv(self.client)
        self.logger.info(f'Received {data} from the wire')
        return data

def serve(service):
    logging.basicConfig(
        format='%(asctime)s - %(levelname)s %(message)s',
        datefmt='%d-%b-%y %H:%M:%S',
        level=logging.INFO
    )
    try:
        service.start()
        service.run()
    finally:
        service.shutdown()

import socket
import logging
from parcs.engine import Engine
from parcs.network import PORT, send, recv, handshake

class Executable:
    def __init__(self):
        self.logger = logging.getLogger('Executable')

    def run(self):
        raise NotImplementedError()

    def start(self):
        self.logger.info('Execution started')

    def shutdown(self):
        self.logger.info('Execution finished')


class Runner(Executable):
    def __init__(self):
        super().__init__()
        self.engine = Engine()
        self.logger = logging.getLogger('Runner')


class Service(Runner):
    def __init__(self):
        super().__init__()
        self.logger = logging.getLogger('Service')

        self.server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.server.bind(('', PORT))
        self.server.listen()

    def start(self):
        super().start()
        self.client, (ip, unused_port) = self.server.accept()
        handshake(self.client, side='server')
        self.logger.info(f'Client from {ip} connected')

    def shutdown(self):
        self.client.close()
        self.server.close()
        super().shutdown()

    def send(self, data):
        self.logger.info(f'Sending {data} over the wire')
        send(self.client, data)

    def recv(self):
        data = recv(self.client)
        self.logger.info(f'Received {data} from the wire')
        return data

def serve(executable):
    logging.basicConfig(
        format='%(asctime)s - %(levelname)s %(message)s',
        datefmt='%d-%b-%y %H:%M:%S',
        level=logging.INFO
    )
    try:
        executable.start()
        executable.run()
    finally:
        executable.shutdown()

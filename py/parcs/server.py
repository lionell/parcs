import socket
import logging
from parcs.engine import Engine
from parcs.network import PORT, send, recv, handshake

class Runner:
    def __init__(self):
        self.engine = Engine()

    def run(self):
        raise NotImplementedError()

    def start(self):
        pass

    def shutdown(self):
        pass


class Service(Runner):
    def __init__(self):
        super().__init__()
        self.server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.server.bind(('', PORT))
        self.server.listen()

    def start(self):
        self.client, (ip, unused_port) = self.server.accept()
        handshake(self.client, side='server')
        logging.info(f'Client from {ip} connected')

    def shutdown(self):
        self.client.close()
        self.server.close()

    def send(self, obj):
        send(self.client, obj)

    def send_all(self, *objs):
        for obj in objs:
            self.send(obj)

    def recv(self):
        return recv(self.client)

def serve(executable):
    logging.basicConfig(
        format='%(asctime)s [%(levelname)s] %(message)s',
        datefmt='%y/%m/%d %H:%M:%S',
        level=logging.INFO
    )
    logging.info('Welcome to PARCS-Python!')
    try:
        executable.start()
        logging.info('Running your program...')
        executable.run()
    finally:
        executable.shutdown()
        logging.info('Bye!')

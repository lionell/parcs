from parcs.server import Service, serve
import logging

class Echo(Service):
    def run(self):
        n = self.recv()
        logging.info(f'Received %d. Sending it back...', n)
        self.send(n)

serve(Echo())

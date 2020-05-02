import time
import logging
from parcs.server import Service, serve

class Sleeper(Service):
    def run(self):
        n = self.recv()
        logging.info(f'Going to sleep for {n} seconds')
        time.sleep(n)
        logging.info('Just woke up')
        self.send(n)

serve(Sleeper())

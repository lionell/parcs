import time
import logging
from parcs.server import Service, serve

class Sleeper(Service):
    def run(self):
        logging.info('Sleeper service beginning')
        n = int(self.recv())
        logging.info(f'Going to sleep for {n} seconds')
        time.sleep(n)
        logging.info('Just woke up')
        self.send(str(n).encode())
        logging.info('Sleeper service ended')

serve(Sleeper())

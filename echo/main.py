from parcs.server import Service, serve
from parcs.data import marshal, unmarshal

class Echo(Service):
    def run(self):
        n = unmarshal(self.recv())
        logging.info(f'Received %d', n)
        self.send(marhsal(n))

serve(Echo())

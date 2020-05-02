import logging
from parcs.server import Service, serve

class Factor(Service):
    def run(self):
        n = self.recv()
        a = self.recv()
        b = self.recv()
        logging.info(f'Looking for factors of {n} between {a} and {b}')
        facts = []
        for i in range(a, b):
            if n % i == 0:
                logging.info(f'{n} is divisible by {i}')
                facts.append(i)
        self.send(facts)

serve(Factor())

from parcs.server import Service, serve
from parcs.data import unpack_int, pack_ints

class Factor(Service):
    def run(self):
        n = unpack_int(self.recv())
        a = unpack_int(self.recv())
        b = unpack_int(self.recv())
        self.logger.info(f'Looking for factors of {n} between {a} and {b}')
        facts = []
        for i in range(a, b):
            if n % i == 0:
                self.logger.info(f'{n} is divisible by {i}')
                facts.append(i)
        self.send(pack_ints(facts))

serve(Factor())

from parcs.server import Service, serve
from parcs.data import unpack_int, pack_ints

class Factor(Service):
    def run(self):
        n = unpack_int(self.recv())
        a = unpack_int(self.recv())
        b = unpack_int(self.recv())
        facts = []
        for i in range(a, b + 1):
            if n % i == 0:
                self.logger.info(f'{n} is divisible by {i}')
                facts.append(i)
        self.send(pack_ints(facts))

serve(Factor())

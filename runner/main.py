import os
import math
from parcs.server import Service, serve
from parcs.data import pack_int, unpack_ints

class Runner(Service):
    def run(self):
        n = int(os.environ.get('N'))
        t = self.engine.run('lionell/parcs-factor')
        t.send(pack_int(n))
        t.send(pack_int(1))
        t.send(pack_int(n))
        facts = unpack_ints(t.recv())
        self.logger.info(f'Factors found: {facts}')

serve(Runner())

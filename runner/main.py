import os
import math
from parcs.server import Runner, serve
from parcs.data import pack_int, unpack_ints

class Example(Runner):
    def run(self):
      n = int(os.environ.get('N'))
      t = self.engine.run('lionell/parcs-factor')
      t.send(pack_int(n))
      t.send(pack_int(2))
      t.send(pack_int(n - 1))
      facts = unpack_ints(t.recv())
      t.shutdown()
      self.logger.info(f'Factors found: {facts}')

serve(Example())
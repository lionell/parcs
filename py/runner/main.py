import os
import math
import logging
from parcs.server import Runner, serve

def split(n, p):
    chunk = math.ceil((n - 2) / p)
    a = 2
    b = a + chunk
    res = []
    while (b < n):
        res.append((a, b))
        a += chunk
        b += chunk
    if a < n:
        res.append((a, n))
    return res


class Example(Runner):
    def run(self):
        n = int(os.environ.get('N'))
        p = int(os.environ.get('P'))
        logging.info(f'Got N = {n}, P = {p} from env vars')
        tasks = []
        for (a, b) in split(n, p):
            t = self.engine.run('lionell/factor-py')
            t.send_all(n, a, b)
            tasks.append(t)
        facts = []
        for t in tasks:
            facts += t.recv()
        for t in tasks:
            t.shutdown()
        logging.info(f'Factors found: {facts}')

serve(Example())

import logging
from parcs.engine import Engine

e = Engine()
tasks = []
for i in range(3, 10):
    t = e.run('lionell/parcs-sleeper:latest')
    t.send(str(i * 10).encode())
    tasks.append(t)
for t in tasks:
    logging.info(f'Waiting on {t.name()} to finish')
    val = t.recv().decode()
    logging.info(f'{t.name()} returned {val}')
    t.shutdown()

from parcs.engine import Engine

# How to set up a cluster
# 1. Create a Docker Swarm master ($ docker swarm init)
# 2. Connect workers ($ docker swarm join ...)
# 3. Create an overlay network called 'parcs'
# 4. Make Leader accessible on some port like 4321

e = Engine()

t = e.run('lionell/parcs-echo:0.1')
t.send('bla'.encode())
t.logger.info(t.recv().decode())
t.shutdown()


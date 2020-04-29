from parcs.server import Service, serve

class Echo(Service):
    def run(self):
        msg = self.recv()
        self.send(msg)

serve(Echo())

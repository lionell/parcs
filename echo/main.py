from parcs.server import Service, serve

class Echo(Service):
    def run(self):
        self.send(self.recv())

serve(Echo())

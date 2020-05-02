from parcs.server import Service, serve

class Factor(Service):
    def run(self):
        self.logger.info('BEFORE N')
        n = self.recv()
        self.logger.info(f'BEFORE A, n = {n}')
        a = self.recv()
        self.logger.info(f'BEFORE B, a = {a}')
        b = self.recv()
        self.logger.info(f'Looking for factors of {n} between {a} and {b}')
        facts = []
        for i in range(a, b):
            if n % i == 0:
                self.logger.info(f'{n} is divisible by {i}')
                facts.append(i)
        self.send(facts)

serve(Factor())

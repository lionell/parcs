.PHONY: runner factor python all

all: runner factor

python:
	cd python; docker build -t lionell/parcs-python .
	docker push lionell/parcs-python:latest

runner: python
	cd runner; docker build -t lionell/parcs-runner .
	docker push lionell/parcs-runner:latest

factor: python
	cd factor; docker build -t lionell/parcs-factor .
	docker push lionell/parcs-factor:latest

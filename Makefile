.PHONY: runner factor

runner:
	cd runner; docker build -t lionell/parcs-runner .
	docker push lionell/parcs-runner:latest

factor:
	cd factor; docker build -t lionell/parcs-factor .
	docker push lionell/parcs-factor:latest

.PHONY: parcs-go parcs-py echo-go echo-py runner-go runner-py factor-py main log clean

PARCS = ${GOPATH}/src/github.com/lionell/parcs
GO = ${PARCS}/go
PY = ${PARCS}/py

all: main

parcs-go:
	docker build --no-cache -t lionell/parcs-go -f ${GO}/parcs/Dockerfile ${GOPATH}

parcs-py:
	docker build --no-cache -t lionell/parcs-py ${PY}/parcs

echo-go: parcs-go
	docker build --no-cache -t lionell/echo-go -f ${GO}/Dockerfile ${GO}/echo

echo-py: parcs-py
	docker build --no-cache -t lionell/echo-py -f ${PY}/Dockerfile ${PY}/echo

runner-go: parcs-go
	docker build --no-cache -t lionell/runner-go -f ${GO}/Dockerfile ${GO}/runner

runner-py: parcs-py
	docker build --no-cache -t lionell/runner-py -f ${PY}/Dockerfile ${PY}/runner

factor-py: parcs-py
	docker build --no-cache -t lionell/factor-py -f ${PY}/Dockerfile ${PY}/factor

main: runner-go factor-py
	docker service create --name runner --network parcs --restart-condition none --env LEADER_URL=tcp://10.0.98.232:4321 lionell/runner-go

log:
	docker service logs -f runner

clean:
	docker service ls | grep lionell | cut -f 1 -d ' ' | xargs docker service rm

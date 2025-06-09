hello:
	echo "Hello world"

run-dev:
	cd /home/nindgabeet/workspace/github.com/sohWenMing/gator/cmd/main_server && \
	ENVPATH="../../.env" go run . $(COMMAND)

run-pg:
	docker run \
	--name dev-postgres \
	-p 5432:5432 \
	-e POSTGRES_PASSWORD=postgres \
	-v pg_data:/var/lib/postgresql/data \
	-d \
	postgres:15.13

stop-pg:
	docker stop dev-postgres
	docker container prune -f


setup-pg:
	make run-pg
	./db_setup.sh
	make stop-pg
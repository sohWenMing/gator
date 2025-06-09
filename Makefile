run-dev:
	cd /home/nindgabeet/workspace/github.com/sohWenMing/gator/cmd/main_server && \
	ENVPATH="../../.env" go run . $(COMMAND)
# to run this, from root directory run make run-dev COMMAND="login nindgabeet"

run-pg-dev:
	docker run \
	--name dev-postgres \
	-p 5432:5432 \
	-e POSTGRES_PASSWORD=postgres \
	-v pg_data:/var/lib/postgresql/data \
	-d \
	postgres:15.13

stop-pg-dev:
	docker stop dev-postgres
	docker container prune -f

setup-pg-dev:
	make run-pg-dev
	./shell_scripts/db_setup.sh
	make stop-pg-dev


down-pg-dev:
	make run-pg-dev
	./shell_scripts/db_down.sh
	make stop-pg-dev

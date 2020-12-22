build_viz:
	docker build -t vis ./visualization

start_viz:
	-docker stop vis
	-docker rm vis
	docker run --env address=:8000 --publish 8000:8000 --detach --name vis vis:latest

start_containers:
	-docker-compose stop
	docker-compose.exe up --build --force-recreate -d

swap_docker:
	cp -f ./networkExamples/dockerfiles/$(name) ./networkExamples/docker-compose.yml

run:
	 $(MAKE) start_viz
	 $(MAKE) start_containers
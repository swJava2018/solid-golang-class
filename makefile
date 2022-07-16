.PHONEY: install-kafka-docker
.PHONEY: install-rabbitmq-docker
.PHONEY: install-kafka-k8s, install-prometheus-k8s, install-elasticsearch-k8s, install-kibana-k8s

install-kafka-docker:
	docker compose -f kafka/docker-compose.yaml up -d

generate-kafka-events:
	./kafka/gen-events.sh

install-rabbitmq-docker:
	docker compose -f rabbitmq/docker-compose.yaml up -d

install-kafka:
	kafka/deploy.sh

install-prometheus:
	prometheus/deploy.sh

install-elasticsearch:
	cd elasticsearch; ./deploy.sh

install-kibana:
	kibana/deploy.sh
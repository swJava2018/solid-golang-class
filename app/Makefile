.PHONEY: up, down, gen-event-kafka
.PHONEY: helm-install

up:
	docker compose -f docker-compose.yaml -f docker-compose-dependencies.yaml up -d
down:
	docker compose -f docker-compose.yaml -f docker-compose-dependencies.yaml down

gen-kafka-event:
	curl -s "https://api.mockaroo.com/api/d5a195e0?count=1000&key=ff7856d0" | kcat -b localhost:9092 -t purchases -P

kubectl-apply:
	kubectl apply -f manifests/configmap.yaml
	kubectl apply -f manifests/deployment.yaml

helm-install:
	helm upgrade --atomic --install event-data-pipeline helm/ -n default
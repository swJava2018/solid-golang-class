.PHONEY: install-kafka, install-prometheus, install-elasticsearch, install-kibana

install-kafka:
	kafka/deploy.sh

install-prometheus:
	prometheus/deploy.sh

install-elasticsearch:
	cd elasticsearch; ./deploy.sh

install-kibana:
	kibana/deploy.sh
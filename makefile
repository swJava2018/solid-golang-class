.PHONEY: install-kafka, install-prometheus, install-elasticsearch

install-kafka:
	kafka/deploy.sh

install-prometheus:
	prometheus/deploy.sh

install-elasticsearch:
	cd elasticsearch; ./deploy.sh

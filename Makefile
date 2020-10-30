.PHONY : test restart restartd kafka-setup s3-setup setup

test:
	docker exec -w /build isengard go test -v ./connectors ./tailing ./config -cover

restart:
	docker-compose rm -svf && docker-compose down && docker-compose build && docker-compose up

stop:
	docker-compose rm -svf && docker-compose down

restartd:
	docker-compose rm -svf && docker-compose down && docker-compose build && docker-compose up -d

kafka-setup:
	docker exec kafka-cluster kafka-topics --create --if-not-exists --zookeeper zookeeper:2181 --partitions 1 --replication-factor 1 --topic test-kafka-topic

s3-setup:
	docker exec localstack awslocal s3 mb s3://local-test-bucket

setup: kafka-setup s3-setup

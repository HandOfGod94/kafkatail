clean-dev-stack:
	docker compose stop
	docker compose rm -f -v
	docker volume rm -f kafkatail_kafka_data kafkatail_zookeeper_data

create-dev-stack:
	docker compose up -d

reset-dev-stack: clean-dev-stack create-dev-stack
	echo "reset dev stack complete"

cucumber-setup:
	ruby --version
	bundle install
	bundle binstub cucumber --path bin

cucumber-test: reset-dev-stack cucumber-setup
	bin/cucumber

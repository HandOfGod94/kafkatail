require "kafka"

kafka = nil

Given '{string} is present on kafka-broker {string} with {int} partition' do |topic, broker_addr, partition|
  kafka = Kafka.new([broker_addr], client_id: "kafkatail-integration-test")
  kafka.create_topic(topic, num_partitions: partition)
end

And '{string} message is pushed to {string} on partition {int}' do |message, topic, partition|
  kafka.deliver_message(message, topic: topic, partition: partition)
end

And 'starwars Human proto message is pushed to {string}' do |topic|
  human = Starwars::Human.new(homePlanet: "earth")
  seralized = Starwars::Human.encode(human)
  kafka.deliver_message(seralized, topic: topic)
end

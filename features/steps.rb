require "kafka"

kafka = nil

Given '{string} is present on kafka-broker {string}' do |topic, broker_addr|
  kafka = Kafka.new([broker_addr], client_id: "kafkatail-integration-test")
  kafka.create_topic(topic)
end

And '{string} message is pushed to {string}' do |message, topic|
  kafka.deliver_message(message, topic: topic)
end


require "aruba/cucumber"
require "aruba/api"
require "kafka"

After 'not @no-kafka' do
  kafka = Kafka.new(["localhost:9093"], client_id: "kafkatail-integration-test")
  kafka.delete_topic("test-topic")

  # we need to sleep for some seconds to ensure topic is deleted
  sleep(1.0)
end

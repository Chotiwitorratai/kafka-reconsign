# kafka-reconsign

cd docker-container
docker-compose -f compose-kafka-2181.yaml up -d

./kafka-topics.sh --bootstrap-server localhost:9092 --list

./kafka-topics.sh --create --topic consume-topic --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1

./kafka-console-producer.sh --broker-list localhost:9092 --topic consume-topic
# {"refId":"REF002"}

kafka-console-producer --bootstrap-server=localhost:9092 --topic=PaymentKNEXT   // for send payment messages 

kafka-console-producer --broker-list localhost:9092 --topic PaymentKNEXT

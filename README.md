# kafka-reconsign

cd docker-container
docker-compose -f compose-kafka-2181.yaml up -d

./kafka-topics.sh --bootstrap-server localhost:9092 --list

./kafka-topics.sh --create --topic consume-topic --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1

./kafka-console-producer.sh --broker-list localhost:9092 --topic consume-topic
# {"refId":"REF002"}

./kafka-console-producer --bootstrap-server=localhost:9092 --topic=PaymentKNEXT   // for send payment messages 

./kafka-console-producer --broker-list localhost:9092 --topic PaymentKNEXT

# {"TransactionRefID":"555test2022","Status":"success","PaymentInfoAmount":1900,"PaymentInfoWebAdditionalInfo":"testing","PartnerInfoName":"TIP","PartnerInfoDeeplinkUrl":"www.example.com","PaymentPlatform":"paotang"}

./kafka-console-producer --bootstrap-server=localhost:9092 --topic=InsuranceCallBack   // for send Insurance callback  

./kafka-console-producer --broker-list localhost:9092 --topic InsuranceCallBack

# {"RefID":"555test2022i","IdCard":"12555995336","PlanCode":"PL01","PlanName":"ประกันโควิด","EffectiveDate":"2022-01-25T12:11:56Z","ExpireDate":"2022-01-25T12:11:56Z","IssueDate":"2022-01-25T12:11:56Z","InsuranceStatus":"success","TotalSumInsured":10000,"ProductOwner":"POtest","PlanType":"Base plan"}


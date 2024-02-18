# go-kafka-notify

https://www.freecodecamp.org/news/build-a-real-time-notification-system-with-go-and-kafka/

curl -X POST http://localhost:8080/send -d "fromID=2&toID=1&message=Raju kutreya."
curl -X POST http://localhost:8080/send -d "fromID=1&toID=2&message=tu banduk kahan se layega Raju"
curl -X POST http://localhost:8080/send -d "fromID=2&toID=3&message=mera paisa kidhar hai"

curl http://localhost:8082/notifications/1
{"notifications":[{"from":{"id":2,"name":"Raju"},"to":{"id":1,"name":"Shyam"},"message":"Raju kutreya."}]}

curl http://localhost:8082/notifications/2
{"notifications":[{"from":{"id":1,"name":"Shyam"},"to":{"id":2,"name":"Raju"},"message":"tu banduk kahan se layega Raju"}]}

curl http://localhost:8082/notifications/3
{"notifications":[{"from":{"id":2,"name":"Raju"},"to":{"id":3,"name":"Anuradha"},"message":"mera paisa kidhar hai"}]}
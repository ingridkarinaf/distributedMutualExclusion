# distributedMutualExclusion
Distributed Systems: Mandatory Hand-in 4

To run the project, open 3 different terminals and run the following commands in each terminal:
`go run main.go 0`
`go run main.go 1`
`go run main.go 2`

Changes
-------

- Return statement inside of the request queue loop
After realising it is in fact possible for multiple nodes to access the ciritcal section as the same time, we identified the issue to be that gRPC calls are processed concurrently (i.e. as threads), meaning that we had to implement a lock in the "AccessRequest" method, ensuring that the different nodes could not execute the method at the same time. We implemented the lock using channels, in which each node will have to wait to receive a value from the channel before it is able to proceed into the AccessRequest gRPC method, in which each node accesses the critical section. 


-------

This will create 3 separate severs running on ports 5001, 5002 and 5003. 

The car represents our critical section. To enter the critical section, you press enter in the terminal (peer) who you want to access the critical section with. 

We are using the logic from Ricart & Agrawala algorithm to make only one peer to access the critical section (car) at once. 
(1) 5001 wants to access the critical section so it sends a request to 5002 and 5003. 5001 stays and waits for the replies before continuing 
(2) neither 5002 nor 5003 wants to access the car so they will reply 5001 immediately and 5001 can drive the car
(3) Simultaneously peer 5001 and 5002 are trying to access the critical section. 5001 will access the critical section first since their id is lower. When 5002 request arrives to 5001, 5002 will put 5001 into a queue to waits its turn until 5001 releases the critical section (car).
(4) When 5001 is releasing the car it iterates through requestQueue and sends the reply for 5002, who has been waiting for the reply from 5001. 
(5) Now when 5002 got reply from both 5001 and 5003, it can access the critical section.

We have attached screenshots of our logs to the Pictures folder, where we demonstrate the scenario described above. Furthermore, we also drew the scenario and this can be found from the pictures folder.



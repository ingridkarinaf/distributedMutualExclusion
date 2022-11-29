# distributedMutualExclusion
Distributed Systems: Mandatory Hand-in 4

To run the project, open 3 different terminals and run the following commands in each terminal:
`go run main.go 0`
`go run main.go 1`
`go run main.go 2`

Changes
-------



The problem with our initial solution was that we didn't take into account the scenario, where a client with a lower port number would had already given a permission to a client with higher port number. The problem here would be when the client with the lower port number would also want to access the critical section while both of them would be waiting a response from the third client. We modeled the problem and the solution in the following graph: https://github.com/ingridkarinaf/distributedMutualExclusion/blob/main/Pictures/mandatory4scenario.drawio.png

From the graph with previous implementation with our wrong logic 5001 would start being in the critical section (t1). 5003 would requeast access first to the critical section and get response from both peers (t2). Right after 5002 would request access and 5003 would give a permission to 5002 since the id is smaller (t3). Next, when 5001 would relaease the critical section it would give the permission simultaneously to both 5002 and 5003 (t4), and since both of them had also approved each others they would both access the critical section at the sametime. 

To solve this we needed to find a way to know if a peer had already given a response to another peer while waiting the critical section to become free. Before we always gave response for a peer if their id was smaller. Now we will also be checking the condition if a peer requesting access has already given us a response. In such situation we would queue this client to the request queue to be allowed to drive the car after. The correct scenario is showed in the graph as "new implementation" where 5003 will access the critical section before 5002, since 5002 had already approved 5003 earlier.

Another problem that we found from our program was that we had a return statment inside of the for loop which job was to iterate over the list of peers who were set into the request queue. Having the return statment in the queue caused that only one peer was replied with a permission. 

Furthermore, we identified one more issue to be gRPC calls being processed concurrently (i.e. as threads), meaning that we had to implement a lock in the "AccessRequest" method, ensuring that the different nodes could not execute the method at the same time. We implemented the lock using channels, in which each node will have to wait to receive a value from the channel before it is able to proceed into the AccessRequest gRPC method, in which each node accesses the critical section. 


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



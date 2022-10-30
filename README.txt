Use gRPC for all messages passing between nodes - check

Use Golang to implement the service and clients - check

Every client has to be deployed as a separate process - check

Log all service calls (Publish, Broadcast, ...) using the log package - check 

Demonstrate that the system can be started with at least 3 client nodes 
Demonstrate that a client node can join the system - check
Demonstrate that a client node can leave the system - check

Optional: All elements of the Chitty-Chat service are deployed as Docker containers - check







R1: Chitty-Chat is a distributed service, that enables its clients to chat. The service is using gRPC for communication. You have to design the API, including gRPC methods and data types.  Discuss, whether you are going to use server-side streaming, client-side streaming, or bidirectional streaming? - check

R2: Clients in Chitty-Chat can Publish a valid chat message at any time they wish.  A valid message is a string of UTF-8 encoded text with a maximum length of 128 characters. A client publishes a message by making a gRPC call to Chitty-Chat. - check

R3: The Chitty-Chat service has to broadcast every published message, together with the current Lamport timestamp, to all participants in the system, by using gRPC. It is an implementation decision left to the students, whether a Vector Clock or a Lamport timestamp is sent.
R4: When a client receives a broadcasted message, it has to write the message and the current Lamport timestamp to the log
R5: Chat clients can join at any time. - check
R6: A "Participant X  joined Chitty-Chat at Lamport time L" message is broadcast to all Participants when client X joins, including the new Participant.
R7: Chat clients can drop out at any time. - check
R8: A "Participant X left Chitty-Chat at Lamport time L" message is broadcast to all remaining Participants when Participant X leaves.
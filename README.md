 ACN Project Report
Team - TechGeeks (ACN08)
 Group Members : 
 Himanshu Gupta (2020H1030121H) 
 Somdatta Sen (2020H1030140H)
 Surbhi Sharma (2020H1030148H)
 Aman Srivastav (2020H1030137H)
1
Summary

 Quick UDP Internet Connection (QUIC) is a recent protocol initiated by Google that combines the functions of HTTP/2, TLS, and TCP directly over UDP, with the goal to reduce the latency
 client-server communication. MP-QUIC is a QUIC extension that enables a QUIC connection to use different paths such as WiFi and LTE on smartphones, or IPv4 and IPv6 on dual-stack hosts.
We have implemented MP-QUIC as an extension of the quic-go implementation as it maintains MP-TCP's benefits (like aggregation benefit, network handover).
We took help of online references available for implementation of “quic-go”.
For file transfer, in our implementation, on the basis of file size, the algorithm decides whether to send the desired file via single path (QUIC) or multipath (MP-QUIC). The file size (or threshold size) is set to 5KB. So, for a file of size greater than 5KB, the transfer takes place using MP-QUIC , else, it uses a single path (QUIC).
We focused on Packet scheduling in MP-QUIC. We were successful in analysing the current scheduling scheme used and its implementation. The research paper that we chose is purely based on Priority based stream scheduling. We noticed that scheduling without the recognition of the stream features can aggravate inter-stream blocking when paths are being shared.
We observed that scheduling without the recognition of the stream features can aggravate inter-stream blocking when sharing paths. So to solve this problem we had allocated paths for each stream by considering the match of stream and path features in the scheduling process (usually done by scheduler).
We proposed a priority-based stream scheduling mechanism for MP-QUIC. PStream provides stream prioritization to prevent time-critical streams from being blocked by non-critical streams.
The Scheduling in MP-QUIC is done on the basis of Path features such as Bandwidth,rtt and delay etc. So our algorithm is based on network assisted packet scheduling.As The client initiates the connection with the server the sender will also send the RM cells which is interspersed with the data cells .With these information such as bandwidth,Feedback from RM cells ,time delay the scheduler then schedules the packets according to the priority assigned. This type of network assisted scheduling algorithm will helps us to reduce the latency and achieve Maximum throughput.
Step-1. Implementation of MP-QUIC
2

 The topology used in our implementation is as below:
 We are considering two paths P1 and P3 which ensures the working of MP-QUIC based on the threshold file size (set to 5KB).
➢ Requirements
1. Install the latest copy of Oracle VM VirtualBox with OS as ubuntu .
2. Setup mininet, miniedit and installed quic-go setup using the references below -
https://multipath-quic.org/ https://multipath-quic.org/conext17-deconinck.pdf https://github.com/lucas-clemente/quic-go
➢ Creating Structures for file transfer using MP-QUIC
1. The client server connection is created as we run the setup-topology.py file.
2. Once the connection is established, we need to create the executable binary files for the client and
server by running “client-multipath.go” and “servermultipath.go” in their respective xterm windows.
3. While starting the server, we need to provide the client storage location (Storage-client) as a parameter, whereas, while starting client, we need to mention the storage location of server (Storage-server) along
with the file name that the client wants to request and the IP address of server.
    3

  ➢ Server Implementation
1. The server binds itself to port number (4242).
2. Session is created with a Session Remote address.
3. Based on the size of the file, stream(s) is/are created which acts as a burst packet that needs to be sent
in parallel/sequentially on the paths of the MP-QUIC network created.
4. The server file is compiled using “go build server-multipath.go”.
➢ Client Implementation
1. Import all the necessary files and create paths for the client and bind it with the server address and the
port number using which the connectivity between client and server will be established using a three
way handshake protocol.
2. A threshold (5KB) is maintained and mode of transmission is decided on the following condition:
● If file size < threshold, then → use single path transmission (QUIC)
● If file size > threshold, then → use multipath transmission (MP-QUIC)
3. Session and streams are maintained in a similar manner as server implementation.
4. The Client file is compiled using command “go build client-multipath.go” which creates a binary
executable file for client.
➢ Command to execute Client and Server
● ./server-multipath storage-client
● ./client-multipath storage-server/{filename-to-request} 100.0.0.1
➢ Observations and Result
4

       S.No.
Parameters
 File Size
 Transfer Time
 1
Rate = 5kbit Burst = 64kbit Latency = 200ms
10.8MB
14mins 12secs
2
Rate = 10kbit Burst = 256kbit Latency = 500ms
10.8MB
5mins 50secs
3
When 1 interface (client-eth1) was turned down using command “ifconfig client-eth1 down”
10.8MB
7mins 45secs
      5

 Step-2. Research Paper Study and Implementation ➢ Background
We observed that scheduling without the recognition of the stream features can aggravate inter-stream blocking when sharing paths. This gap is filled by a priority-based online stream scheduling mechanism used in MP-QUIC, which performs path scheduling based on the stream features.
➢ Findings from Research Paper
1. Usually packets compete to get scheduled and to get on a path with higher bandwidth or low
congestion.
2. Instead of this, we allocate paths for each stream by considering the match of stream and path
features in the scheduling process (usually done by scheduler).
3. The streams can utilize the allocated paths concurrently from the beginning, thus reducing the
burst transmission on the fast path.
4. Streams can be prioritized by marking them as dependent(assigning priority value) and other path
features(bandwidth, RTT, Completion time, etc.) and then the scheduler allocates the new stream
to each path with a calculated amount of data.
5. To further optimize the completion time of the time-critical stream, it schedules a stream to a single
path.
6. Each dependent stream is allocated a priority value between 1 and 255 (inclusive) and a higher
value means a higher priority.
7. The scheduler allocates the new stream to each path with a calculated amount of data. After the
scheduling process, the stream is added to the corresponding paths for sending.
 6

 ➢ Algorithm - Priority Stream
7
1: 2: 3: 4: 5: 6: 7: 8: 9: 10: 11: 12: 13: 14: 15: 16: 17: 18: 19: 20: 21:
procedure SCHEDULESTREAM
sortedList ← sort paths by ascending order of estimated one way delay of path i initialize di = 0, 1, 2 [1,2,....,n]
D’ ← D ,D is the total data volume of stream, D’ is a temporary variable
for path i in sortedList[1,2,..., n-1] do. Fill the gap
Ogap = Oi+1- Oi (find the gap between each consecutive paths)
!= 0 and D’ > 0 then
if D’ >= Ogap x 𝛴k=1 to i bsk then
for j ←1 to i do
inc ← Ogap x bsj ,where bsj is the bandwidth share of stream j.
dj ← dj + inc where dj is the fraction of data assigned to path j. D’ ← D’- inc
end for
else
for j ← 1 to i do
inc ← D’ x (bsj / 𝛴k=1 to i bsk )
dj ← dj + inc end for
D’← 0
end if end if
if Ogap

 22: 23: 24: 25: 26: 27: 28:
end for
if D’ > 0 then share the rest volume proportionally
for path i in sortedList do
di ← di + D’ x (bsi / 𝛴k=1 to i bsk ) , di is the fraction of data assigned to path i
end for end if
end procedure
  8

 Step-3. Enhancement(s) and Innovation Network Assisted Path scheduling
➢ Issues (as per our observation)
1. The scheduling algorithm which is used considers the path features, RTT and many other things, on
the sender side.
2. This type of scheduling doesn’t take care about the network's delay or network path or network
features which may cause issues in the coming paths or deviations.
3. These features should also be considered when selecting the best path to send data, so that packets
can flow through fast paths and provide low latency.
4. However, if we take network features into consideration while performing packet scheduling, we can
maximize the transmission and can reduce the packet loss and packet retransmission rate to great extent, as network features play an important role in Packet Scheduling.
➢ Solution Approach (by Team - TechGeeks)
1. During the initiation of the connection, we can send fixed length cells when doing handshake.
2. These RM (Resource Management) cells are interspersed with the data cells. These cells will be
passed along with client hello messages during the initiation of the connection.
3. These cells will be explicitly sent back by the receiver .Each data cell contains an EFCI (Explicit
Forward Congestion Indication) bit. A congested network switch can set the EFCI bit in a data cell to
1 when signal congestion to the destination host.
4. Using the EFCI in data cells and the CI bit in RM cells, a sender can thus be notified about congestion
at paths.
5. A switch can also set the NI (No Increase) bit in a passing RM cell to 1 under mild congestion and
can set the CI(Congestion Indication) bit to 1 under severe congestion conditions.
6. When a destination host receives an RM cell, it will send the RM cell back to the sender with its CI
and NI bits intact.
7. With all the information indicated by bit indicators(such as path delay,path congestion,bandwidth
etc), packets can be scheduled on the path.
8. The data cells will be triggered after every constant time so that the scheduler can schedule the
data in the best path.
 9

 ➢Algorithm:
1. Procedure SCHEDULE:
2. During the Connection Initiation:
a. Estimate the one way delay of each path.
b. Send the data cells(Resource Management) in each path.
3. Get the Feedback from the Data cells.
4. While(Data is not send )
5. Begin:
6. Sort the path[1,2,....,n-1] in the increasing order of time delay.
7. Assign the priority to each path based on:
a. Feedback from RM cells
b. Time delay
c. Bandwidth
d. Completion time
8. Send the Packets in the increasing order of priority .
9. Update the priority of each path by sending again RM cells after some constant time.
10. End.
10

 References
https://multipath-quic.org/
https://multipath-quic.org/conext17-deconinck.pdf https://github.com/lucas-clemente/quic-go
PStream: Priority-Based Stream Scheduling for Heterogeneous Paths in Multipath-QUIC
    11


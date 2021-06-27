# Filetransfer-using-MPQUIC-protocol
Setup MP-QUIC structure. Create a topology as below :3. Create the files client-multipath.go and server-multipath.go 4. The files are stored in location “storage-server” and are accessible by server. 5. To execute the code, follow the below steps : Run the python file exported from topology using command : - sudo python2 setup-topology.py Once the terminal enters mininet, open xterm of server and client : - xterm server client In xterm of server, run the below commands : - go build server-multipath.go - ./server-multipath storage-client Here, storage-client is the destination location where client stores its files. In xterm of client, run the below commands : - go build client-multipath.go - ./client-multipath storage-server/abc.txt 100.0.0.1 Here, 100.0.0.1 is the IP address of node - server and abc.txt is the file requested by client which is stored in location “storage-server”...Priority-Based Stream SchedulingIn MP-QUIC, when multiple streams share common path, there are chances of Inter-stream blocking, which may have severe consequences. Also, stream features such as Bandwidth,RTT Delay etc are not taken into consideration when streams are directed on paths in network. In Priority Based Scheduling, instead of making all streams competing for the fast path in a greedy fashion, we allocate paths for each stream by considering the match of stream and path features in the scheduling process. In this, streams can be prioritized by giving them priority value based on path features(bandwidth, RTT, Completion time, etc.) and then the scheduler allocates the new stream to each path with a calculated amount of data.  This type of scheduling reduces the burst transmission of packets in congested path or paths with having low delay      .Implementation To implement Priority based scheduling, we have assigned a priority to each stream that is created for file transfer / transfer of all the packets created. The streams are scheduled on the basis of decreasing priority of streams on a common path. For example -  If 3 streams are created, and we have 2 paths available in network, then, 2 streams will be redirected to 2 paths available, and 3rd stream will be sharing the path with either Stream1 / Stream2. In this case, the stream scheduling is done on the basis of priority. The probability of a stream being selected is calculated by dividing its priority by the priority sum of all the scheduled streams on this path.  Network Assisted Path schedulingThe scheduling algorithm we have now, considers the path features, RTT and many things, on the sender side. During the initiation of connection between client and server, the RTT delay is taken and fixed length cells are sent along with handshake. When a  destination host receives an RM cell, it will send the RM cell back to the sender with its CI and NI bits intact. With all the information indicated by bit indicators(such as path delay,path congestion,bandwidth etc), packets are scheduled on the path  The data cells will be triggered after every constant time so that scheduler can schedule the data in the best path.
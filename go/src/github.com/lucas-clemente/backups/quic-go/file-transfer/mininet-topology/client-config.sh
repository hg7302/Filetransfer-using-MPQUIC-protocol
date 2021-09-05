#!/bin/bash

ifconfig client-eth0 10.0.0.1
ifconfig client-eth1 11.0.0.1

ip route add 10.0.0.0/8 dev client-eth0 src 10.0.0.1 table rt_client_eth0
ip route add default via 10.0.0.2 dev client-eth0 table rt_client_eth0

ip rule add from 10.0.0.1/32 table rt_client_eth0
ip rule add to 10.0.0.1/32 table rt_client_eth0


ip route add 11.0.0.0/8 dev client-eth1 src 11.0.0.1 table rt_client_eth1
ip route add default via 11.0.0.2 dev client-eth1 table rt_client_eth1

ip rule add from 11.0.0.1/32 table rt_client_eth1
ip rule add to 11.0.0.1/32 table rt_client_eth1

route add default gw 10.0.0.2 client-eth0
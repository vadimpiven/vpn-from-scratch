# Network connection basics, parts 0-1 (2019/07/11)

## OSI model

### Definition
The Open Systems Interconnection model is standardizes the communication functions of a telecommunication or computing system without regard to its underlying internal structure and technology.

### Layers
Each layer of OSI model is represented by different protocols. Each protocol could interact only on layer it belongs to, with one layer up and one layer down. Each layer taggs the data with headers that will help to transform data into its original look on the other end of connection.
1. Physical layer (USB, RS-232) – interacts with binary data transmitted through the physical cable
2. Data link layer (ARP, IEEE 802.3) – performs physical addressing and data integrity checks
3. Network layer (IPv4, IPv6, ICMP) – performs logical addressing, interacts with queues of data frames, performs translation of logical addresses and names in physical addresses (DNS functionality)
4. Transport layer (TCP, UDP) – performs data segmentation, transferring variable-length data sequences and error control while maintaining the Quality of Service functionality
5. Session layer (Sockets, L2TP, PPTP) – manages the creation and closing of session, exchange of data and maintaining the session during periods of inactivity
6. Presentation layer (ASCII, SSL, TLS, MIME) – performs data encoding and decoding, enctyption and decryption, compresing and decompressing, convertion between different protocols
7. Application layer (HTTP, FTP, SMTP, TELNET) – provides the ability to interact with network for userland application

### TCP/IP
The Internet protocol suite is the conceptual model and set of communications protocols used in the Internet and similar computer networks – HTTP on application layer, TCP on transport layer, IPv4 and ARP on network layer and Ethernet on data link layer are the most well-known examples. The three top layers in the OSI model, i.e. the application layer, the presentation layer and the session layer, are not distinguished separately in the TCP/IP model which only has an application layer above the transport layer.

## IP address

### Definition
An Internet Protocol address is a numerical label assigned to each device connected to a computer network that uses the Internet Protocol for communication. This label allows to uniqly identify each node (host) in the composite network.

### Address types
- local/physical address (layer 1) – MAC (Media Access Control) address, unique for each physical device, distributed centrally
- network/IP address (layer 2) – 4 bytes written as decimal and separated by dots, identifies the unique connection, assigned independently from MAC
- domain name (layer 3) – symbolic name uniquely corresponding to the certain IP address, translation is performed by DNS (Domain Name System)

### Reserved IP addresses
IP address consists of subnet address and host address, range of host addresses available is defined by mask (for example 255.255.255.0 mask means that for IP address 192.168.1.1 subnet address will be 192.168.1.0 and host address will be 0.0.0.1). Subnet address is reserved and couldn't be assigned to any node in this subnet. Host address with all ones in binary representation is also reserved, it's called a broadcast address. All packages sent to broadcast address will be delivered to every node in the subnet.

## VLAN

### Definition
Virtual Local Area Network is a group of devices which can communicate with each other on data link layer regardless of their physical location. In the same time, devices located in different VLANs can interact with each other only on network layer even if they are physically connected to the same switch.

### Communication
By default, each device connected to swith is identified by switch port it's connected to, MAC address and VLAN number. In order to exchange information between devices located in the same VLAN but physically connected to different switches, it is necessary that there exists a channel between these switches also related to this particular VLAN. Such cnannel is called untagged. When there are two VLANs both represented on two different switches – two untagged channels are required to make data exchange possible. The same time, it's possible to configure a tagged channel called 'trunk'. When data is passed through a trunk it is labeled with tag, marking to which VLAN the passed data relates. When a tagged data is passed through an untagged channel or when a data passed is intended to be delivered to the same swith port it came from - frame will be dropped by switch.

### Transferring tagged data over the Internet
By default tagged frames wouldn't be transfered over the internet. The easiest way to save tags is to transfer data over the VPN bridge (VTun or OpenVPN).

## Sources
- [Сетевая модель OSI](https://ru.wikipedia.org/wiki/Сетевая_модель_OSI)
- [TCP/IP](https://ru.wikipedia.org/wiki/TCP/IP)
- [IP-адрес](http://xgu.ru/wiki/IP-адрес)
- [IP address](https://en.wikipedia.org/wiki/IP_address)
- [VLAN](http://xgu.ru/wiki/VLAN)
- [OpenVPN Bridge](http://xgu.ru/wiki/OpenVPN_Bridge)
- [Сети для самых маленьких. Часть нулевая. Планирование](https://linkmeup.ru/blog/11.html)
- [Сети для самых маленьких. Часть первая (которая после нулевой). Подключение к оборудованию cisco](https://linkmeup.ru/blog/12.html)
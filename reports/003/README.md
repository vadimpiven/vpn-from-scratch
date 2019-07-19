# Network connection basics, parts 2-3 (2019/07/16)

## Public-key cryptography

### Concept
There are three basic tenets of secure transactions:
- data encryption – makes reading the message impossible for the third party
- integrity check – a way to detect whether the message was tampered
- user authentication – identifies the sender and the recipient

### Definition
Public-key cryptography is an asymmetric encryption method. It is assumed that each data exchange participant has two keys - public and private. Data encrypted with public key could be decrypted by corresponding private key and vice versa. Each of the data exchange participants knows the public keys of all the others, private keys are kept in secret.

### Method
1. Sender runs the message through a message-digest function and calculates a checksum (integrity check).
2. This checksum is encrypted with a sender's privete key and attached to the message (sender authentication).
3. The message with attached checksum is encrypted by recipient's public key (recipient authentication, data encryption).
4. Recipient decrypts the message using its private key and runs it through the same message-digest function, then compares the result with checksum, decrypted using sender's public key. If they match then the message wasn't tampted.

## SSL/TLS

### Authentication
A server authenticates itself to a client by sending an unencrypted ASCII-based digital sertificate (SSL sertificate) containing information about the company operating the server, including the server's public key. This sertificate is signed by a trusted issuer, which means that the issuer has investigated the company operating the server and believes it to be legitimate. If the client trusts the issuer, it can trust the server. The issuer signs the sertificate by generating check-sum for it and encrypting this checksum with issuer's private key. If the client trusts the issuer, then it already knows the issuer's public key.

### Encryption
In order for HTTPS encrypted connection to be established, sercer and client have to define a secret encryption key (HTTPS uses symmetric encryption instead of assymmetric as it's performed faster, so only one key is needed). There are two procedures that allow secure exchange: RSA in TLS 1.2 and DH (DH, DHE, ECDH, ECDHE) in TLS 1.3. In both cases, the client and server randomly select two large prime numbers and exchange them in the open. These two primes (![x](https://latex.codecogs.com/png.latex?x) and ![y](https://latex.codecogs.com/png.latex?y)), as well as one additional number ![z](https://latex.codecogs.com/png.latex?z), will then be used to calculate the symmetric encryption key. If RSA is used, client generates so-called pre-master secret ![z](https://latex.codecogs.com/png.latex?z), encrypts it with server's public key, then server receives it and decrypts using server's private key. If DH is used, client chooses pre-master secret ![a](https://latex.codecogs.com/png.latex?a) and server chooses pre-master secret ![b](https://latex.codecogs.com/png.latex?b), then client calculates ![A = x^a \mod y](https://latex.codecogs.com/png.latex?A%20%3D%20x%5Ea%20%5Cmod%20y), server calculates ![B = x^b \mod y](https://latex.codecogs.com/png.latex?B%20%3D%20x%5Eb%20%5Cmod%20y) and they exchange these numbers. Then client calculates ![z = B^a \mod y](https://latex.codecogs.com/png.latex?z%20%3D%20B%5Ea%20%5Cmod%20y), server calculates ![z = A^b mod y](https://latex.codecogs.com/png.latex?z%20%3D%20A%5Eb%20mod%20y) (the resulting ![z](https://latex.codecogs.com/png.latex?z) will match). Nowadays DH is preferred over RSA because of Perfect Forward Secrecy (PFS). This means that when DH key exchange is compromised, hacker can only decrypt one session whereas if RSA is used and server's secret key is compromised – all user sessions could be decrypted.

## VPN

### Definition
Virtual Private Network (VPN) allows user to connect to the internet through a different computer running VPN server. This means that all transfered data is encrypted and readable only by sender and server (noone can steal your password while using public wifi), sender and server see each other as if situated in same local network, sender is identified by websites with server's IP address and physical location, local Internet Service Provider (ISP) could dump but couldn't read the traffic, access to the websites couldn't be blocked.

### Placement
- browser extension – only browser traffic is encrypted
- standalone – every application on the computer connects to network through the VPN
- router VPN - traffic from all devices in local network is encrypted

### Protocols
- Point-to-Point Tunneling Protocol (PPTP) – the oldest one, encryption is weak
- Layer 2 Tunneling Protocol (L2TP/IPSec) – creates a tunnel but uses no encryption, that's why commonly used with IPSec (IP security – suite of protocols between 2 communication points across the IP network that provide data authentication, integrity, and confidentiality)
- Secure Socket Tunneling Protocol (SSTP) – newer version of PPTP, equivalent to HTTPS encryption
- Internet Key Exchange, version 2 (IKEv2) – newer version of L2TP, also should be used with IPSec
- OpenVPN – the open source technology, most popular

## IP addresses
IP addresses are divided into private and public. Public IP addresses are distributed by IANA (Internet Assigned Numbers Authority), private IP addresses can be used without any control. Private IP ranges are:
- 10.0.0.0-10.255.255.255
- 172.16.0.0-172.31.255.255
- 192.168.0.0-192.168.255.255
All packages addressed to private IP will be discarded by ISP.

## MAC address
MAC address FF:FF:FF:FF:FF:FF means that frame is addressed to all devices in local VLAN. This is used by ARP requests. If recipient is located not in local VLAN, then MAC address of switch is used instead of recipient's MAC address and frame is sent to 'default gateway'.

## Sources
- [Что такое TLS-рукопожатие и как оно устроено](https://tproger.ru/articles/tls-handshake-explained/)
- [How does VPN work?](https://www.namecheap.com/vpn/how-does-vpn-virtual-private-network-work/)
- [Сети для самых маленьких. Часть вторая. Коммутация](https://linkmeup.ru/blog/13.html)
- [Сети для самых маленьких. Часть третья. Статическая маршрутизация
](https://linkmeup.ru/blog/14.html)
- [Еще раз про IP-адреса, маски подсетей и вообще](https://habr.com/ru/post/129664/)
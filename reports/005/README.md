# Network connection basics, parts 6-7 (2019/07/18)

## IPsec

### Definition
Internet Protocol Security (IPsec) is a secure network protocol suite that authenticates and encrypts the packets of data sent over an Internet Protocol network. It is used in virtual private networks (VPNs). IPsec can protect data flows between a pair of hosts (host-to-host), between a pair of security gateways (network-to-network), or between a security gateway and a host (network-to-host).

### Protocols
IPsec is a network layer set of protocols:
- ESP (Encapsulating Security Payload) - encrypts and encapsulates data
- AH (Authentication Header) - authentification, integrity checks
- IKE (Internet Key Exchange protocol) - is used to set up IPSec SA (Security Association, params of given network tunnel)

### Specifications
The first stage of initiating IPsec tunnel requires to start up ISAKMP tunnel which is used to exchange internal service data (encryption method, keys and so on). This is done by IKE protocol which exists in two versions:
- IKEv1 - an old one, specifies two ways of function: main mode (default, slow) and aggressive mode (two times faster, has security leaks)
- IKEv2 - the new one, recommended, uses DH procedure
When two sides exchanged all required parameters the main tunnel starts up; ISAKMP tunnel is not removed as it will be used to renew encryption keys each 4608000kb/3600s (by default).

### Modes
IPSec modes are closely related to the function of AH and ESP protocols. Both of these protocols provide protection by adding to a datagram a header (and possibly other fields) containing security information. The choice of mode does not affect the method by which each generates its header, but rather, changes what specific parts of the IP datagram are protected and how the headers are arranged to accomplish this.
- transport mode - only the payload of the IP packet is encrypted or authenticated, the IP header is neither modified nor encrypted
- tunnel mode - the entire IP packet is encrypted and authenticated (used to create VPN)

## Git - preparations for code review

### Create and switch to NewFeature branch
~~~
git checkout -b NewFeature
~~~

### Confirm changes
~~~
git status
git diff src/controllers/v1/comments.js
git add src/controllers/v1/comments.js
~~~

### Commit & push
~~~
git commit
git push origin NewFeature
~~~

### Generate definition for pull-request (squash and merge)
~~~
git log --pretty='%h: %B' --first-parent --no-merges --reverse
~~~

### Unite small steps in a whole story (manually, pick->squash)
~~~
git rebase --interactive master
git push origin NewFeature --force
~~~

### Delete NewFeature branch after merge
~~~
git checkout master
git pull origin master
git branch -D NewFeature
~~~

## Sources
- [Сети для самых маленьких. Часть шестая. Динамическая маршрутизация](https://linkmeup.ru/blog/33.html)
- [Сети для самых маленьких. Часть седьмая. VPN](https://linkmeup.ru/blog/50.html)
- [IPsec](https://en.wikipedia.org/wiki/IPsec)
- [Анатомия IPsec. Проверяем на прочность легендарный протокол](https://habr.com/ru/company/xakep/blog/256659/)
- [IPSec Modes: Transport and Tunnel](http://www.tcpipguide.com/free/t_IPSecModesTransportandTunnel.htm)
- [Рецепт полезного код-ревью от разработчика из Яндекса](https://m.habr.com/ru/company/yandex/blog/422143/)
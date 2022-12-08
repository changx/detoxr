# DNS Detoxifier

利用 GFW 投毒抢答的工作方式，自动判别域名是否被投毒。

被投毒的域名，通过 DNS over HTTPS 解析，安全域名用指定的 name server 解析。

``` udp/1053 ```  name server port 

``` tcp/3000 ``` web admin http port

作为 dnsmasq 上游的 forwarder ns，在 dnsmasq.conf 里配置

```server=/#/127.0.0.1#1053```

Web 管理界面访问方法：

```http://ip-of-the-router-detoxr-running-on:3000/```

<img width="1336" alt="Screenshot 2022-12-08 at 16 43 40" src="https://user-images.githubusercontent.com/620665/206399985-fd008737-983b-469f-bd63-d12894011c13.png">

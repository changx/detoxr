# DNS Detoxifier

利用 GFW 投毒抢答的工作方式，自动判别域名是否被投毒。

被投毒的域名，通过 DNS over HTTPS 解析，安全域名用指定的 name server 解析。

``` udp/1053 ```  name server port 

``` tcp/3000 ``` web admin http port

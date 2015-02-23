vulcand-elb-rewrite-ssl
============

```
Client ---> ELB:80  (X-Forwarded-Port:80)  ---> Vulcand:80 (rewrite SSL)
Client ---> ELB:443 (X-Forwarded-Port:443) ---> Vulcand:80 (OK)
```

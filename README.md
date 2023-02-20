
## The application

This code is a DNS proxy that listens on both TCP and UDP ports and forwards DNS requests to a remote DNS server. The proxy supports both TCP and UDP protocols and uses a secure connection (TLS) to communicate with the remote DNS server when using the TCP protocol.

### Running the application 
``` 
# keep in mind that since the application runs on port 53, it need root privilege to run any process in this port.
# Also, the port 53 normally is already in use, so you need to stop the resolve service .to be able to run the application. 
# To stop resolve service (ubuntu): sudo systemctl stop systemd-resolved.service.

sudo go run go/main.go

# Don't forget to restart the resolve services (ubuntu): sudo systemctl restart systemd-resolved.service
``` 

### Running the application using docker
```
docker build -t proxy:v0.1 .
docker run --rm -d --name proxy -p 5333:53/tcp -p 5333:53/udp proxy:v0.1
```

### Testing name resolution
``` 
# tcp 
dig -p 5333 +tcp +short @localhost google.com

# udp
dig -p 5333 +tcp +short @localhost google.com
``` 

## Security 
* Access to the proxy should be restricted to authorized personnel only.
* The proxy should be configured to log all relevant events and activities, such as incoming and outgoing traffic, errors, and exceptions.
* The code should be reviewed for security vulnerabilities, and regular security updates should be applied.
* The communication is only encrypted after the proxy connection, all communication before this point is still in plain text.

## Running in a microservice environment 
Assuming we want to deploy this proxy in a Kubernetes cluster, there are a few different strategies that we can follow:

* Deploy the proxy (as a daemonset) in the cluster and create a service that forwards DNS messages to it.
* Deploy as a sidecar with specific pods that need this solution.

## Improvements  
* Logging and monitoring
* Caching
* Unit tests
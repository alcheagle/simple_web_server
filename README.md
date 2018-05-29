# Description

This is an example deployment of a web application for k8s. The web application
answers on "0.0.0.0:8080/" path with the hostname of the machine on which it is running.

The k8s deployment creates:
+ a replicaSet of 4 pods, by default
+ a service listening on port 8080
+ an Ingress resource answering on path
"INGRESS_CONTROLLER_ADDRESS:INGRESS_CONTROLLER_PORT/simple_web_server"

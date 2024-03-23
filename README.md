# IoTA-Admin
The Fiware IoT-Agent Admin will be a service to manage IoT-Agents with this [IoT Agent API - Fiware-iotagent-node-lib](https://iotagent-node-lib.readthedocs.io/en/3.3.0/api.html).

As it's right now, creating, deleting and listing for devices as well as config groups works. 

There is also support for managing multiple iot-agents.

## Next up
Adding simple authentication via username and password, as well as protecting "public" routes with brute force protection and rate limiting.

## Features tha will come
- LDAP authentication
  
## Features that may come
- adding Orion-Context Broker sdk
- live monitoring
- mangage Orion-Context Broker
- managing and deploying other IoT Agents


## Why is the release version so high?
I made a mistake and released a v1.4.x which lead to go thinking an old version being the latest version.

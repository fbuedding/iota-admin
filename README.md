# IoTA-Admin
The Fiware IoT-Agent Admin is a service to manage IoT-Agents using the [IoT Agent API - Fiware-iotagent-node-lib](https://iotagent-node-lib.readthedocs.io/en/3.3.0/api.html).

It supports creating, deleting, and listing devices as well as config groups. It also supports managing multiple IoT-Agents.

## Features
- **Device Management:** Create, delete, and list devices and config groups.
- **Multi-Agent Support:** Manage multiple IoT-Agents.
- **Security:** Simple authentication via username and password.
- **Protection:** Public routes are secured with brute force protection and rate limiting.

## Planned Features
- LDAP authentication
  
## Potential Future Features
- Adding Orion-Context Broker SDK
- Live monitoring
- Manage Orion-Context Broker
- Managing and deploying other IoT Agents

## Why is the release version so high?
I made a mistake and released a v1.4.x which led Go to think an old version was the latest version.

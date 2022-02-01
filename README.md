# Crypto Dashboard
**An experiment in microservice architecture**

## Local Testing
**Prerequisites:**
- Kubectl
- Minikube (running/started)

```bash
git clone github.com/jmoussa/crypto-dashboard
cd crypto-dashboard
go get .
make build
make deploy
minikube tunnel # expose API port to your local machine
# Additional Commands
# clean up kubernetes deployment
make clean

# open minikube dashboard for monitoring
minikube dashboard

# test endpoint in browser
localhost:3000/coindesk (test coindesk REST API + gRPC endpoint)
```

## Architecture
### UI
- Single Dashboard view
    - Monitoring top 20 crypto prices (scrape from etherscan.io)
    - Twitter NFT News Feed (API)
    - CoinDesk, TodayOnChain, CCN, and CoinTelegraph
        - Crypto/Finance News Feed (API or scrape)
    - Companies that Accept Crypto (scrape from somewhere?)

### API
- REST at the high level
- gRPC for microservice communications

Top level REST API which connects to gRPC Servers that will perform the content gathering jobs and pipe the input back.


## CI/CD
Using Kubernetes and Helm for management and monitoring, (minikube for local development)
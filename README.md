# Crypto Dashboard
**An experiment in microservice architecture**

## Local Testing
**Prerequisites:**
- Kubectl
- Minikube (running/started)

```bash
# set this at the beginning of your terminal session to reference the local docker images
eval $(minikube docker-env)

git clone github.com/jmoussa/crypto-dashboard
cd crypto-dashboard
cp $(pwd)/config/config.json.template $(pwd)/config/config.json # fill in the correct values
export CONFIG_LOCATION=$(pwd)/config
go get .
make build
make deploy # will start minikube dashboard to visualize the deployment 

# in a separate terminal window
minikube tunnel # expose API port to your local machine

# Additional Commands
# clean up go cache, docker, and kubernetes deployment
make clean

# open minikube dashboard for monitoring
minikube dashboard

# test endpoint in browser
localhost:3000/coindesk (test coindesk REST API + CoinDesk Article Scraping gRPC endpoint)
localhost:3000/twitter (test coindesk REST API + Twitter Scraping gRPC endpoint)
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
Using Kubernetes and Helm for environment management and monitoring, and minikube for local development.
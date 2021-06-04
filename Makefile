GOCMD=go
BINARY_NAME=student
GOBUILD=$(GOCMD) build
PATH_TO_CHAINCODE=./chaincode/student
# generates Genesis Block
genesis-block:
	@./genesis-block.sh
	
# starts the orderer in devmode 
orderer:
	@ORDERER_GENERAL_GENESISPROFILE=SampleDevModeSolo orderer	

# starts peer
peer:
	@FABRIC_LOGGING_SPEC=chaincode=debug CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052 peer node start --peer-chaincodedev=true

#  channel configurations
channel:
	@configtxgen -channelID ch1 -outputCreateChannelTx ch1.tx -profile SampleSingleMSPChannel -configPath ${FABRIC_CFG_PATH}
	@peer channel create -o 127.0.0.1:7050 -c ch1 -f ch1.tx
	@peer channel join -b ch1.block
 
 # build the chaincode binary
chaincode1: 
	@echo "building chaincode..." 
	@$(GOBUILD) -v -o ./bin/$(BINARY_NAME) $(PATH_TO_CHAINCODE)
	@CORE_CHAINCODE_LOGLEVEL=debug CORE_PEER_TLS_ENABLED=false CORE_CHAINCODE_ID_NAME=mycc:1.0 ./bin/student -peer.address 127.0.0.1:7052


approve-chaincode:
	@peer lifecycle chaincode approveformyorg  -o 127.0.0.1:7050 --channelID ch1 --name mycc --version 1.0 --sequence 1 --init-required --signature-policy "OR ('SampleOrg.member')" --package-id mycc:1.0
	@peer lifecycle chaincode checkcommitreadiness -o 127.0.0.1:7050 --channelID ch1 --name mycc --version 1.0 --sequence 1 --init-required --signature-policy "OR ('SampleOrg.member')"
	@peer lifecycle chaincode commit -o 127.0.0.1:7050 --channelID ch1 --name mycc --version 1.0 --sequence 1 --init-required --signature-policy "OR ('SampleOrg.member')" --peerAddresses 127.0.0.1:7051


ipfs-bootstrap:
	@docker-compose -f ./ipfs/docker-compose.yaml up -d

ipfs-destroy:
	@docker-compose -f ./ipfs/docker-compose.yaml down --volumes

ipfs-cors:
	@docker exec ipfs_host ipfs config --json API.HTTPHeaders.Access-Control-Allow-Origin '["http://0.0.0.0:5001", "http://localhost:3000", "http://127.0.0.1:5001", "https://webui.ipfs.io"]'
	@docker exec ipfs_host ipfs config --json API.HTTPHeaders.Access-Control-Allow-Methods '["PUT", "POST"]'


#  docker exec ipfs_host ipfs dag get bafyreicrur7nilot43ds2zso6hiceqbbh3sv3erwjnjdcklwsrbez7onky
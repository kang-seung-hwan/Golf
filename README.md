# myFabricNetwork

fabric network 구축하기

## 전체 구성

4개의 Organizations : OrdererOrg, ScoreOrg, ReservationOrg

각 Org 별 2개의 peer
s
- OrdererOrg : raft1, raft2, raft3
- ScoreOrg : peer0
- ReservationOrg : peer0

한개의 channel 생성

- mychannel1 : ScoreOrg, ReservationOrg

채널의 anchor peer

- mychannel1 : peer0.score, peer0.reservation

각 peer 별로 couchdb 연결

- peer0.score : couchdb0 (localhost:5984)
- peer0.reservation : couchdb4 (localhost:6984)

각 채널로 fabcar chaincode 배포

## 구동

1. clone

```bash
git clone https://github.com/AdoreJE/myFabricNetwork
```

2. network up

```bash
cd myFabricNetwork/test-network
./network.sh up -ca -s couchdb
```

3. create channel

```bash
./network.sh createChannel
```

4. deploy chaincode

```bash
./network.sh deployCC -cci initLedger -ccn fabcar -ccp ../chaincode/go
```

5. 환경변수 설정

```bash
# 모든 peer 에서 공통으로 적용
export CORE_PEER_TLS_ENABLED=true
export FABRIC_CFG_PATH=${PWD}/../config

### peer0.score
export CORE_PEER_LOCALMSPID="ScoreMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/score.example.com/peers/peer0.score.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/score.example.com/users/Admin@score.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

#### 다른 peer 의 환경변수는 아래 부록 참고
```

6. smart contract 실행

peer0.score에 대한 환경변수를 설정한 경우

```bash
peer chaincode query -C mychannel1 -n fabcar-ch1 -c '{"function":"queryAllCars","Args":[""]}'
```


7. couchdb 확인
   브라우저에서 couchdb0 에 접속

```
localhost:5984/_utils/#login
```

Username : admin

Password : adminpw

mychannel1_fabcar-ch1 확인

## 부록

환경변수

```bash
# 모든 peer 에서 공통으로 적용
export CORE_PEER_TLS_ENABLED=true
export FABRIC_CFG_PATH=${PWD}/../config

### peer0.score
export CORE_PEER_LOCALMSPID="ScoreMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/score.example.com/peers/peer0.score.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/score.example.com/users/Admin@score.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

### peer0.reservation
export CORE_PEER_LOCALMSPID="ReservationMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/reservation.example.com/peers/peer0.reservation.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/reservation.example.com/users/Admin@reservation.example.com/msp
export CORE_PEER_ADDRESS=localhost:8051
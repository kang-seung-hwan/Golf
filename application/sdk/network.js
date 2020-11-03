// Setting for Hyperledger Fabric
const { Wallets, Gateway } = require('fabric-network');

const path = require('path');
const fs = require('fs');

exports.getReservationNetwork = async () => {
  // 경로 파일명 맞추기
  const ccpPath = path.resolve(
    __dirname,
    '..',
    '..',
    'test-network',
    'organizations',
    'peerOrganizations',
    'reservation.example.com',
    'connection-reservation.json'
  );

  const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

  // Create a new file system based wallet for managing identities.
  const walletPath = path.join(process.cwd(), 'wallet');
  const wallet = await Wallets.newFileSystemWallet(walletPath);
  console.log(`Wallet path: ${walletPath}`);

  // Check to see if we've already enrolled the user.
  const identity = await wallet.get('appUser');
  if (!identity) {
    console.log(
      'An identity for the user "appUser" does not exist in the wallet'
    );
    console.log('Run the registerUser.js application before retrying');
    return;
  }

  // Create a new gateway for connecting to our peer node.
  const gateway = new Gateway();
  await gateway.connect(ccp, {
    wallet,
    identity: 'appUser',
    discovery: { enabled: true, asLocalhost: true },
  });

  // 체널 이름 맞추기
  const network = await gateway.getNetwork('mychannel1');

  // 체인코드 이름 맞추기. deployCC 했을 때 -ccn 옵션으로 준 이름
  const contract = network.getContract('reservation');

  return [contract, gateway];
};

exports.getScoreNetwork = async () => {
  // Connection profile
  const ccpPath = path.resolve(
    __dirname,
    '..',
    '..',
    'test-network',
    'organizations',
    'peerOrganizations',
    'reservation.example.com',
    'connection-reservation.json'
  );

  const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

  // Create a new file system based wallet for managing identities.
  const walletPath = path.join(process.cwd(), 'wallet');
  const wallet = await Wallets.newFileSystemWallet(walletPath);
  console.log(`Wallet path: ${walletPath}`);

  // Check to see if we've already enrolled the user.
  const identity = await wallet.get('appUser');
  if (!identity) {
    console.log(
      'An identity for the user "appUser" does not exist in the wallet'
    );
    console.log('Run the registerUser.js application before retrying');
    return;
  }

  // Create a new gateway for connecting to our peer node.
  const gateway = new Gateway();
  await gateway.connect(ccp, {
    wallet,
    identity: 'appUser',
    discovery: { enabled: true, asLocalhost: true },
  });

  // Get the network (channel) our contract is deployed to.
  const network = await gateway.getNetwork('mychannel1');

  // Get the contract from the network.
  const contract = network.getContract('score');

  return [contract, gateway];
};

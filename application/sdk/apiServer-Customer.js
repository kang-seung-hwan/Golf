const express = require('express');
const cors = require('cors');
const bodyParser = require('body-parser');

let app = express();
app.use(bodyParser.json());

let network = require('./network.js');
// CORS 설정
app.use(cors())

// queryAllGround
app.get('/api/queryallground', async function (req, res) {
  try {
    // queryAllGround가 reservation에 저장되있으므로 getReservationNetwork()
    const [contract, gateway] = await network.getReservationNetwork();
    // query 함수니까 evaluationTransaction()
    // 함수명 맞출것.(첫번째 문자는 소문자로)
    // queryAllGround 의 반환형이 []*Ground -> result에 저장됨.
    const result = await contract.evaluateTransaction('queryAllGround');
    console.log(
      `Transaction has been evaluated, result is: ${result.toString()}`
    );

    // result를 string으로 바꿔서 클라이언트로 전달
    res.status(200).json({ response: result.toString() });
    // Disconnect from the gateway.
    await gateway.disconnect();
  } catch (error) {
    console.error(`Failed to evaluate transaction: ${error}`);
    // process.exit(1);
  }
});

// queryGround
// queryGround는 groundID를 받아야함
// groundID는 get 방식으로 받을거고 그럼 url에 포함됨.
app.get('/api/query/:ground_ID', async function (req, res) {
  try {
    const [contract, gateway] = await network.getReservationNetwork();
    
    const result = await contract.evaluateTransaction(
      'queryGround',
      // url 에 포함된 ground_ID는 여기에 저장되어 있음.
      req.params.ground_ID
    );
    console.log(
      `Transaction has been evaluated, result is: ${result.toString()}`
    );
    // 클라이언트로 전송
    res.status(200).json({ response: result.toString() });
    await gateway.disconnect();
  } catch (error) {
    console.error(`Failed to evaluate transaction: ${error}`);
    res.status(500).json({ error: error });
    // process.exit(1);
  }
});

// createGround
// 인자가 많기 때문에 post로  받음
// groundID, groundID, startTime, endTime, totalHole
app.post('/api/createGround/', async function (req, res) {
  try {
    const [contract, gateway] = await network.getReservationNetwork();
    // invoke 함수이므로 submitTransaction
    // 첫번째 인자로 함수명을 넣고
    // 두번째로 파라미터들을 차례대로 입력
    // ctx는 않넣어도됨
    await contract.submitTransaction(
      'createGround',
      req.body.groundID,
      req.body.groundName,
      req.body.availableTimeStart,
      req.body.availableTimeEnd,
      req.body.totalHole
    );
    console.log('Transaction has been submitted');
    
    // 결과를 클라이언트로 전송하는 부분
    res.send('Transaction has been submitted');
    // Disconnect from the gateway.
    await gateway.disconnect();
  } catch (error) {
    console.error(`Failed to submit transaction: ${error}`);
    // process.exit(1);
  }
});

// reserveGround
app.post('/api/reserveGround/', async function (req, res) {
  try {
    const [contract, gateway] = await network.getReservationNetwork();
    contract.addContractListener(listener);

    // Submit the specified transaction.
    await contract.submitTransaction(
      'reserveGround',
      req.body.groundID,
      req.body.userID,
      req.body.begin,
      req.body.end
    );
    console.log('Transaction has been submitted');
    res.send('Transaction has been submitted');
    // Disconnect from the gateway.
    await gateway.disconnect();
  } catch (error) {
    console.error(`Failed to submit transaction: ${error}`);
    // process.exit(1);
  }
});

app.get('/api/confirmReservation/:groundID/:userID', async function (req, res) {
  try {
    const [contract, gateway] = await network.getReservationNetwork();
    // Evaluate the specified transaction.
    // queryGround transaction - requires 1 argument, ex: ('queryGround', 'Ground01')
    // queryAllGround transaction - requires no arguments, ex: ('queryAllGround')
    const result = await contract.evaluateTransaction(
      'confirmReservation',
      req.params.groundID,
      req.params.userID
    );
    console.log(
      `Transaction has been evaluated, result is: ${result.toString()}`
    );
    res.status(200).json({ response: result.toString() });
    await gateway.disconnect();
  } catch (error) {
    console.error(`Failed to evaluate transaction: ${error}`);
    res.status(500).json({ error: error });
    // process.exit(1);
  }
});


app.get('/api/userConfirmReservation/:userID', async function (req, res) {
  try {
    const [contract, gateway] = await network.getReservationNetwork();

    // result에 []*Reservation이 반환되어 저장됨.
    const result = await contract.evaluateTransaction(
      'userConfirmReservation',
      req.params.userID
    );
    console.log(
      `Transaction has been evaluated, result is: ${result.toString()}`
    );
    // 클라이언트로 전송
    res.status(200).json({ response: result.toString() });
    await gateway.disconnect();
  } catch (error) {
    console.error(`Failed to evaluate transaction: ${error}`);
    res.status(500).json({ error: error });
    // process.exit(1);
  }
});

// 여기서부터 score 부분

// reserveGround
app.post('/api/startGame/', async function (req, res) {
  try {
    const [contract, gateway] = await network.getScoreNetwork();

    // Submit the specified transaction.
    await contract.submitTransaction(
      'startGame',
      req.body.groundID,
      req.body.user,
      req.body.userNumber,
      req.body.gameCode
    );
    console.log('Transaction has been submitted');
    res.send('Transaction has been submitted');
    // Disconnect from the gateway.
    await gateway.disconnect();
  } catch (error) {
    console.error(`Failed to submit transaction: ${error}`);
    // process.exit(1);
  }
});

// queryGameInfo
app.get('/api/queryGameInfo/:groundID/:gameNumber', async function (req, res) {
  try {
    const [contract, gateway] = await network.getScoreNetwork();

    // result에 []*Reservation이 반환되어 저장됨.
    const result = await contract.evaluateTransaction(
      'queryGameInfo',
      req.params.groundID,
      req.params.gameNumber
    );
    console.log(
      `Transaction has been evaluated, result is: ${result.toString()}`
    );
    // 클라이언트로 전송
    res.status(200).json({ response: result.toString() });
    await gateway.disconnect();
  } catch (error) {
    console.error(`Failed to evaluate transaction: ${error}`);
    res.status(500).json({ error: error });
    // process.exit(1);
  }
});

// setScore
app.post('/api/setScore/', async function (req, res) {
  try {
    const [contract, gateway] = await network.getScoreNetwork();

    // Submit the specified transaction.
    await contract.submitTransaction(
      'setScore',
      req.body.gameNumber,
      req.body.holeNumber,
      req.body.userNumber,
      req.body.score
    );
    console.log('Transaction has been submitted');
    res.send('Transaction has been submitted');
    // Disconnect from the gateway.
    await gateway.disconnect();
  } catch (error) {
    console.error(`Failed to submit transaction: ${error}`);
    // process.exit(1);
  }
});


// queryScore
app.get('/api/queryScore/:gameNumber/:holeNumber', async function (req, res) {
  try {
    const [contract, gateway] = await network.getScoreNetwork();

    // result에 []*Reservation이 반환되어 저장됨.
    const result = await contract.evaluateTransaction(
      'queryScore',
      req.params.gameNumber,
      req.params.holeNumber
    );
    console.log(
      `Transaction has been evaluated, result is: ${result.toString()}`
    );
    // 클라이언트로 전송
    res.status(200).json({ response: result.toString() });
    await gateway.disconnect();
  } catch (error) {
    console.error(`Failed to evaluate transaction: ${error}`);
    res.status(500).json({ error: error });
    // process.exit(1);
  }
});


// agreeScore
app.post('/api/agreeScore/', async function (req, res) {
  try {
    const [contract, gateway] = await network.getScoreNetwork();

    // post로 요청할 때 isAgreed 에 "agree" 로 저장되어야함 
    // Submit the specified transaction.
    await contract.submitTransaction(
      'agreeScore',
      req.body.gameNumber,
      req.body.holeNumber,
      req.body.userNumber,
      req.body.isAgreed
    );
    console.log('Transaction has been submitted');
    res.send('Transaction has been submitted');
    // Disconnect from the gateway.
    await gateway.disconnect();
  } catch (error) {
    console.error(`Failed to submit transaction: ${error}`);
    // process.exit(1);
  }
});


// queryAgreement
app.get('/api/queryScore/:gameNumber/:holeNumber', async function (req, res) {
  try {
    const [contract, gateway] = await network.getScoreNetwork();

    // result에 []*Reservation이 반환되어 저장됨.
    const result = await contract.evaluateTransaction(
      'queryAgreement',
      req.params.gameNumber,
      req.params.holeNumber
    );
    console.log(
      `Transaction has been evaluated, result is: ${result.toString()}`
    );
    // 클라이언트로 전송
    res.status(200).json({ response: result.toString() });
    await gateway.disconnect();
  } catch (error) {
    console.error(`Failed to evaluate transaction: ${error}`);
    res.status(500).json({ error: error });
    // process.exit(1);
  }
});




// validateScore
app.post('/api/validateScore/', async function (req, res) {
  try {
    const [contract, gateway] = await network.getScoreNetwork();

    // post로 요청할 때 isAgreed 에 "agree" 로 저장되어야함 
    // Submit the specified transaction.
    await contract.submitTransaction(
      'validateScore',
      req.body.gameNumber,
      req.body.holeNumber
    );
    console.log('Transaction has been submitted');
    res.send('Transaction has been submitted');
    // Disconnect from the gateway.
    await gateway.disconnect();
  } catch (error) {
    console.error(`Failed to submit transaction: ${error}`);
    // process.exit(1);
  }
});


// queryTotalGameScore
app.get('/api/queryTotalGameScore/:gameNumber', async function (req, res) {
  try {
    const [contract, gateway] = await network.getScoreNetwork();

    // result에 []*Reservation이 반환되어 저장됨.
    const result = await contract.evaluateTransaction(
      'queryTotalGameScore',
      req.params.gameNumber
    );
    console.log(
      `Transaction has been evaluated, result is: ${result.toString()}`
    );
    // 클라이언트로 전송
    res.status(200).json({ response: result.toString() });
    await gateway.disconnect();
  } catch (error) {
    console.error(`Failed to evaluate transaction: ${error}`);
    res.status(500).json({ error: error });
    // process.exit(1);
  }
});

app.listen(8080);

// // 나중에 쓸 수 도?
let details = '';
const listener = async (event) => {
  if (event.eventName === 'newReservation') {
    console.log('\n\nnewReservation');
    details = event.payload.toString('utf8');
    // Run business process to handle orders
    console.log('===============');
    console.log(details);
  }
};

// app.get('/api/events', async function (req, res) {
//   console.log('request received');
//   res.writeHead(200, {
//     Connection: 'keep-alive',
//     'Content-Type': 'text/event-stream',
//     'Cache-Control': 'no-cache',
//   });
//   res.write('\n');
//   setInterval(() => {
//     res.write('event: reservation\n'); // added these
//     res.write(`data: ${JSON.stringify(details)}`);
//     res.write('\n\n');
//   }, 5000);
// });

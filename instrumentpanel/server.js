// server.js
const WebSocket = require('ws');
const wss = new WebSocket.Server({ port: 8080 });

wss.on('connection', (ws) => {
  ws.on('message', (message) => {
    console.log('received:', message);
    // Handle control commands
  });

  // Simulate sending telemetry data
  setInterval(() => {
    const data = JSON.stringify({
      attitude: { pitch: Math.random() * 10, roll: Math.random() * 10 },
      airspeed: Math.random() * 200,
      altitude: Math.random() * 10000,
      verticalSpeed: (Math.random() - 0.5) * 2000,
      heading: Math.random() * 360,
      turnRate: (Math.random() - 0.5) * 5,
      engineRPMs: [Math.random() * 6000, Math.random() * 6000, Math.random() * 6000, Math.random() * 6000],
      fuel: Math.random() * 100,
    });
    ws.send(data);
  }, 1000);
});

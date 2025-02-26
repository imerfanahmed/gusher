import ws from 'k6/ws';
import { check } from 'k6';

export const options = {
  vus: 500, // Number of virtual users
  duration: '60s', // Test duration
};

export default function () {
  //const url = 'wss://bifrost.magicoffice.co.uk:6001/app/b3Asdg5EChRS1dn?protocol=7&client=js&version=4.4.0'; // Replace with your Soketi server URL and credentials
  // const url1 = 'ws://127.0.0.1:6001/app/Ex823XN5GxGhY2q'; // Replace with your Soketi server URL and credentials
  const url = 'ws://83.136.252.254:3000/app/3uGbZsHavJZVz37'; // Replace with your Soketi server URL and credentials

  const res1 = ws.connect(url, {}, (socket) => {
    // WebSocket connection established
    socket.on('open', () => {
      console.log('Connected');

      // Subscribe to a channel
      const subscribeMessage = JSON.stringify({
        event: 'pusher:subscribe',
        data: {
          channel: 'test-channel', // Replace with your channel name
        },
      });
      socket.send(subscribeMessage);

      // Log incoming messages
      socket.on('message', (message) => {
        console.log('Received message:', message);
      });
    });

    // Check for errors
    socket.on('error', (e) => {
      console.log('Error:', e.error());
    });

    // Disconnect after 10 seconds
    socket.setTimeout(() => {
      socket.close();
    }, 10000);
  });


  // Validate connection
  check(res1, {
    'connected successfully': (r) => r && r.status === 101,
  });


}

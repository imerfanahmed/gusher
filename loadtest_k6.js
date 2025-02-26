import ws from 'k6/ws';
import { Trend } from 'k6/metrics';

/**
 * This script tests the WebSocket server with two scenarios:
 * 1. soakTraffic: Sustained connections to simulate stable users.
 * 2. highTraffic: High connection/disconnection rate to simulate active traffic.
 *
 * Prerequisites:
 * - Ensure your Go WebSocket server is running:
 *   go run cmd/main.go
 * - Optionally, simulate message sending (e.g., via a separate tool or script).
 *
 * To simulate varying message rates, you could use a separate client:
 * - Low: 1 msg/sec
 * - Mild: 2 msg/sec
 * - Overkill: 10 msg/sec
 */

const delayTrend = new Trend('message_delay_ms'); // Tracks delay from send to receive

// Default performance thresholds
let maxP95 = 100; // 95th percentile delay in ms
let maxAvg = 100; // Average delay in ms

// Adjust thresholds based on external dependencies (if applicable)
if (['mysql'].includes(__ENV.DB_DRIVER)) {
    maxP95 += 500; // MySQL adds latency
    maxAvg += 100;
}

if (['redis'].includes(__ENV.CACHE_DRIVER)) {
    maxP95 += 20;  // Redis caching overhead
    maxAvg += 20;
}

export const options = {
    thresholds: {
        message_delay_ms: [
            { threshold: `p(95)<${maxP95}`, abortOnFail: false },
            { threshold: `avg<${maxAvg}`, abortOnFail: false },
        ],
    },

    scenarios: {
        // Soak Traffic: Sustained connections
        soakTraffic: {
            executor: 'ramping-vus',
            startVUs: 0,
            startTime: '0s',
            stages: [
                { duration: '50s', target: 250 }, // Ramp up to 250 VUs
                { duration: '110s', target: 250 }, // Hold at 250 VUs
            ],
            gracefulRampDown: '40s',
            env: {
                SLEEP_FOR: '160', // Total test duration in seconds
                WS_HOST: __ENV.WS_HOST || 'ws://83.136.252.254:3000/app/app_key',
            },
        },

        // High Traffic: Connection churn
        highTraffic: {
            executor: 'ramping-vus',
            startVUs: 0,
            startTime: '50s', // Start after soakTraffic stabilizes
            stages: [
                { duration: '50s', target: 250 }, // Ramp up to 250 VUs
                { duration: '30s', target: 250 }, // Hold
                { duration: '10s', target: 100 }, // Ramp down
                { duration: '10s', target: 50 },  // Further down
                { duration: '10s', target: 100 }, // Spike back up
            ],
            gracefulRampDown: '20s',
            env: {
                SLEEP_FOR: '110', // Shorter duration for churn
                WS_HOST: __ENV.WS_HOST || 'ws://83.136.252.254:3000/app/app_key',
            },
        },
    },
};

export default () => {
    ws.connect(__ENV.WS_HOST, null, (socket) => {
        // Close connection after specified duration
        socket.setTimeout(() => {
            socket.close();
        }, __ENV.SLEEP_FOR * 1000);

        socket.on('open', () => {
            // Keep connection alive with periodic pings
            socket.setInterval(() => {
                socket.send(JSON.stringify({
                    event: 'pusher:ping',
                    data: {},
                }));
            }, 30000); // Every 30 seconds

            // Subscribe to a unique channel per VU
            const channel = `test-channel-${__VU}`;
            socket.send(JSON.stringify({
                event: 'pusher:subscribe',
                data: { channel: channel },
            }));

            socket.on('message', (message) => {
                const receivedTime = Date.now();
                let msg;

                try {
                    msg = JSON.parse(message);
                } catch (e) {
                    console.error(`VU ${__VU}: Failed to parse message: ${e}`);
                    return;
                }

                // Measure delay for timed messages (if sent externally)
                if (msg.event === 'timed-message') {
                    const data = JSON.parse(msg.data);
                    delayTrend.add(receivedTime - data.time);
                } else {
                    console.log(`VU ${__VU}: Received event: ${msg.event} on ${channel}`);
                }
            });

            socket.on('error', (e) => {
                console.error(`VU ${__VU}: WebSocket error: ${e.error()}`);
            });
        });

        socket.on('close', () => {
            console.log(`VU ${__VU}: Connection closed`);
        });
    });
};
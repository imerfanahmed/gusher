document.addEventListener('DOMContentLoaded', () => {
    const tabs = {
        apps: document.getElementById('appsSection'),
        channels: document.getElementById('channelsSection'),
        events: document.getElementById('eventsSection'),
        debug: document.getElementById('debugSection')
    };

    const buttons = {
        appsBtn: document.getElementById('appsBtn'),
        channelsBtn: document.getElementById('channelsBtn'),
        eventsBtn: document.getElementById('eventsBtn'),
        debugBtn: document.getElementById('debugBtn')
    };

    // Tab switching
    Object.keys(buttons).forEach(key => {
        buttons[key].addEventListener('click', () => {
            Object.values(tabs).forEach(tab => tab.classList.remove('active'));
            tabs[key.replace('Btn', '')].classList.add('active');
        });
    });

    // Initialize Pusher for real-time updates
    const pusher = new Pusher('app_key', {
        cluster
: 'eu',
        wsHost: ' 83.136.252.254',
        wsPort: 3000,
        forceTLS: false,
        disableStats: true
    });

    // App Management
    const createAppForm = document.getElementById('createAppForm');
    const appList = document.getElementById('appList');

    createAppForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const appKey = document.getElementById('appKey').value;
        const appSecret = document.getElementById('appSecret').value;
        const appID = document.getElementById('appID').value;

        try {
            const response = await fetch('http://127.0.0.1:8080/apps', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ key: appKey, secret: appSecret, id: appID })
            });
            if (response.ok) {
                loadApps();
                createAppForm.reset();
            }
        } catch (error) {
            console.error('Failed to create app:', error);
        }
    });

    async function loadApps() {
        try {
            const response = await fetch('http://127.0.0.1:8080/apps');
            const apps = await response.json();
            appList.innerHTML = apps.map(app => `<li>App ID: ${app.id}, Key: ${app.key}</li>`).join('');
        } catch (error) {
            console.error('Failed to load apps:', error);
        }
    }

    // Channel Debugging
    const channelFilter = document.getElementById('channelFilter');
    const channelList = document.getElementById('channelList');

    channelFilter.addEventListener('input', () => {
        const filter = channelFilter.value.toLowerCase();
        // Simulate channel data (replace with real API call)
        const channels = ['test-channel-1', 'test-channel-2', 'private-channel'];
        channelList.innerHTML = channels
            .filter(ch => ch.toLowerCase().includes(filter))
            .map(ch => `<li>${ch} (Occupancy: ${Math.floor(Math.random() * 10)})</li>`).join('');
    });

    // Event Logger
    const eventFilter = document.getElementById('eventFilter');
    const eventLog = document.getElementById('eventLog');

    const eventChannel = pusher.subscribe('event-log');
    eventChannel.bind('new-event', data => {
        const event = data.data;
        const filtered = eventFilter.value.toLowerCase();
        if (event.event.toLowerCase().includes(filtered) || event.channel.toLowerCase().includes(filtered)) {
            eventLog.insertAdjacentHTML('afterbegin', `<li>[${new Date().toLocaleTimeString()}] ${event.event} on ${event.channel}: ${JSON.stringify(event.data)}</li>`);
            if (eventLog.children.length > 100) eventLog.removeChild(eventLog.lastChild); // Limit log size
        }
    });

    // Debug Console
    const debugChannelInput = document.getElementById('debugChannel');
    const debugEventInput = document.getElementById('debugEvent');
    const debugDataInput = document.getElementById('debugData');
    const triggerEventBtn = document.getElementById('triggerEvent');
    const debugOutput = document.getElementById('debugOutput');

    triggerEventBtn.addEventListener('click', () => {
        const channel = debugChannelInput.value;
        const event = debugEventInput.value;
        let data = {};
        try {
            data = JSON.parse(debugDataInput.value || '{}');
        } catch (e) {
            debugOutput.innerText = `Invalid JSON: ${e.message}`;
            return;
        }

        pusher.trigger(channel, event, data).then(() => {
            debugOutput.innerText = `Triggered ${event} on ${channel} with data: ${JSON.stringify(data)}`;
        }).catch(error => {
            debugOutput.innerText = `Error triggering event: ${error.message}`;
        });
    });

    // Load initial data
    loadApps();
    document.getElementById('appsSection').classList.add('active');
});
from flask import Flask, render_template
from flask_socketio import SocketIO, emit
import threading
import time
import random

app = Flask(__name__)
socketio = SocketIO(app)

gps_data = {
    'latitude': 30.2672,
    'longitude': -97.7431,
    'altitude': 0,
    'speed': 0,
    'heading': 0,
    'verticalSpeed': 0,
    'pitch': 0,
    'roll': 0
}

@app.route('/')
def index():
    return render_template('index.html')

@socketio.on('connect')
def handle_connect():
    print('Client connected')

@socketio.on('disconnect')
def handle_disconnect():
    print('Client disconnected')

def gps_simulation():
    while True:
        gps_data['latitude'] = 30.2672 + random.uniform(-0.001, 0.001)
        gps_data['longitude'] = -97.7431 + random.uniform(-0.001, 0.001)
        gps_data['altitude'] = random.uniform(100, 1000)
        gps_data['speed'] = random.uniform(50, 200)
        gps_data['heading'] = random.uniform(0, 360)
        gps_data['verticalSpeed'] = random.uniform(-10, 10)
        gps_data['pitch'] = random.uniform(-10, 10)
        gps_data['roll'] = random.uniform(-30, 30)
        socketio.emit('gps_update', gps_data)
        time.sleep(1)

if __name__ == '__main__':
    threading.Thread(target=gps_simulation).start()
    socketio.run(app, debug=True)

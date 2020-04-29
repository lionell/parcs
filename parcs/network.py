import struct
import socket
import time

SLEEP_DURATION = 0.5
PORT = 4444

def send(sock, msg):
    data = struct.pack('>I', len(msg)) + msg
    sock.sendall(data)

def recv(sock):
    raw_len = _recv(sock, 4)
    size = struct.unpack('>I', raw_len)[0]
    return _recv(sock, size)

def _recv(sock, size):
    data = b''
    while len(data) < size:
        chunk = sock.recv(size - len(data))
        if not chunk:
            return None
        data += chunk
    return data

def handshake(sock, side):
    if side == 'server':
        assert(_recv(sock, 3) == b'SYN')
        sock.sendall(b'ACK')
    else:
        sock.sendall(b'SYN')
        assert(_recv(sock, 3) == b'ACK')

def dns_lookup(hostname):
    while True:
        try:
            return socket.gethostbyname(hostname)
        except:
            time.sleep(SLEEP_DURATION)


def connect(sock, ip):
    while True:
        try:
            sock.connect((ip, PORT))
            return
        except:
            time.sleep(SLEEP_DURATION)

import struct
import socket
import time
from parcs.data import BYTES_IN_INT, pack_int, unpack_int

SLEEP_DURATION = 0.5
PORT = 4444

def send(sock, msg):
    data = pack_int(len(msg)) + msg
    sock.sendall(data)

def recv(sock):
    raw_len = _recv(sock, BYTES_IN_INT)
    size = unpack_int(raw_len)
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

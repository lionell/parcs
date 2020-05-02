import struct
import socket
import time
import logging
from parcs.data import *

SLEEP_DURATION = 0.5
PORT = 4444

def send(sock, obj):
    send_bytes(sock, marshal(obj))

def send_bytes(sock, data):
    sock.sendall(encode_uint64(len(data)))
    sock.sendall(data)

def recv(sock):
    return unmarshal(recv_bytes(sock))

def recv_bytes(sock):
    size = decode_uint64(_recv_all_bytes(sock, BYTES_IN_INT))
    return _recv_all_bytes(sock, size)

def _recv_all_bytes(sock, size):
    data = b''
    while len(data) < size:
        chunk = sock.recv(size - len(data))
        if not chunk:
            return None
        data += chunk
    return data

def handshake(sock, side):
    if side == 'server':
        assert(_recv_all_bytes(sock, 3) == b'SYN')
        sock.sendall(b'ACK')
    else:
        sock.sendall(b'SYN')
        assert(_recv_all_bytes(sock, 3) == b'ACK')

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

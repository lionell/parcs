import json

BYTES_IN_INT = 8
BYTEORDER = 'big'

def marshal(obj):
    return json.dumps(obj).encode()

def unmarshal(bs):
    return json.loads(bs.decode())

def encode_uint64(n):
    return n.to_bytes(length=BYTES_IN_INT, byteorder=BYTEORDER)

def decode_uint64(data):
    return int.from_bytes(data, byteorder=BYTEORDER)

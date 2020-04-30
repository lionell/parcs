BYTES_IN_INT = 8
BYTEORDER = 'big'

def pack_int(n):
    return n.to_bytes(length=BYTES_IN_INT, byteorder=BYTEORDER)

def unpack_int(bs):
    return int.from_bytes(bs, byteorder=BYTEORDER)

def pack_ints(l):
    return b''.join(map(pack_int, l))

def unpack_ints(bs):
    return [
        unpack_int(bs[i:i + BYTES_IN_INT])
        for i in range(0, len(bs), BYTES_IN_INT)
    ]

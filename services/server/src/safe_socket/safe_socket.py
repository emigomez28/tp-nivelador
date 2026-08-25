import socket


def recv_all(sock: socket.socket, size):
    total_read = b""

    while len(total_read) < size:
        read = sock.recv(size - len(total_read))

        if read == b"":
            raise RuntimeError("Connection closed before receiving all data")

        total_read += read

    return total_read


def send_all(sock: socket.socket, data: bytes):
    total_sent = 0

    while total_sent < len(data):
        sent = sock.send(data[total_sent:])
        total_sent += sent

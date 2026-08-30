import socket

import safe_socket

from .message import (
    _HEADER_SIZE,
    _LENGTH_FIELD_SIZE,
    _MAX_PAYLOAD_SIZE,
    _TYPE_FIELD_SIZE,
    Message,
)


def send_message(sock: socket.socket, message: Message) -> None:
    if len(message.payload) > _MAX_PAYLOAD_SIZE:
        raise RuntimeError(f"payload exceeds {_MAX_PAYLOAD_SIZE} bytes")

    header = get_header(message)
    safe_socket.send_all(sock, header + message.payload)


def recv_message(sock: socket.socket) -> Message:
    header = safe_socket.recv_all(sock, _HEADER_SIZE)

    message_type = header[0]
    payload_size = int.from_bytes(header[_TYPE_FIELD_SIZE:], "big")

    try:
        payload = safe_socket.recv_all(sock, payload_size)
    except RuntimeError:
        raise RuntimeError("connection closed mid-message")

    return Message(message_type, payload)


def get_header(message):
    return bytes([message.type]) + len(message.payload).to_bytes(
        _LENGTH_FIELD_SIZE, "big"
    )

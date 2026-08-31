_TYPE_FIELD_SIZE = 1
_LENGTH_FIELD_SIZE = 2
_HEADER_SIZE = _TYPE_FIELD_SIZE + _LENGTH_FIELD_SIZE

_MAX_UINT_16 = 65535
_MAX_PAYLOAD_SIZE = _MAX_UINT_16

_UNKNOWN_TYPE_NAME = "UNKNOWN"


class MessageType:
    START_TRANSMISSION = 0x01
    BETS = 0x02
    END_TRANSMISSION = 0x03
    OK = 0x10
    WINNERS = 0x11
    ERROR = 0x12

    @classmethod
    def get_name_from_code(cls, code: int) -> str:
        return _TYPE_NAMES.get(code, _UNKNOWN_TYPE_NAME)


_TYPE_NAMES = {
    MessageType.START_TRANSMISSION: "START_TRANSMISSION",
    MessageType.BETS: "BETS",
    MessageType.END_TRANSMISSION: "END_TRANSMISSION",
    MessageType.OK: "OK",
    MessageType.WINNERS: "WINNERS",
    MessageType.ERROR: "ERROR",
}


class Message:
    def __init__(self, message_type: int, payload: bytes = b"") -> None:
        self.type = message_type
        self.payload = payload

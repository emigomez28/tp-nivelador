from lottery import Bet

_AGENCY_ID_SIZE = 2
_TEXT_LENGTH_SIZE = 1
_DOCUMENT_SIZE = 4
_NUMBER_SIZE = 4

_MAX_TEXT_LENGTH = 255
_ENDIANNESS = "big"


def decode_agency_id(payload: bytes) -> int:
    if len(payload) != _AGENCY_ID_SIZE:
        raise RuntimeError(
            f"agency id must be {_AGENCY_ID_SIZE} bytes, got {len(payload)}"
        )

    return int.from_bytes(payload, _ENDIANNESS)


def decode_bets(payload: bytes, agency_id: int) -> list[Bet]:
    if not payload:
        raise RuntimeError("empty bet batch")

    bets = []
    offset = 0

    while offset < len(payload):
        bet, offset = _decode_bet(payload, offset, agency_id)
        bets.append(bet)

    return bets


def encode_winners(winners: list[Bet]) -> bytes:
    payload = b""

    for bet in winners:
        payload += _encode_bet(bet)

    return payload


def _decode_bet(payload: bytes, offset: int, agency_id: int) -> tuple[Bet, int]:
    first_name, offset = _read_text(payload, offset)
    last_name, offset = _read_text(payload, offset)
    document, offset = _read_int(payload, offset, _DOCUMENT_SIZE)
    birthdate, offset = _read_text(payload, offset)
    number, offset = _read_int(payload, offset, _NUMBER_SIZE)

    bet = Bet(agency_id, first_name, last_name, document, birthdate, number)

    return bet, offset


def _encode_bet(bet: Bet) -> bytes:
    return (
        _encode_text(bet.first_name)
        + _encode_text(bet.last_name)
        + _to_bytes(bet.document, _DOCUMENT_SIZE, "document")
        + _encode_text(bet.birthdate)
        + _to_bytes(bet.number, _NUMBER_SIZE, "number")
    )


def _read_int(payload: bytes, offset: int, size: int) -> tuple[int, int]:
    end = offset + size

    if end > len(payload):
        raise RuntimeError("truncated bet record")

    return int.from_bytes(payload[offset:end], _ENDIANNESS), end


def _read_text(payload: bytes, offset: int) -> tuple[str, int]:
    length, offset = _read_int(payload, offset, _TEXT_LENGTH_SIZE)
    end = offset + length

    if end > len(payload):
        raise RuntimeError("truncated bet record")

    try:
        return payload[offset:end].decode("utf-8"), end
    except UnicodeDecodeError:
        raise RuntimeError("bet record field is not valid utf-8")


def _encode_text(value: str) -> bytes:
    encoded = value.encode("utf-8")

    if len(encoded) > _MAX_TEXT_LENGTH:
        raise RuntimeError(f"field exceeds {_MAX_TEXT_LENGTH} bytes: {value}")

    return len(encoded).to_bytes(_TEXT_LENGTH_SIZE, _ENDIANNESS) + encoded


def _to_bytes(value: int, size: int, field_name: str) -> bytes:
    try:
        return value.to_bytes(size, _ENDIANNESS)
    except OverflowError:
        raise RuntimeError(f"{field_name} does not fit in {size} bytes: {value}")

from lottery import Bet

_BET_FIELDS_AMOUNT = 5
_FIELD_SEPARATOR = ","
_BET_SEPARATOR = "\n"


def decode_agency_id(payload: bytes) -> int:
    try:
        return int(payload.decode("utf-8"))
    except ValueError:
        raise RuntimeError(f"invalid agency id: {payload}")


def decode_bets(payload: bytes, agency_id: int) -> list[Bet]:
    if not payload:
        raise RuntimeError("empty bet batch")

    lines = payload.decode("utf-8").split(_BET_SEPARATOR)
    return [_decode_bet(line, agency_id) for line in lines]


def _decode_bet(line: str, agency_id: int) -> Bet:
    fields = line.split(_FIELD_SEPARATOR)

    if len(fields) != _BET_FIELDS_AMOUNT:
        raise RuntimeError(f"expected {_BET_FIELDS_AMOUNT} fields, got {len(fields)}")

    first_name, last_name, document, birthdate, number = fields

    return Bet(
        agency_id,
        first_name,
        last_name,
        _cast_to_int(document, "document"),
        birthdate,
        _cast_to_int(number, "number"),
    )


def encode_winners(winners: list[Bet]) -> bytes:
    lines = []
    for bet in winners:
        fields = _get_fields_as_list(bet)
        line = _FIELD_SEPARATOR.join(fields)
        lines.append(line)

    return "\n".join(lines).encode("utf-8")


def _cast_to_int(value: str, field_name: str) -> int:
    try:
        return int(value)
    except ValueError:
        raise RuntimeError(f"could not cast {field_name} field: '{value}' to int")


def _get_fields_as_list(bet: Bet):
    return [
        bet.first_name,
        bet.last_name,
        str(bet.document),
        bet.birthdate,
        str(bet.number),
    ]

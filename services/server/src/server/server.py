import socket
import logger
import protocol
from lottery import Lottery


class Server:
    def __init__(self, server_host: str, server_port: int, storage_path: str) -> None:
        self.server_host = server_host
        self.server_port = server_port
        self.lottery = Lottery(storage_path)

    def _handle_client(self, client_socket):
        action = "handle-client"

        try:
            agency_id = self._recv_start_transmission(client_socket)
            bets_amount = self._recv_bets(client_socket, agency_id)
            self._send_winners(client_socket, agency_id)
        except Exception as e:
            logger.error(action, logger.LogResult.fail, "err", e)
            return

        logger.info(
            action,
            logger.LogResult.success,
            "agency-id",
            agency_id,
            "bets-amount",
            bets_amount,
        )

    def run(self):
        action = "accept-connection"
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as server_socket:
            server_socket.bind((self.server_host, self.server_port))
            server_socket.listen()
            while True:
                try:
                    logger.info(action, logger.LogResult.in_progress)
                    client_socket, _ = server_socket.accept()
                except Exception as e:
                    logger.error(action, logger.LogResult.fail)
                    raise e
                logger.info(action, logger.LogResult.success)

                with client_socket:
                    self._handle_client(client_socket)

    def _recv_expecting(self, client_socket, expected: int) -> protocol.Message:
        message = protocol.recv_message(client_socket)

        if message.type != expected:
            expected_name = protocol.MessageType.get_name_from_code(expected)
            received_name = protocol.MessageType.get_name_from_code(message.type)
            raise RuntimeError(f"expected {expected_name}, got {received_name}")

        return message

    def _recv_start_transmission(self, client_socket) -> int:
        message = self._recv_expecting(
            client_socket, protocol.MessageType.START_TRANSMISSION
        )

        agency_id = protocol.decode_agency_id(message.payload)
        protocol.send_message(client_socket, protocol.Message(protocol.MessageType.OK))
        return agency_id

    def _recv_bets(self, client_socket, agency_id) -> int:
        bets_amount = 0

        while True:
            message = protocol.recv_message(client_socket)

            if message.type == protocol.MessageType.END_TRANSMISSION:
                break

            if message.type != protocol.MessageType.BETS:
                received_name = protocol.MessageType.get_name_from_code(message.type)
                raise RuntimeError(
                    f"expected BETS or END_TRANSMISSION, got {received_name}"
                )

            bets = protocol.decode_bets(message.payload, agency_id)
            self.lottery.store_bets(bets)

            bets_amount += len(bets)

            protocol.send_message(
                client_socket, protocol.Message(protocol.MessageType.OK)
            )

        return bets_amount

    def _send_winners(self, client_socket, agency_id: int) -> None:
        winners = self._get_winners(agency_id)

        protocol.send_message(
            client_socket,
            protocol.Message(
                protocol.MessageType.WINNERS, protocol.encode_winners(winners)
            ),
        )

    def _get_winners(self, agency_id):
        winners = []

        for bet in self.lottery.load_bets():
            if bet.agency_id == agency_id and self.lottery.has_won(bet):
                winners.append(bet)

        return winners

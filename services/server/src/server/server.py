import socket
import threading

import logger
import protocol
from lottery import Lottery

CLIENT_THREAD_JOIN_TIMEOUT_SECONDS = 0.5


def _log_draw() -> None:
    logger.info("lottery-draw", logger.LogResult.success)


class Server:
    def __init__(
        self,
        server_host: str,
        server_port: int,
        storage_path: str,
        agency_quorum_min: int,
    ) -> None:
        self.server_host = server_host
        self.server_port = server_port
        self.lottery = Lottery(storage_path)
        self._lottery_lock = threading.Lock()
        self._quorum = threading.Barrier(agency_quorum_min, action=_log_draw)
        self._is_running = True
        self._server_socket = None
        self._client_sockets_by_thread = {}

    def shutdown(self) -> None:
        self._is_running = False
        self._unblock_accept()

    def _unblock_accept(self) -> None:
        if self._server_socket is None:
            return
        try:
            self._server_socket.shutdown(socket.SHUT_RDWR)
        except OSError:
            pass

    def _unblock_clients(self) -> None:
        for client_socket in self._client_sockets_by_thread.values():
            try:
                client_socket.shutdown(socket.SHUT_RDWR)
            except OSError:
                pass

    def _handle_client(self, client_socket):
        action = "handle-client"
        agency_id = None

        try:
            with client_socket:
                agency_id = self._recv_start_transmission(client_socket)
                threading.current_thread().name = f"thread-agency-{agency_id}"
                bets_amount = self._recv_bets(client_socket, agency_id)
                self._quorum.wait()
                self._send_winners(client_socket, agency_id)

            logger.info(
                action,
                logger.LogResult.success,
                "agency-id",
                agency_id,
                "bets-amount",
                bets_amount,
            )
        except Exception as e:
            if self._is_running:
                logger.error(action, logger.LogResult.fail, "err", e)
            else:
                logger.info(
                    "handle-client-shutdown",
                    logger.LogResult.success,
                    "agency-id",
                    agency_id,
                )

    def run(self):
        action = "accept-connection"
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as server_socket:
            self._server_socket = server_socket
            server_socket.bind((self.server_host, self.server_port))
            server_socket.listen()

            while self._is_running:
                try:
                    logger.info(action, logger.LogResult.in_progress)
                    client_socket, _ = server_socket.accept()
                except OSError as e:
                    if not self._is_running:
                        break
                    logger.error(action, logger.LogResult.fail)
                    raise e
                logger.info(action, logger.LogResult.success)

                self._untrack_finished_threads()
                self._spawn_client_thread(client_socket)

        self._shutdown_clients()

    def _shutdown_clients(self) -> None:
        logger.info("shutdown", logger.LogResult.in_progress)
        self._quorum.abort()
        self._unblock_clients()
        self._join_client_threads()
        logger.info("shutdown", logger.LogResult.success)

    def _spawn_client_thread(self, client_socket) -> None:
        thread = threading.Thread(target=self._handle_client, args=(client_socket,))
        self._client_sockets_by_thread[thread] = client_socket
        thread.start()

    def _untrack_finished_threads(self) -> None:
        alive_clients = {}

        for thread, client_socket in self._client_sockets_by_thread.items():
            if thread.is_alive():
                alive_clients[thread] = client_socket

        self._client_sockets_by_thread = alive_clients

    def _join_client_threads(self) -> None:
        threads = list(self._client_sockets_by_thread.keys())

        for thread in threads:
            thread.join(timeout=CLIENT_THREAD_JOIN_TIMEOUT_SECONDS)

        pending = [thread for thread in threads if thread.is_alive()]

        for pending_thread in pending:
            logger.error(
                "join-client-threads",
                logger.LogResult.fail,
                "pending-thread",
                pending_thread.name,
            )

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
            with self._lottery_lock:
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

        with self._lottery_lock:
            for bet in self.lottery.load_bets():
                if bet.agency_id == agency_id and self.lottery.has_won(bet):
                    winners.append(bet)

        return winners

import os
import signal
import sys

import logger
import server

SERVER_HOST = os.environ["SERVER_HOST"]
SERVER_PORT = int(os.environ["SERVER_PORT"])
STORAGE_PATH = os.environ.get("STORAGE_PATH", "bets_stored.csv")
AGENCY_QUORUM_MIN = int(os.environ.get("AGENCY_QUORUM_MIN", 1))


def main():
    logger.init()
    s = server.Server(SERVER_HOST, SERVER_PORT, STORAGE_PATH, AGENCY_QUORUM_MIN)
    signal.signal(signal.SIGTERM, lambda signal_number, curr_stack_frame: s.shutdown())
    try:
        s.run()
    except Exception as e:
        logger.error("server-run", logger.LogResult.fail, "err", e)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())

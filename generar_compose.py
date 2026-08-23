import sys

EXPECTED_PARAMS = 2
NOMBRE_SCRIPT = 0
CANT_CLIENTES = 1

SERVER_HOST = "server"
SERVER_PORT = 5678
ARCHIVO_SALIDA = "docker-compose.yaml"

SERVER_TEMPLATE = """\
services:
  server:
    build:
      context: ./services/server
      dockerfile: Dockerfile
    container_name: server
    environment:
      - PYTHONUNBUFFERED=1
      - SERVER_HOST={server_host}
      - SERVER_PORT={server_port}
    ports:
      - "{server_port}:{server_port}"
"""

CLIENT_TEMPLATE = """\
  client_{id}:
    build:
      context: ./services/client
      dockerfile: Dockerfile
    container_name: client_{id}
    depends_on:
      - server
    environment:
      - AGENCY_ID={id}
      - SERVER_HOST={server_host}
      - SERVER_PORT={server_port}
      - INPUT_FILE=./input/input-{id}.csv
      - OUTPUT_FILE=./output/output_file-{id}.csv
"""


def generar_compose(cant_clientes):
    contenido = [
        SERVER_TEMPLATE.format(server_host=SERVER_HOST, server_port=SERVER_PORT)
    ]
    for i in range(cant_clientes):
        contenido.append(
            CLIENT_TEMPLATE.format(
                id=i, server_host=SERVER_HOST, server_port=SERVER_PORT
            )
        )

    with open(ARCHIVO_SALIDA, "w") as archivo:
        for i, servicio in enumerate(contenido):
            archivo.write(servicio)
            if i != len(contenido) - 1:
                archivo.write("\n")


def leer_cant_clientes(argv):
    if len(argv) != EXPECTED_PARAMS:
        print(f"Uso: python3 {argv[NOMBRE_SCRIPT]} <cant_clientes>")
        return None

    try:
        cant_clientes = int(argv[CANT_CLIENTES])
    except ValueError:
        print(f"'{argv[CANT_CLIENTES]}' no es un entero valido")
        return None

    if cant_clientes < 1:
        print("La cantidad de clientes debe ser mayor a 0")
        return None

    return cant_clientes


def main():
    cant_clientes = leer_cant_clientes(sys.argv)
    if cant_clientes is None:
        return

    generar_compose(cant_clientes)


if __name__ == "__main__":
    main()

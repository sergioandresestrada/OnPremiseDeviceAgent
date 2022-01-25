import sys
import socket

if len(sys.argv) != 2:
    puerto = 9999

else:
    puerto = int(sys.argv[1])

print(f"Creando socket escuchando en el puerto {puerto}")
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

s.bind(("", puerto))

s.listen()

while True:
    print("\nEsperando por un mensaje...")
    
    sd, origen = s.accept()
    print("Se ha recibido una conexión")

    mensaje = sd.recv(1024).decode("utf-8")

    print(f"Se ha recibido el mensaje: {mensaje}")

    print("Se cierra el socket de datos")
    sd.close()
    

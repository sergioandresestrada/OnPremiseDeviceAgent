import sys
import socket

if len(sys.argv) != 2:
    puerto = 9999

else:
    puerto = int(sys.argv[1])

print(f"Creating listening socket in port: {puerto}")
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

s.bind(("", puerto))

s.listen()

while True:
    print("\nWaiting for a message...")
    
    sd, origen = s.accept()
    print("A connection has been received")

    mensaje = sd.recv(1024).decode("utf-8")

    print(f"Message received: {mensaje}")

    print("Closing data socket")
    sd.close()
    

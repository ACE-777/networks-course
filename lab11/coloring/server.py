import socket
import pickle
import threading

import tkinter as tk

server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server_socket.bind(('localhost', 9999))
server_socket.listen(1)

client_socket, client_address = server_socket.accept()
print(f"Establish connection with {client_address}")

window = tk.Tk()
window.title("Server")

canvas = tk.Canvas(window, width=500, height=500)
canvas.pack()


def draw_line(data):
    x1, y1, x2, y2 = data
    canvas.create_line(x1, y1, x2, y2, fill="black")


def receive_data():
    while True:
        try:
            data = client_socket.recv(1024)
            if data:
                data = pickle.loads(data)
                draw_line(data)
        except Exception as e:
            print(e)

            break


receive_thread = threading.Thread(target=receive_data)
receive_thread.daemon = True
receive_thread.start()

window.mainloop()

client_socket.close()
server_socket.close()

import socket
import pickle

import tkinter as tk

client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
client_socket.connect(('localhost', 9999))

window = tk.Tk()
window.title("Client")

canvas = tk.Canvas(window, width=500, height=500)
canvas.pack()

x_prev, y_prev = None, None


def send_data(event):
    global x_prev, y_prev
    x, y = event.x, event.y
    if x_prev and y_prev:
        data = (x_prev, y_prev, x, y)
        client_socket.send(pickle.dumps(data))
        draw_line((x_prev, y_prev, x, y), 'green')

    x_prev, y_prev = x, y


canvas.bind("<B1-Motion>", send_data)


def start_new_line(event):
    global x_prev, y_prev
    x_prev, y_prev = None, None


canvas.bind("<ButtonPress-1>", start_new_line)


def draw_line(data, color):
    x1, y1, x2, y2 = data
    canvas.create_line(x1, y1, x2, y2, fill=color)


window.mainloop()

client_socket.close()

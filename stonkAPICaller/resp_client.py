import asyncio
import json
import pickle
import socket
import sys
import time

import schedule as schedule

# from main import r


# from resp import Client

def start_resp(redis_client):
    # how it is supposed to run
    # loop = asyncio.get_event_loop()
    # client = ClientMain()
    # loop.run_until_complete(client.main_client('$pyredis>>', loop))
    # loop.close()
    #
    schedule.every(1).minutes.do(start_sending_responses, redis_client=redis_client)
    start_time = time.time()
    while 1:
        schedule.run_pending()
        time.sleep(1)
        if time.time() - start_time > 3600:
            break
    print("PROGRAM HAS STOPPED")


def encode(data, simple_str=False):
    if isinstance(data, ValueError):
        return f'-{str(data)}\r\n'
    elif isinstance(data, str) and simple_str:
        return f'+{data}\r\n'
    elif isinstance(data, int):
        return f':{data}\r\n'
    elif isinstance(data, str):
        return f'${len(data)}\r\n{data}\r\n'
    elif isinstance(data, list) or isinstance(data, tuple):
        enc = f'*{len(data)}\r\n'
        for itm in data:
            enc += encode(itm)
        return enc


def my_decode(data, r):
    marker = data[0]
    if marker == '-':
        return str(data[1:-2])
    elif marker == '+':
        return str(data[1:-2])
    elif marker == ':':
        return int(data[1:-2])
    elif marker == '$':
        parts = data[1:].split('\r\n')
        ln = int(parts[0])
        if ln == -1:
            return None
        elif ln == 0:
            return ''
        else:
            return str(parts[1])
    elif marker == '*':
        parts = data[1:].split('\r\n')
        ln = int(parts[0])
        items = [None] * ln
        for i in range(1, 1 + ln):
            items[i - 1] = my_decode(parts[i - 1], r)
        return items
    else:
        raise Exception('resp decoding error')


class ClientMain:
    def __init__(self):
        self.client = None

    async def main_client(self, prompt, evt_loop, redis_client):
        print(f'{prompt} Welcome to pyRedis')
        running = True
        cl = Client(evt_loop, 0, "hi")
        self.client = cl
        while running:
            cmd = input(f'{prompt} ')
            if cmd == '.exit':
                running = False
            else:
                try:
                    reply = await cl.execute(cmd, redis_client)
                    print(prompt, reply)
                except Exception as e:
                    print(prompt, 'errr:', e)
        cl.close()


def start_sending_responses(redis_client):
    client = Client(None, 6379, "127.0.0.1")
    # client.execute(".connect 127.0.0.1:6379")
    # print("WORKED CONNECTING")
    # client.execute("get hi")
    # print("WORKED SET")
    # client.execute(".exit")
    # client.execute("ping")
    # client.execute("arrset hello 7", redis_client)
    # client.execute("arrget hello", redis_client)
    # client.execute("arrset hello time random 2", redis_client)
    client.execute("arrget hello time random", redis_client)
    print("I'm working...")


class Client:
    def __init__(self, evt_loop, port_num, ip):
        self.event_loop = evt_loop
        self.reader = None
        self.writer = None
        self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.socket.connect((ip, port_num))

    # def send(self, message):
    #     self.socket.sendto(message.encode(), ('127.0.0.1', 6379))

    # supposed to be asynced
    def execute(self, cmd, redis_client):
        # if self.socket is None:
        #     if cmd.startswith('.connect'):
        #         cmd = cmd.replace('.connect', '').strip()
        #         parts = cmd.split(':')
        #         # Supposed to be an await here
        #         # and , loop=self.event_loop
        #         #     s.send(bytes((encode(["set", "others", "7"])), encoding='utf8'))
        #         # (r, w) = await asyncio.open_connection(parts[0], parts[1], loop=self.event_loop)
        #         # self.reader = r
        #         # self.writer = w
        #     else:
        #         return 'Please connect with <.connect ip:port>'
        # else:
        cmd_split = cmd.split(' ')
        print(cmd_split)
        byte_string_cmd = bytes((encode(cmd_split)), encoding='utf8')
        self.socket.send(byte_string_cmd)
        # self.writer.write(encode(cmd_split).encode())
        # Supposed to be an await here
        # data = self.reader.read(1024)
        data1 = self.socket.recv(4096)
        print("Response: ", my_decode(data1.decode(), redis_client))
        return my_decode(data1.decode(), redis_client)

    def close(self):
        if self.writer is not None:
            self.writer.close()

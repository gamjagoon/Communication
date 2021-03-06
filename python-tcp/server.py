import socket, json

HOST = ""
PORT = 0

json_file = open('address.json')
json_data = json.load(json_file)
HOST = json_data["HOST"]
PORT = json_data["PORT"]
json_file.close()
# 소켓 객체를 생성합니다. 
# 주소 체계(address family)로 IPv4, 소켓 타입으로 TCP 사용합니다.  
server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
server_socket.bind((HOST, PORT))
server_socket.listen()
print("start server!!")
client_socket, addr = server_socket.accept()
print('Connected by', addr)

while True:
    data = client_socket.recv(128)

    # 데이터가 없거나 exit 면 종료
    if not data or data == "exit":
        break
    print('Received from', addr, data.decode())
    client_socket.sendall("ok".encode())
# 소켓을 닫습니다.
client_socket.close()
server_socket.close()
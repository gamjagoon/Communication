import socket,json
HOST = "" 
PORT = 0
json_file = open('address.json')
json_data = json.load(json_file)
HOST = json_data["HOST"]
PORT = json_data["PORT"]
json_file.close()
# 소켓 객체를 생성합니다. 
# 주소 체계(address family)로 IPv4, 소켓 타입으로 TCP 사용합니다.  
client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
print("make object ")
client_socket.connect((HOST, PORT))
print("start client")

while True:
    data = input()
    client_socket.sendall(data.encode())
    # 데이터가 없거나 exit 면 종료
    if not data or data == "exit":
        break
    data = client_socket.recv(128)
    print('Received', repr(data.decode()))

# 소켓을 닫습니다.
client_socket.close()

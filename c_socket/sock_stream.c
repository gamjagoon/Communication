#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/un.h>

#define SOCK_PATH "sock_stream_un"

static void print_usage(const char *progname)
{
    printf("%s (server | client)\n",progname);
}

static int do_server(void)
{
    int sock = socket(AF_UNIX, SOCK_STREAM, 0);
    struct sockaddr_un addr;
    char buf[1024];
    if(sock == -1){
        perror("socket()");
        return -1;
    }
    memset(&addr, 0, sizeof(addr));
    addr.sun_family = AF_UNIX;
    strncpy(addr.sun_path, SOCK_PATH,sizeof(addr.sun_path) - 1);

    int ret = bind(sock, (struct sockaddr*)&addr, sizeof(struct sockaddr_un));
    if(ret == - 1){
        perror("bind()");
        close(sock);
        return -1;
    }

    listen(sock,5);

    ret = accept(sock,NULL, NULL);
    if(ret == - 1){
        perror("accept()");
        close(sock);
        return -1;
    }

    memset(buf,0,sizeof(buf));
    size_t err = recv(ret, buf, sizeof(buf), 0 );
    if(err == -1){
        perror("recv()");
        close(sock);
        return -1;
    }
    printf("client said [%s]\n",buf);
    close(ret);
    close(sock);

    return 0;
}

static int do_client(void)
{
    int sock =socket(AF_UNIX,SOCK_STREAM,0);
    if(sock == -1){
        perror("socket()");
        return -1;
    }
    struct sockaddr_un addr;
    memset(&addr, 0, sizeof(addr));
    addr.sun_family = AF_UNIX;
    strncpy(addr.sun_path, SOCK_PATH,sizeof(addr.sun_path) - 1);
    int err = connect(sock, (struct sockaddr* )&addr, sizeof(struct sockaddr_un));
    if(err == -1){
        perror("connect()");
        return -1;
    }
    char buf[128];
    memset(buf,0,sizeof(buf));
    snprintf(buf,sizeof(buf),"this is msg from sock_strea");
    send(sock, buf, sizeof(buf), 0);

    close(sock);
    return 0;
}

int main(int argc, char **argv)
{
    int ret;
    if (argc < 2){
        print_usage(argv[0]);
        return -1;
    }

    if (!strcmp(argv[1], "server")){
        ret = do_server();
    } else if(!strcmp(argv[1],"client")){
        ret = do_client();
    } else {
        print_usage(argv[0]);
        return -1;
    }
    return ret;
}
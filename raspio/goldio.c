#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <sys/types.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/mman.h>
#include <wiringPi.h>
#include <time.h>
#include <signal.h>

#define MAPSIZE 20
#define FORMAT "/gosprints%d"

struct Pin {
    int port_num;
    char *map;
    int fd;
    int cycle;
    char val;
    int dist;
};

struct Pin *pins;
char name[] = FORMAT;
int len = 2;
int _H = HIGH;
int _L = LOW;

char simulate=0;
int simulate_wait_secs = 3;
char simulate_wait_now = 0;

const char* shm_name(int i) {
    sprintf(name, FORMAT, i);
    return name;
}

void cleanup() {
    int i;
    for(i=0; i<len; i++) {
        munmap(pins[i].map, MAPSIZE);
        shm_unlink(shm_name(i));
    }
    free(pins);
    printf("closed.\n");

    exit(EXIT_SUCCESS);
}

void prevent_false_start() {
    char wait_str[11];
    char *env = getenv("GOSPRINTS_GOLDIO_WAIT");
    if(env != NULL) {
        strncpy(wait_str, env, 10); 
        sscanf(wait_str, "%d", &simulate_wait_secs);
    }
    simulate_wait_now=1;
}

void reset() {
    if(simulate) {
        prevent_false_start();
    }
    int i;
    for(i=0; i<len; i++) {
        pins[i].dist=0;
        pins[i].cycle=0;
        memset(pins[i].map, '\0', sizeof(*(pins[i].map)));
        sprintf((char*) pins[i].map, "%d", pins[i].dist);
    }
    printf("reset\n");
}



void init(char *ports, const char pull_up) {
    char* token;
    int port_num, i;
    struct Pin* pin;
    for (i=0; (token = strtok(ports, ",")) != NULL; i++) {
        if(strlen(token) > 0) {
            port_num = atoi(token);
            if(port_num>0) {
                if(!simulate)
                    pinMode(port_num, INPUT);
                if(pull_up && !simulate) {
                    pullUpDnControl(port_num, PUD_UP);
                }
                pins = realloc(pins, (i+1)*sizeof(*pin));
                pins[i].port_num = port_num;
                pins[i].fd = shm_open(shm_name(i), O_TRUNC|O_RDWR|O_CREAT, 0666);
                pins[i].map = mmap(NULL, MAPSIZE, PROT_WRITE, MAP_SHARED, pins[i].fd, 0);

                ftruncate(pins[i].fd, MAPSIZE);
            }
        }
        ports = NULL;
    }
    len=i;
}

void race(char pull_up, char threshold, int wait) {
    if(pull_up) {
        _H = LOW;
        _L = HIGH;
    }

    char val=_L;
    int i;

    while(1) {
        for(i=0; i<len; i++) {
            val = digitalRead(pins[i].port_num);
            if(val!=pins[i].val && val==_H) {
                pins[i].cycle++;
                if(pins[i].cycle==threshold) {
                    pins[i].dist++;
                    pins[i].cycle=0;
                    sprintf((char*) pins[i].map, "%i", pins[i].dist);
                }
            }

            pins[i].val = val;
        }
        if(wait)
            usleep(wait);
    } 
}

void simulation(char threshold, int wait) {
    int i;
    while(1) {
        if(simulate_wait_now) {
            sleep(simulate_wait_secs);
            simulate_wait_now=0;
        }
        for(i=0; i<len; i++) {
            pins[i].dist+=(int)random()%threshold;
            sprintf((char*) pins[i].map, "%i", pins[i].dist);
            /* fprintf(stderr, "updating #%d with distance: %d\n", i, pins[i].dist); */
        }
        if(wait)
            usleep(random()%wait);
    }
}

void usage(const char exec[]) {
    fprintf(stderr, "Usage: %s [-p] [-s] [-w <sleep usec>] [-t <threshold>] <port 1>[,<port 2>,..]\n", exec);
    exit(EXIT_FAILURE);
}

void handlesig(int sig_num, void f()) {
    if(signal(sig_num, f)==SIG_ERR) {
        fprintf(stderr, "cant catch signal %d\n", sig_num);
        exit(EXIT_FAILURE);
    }
}

int main(int argc, char *const* argv)
{
    unsigned char pull_up=0, threshold=1;
    int opt, wait=0;

    while ((opt = getopt(argc, argv, "pst:w:")) != -1) {
        switch (opt) {
            case 'p':
                pull_up = 1;
                break;
            case 't':
                threshold = atoi(optarg);
                break;
            case 'w':
                wait = atoi(optarg);
                break;
            case 's':
                simulate=1;
                break;
            default: usage(argv[0]);
            }
    }

    if(optind >= argc) 
        usage(argv[0]);

    handlesig(SIGINT, cleanup);
    handlesig(SIGTERM, cleanup);
    handlesig(SIGHUP, reset);
    handlesig(SIGABRT, reset);

    if(!simulate)
        wiringPiSetupPhys();

    init(argv[optind], pull_up);

    if(simulate)
        simulation(threshold, wait);
    else
        race(pull_up, threshold, wait);

    cleanup();
    return 0;
}

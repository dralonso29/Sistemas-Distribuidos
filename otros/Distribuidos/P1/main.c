//Eva Gordo Calleja
//Grado telemática

#include <stdio.h>
#include "threads.h"
#include <time.h>
#include <stdlib.h>
#include <unistd.h>

int glob;

void f1(void *arg){
	int i;
	char *arg1 = arg;
	for (i=0; i<10; i++){
		usleep(100000);
		glob++;
		printf("f1 %i %s\n", glob, arg1);
		yieldthread();
	}
	printf("EXIT 1\n");
	exitsthread();
}


void f2(){
	int i;
	for (i=0; i<10; i++){		
		usleep(100000);
		glob++;
		printf("f2 %i\n", glob);
		yieldthread();
	}
	printf("EXIT 2\n");
	exitsthread();
}

void f3(){
	int i;
	for (i=0; i<10; i++){
		usleep(100000);
		glob++;
		printf("f3 %i\n", glob);
		yieldthread();
	}
	printf("EXIT 3\n");
	exitsthread();
}

int main(int argc,char **argv){
	int i;
	char *arg1;
	arg1 = malloc(sizeof(char));
	arg1 = "ARGUMENTO";
	initthreads();
	createthread(f1, arg1, 1024*4);
	createthread(f2, NULL, 1024*4);
	createthread(f3, NULL, 1024*4);

	for(i=0; i < 20; i++){
		usleep(100000);
		glob++;
		printf("im am main %d\n",  glob);
		yieldthread();
	}
	printf("EXIT MAIN\n");
	
	exitsthread();
	return 0;
}

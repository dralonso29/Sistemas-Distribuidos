//Jorge Vela Peña
//Grado en Ingeniería Telemática


#include <stdio.h>
#include "threads.h"
#include <time.h>


#include <unistd.h>

void f2();
void f3();
void f4();


void f2(){
	int i;
	for(i=0; i < 20; i++){
		printf("F2 %d\n", i);
		usleep(100000);
		yieldthread();
	}
	printf("FIN F2-------------------------------------------------------------\n");
	exitsthread();
}

void f3(){
	int i;
	for(i=0; i < 20; i++){
		printf("F3 %d\n", i);
		usleep(100000);
		if(i==4){
			createthread(f4, NULL ,64*1024);
		}
		yieldthread();
	}
	printf("FIN F3-------------------------------------------------------------\n");
	exitsthread();
}

void f4(){
	int i;
	for(i=0; i < 20; i++){
		printf("F4 %d\n", i);
		usleep(100000);
		yieldthread();
	}
	printf("FIN F4-------------------------------------------------------------\n");
	exitsthread();
}

int main (int argc,char **argv)
{
	int i;
	initthreads();
	createthread(f2, NULL ,64*1024);
	createthread(f3, NULL ,64*1024);

	for (i=0; i<20;i++){
		printf("Main %d\n", i);
		yieldthread();
	}
	printf("%s\n", "Main FIN ------------------------------------------------------");

	exitsthread();
	return 0;
}
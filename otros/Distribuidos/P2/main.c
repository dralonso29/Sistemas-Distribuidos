//Eva Gordo Calleja
//Grado en Ingeniería Telemática


#include <stdio.h>
#include "threads.h"
#include <time.h>
#include <unistd.h>
#include <stdlib.h>

void f1();
void f2();
void f4();
void f5();

void f1(){
	int i;
	//sleepthread(2000);

	for(i=0; i < 20; i++){
		printf("F1 %d\n", i);
		usleep(100000);
		
		if(i==5){
			suspendthread();
			//aux();
			//sleep(10);
		}
		
		if(i == 7){
			aux();
			sleep(5);		
		}
		
		yieldthread();
		/*
		if(i==16){
			resumethread(2);
		}
		*/
	}
	//sleepthread(2000);
	exitsthread();
}

void f2(){
	int i;
	//suspendthread();

	for(i=0; i < 20; i++){
		printf("F2 %d\n", i);
		usleep(100000);
		
		if(i==10){
			resumethread(1);
			//killthread(3);
			//createthread(f4, NULL ,64*1024);
			//sleepthread(2000);
		}
		//killthread(1);
		yieldthread();
	}

	exitsthread();
}
/*
void f5(){
	for(int i=0; i < 20; i++){
		printf("F5 \n");
		usleep(100000);
		yieldthread();
	}
	printf("FIN F5-------------------------------------------------------------\n");
	exitsthread();
}
*/
void f3(){
	int i;

	//killthread(1);

	for(i=0; i < 20; i++){
		printf("F4 %d\n", i);
		usleep(100000);
		yieldthread();
	}
	//suspendthread();
	printf("FIN F4-------------------------------------------------------------\n");
	exitsthread();
}

int main (int argc,char **argv)
{
	int i, numero;
	initthreads();

	//int *a;


	createthread(f1, NULL ,64*1024);
	createthread(f2, NULL ,64*1024);
	createthread(f3, NULL ,64*1024);
	//createthread(f5,NULL,64*1024);

	
	for (i=0; i<20;i++){
		printf("Main %d\n", i);
		usleep(100000);
		yieldthread();
	}

	
	//numero = suspendedthreads(&a);

	//printf("NUUUUMEROOOOOSSS %d\n", numero);
	/*
	for(int j=0;j<numero;j++){
		printf("%d ", a[j]);
	}

	printf("%s\n", "Main FIN ------------------------------------------------------");
	
	free(a);
	*/
	exitsthread();
	return 0;
}

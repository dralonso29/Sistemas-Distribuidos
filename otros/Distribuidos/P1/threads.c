//Eva Gordo Calleja
//Grado telemática
//gcc -c -g -Wall -Wshadow main.c && gcc -c -g -Wall -Wshadow threads.c && gcc -o main threads.o main.o
//./main
#include <ucontext.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/time.h>
#include <err.h>
#include <unistd.h>


enum {
	sizeStack=64*1024,
	maxThreads = 5
};

enum status{
	Ending,
	Using,
	Waiting
};

int idThread = 1;
double toMilisg();
int current = 0;

struct thread {   
	int id;
	char *stack;     
	double milsg;
	enum status estado;
	ucontext_t context;
};

static struct thread arrayThread[maxThreads];

void initthreads(void){
	toMilisg();
	char *puntero;
	puntero = malloc(sizeStack);
	arrayThread[0].id = idThread;
	arrayThread[0].stack = puntero;
	arrayThread[0].milsg = toMilisg();
	arrayThread[0].estado = Using;
}

int createthread(void (*mainf)(void*), void *arg, int stacksize){
	int i;
	char *punteroStack;
	punteroStack = malloc(stacksize);

	for(i=1;i<maxThreads; i++){
		if(arrayThread[i].estado==Ending){
			//printf("CREANDO HILO NÚMERO: %d\n",i);
			idThread++;
			getcontext(&(arrayThread[i].context));
			arrayThread[i].id=idThread;
			arrayThread[i].milsg = toMilisg();
			arrayThread[i].estado = Waiting;
			arrayThread[i].stack = punteroStack;
			
			//hacer malloc porque nos va a dar la dirección de donde empieza esa memoria reservada.
			arrayThread[i].context.uc_stack.ss_sp = punteroStack;
			arrayThread[i].context.uc_stack.ss_size = stacksize;
			arrayThread[i].context.uc_link = NULL;
			makecontext(&(arrayThread[i].context), (void(*)(void))mainf, 1,arg);
			break;
		}
	}
	return arrayThread[i].id;
}

int nextposition(void){
	for(int i=current+1; i%maxThreads != current; i++){
		if(arrayThread[i%maxThreads].estado==Ending && arrayThread[i%maxThreads].context.uc_stack.ss_size != 0){
			arrayThread[i%maxThreads].context.uc_stack.ss_size = 0;
			free(arrayThread[i%maxThreads].stack);
		}
		
		if(arrayThread[i%maxThreads].estado==Waiting){
			arrayThread[i%maxThreads].estado = Using;
			arrayThread[i%maxThreads].milsg = toMilisg();
			return i%maxThreads;
		}
	}
	return -1;
}

int gestor(void){
    	if(toMilisg() - arrayThread[current].milsg < 200){
		printf("TIEMPO %g\n",toMilisg() - arrayThread[current].milsg);
		return current;
   	}else{
		return nextposition();
    	}
}

void yieldthread(void){
	int nextThread = 0;
	int present=current;

	if(arrayThread[current].estado != Using){
		err(1, "PROBLEMAS EN YIELD");
	}else{
		nextThread = gestor();
		if(nextThread != -1){
			arrayThread[current].estado = Waiting;
			current = nextThread;
			if(swapcontext(&(arrayThread[present].context),&(arrayThread[nextThread].context))==-1){
				err(1, "PROBLEMAS EN SWAP YIELD");
			}
		}
	}
}

void exitsthread(void){
	int nextThread = 0;
	int present = current;

	nextThread = nextposition();

	if(nextThread == -1){
		exit(0);
	}

	arrayThread[current].estado = Ending;
	arrayThread[nextThread].estado = Using;
	current = nextThread;
	if(swapcontext(&(arrayThread[present].context),&(arrayThread[nextThread].context))==-1){
		err(1, "PROBLEMAS EN SWAP EXIT");
	}
}

double toMilisg(void){
	struct timeval t;
	double doneMil;
    	
    	if (gettimeofday(&t, NULL) < 0 )
    		return 0.0;
	
	doneMil = (t.tv_usec + t.tv_sec * 1000000.0);

    	return doneMil;
}

int curidthread(void){
	int i;
	for(i=0;i<maxThreads; i++){
		if(arrayThread[i].estado==Using){
			return arrayThread[i].id;
		}
	}
	return 0;	
}


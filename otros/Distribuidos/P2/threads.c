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



















//CAMBIIIAAARRRR A 32




enum {
	sizeStack=64*1024,
	maxThreads = 4
};

enum status{
	Ending,
	Using,
	Waiting,
	Suspending,
	Sleeping
};

int idThread = 1;
double toMilisg();
int current = 0;

struct thread {   
	int id;
	char *stack;     
	double milsg;
	double milsgSleep;
	enum status estado;
	ucontext_t context;
};

static struct thread arrayThread[maxThreads];

void initthreads(void){
	toMilisg();
	char *puntero;
	puntero = malloc(sizeStack);
	arrayThread[0].id = 0;
	arrayThread[0].stack = puntero;
	arrayThread[0].milsg = toMilisg();
	arrayThread[0].estado = Using;
	arrayThread[0].milsgSleep = 0;

}

int createthread(void (*mainf)(void*), void *arg, int stacksize){
	int i;
	char *punteroStack;
	punteroStack = malloc(stacksize);

	for(i=1;i<maxThreads; i++){
		if(arrayThread[i].estado==Ending){
			//printf("CREANDO HILO NÚMERO: %d\n",i);
			
			getcontext(&(arrayThread[i].context));
			arrayThread[i].id=idThread;
			arrayThread[i].milsg = toMilisg();
			arrayThread[i].estado = Waiting;
			arrayThread[i].stack = punteroStack;
			arrayThread[i].milsgSleep = 0;
			
			//hacer malloc porque nos va a dar la dirección de donde empieza esa memoria reservada.
			arrayThread[i].context.uc_stack.ss_sp = punteroStack;
			arrayThread[i].context.uc_stack.ss_size = stacksize;
			arrayThread[i].context.uc_link = NULL;
			makecontext(&(arrayThread[i].context), (void(*)(void))mainf, 1,arg);
			idThread++;
			break;
		}
	}
	return arrayThread[i].id;
}

int nextposition(void){
	for(int i=current+1; i%maxThreads != current; i++){
		if(arrayThread[i%maxThreads].estado==Waiting){
			arrayThread[i%maxThreads].estado = Using;
			arrayThread[i%maxThreads].milsg = toMilisg();
			return i%maxThreads;
		}
		if(arrayThread[i%maxThreads].estado==Ending && arrayThread[i%maxThreads].context.uc_stack.ss_size != 0){
			arrayThread[i%maxThreads].context.uc_stack.ss_size = 0;
			free(arrayThread[i%maxThreads].stack);
		}
		
		//new
		/*
		if(arrayThread[i%maxThreads].estado==Sleeping){
			if(toMilisg() - arrayThread[current].milsg > arrayThread[current].milsgSleep){
				arrayThread[i].estado = Waiting;
			}
		}
		*/
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

void aux(){
	int i;
	int devuelveEstado = 0;
	int devuelveID = 0;
	for(i=0;i<maxThreads; i++){
		devuelveID = arrayThread[i].id;
		devuelveEstado = arrayThread[i].estado;
		printf("pos_a:  %d ID %d ESTADO %d \n",i,devuelveID, devuelveEstado);
	}
}

void yieldthread(void){
	int nextThread = 0;
	int present=current;

	//aux();
	printf("CURRENT %d\n", current);
	if(arrayThread[current].estado != Using){
		fprintf(stderr, "%s\n", "Error: yieldthread");
	}else{
		nextThread = gestor();
		if(nextThread != -1){
			arrayThread[current].estado = Waiting;
			//arrayThread[current].milsg = toMilisg();
			current = nextThread;
			if(swapcontext(&(arrayThread[present].context),&(arrayThread[nextThread].context))==-1){
				fprintf(stderr, "%s\n", "Error: yieldthread. Swap problems");
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
	
	//si el nexthread es = que el actual y tengo algun hilo dormido, tengo que esperar el tiempo que ese hilo este dormido.
	//si solo tiene un hilo y hay alguno suspendido, tiene que terminar el programa  y notificar que tiens hilos suspendidos.

	if(swapcontext(&(arrayThread[present].context),&(arrayThread[nextThread].context))==-1){
		fprintf(stderr, "%s\n", "Error: exitsthread. Swap problems");
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

void suspendthread(void){
	int nextThread = 0;
	int present = current;
	nextThread = nextposition();
	
	if(arrayThread[current].estado==Using){
		if(nextThread == -1){
			fprintf(stderr, "%s\n", "Error: suspendthread. No threads Waiting");
		}else{
			arrayThread[current].estado = Suspending;
			arrayThread[nextThread].estado = Using;
			current = nextThread;
			if(swapcontext(&(arrayThread[present].context),&(arrayThread[nextThread].context))==-1){
				fprintf(stderr, "%s\n", "Error: suspendthread. Swap problems");
			}
		}
	}else{
		fprintf(stderr, "%s\n", "Error: suspendthread. No using");
	}
}

int resumethread(int id){
	int i;
	for (i=0; i<maxThreads; i++){
		printf("NUMERITOS %d\n", arrayThread[i].id);
		if(arrayThread[i].id == id && arrayThread[i].estado == Suspending){
			printf("NUMERITOS SUSPENDING %d\n", arrayThread[i].id);
			arrayThread[i].estado = Waiting;
			aux();
			sleep(5);
			return 0;
		}
	}
	return -1;
}

int killthread(int id){
	int i;

	for(i=0; i<maxThreads; i++){
		if(arrayThread[i].id==id){
			if (arrayThread[i].estado == Ending){
				return -1;			
			}else if (arrayThread[i].estado==Using){
				exitsthread();
				return 0;
			}else{
				arrayThread[i].estado = Ending;
				return 0;
			}
		}
	}
	return -1;
}

/*
int suspendedthreads(int **list){
	int i, j = 0;
	int suspend = 0;

	for(i=0; i < maxThreads; i++){
		if(arrayThread[i].estado == Suspending){
			suspend++;
		}
	}

	if(suspend > 0){
		*list = malloc(sizeof(int)*suspend);
		for(i=0; i <maxThreads; i++){ 
			if(arrayThread[i].estado == Suspending){
				*list[j]= arrayThread[i].id;
				j++;
			}
		}
	}

	return suspend;

}




//SLEEP: IF(toMilisg() - arrayThread[current].milsg > arrayThread[current].milsgSleep){arrayThread[i].estado = Waiting;}

void sleepthread(int msec){
      int i;
      for(i=0; i <maxThreads; i++){
          if(arrayThread[i].estado == Using){
              arrayThread[i].estado = Sleeping;
              arrayThread[i].milsg = toMilisg();
              arrayThread[i].milsgSleep = msec;
              gestor();
          }
      }
}
*/

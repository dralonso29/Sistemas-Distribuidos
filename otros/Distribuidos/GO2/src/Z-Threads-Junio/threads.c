//Jorge Vela Peña
//Grado en Ingeniería Telemática

#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <ucontext.h>
#include <sys/time.h>

enum{
	sizeStack=64*1024,
	MaxThreads = 32
};


enum {
    Free=0,
    Ready=1, 
    Delete=2, 
    Using = 3
};


static int idthread=0;

static long long initTime();
static void scheduler (int j);


static struct thread{
    int id;
    long long t;
    char* stack;
    ucontext_t context;
    int usoThread;
}threadInit;

static struct thread arrayTh[31];


void
initthreads(void) {
    initTime();
    threadInit.id=idthread;
    threadInit.t=initTime();
    threadInit.stack = malloc(sizeStack);
    threadInit.usoThread = Using;
    idthread++;
    /*  */
    threadInit.context.uc_stack.ss_sp = threadInit.stack;
    threadInit.context.uc_stack.ss_size = sizeStack;
    threadInit.context.uc_link = NULL; 
    arrayTh[0]= threadInit;
}


int
createthread(void (*mainf)(void*), void *arg, int stacksize){
    int id;
    id = -1;
    int i;
    for(i=0; i<MaxThreads; i++){
        if(arrayTh[i].usoThread == Free || arrayTh[i].usoThread == Delete){
            struct thread newthread;
            getcontext(&(newthread.context));

            newthread.id=idthread;
            id=newthread.id;
            newthread.usoThread = Ready;
            idthread++;

            newthread.t=initTime();

            if(arrayTh[i].usoThread == Delete){
                free(arrayTh[i].stack);
                newthread.stack = malloc(stacksize);
            }else{
                newthread.stack = malloc(stacksize);
            }

            newthread.context.uc_stack.ss_sp = newthread.stack;
            newthread.context.uc_stack.ss_size = stacksize;
            newthread.context.uc_link = NULL;

            arrayTh[i]= newthread;
            makecontext(&(arrayTh[i].context), (void(*)(void))mainf, 1,arg);
            break;
        } 
      }
    return id;  
}

void
yieldthread(void){
    int i;
    for (i=0;i<MaxThreads;i++){
        if(arrayTh[i].usoThread==Using){
            initTime();
            if(initTime() - arrayTh[i].t < 200){
                break;
            }else{
                arrayTh[i].usoThread= Ready;
                scheduler(i);
                break;
            }
        }
    }
}

static void
funcDelete(int i){
    free(arrayTh[i].stack);
    arrayTh[i].id=0;
    arrayTh[i].usoThread=Free;
    arrayTh[i].t=0;
}
static void
funcScheduler(int i, int j){
    arrayTh[i].usoThread = Using;
    arrayTh[i].t= initTime();
    swapcontext(&(arrayTh[j].context), &(arrayTh[i].context));
}

static void
scheduler (int j){
    int i;
    int changeSwap;
    changeSwap=0;


	int p= j+1;
	int q= MaxThreads;

	for(i = p; i < q;i++){
        if(arrayTh[i].usoThread == Delete){
           funcDelete(i);
        }
        if(arrayTh[i].usoThread == Ready){
          funcScheduler(i,j);
          break;
        }
		if(i==MaxThreads-1){
			printf("ENTRA");
			p=0;
			q=j;
		}
    }

    /*for(i = j+1; i < MaxThreads;i++){
        if(arrayTh[i].usoThread == Delete){
           funcDelete(i);
        }
        if(arrayTh[i].usoThread == Ready){
          funcScheduler(i,j);
          changeSwap=1;
          break;
        }
    }
    for(i=0; i<j; i++){
        if(arrayTh[i].usoThread == Delete){
           funcDelete(i);
        }

        if(arrayTh[i].usoThread == Ready && changeSwap==0){
          funcScheduler(i,j);
          break;
        }
    }*/
}

void
exitsthread(void){
    int i;
    for (i=0;i<MaxThreads;i++){
        if (arrayTh[i].usoThread == Using){
            arrayTh[i].usoThread=Delete;
            scheduler(i);
        }
    }
}

int
curidthread(void){
    int i;
    for (i=0;i<MaxThreads;i++){
        if(arrayTh[i].usoThread==Using)
            return arrayTh[i].id;
        }
    return 0;
}


static long long 
initTime(){
    struct timeval t_ini;
    gettimeofday(&t_ini, NULL);
    long long int r= (((long long)(t_ini.tv_usec*0.001))+ ((long long)t_ini.tv_sec)*1000);
    return r;
}

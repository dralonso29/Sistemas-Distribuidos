/*#include <stdio.h>
#include <sys/time.h>
int main(void)
{
    struct timeval ti, tf;
    double tiempo;
    gettimeofday(&ti, NULL);   // Instante inicial
    printf("Lee este mensaje y pulsa ENTER\n");
    getchar();
    gettimeofday(&tf, NULL);   // Instante final
    tiempo= (tf.tv_sec - ti.tv_sec)*1000 + (tf.tv_usec - ti.tv_usec)/1000.0;
    printf("Has tardado: %g milisegundos\n", tiempo);
}
*/


void aux(){
	int i;
	int devuelveEstado = 0;
	int devuelveID = 0;
	for(i=0;i<maxThreads; i++){
		devuelveID = arrayThread[i].id;
		devuelveEstado = arrayThread[i].estado;
		printf("ID %d ESTADO %d \n", devuelveID, devuelveEstado);
	}
}

void exitsthread(void){
	int nextThread = 0;
	int actual = current;

	nextThread = siguiente();

	printf("NEXT THREAD %d\n", nextThread);
	arrayThread[current].estado = Ending;
	arrayThread[nextThread].estado = Using;
	current = nextThread;
	if(current == -1){
		exit(0);
	}
	swapcontext(&(arrayThread[actual].context),&(arrayThread[nextThread].context));
	/*for(i=current+1; i%maxThreads != current; i++){
		if(arrayThread[i].estado==Waiting){		
			arrayThread[i].estado=Ending;
			nextid = gestor();
			arrayThread[nextid].estado=Using;
			current = nextid;
			
			printf("ME QUIERO MORIR %d\n", i);
			swapcontext(&(arrayThread[actual].context),&(arrayThread[nextid].context));
			
		}
	}
	*/
}




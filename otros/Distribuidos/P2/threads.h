//Eva Gordo Calleja
//Grado telemática

void initthreads(void);
int createthread(void (*mainf)(void*), void *arg, int stacksize);
void exitsthread(void);
void yieldthread(void);
int curidthread(void);
void aux();
void suspendthread(void);
int resumethread(int id);
int suspendedthreads(int **list);
int killthread(int id);
void sleepthread(int msec);
void aux();

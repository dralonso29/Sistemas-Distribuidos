## Practica final 2019-2020 ##
Realizada por Daniel Pulido y David Rodríguez

- Pregunta:  	¿Puedes justificar, con datos objetivos (a partir de las simulaciones) por qué una estrategia de write
		back + invalidación es más eficiente que una de write through con propagación de cambios?

- Respuesta: 	A partir de los resultados obtenidos gracias a la realizacion de la adicional 2, podemos observar que es
		la suma de tiempo invertida en invalidar otra cache más en actualizar esa misma linea en la otra cache es
		menor que actualizar la memoria principal con el nuevo valor de la linea. (adjuntamos fichero con los valo_
		res tomados en fichero: "Estadisticas.txt")
		
		----STATS---------
		En invalidar otra línea de una caché diferente tardo  7 u.t en promedio
		En actualizar esa línea en otra caché diferente tardo  12 u.t en promedio
		en actualizar la memoria principal con el nuevo valor de una línea tardo  32 u.t en promedio
		DATOS TOMADOS CON  10 EJECUCIONES DEL PROGRAMA
		-------------------

- Pregunta:	Calcula cuántas peticiones de actualización tiene que servir cada estadio al resto. Si esto fuese un
		sistema de coherencia entre cachés, y en el paso directo de mensajes entre cachés tardase la mitad de
		tiempo por petición que al actualizar los datos de caché a memoria principal (el “servicio central” en
		nuestro simulador): ¿qué estrategia crees que es más efectiva? ¿Por qué?

- Respuesta:	Estadio 1: 4 veces
		Estadio 2: 4 veces
		Estadio 3: 8 veces
		Estadio 4: 4 veces
		
		La estrategia Snoopy seria mas efectiva porque los accesos a cache son mas rapido (el doble) y ademas solo
		se realizan cuando el estadio lo necesita. Solo se comunican con sistema central al principio y al final 
		de la ejecucion.

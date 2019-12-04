En estos dos ejemplos se podia dar el caso de que los dos primeros coches que llegasen, si eran de sentidos diferentes, se chocasen.
El problema esta en el primer if de la funcion carRunning:

func carRunning(c *car, n *sync.WaitGroup)  {
			.
			.
			.

	if c.side == "east"{
		brg.eastch <- car{c.id, c.side}
	}else{
		brg.westch <- car{c.id, c.side}
	}
			.
			.
			.
Si nos fijamos, dos gorutinas podrian acceder a ese if antes de que una de ellas estableciese el valor del sentido del puente. Es por eso que en ese caso se podrian chocar dos coches. Recordemos que el mecanismo para que un coche pueda pasar viene por el hecho de poder meter al coche en el semaforo.

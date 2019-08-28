# com.mercadolibre.itacademy.apigo
Api GO with simple. waitgroup and channel implementation

EndPoints: 
  "/users/:userId" -> get user with id
	"/countries/:countryId" -> get country with id
	"/sites/:siteId" -> get site with id
	"/results/:userId" get Result object from user id

	endpoint para result workgroup
	"/waitgroup/results/:userId" -> result object with user id

	endpoint para result con channels
	"/channel/results/:userId" -> result object with user id
  
  IMPORTANTE: el endpoint de CHANNEL internamente hace 3 calls a Users, Countries y Sites, pero
  estam cambiados para que hagan call al mock server: localhost:8085
  
  Tiene implementado un CircuitBreaker:
  al haber 3 pedidos SEGUIDOS con error (respuesta 500, timeout o servidor caido) circuitBreaker 
  cambia su estado a OPEN, iniciando un timeout de 10 segundos, pasado ese tiempo se cambia el estado
  a HalfOpen donde pingea a los 3 endpoints (actualmente el Mockserver) y si no recibe Status Code 200
  de los 3 vuelve al estado Open, iniciando el timeout de nuevo.
  
  TimeOut es igual a 10 segundos por metodos practicos, facilmente se puede realizar un timeout mas largo
  igualmente al timeout de la conexion que esta determinado en 5 segundos

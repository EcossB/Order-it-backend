package websockets

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// upgrader es la herramienta que nos ayuda a convertir la peticion http (get) en un websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	//CheckOrigin es una funcion que evita ataques csrf
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// writePump es una funcion la cual saca los mensajes del buzon y los envia a traves de la red (websocket)
func writePump(hub *Hub, client *Client) {

	/*
		aseguramos que si esta funcion termina desconectamos limpiamente del hub y cerramos la conexion.
	*/
	defer func() {
		hub.RemoveClient(client)
		client.Conn.Close()
	}()

	// Este bucle se queda esperando (bloqueado) hasta que aparezca un mensaje en el buzón 'Send'
	for {
		message, ok := <-client.Send

		// Si !ok significa que el Hub cerró este buzón (ej: usamos RemoveClient). Terminamos.
		if !ok {
			client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		//si encontramos un mensaje, lo mandamos por tunel tcp/ip
		err := client.Conn.WriteMessage(websocket.TextMessage, message)

		if err != nil {
			//Si falla al escribir, salimos pacificamente
			return
		}
	}

}

func ServeWS(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		//1. extraemos quien se esta conectando con los query params
		// Ej: /api/ws?tenant_id=123&role=CHEF&kitchen_id=456

		tenantId := c.Query("tenant_id")
		role := c.Query("role")

		if tenantId == "" || role == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tenantId y role son requeridos"})
			return
		}

		//2 definimos a que sala (topic) va a pertenecer
		var topic string

		if role == "CHEF" {
			kitchenId := c.Query("kitchen_id")
			if kitchenId == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "kitchenId es requerido para el rol CHEF"})
				return
			}

			topic = fmt.Sprintf("tenant:%s:kitchen:%s", tenantId, kitchenId)

		} else if role == "WAITER" {
			userID := c.Query("user_id")
			topic = fmt.Sprintf("tenant:%s:waiter:%s", tenantId, userID)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "rol no soportado para los websocket"})
			return
		}

		// 3. convertimos la peticion HTTP en una conexion websocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

		if err != nil {
			log.Println("Error actualizando a websocket:", err)
			return
		}

		//4 creamos nuestro objecto cliente

		client := &Client{
			Conn:  conn,
			Send:  make(chan []byte, 256), // un buzon con capacidad para 256 mensajes en cola
			Topic: topic,
		}

		//5. lo registramo en una sala de nuestro hub
		hub.AddClient(client)

		//6 iniciamos un "bucle" en un hilo separado que se encargara de enviar cualquier cosa que caiga en el buzon.
		go writePump(hub, client)

	}
}

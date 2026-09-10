package websockets

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Client representa un cliente conectado (mesero, chef)
type Client struct {
	//Conn Es el tunel por donde viajan los datos a traves de tcp/ip
	Conn *websocket.Conn

	//Send es un canal de go, por donde enviamos los datos a traves de Conn, es como un buzon.
	Send chan []byte

	//Topic Es el canal al que pertenece el mensaje (ej. un tenant)
	Topic string
}

// Hub es el cerebro central que lleva registro de quien esta conectado.
type Hub struct {
	// Clients es un mapa. La llave es el Topic (sala) y el valor es un mapa de Clientes.
	// Esto nos permite buscar rápidamente a todos los clientes de una "sala" específica.
	Clients map[string]map[*Client]bool

	// Mutex nos protege de condiciones de carrera (Race Conditions).
	// Asegura que dos peticiones no intenten agregar o borrar clientes al mismo tiempo.
	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[string]map[*Client]bool),
	}
}

func (h *Hub) AddClient(client *Client) {

	//Hacemos lock para evitar race conditions
	h.mu.Lock()
	defer h.mu.Unlock() //nos aseguramos de desbloquear cuando termine de correr.

	//si la sala no existe, la creamos.
	if h.Clients[client.Topic] == nil {
		h.Clients[client.Topic] = make(map[*Client]bool)
	}

	//agregamos el cliente a la sala
	h.Clients[client.Topic][client] = true

}

// RemoveClient es un metodo que elimina algun dispositivo de la sala cuando se desconecte
func (h *Hub) RemoveClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.Clients[client.Topic]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client) //lo borramos del mapa
			close(client.Send)      //cerramos el canal

			//si la sala se quedo vacia, la eliminamos para liberar la RAM
			if len(clients) == 0 {
				delete(h.Clients, client.Topic)
			}

		}
	}
}

// Broadcast envía un mensaje (payload JSON) a todos los clientes de una sala
func (h *Hub) Broadcast(topic string, msg []byte) {
	h.mu.RLock() //se usa rlock porque solo se va a leer el mapa
	defer h.mu.RUnlock()

	//buscamos los clientes inscritos en la sala
	if clientsInTopic, ok := h.Clients[topic]; ok {
		for client := range clientsInTopic {

			// intentamos mandar un mensaje al cliente a traves del canal
			select {
			case client.Send <- msg:

			default:
				// si el buzon esta lleno o bloqueado, asumimos que el cliente murio.
				close(client.Send)
				delete(clientsInTopic, client)
			}

		}
	}
}

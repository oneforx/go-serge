Le type ServerEngine contient les champs suivants :

    TcpListener : un pointeur vers une instance de net.TCPListener.
    TcpInit : un booléen qui indique si la connexion TCP est initialisée.
    TcpConnexions : une map qui associe chaque adresse TCP à un pointeur vers une instance de net.TCPConn.
    UdpListener : un pointeur vers une instance de net.UDPConn.
    UdpInit : un booléen qui indique si la connexion UDP est initialisée.
    UdpConnexions : une map qui associe chaque adresse UDP à un pointeur vers une instance de net.UDPAddr.

La méthode Start démarre le serveur sur l'adresse spécifiée et gère les connexions TCP et UDP entrantes en appelant les gestionnaires de connexions correspondants passés en paramètre. La méthode utilise une goroutine pour chaque gestionnaire.

La méthode USendToAll envoie un message à tous les clients connectés via UDP.

La méthode USendToAddr envoie un message à un client connecté spécifique via UDP.

Le package engine utilise également les packages standard encoding/json, fmt, log, net et sync, ainsi que le package github.com/google/uuid.
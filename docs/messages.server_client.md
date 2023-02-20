# Documentation du fichier server_client.go

Le fichier server_client.go du package messages contient des fonctions qui permettent de créer des messages à destination des clients du serveur. Ces fonctions prennent en paramètre des données à inclure dans le message et renvoient une structure engine.Message qui représente le message.

Voici les différentes fonctions disponibles :

    SC_CONNECT_FAILED(data interface{}) engine.Message : cette fonction permet de créer un message indiquant une erreur de connexion. La variable data doit contenir un message décrivant l'erreur. Le message créé a le type "CONNECT_FAILED", le mode réseau "NET_HYB" et la cible "CLIENT_TARGET".

    SC_CONNECT_SUCCESS(data interface{}) engine.Message : cette fonction permet de créer un message indiquant une connexion réussie. La variable data doit contenir un message de bienvenue ou toute autre donnée utile. Le message créé a le type "CONNECT_SUCCESS", le mode réseau "NET_HYB" et la cible "CLIENT_TARGET".

    SC_PING(latence int) engine.Message : cette fonction permet de créer un message de ping. La variable latence doit contenir le temps de latence entre le serveur et le client en millisecondes. Le message créé a le type "PING", le mode réseau "NET_TCP" et la cible "CLIENT_TARGET".

    SC_CREATE_ENTITY(entityData ecs.Entity) engine.Message : cette fonction permet de créer un message pour créer une entité côté client. La variable entityData doit contenir les données de l'entité à créer. Le message créé a le type "CREATE_ENTITY", le mode réseau "NET_TCP" et la cible "CLIENT_TARGET".

    SC_UPDATE_ENTITY(entityData ecs.Entity) engine.Message : cette fonction permet de créer un message pour mettre à jour une entité côté client. La variable entityData doit contenir les nouvelles données de l'entité. Le message créé a le type "UPDATE_ENTITY", le mode réseau "NET_UDP" et la cible "CLIENT_TARGET".

    SC_DELETE_ENTITY(id string) engine.Message : cette fonction permet de créer un message pour supprimer une entité côté client. La variable id doit contenir l'identifiant de l'entité à supprimer. Le message créé a le type "DELETE_ENTITY", le mode réseau "NET_TCP" et la cible "CLIENT_TARGET".
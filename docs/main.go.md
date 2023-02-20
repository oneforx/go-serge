Le code définit une structure de données Client qui stocke les détails de chaque client connecté au serveur, tels que l'adresse UDP et la connexion TCP, l'ID client, le jeton, l'adresse e-mail et le mot de passe.

Le serveur a deux modes de connexion: TCP et UDP, et un mode de connexion est défini en tant que type de constante NET_MODE.

Le main fonction est la fonction principale qui démarre le serveur et initialise tous les paramètres nécessaires.

La fonction commence par créer une instance du moteur de serveur engine.ServerEngine qui va être utilisé pour démarrer et gérer les connexions TCP et UDP.

Il crée également un map nommé clients pour stocker les informations de chaque client connecté.

Un ID unique est généré pour chaque client connecté en utilisant la bibliothèque github.com/google/uuid, puis le client est stocké dans le map clients.

Ensuite, une liste d'entités de jeu gameEntities est créée, et chaque entité est ajoutée à un monde de jeu gameWorld à l'aide de la méthode AddEntity.

Enfin, la méthode Start de l'instance serverEngine est appelée, en passant l'adresse et le canal de messages pour les connexions UDP, ainsi que deux gestionnaires pour les connexions TCP et UDP.

Le gestionnaire pour les connexions TCP commence par envoyer un message de ping au client pour tester la connexion, puis attend de lire des données à partir de la connexion.

Il essaie ensuite de désérialiser les données lues dans une structure de message engine.Message en utilisant la fonction json.Unmarshal.

Si le message est de type DISCONNECT, le serveur envoie un message DISCONNECT de retour pour déconnecter l'utilisateur.

Si le message est de type CONNECT, le serveur essaie de trouver un client correspondant à l'e-mail et au mot de passe envoyés par le client.

Si un client correspondant est trouvé, le serveur envoie un message de connexion réussie au client avec un jeton unique.

Si un client correspondant n'est pas trouvé, le serveur envoie un message de connexion échouée au client.

Si le message est de type CONNECT_TOKEN, le serveur lie la connexion au client correspondant.

C'est en gros une vue d'ensemble du fichier main, j'espère que cela t'aide. N'hésite pas à me poser des questions si tu en as.
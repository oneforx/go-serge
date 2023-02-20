Ce fichier définit plusieurs fonctions qui retournent des messages destinés au serveur. 

Chaque message a un type spécifique qui correspond à une action à réaliser sur le serveur, ainsi qu'une cible et un mode de réseau.

Par exemple, la fonction CS_CONNECT_TOKEN retourne un message de type "CONNECT_TOKEN" qui contient des données de connexion à envoyer au serveur. 

La fonction CS_DELETE_CHARACTER retourne un message de type "DELETE_CHARACTER" qui contient l'ID du personnage à supprimer sur le serveur.
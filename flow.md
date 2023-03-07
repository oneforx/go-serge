Idée :

Jeux tcp / udp -> OK
Se déplacer -> OK
Meilleur gestion de l'écriture du message sur le client Split(scanner) -> OK
Gérer la connexion de plusieurs client différent -> OK
Envoyer un message -> CREATE_WORLD OK
OPTIMISER les message -> UPDATE
    ECS
        -> Entities filter methods
            -> Get > Only for read, mutex.RLock Update -> Can Read and Write Lock
Entities can have states, and its not sent to the client

Envoyer des glands
    -> Client send CS_SHOOT(orientation radiant float64)
        -> A la réception du message, on détecte la balle viens de qui et on crée une EntityWithOwner
    -> Entity have state and isn't sent to client

    -> Fixer la vitesse d'envois
    -> Quand un glands touche une entité
        -> Enlever de la vie
        -> Tracer avec un tableau Entity.States { touchs: []entityIds }
            -> Quand la vie de l'entité est à 0 ou inférieur il est considérer comme mort
                ON_DEAD:
                    -> SEND_MESSAGE { entityId:  content: "X a mis ko y avec l'aide de [a,b,c] " }
                    -> DEAD_SCREEN { bouton Réssuciter }
            -> Le dernier qui l'a touché et que l'entité component life est à 
            

Faire un système de camera
    UPDATE
        - READ new position
        - SET camera.target_position = new_position
        - IF camera.position != target_position
            -> Incrémenter / décrémenter de 1 vers la target position


Le système de mort
    -> A la mort > +18 ? Afficher le bouton quitter le corp sinon afficher "Rejouer"
    -> Le corp ou sac reste sur la carte
    -> Si le joueur a saisi -18 ans alors on affiche des sacs à la place
    -> Possibilité de fouiller le corp -> Pourcentage de chance d'avoir les items que le corp contient
    -> COMPONENT POSSESSION, remove that component when 
    -> COMPONENT ITEM_CONTAINER
    -> [TCHAT] Envoyer un message system { ownerId: systemId, characterId: system.character[id], content }
        -> Character { id, name }
            -> Créer un character "La faucheuse"


Faire un tchat
    -> MESSAGE_SEND_MESSAGE
        -> Créer un utilisateur systéme, utilisé comme identifiant, comme compte du serveur
        -> { ownerId, content, type } Seulement le système qui écrit un message avec un type différent
        -> Les messages défilent du bas vers le haut


ECRANS:
GAME_SCREEN: Information sur l'écran de jeux
    -> Tchat
    -> Vie
    -> ButtonInventaire -> Inventaire -> Appuyer sur I
    -> Slots
OPTION:
    -> Modifier les commandes

CHARACTERS: Faire un écran de création de personnage
    -> MESSAGE_REQUEST_LIST_CHARACTER {NET_TCP}
    -> MESSAGE_RESPONSE_LIST_CHARACTER {[]Character}
    -> Interface, bouton créer, input name, bouton droite / gauche pour sélectionner un personnage

ITEMS:
    -> Livre de résurection: Utiliser l'item permet de réssusité les entités KO proches
        -> ITEM_HOLDING: Afficher un cercle autour du personnage
    -> Potion de vitesse
    -> Potion de force
    -> Potion de agilité
    -> Potion de dexterité

Prendre des objets par terre
-> Créer des objets
-> Ajouter le composant Characteristic
-> Seul les entités avec le composant Characteristic peuvent boire la potion
-> Faire une zone de capture
    -> Créer un système de vague

Faire une base de donnée:
    GameServer read at launch / save at exit
    USERS, ENTITIES, USER_HAS_ENTITY


Messages.AREA -> Envoyer à tous ceux qui doivent voir, Envoyer au monde entier

Répertorier les components qui ont besoin d'être vue par les autres
    ComponentType: 
        LOCAL -> Tout ceux qui voit l'entité sont au courant du changement, comme par exemple la position
        WORLD -> Peut être envoyer à tout le monde
        PRIVATE -> Peut être envoyer qu'au propriétaire

Système de mod
ModFolder
    library
        - Components
        - Entities
        - Systems

Architecture
    Redis -> World Database
    WorldServer -> Gestion du monde
    ApiServer -> Gestion de l'authentification / Management des ressources


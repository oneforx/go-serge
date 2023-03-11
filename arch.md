AuthServer
    - Redis


WorldServer
    - Config
    - *LibraryManager
        - Libraries - []ILibrary
            - OK | SystemTypes
            - EntityTypes
            - OK | ComponentTypes
            - EXPERIMENTAL MessageTypes
            - EXPERIMENTAL Listen(Identifier, func(*Client, *WorldServer))
            - Method - RegisterComponent
            - Method - RegisterComponents
            - Method - RegisterSystem
            - Method - RegisterSystems
            - Method - RegisterMessageHandler
            - Method - 
            - 
    - World
        - *Entities / *Systems
    - *Tcp
    - *Udp
    - Method - Implemented - LoadLibrary
    - Method - RegisterMessage(message, handler)
    - Method - Init(world)
    - Method - Start() // Initialise TCP & UDP server

Quand je veux créer un jeux multijoueur
- Je voudrais définir les entités / composants et systems -> OK

Pour permettre aux autres développeur de créer un world server
- Je dois les permettres d'enregistrer des messages

Quand une entité se met à jour, envoyer la mise à jour aux Client qui sont à moins 1000px 
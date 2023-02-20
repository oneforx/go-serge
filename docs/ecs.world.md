Le fichier world.go du package ecs définit une interface IWorld et sa structure d'implémentation World. IWorld décrit les méthodes que doit implémenter une instance de World pour être considérée comme un monde (un conteneur d'entités et de systèmes).

La structure World dispose de quatre champs : Id (l'identifiant du monde), Entities (la liste des entités du monde), Systems (la liste des systèmes du monde) et tags (une carte associant des identifiants de tags à des identifiants d'entités).

Les méthodes de l'interface IWorld comprennent :

    - GetId() : renvoie l'identifiant du monde.
    
    - AddEntity(entity *IEntity) error : ajoute une entité à la liste des entités du monde. Renvoie une erreur si une entité avec le même identifiant existe déjà.
    
    - GetEntity(id string) *IEntity : renvoie l'entité correspondant à l'identifiant donné.
    
    - GetEntities() []*IEntity : renvoie toutes les entités du monde.
    
    - GetEntitiesPossessedBy(ownerID string) []*IEntity : renvoie toutes les entités possédées par l'entité ayant l'identifiant donné.
    
    - GetEntitiesByComponentId(id string) []*IEntity : renvoie toutes les entités ayant le composant avec l'identifiant donné.
    
    - GetEntitiesWithComponents(v ...string) []*IEntity : renvoie toutes les entités possédant tous les composants avec les identifiants donnés.
    
    - GetEntitiesWithComposition(composition []string) []*IEntity : renvoie toutes les entités possédant une composition contenant tous les composants avec les identifiants donnés.
    
    - GetEntitiesWithStrictComposition(composition []string) []*IEntity : renvoie toutes les entités possédant exactement les composants avec les identifiants donnés.
    
    - RemoveEntity(id string) error : supprime l'entité correspondant à l'identifiant donné. Renvoie une erreur si l'entité n'existe pas.

    - AddSystem(sys *ISystem): Cette méthode prend en paramètre un pointeur vers une structure ISystem et l'ajoute à la liste des systèmes de l'instance du monde. Elle ne renvoie rien.

    - GetSystemById(id string) *ISystem: Cette méthode prend en paramètre un identifiant de système en tant que chaîne de caractères et renvoie un pointeur vers la structure ISystem correspondante, s'il existe. Si aucun système correspondant n'est trouvé, la méthode renvoie nil.

    - RemoveSystem(id string) (err error): Cette méthode prend en paramètre un identifiant de système en tant que chaîne de caractères et supprime le système correspondant de la liste des systèmes de l'instance du monde. Si aucun système correspondant n'est trouvé, la méthode renvoie une erreur Cannot delete entity because it doesn't exist. Sinon, elle renvoie nil.

    - Update(): Cette méthode parcourt tous les systèmes de l'instance du monde et appelle la méthode Update() de chaque système.

Il convient de noter que la plupart des méthodes définies dans la structure World sont des méthodes de recherche qui renvoient des tranches de pointeurs vers des entités ou des systèmes. Les entités et les systèmes sont stockés sous forme de pointeurs pour éviter de copier inutilement les données et pour permettre leur modification dans la liste d'entités ou de systèmes.
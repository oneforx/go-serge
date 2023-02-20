Le fichier entity.go contient la définition de l'interface IEntity et de son implémentation concrète Entity.

L'interface IEntity définit un ensemble de méthodes qui doivent être implémentées par toutes les entités du jeu. Elle permet de récupérer l'ID et le propriétaire de l'entité, d'ajouter des composants à l'entité, de vérifier si elle possède un certain composant, de récupérer un composant par son ID, de récupérer tous les composants de l'entité, de récupérer la composition de l'entité (une liste de noms de composants) et de vérifier si elle possède une certaine composition.

La structure Entity représente une entité du jeu, elle contient un ID, un ID de propriétaire (qui peut être une chaîne vide si l'entité n'est pas possédée par un client), un pointeur vers le monde dans lequel l'entité existe, et une liste de pointeurs vers les composants qui sont attachés à l'entité. Elle implémente toutes les méthodes de l'interface IEntity.

Les méthodes de l'interface IEntity et de l'implémentation concrète Entity utilisent des pointeurs vers des interfaces IComponent et IWorld, qui sont définies dans d'autres fichiers du package ecs.

Le fichier contient également des fonctions de création d'entité (CEntity et CEntityPossessed), qui prennent en paramètre le monde dans lequel l'entité existe, son ID, un ID de propriétaire (dans le cas où l'entité est possédée par un client) et une liste de pointeurs vers les composants qui sont attachés à l'entité.
Ce fichier déclare une interface appelée ISystem dans le package ecs. Cette interface contient deux méthodes :

    - GetId() qui renvoie un identifiant sous forme de chaîne de caractères pour le système.
    
    - Update() qui est utilisée pour mettre à jour l'état du système.

Toute structure qui implémente cette interface doit fournir une implémentation pour ces deux méthodes. Cela permet à différents systèmes de partager une même interface, ce qui facilite l'utilisation et l'interopérabilité entre ces systèmes dans le cadre d'une architecture orientée objet.

Comment créer un system :
    -> Quand un system peut affecte un seul composant alors appelé le fichier nom_du_composant.system.go
    -> 
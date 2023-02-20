Le fichier types.go du package engine contient les types et variables suivants :

Type TargetType : Cette variable de type chaine de caractères (string) est utilisée pour identifier le type de cible (CLIENT_TARGET ou SERVER_TARGET).

Variables CLIENT_TARGET et SERVER_TARGET : Ces variables définissent les types de cibles possibles (client ou serveur).

Type Message : Cette structure contient les champs suivants :

    MessageType : une chaine de caractères qui décrit le type du message
    Data : une interface{} qui contient les données du message
    Target : une variable de type TargetType qui définit la cible du message
    NetMode : une variable de type NET_MODE qui indique le mode de connexion

Type FeedBack : Cette structure est utilisée pour fournir un retour d'information. Elle contient les champs suivants :

    Host : une chaine de caractères qui indique le nom du bloc ou de la fonction qui héberge les jobs
    Job : une chaine de caractères qui indique le nom de la fonction qui a déclenché le retour d'information
    Label : une chaine de caractères pour donner un label au feedback
    Comment : une chaine de caractères pour commenter le feedback
    Data : une interface{} pour contenir les données supplémentaires

Type NET_MODE : Cette variable de type chaine de caractères (string) est utilisée pour définir le mode de connexion. Les valeurs possibles sont TCP, UDP et HYB (hybride, pouvant être envoyé sur le réseau TCP ou UDP).

Variables NET_TCP, NET_UDP et NET_HYB : Ces variables définissent les modes de connexion possibles (TCP, UDP ou hybride).
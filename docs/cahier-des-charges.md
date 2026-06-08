# Cahier des charges - Forum

## Objectif du projet

Créer une application web de forum permettant aux utilisateurs de lire, publier, commenter et interagir avec des contenus.

## Utilisateurs

- Visiteur non connecté
- Utilisateur connecté

## Fonctionnalités principales

- Inscription
- Connexion
- Déconnexion
- Création de posts
- Commentaires
- Catégories
- Likes et dislikes
- Filtres

## Contraintes techniques

- Back-end en Go
- Base de données SQLite
- Pages en HTML et CSS
- JavaScript utilisé pour des interactions simples
- Gestion des sessions avec cookies
- Gestion des erreurs HTTP
- Code organisé avec une arborescence claire
- Utilisation de Git avec des branches et des commits explicites

## Données à gérer

Le forum devra stocker :
- les utilisateurs ;
- les sessions ;
- les posts ;
- les commentaires ;
- les catégories ;
- les likes et dislikes.

## Livrables attendus

- Lien du dépôt GitHub
- Capture du Trello
- Application fonctionnelle
- Documentation du projet
- Support de présentation orale

## Droits des utilisateurs

### Visiteur non connecté

Un visiteur non connecté peut :
- lire les posts ;
- lire les commentaires ;
- voir le nombre de likes et dislikes ;
- filtrer les posts par catégorie.

Un visiteur non connecté ne peut pas :
- créer un post ;
- écrire un commentaire ;
- liker ou disliker ;
- accéder aux pages réservées aux utilisateurs connectés.

### Utilisateur connecté

Un utilisateur connecté peut :
- créer un post ;
- modifier ses propres posts ;
- supprimer ses propres posts ;
- écrire des commentaires ;
- liker ou disliker des posts et commentaires ;
- filtrer ses propres posts ;
- filtrer les posts qu'il a likés.

## Critères de validation

Le projet sera considéré comme conforme si :
- un visiteur peut lire les posts et commentaires ;
- un utilisateur peut s'inscrire, se connecter et se déconnecter ;
- un utilisateur connecté peut créer, modifier et supprimer ses propres posts ;
- un utilisateur connecté peut commenter ;
- un utilisateur connecté peut liker ou disliker ;
- les posts peuvent être filtrés par catégorie ;
- un utilisateur connecté peut voir ses posts et ses posts likés ;
- les erreurs web sont gérées proprement ;
- l'interface est lisible sur mobile, tablette et ordinateur.

## Décisions à valider en équipe

- Liste finale des catégories
- Schéma de base de données
- Répartition des tâches
- Style visuel du site
- Organisation des branches Git
- Cartes Trello et priorités

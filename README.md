# Forum Project

Projet de forum web réalisé en Go, SQLite, HTML, CSS et JavaScript.

L'objectif est de créer une application de type Reddit permettant aux utilisateurs de s'inscrire, se connecter, publier des posts, commenter, liker/disliker et filtrer les contenus.

## Etat actuel

Fonctionnalités commencées :

- Page d'accueil
- Page d'inscription
- Page de connexion
- Déconnexion
- Création de compte avec mot de passe hashé avec `bcrypt`
- Création de session avec UUID
- Cookie de session
- Affichage différent selon utilisateur connecté ou non

Fonctionnalités encore à faire :

- Création de posts
- Edition et suppression de ses propres posts
- Commentaires
- Catégories
- Likes et dislikes
- Filtres par catégorie, posts créés et posts likés
- Gestion complète des erreurs HTTP
- Tests automatisés
- Initialisation automatique des migrations SQL

## Stack technique

- Go
- SQLite
- HTML
- CSS
- JavaScript pour interactions simples
- `github.com/mattn/go-sqlite3`
- `golang.org/x/crypto/bcrypt`
- `github.com/google/uuid`

## Structure du projet

- `cmd/server/` : point d'entrée du serveur Go
- `internal/db/` : connexion à la base SQLite
- `internal/db/migrations/` : scripts SQL de création et évolution des tables
- `internal/handlers/` : logique des routes HTTP
- `internal/templates/` : fonction de rendu des templates HTML
- `web/templates/` : pages HTML
- `web/static/css/` : fichiers CSS
- `docs/` : cahier des charges et documentation équipe

## Routes actuelles

- `GET /` : page d'accueil
- `GET /login` : formulaire de connexion
- `POST /login` : traitement de la connexion
- `GET /register` : formulaire d'inscription
- `POST /register` : traitement de l'inscription
- `GET /logout` : déconnexion

## Lancer le projet

Sur Windows avec `go-sqlite3`, `cgo` et `gcc` doivent être disponibles.

```powershell
$env:Path += ";C:\msys64\ucrt64\bin"
$env:CGO_ENABLED="1"
go run .\cmd\server
```

Le serveur démarre sur :

```text
http://localhost:8080
```

## Tester le projet

```powershell
$env:Path += ";C:\msys64\ucrt64\bin"
$env:CGO_ENABLED="1"
go test ./...
```

## Base de données

La base locale s'appelle actuellement :

```text
forum.db
```

Les fichiers `.db` sont ignorés par Git car ils sont générés localement et peuvent contenir des données privées.

Les migrations SQL sont dans :

```text
internal/db/migrations/
```

Point à améliorer :

- le code ouvre la base SQLite ;
- les migrations existent ;
- mais l'application ne les exécute pas encore automatiquement.

## Cahier des charges

Le cahier des charges est disponible ici :

```text
docs/cahier-des-charges.md
```

Ce fichier sert de checklist pour vérifier que le projet respecte l'énoncé.

## Workflow Git

Chaque membre travaille sur sa branche.

Avant de commencer une tâche :

```powershell
git checkout main
git pull
git checkout -b feat/nom-de-la-tache
```

Avant de merge :

```powershell
go test ./...
git status
```

Les commits suivent ce format :

```text
[TAG] description courte en français
```

Tags utilisés :

- `[ADD]` : ajout de fichier ou fonctionnalité
- `[MODIF]` : modification d'un fichier existant
- `[FIX]` : correction de bug
- `[INIT]` : initialisation du projet ou d'une section
- `[WIP]` : travail en cours
- `[MERGE]` : fusion de branche

## Prochaines priorités

1. Vérifier et stabiliser l'authentification existante
2. Automatiser l'exécution des migrations SQL
3. Ajouter les tests sur inscription, connexion et sessions
4. Créer le modèle des posts
5. Ajouter création, lecture, modification et suppression de posts
6. Ajouter commentaires
7. Ajouter catégories
8. Ajouter likes/dislikes
9. Ajouter filtres
10. Améliorer accessibilité, responsive et erreurs HTTP

## Points à clarifier en équipe

- Garder `mail` ou renommer en `email`
- Garder `password` ou renommer en `password_hash`
- Supprimer ou justifier `package.json`
- Choisir la méthode officielle pour les migrations
- Adapter le cookie `Secure` entre local et production
- Définir les catégories initiales du forum

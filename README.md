# Forum Project

Forum web musical realise en Go, SQLite, HTML, CSS et JavaScript.

L'objectif est de creer une application de type Reddit autour de la musique : inscription, connexion, posts, commentaires, tags, likes/dislikes et filtres.

## Etat actuel

Fonctionnalites deja presentes :

- inscription utilisateur ;
- connexion utilisateur ;
- mot de passe chiffre avec `bcrypt` ;
- session avec UUID ;
- cookie de session `HttpOnly` ;
- deconnexion ;
- page d'accueil avec posts ;
- creation de tags ;
- creation de posts avec un ou plusieurs tags ;
- affichage du detail d'un post ;
- commentaires et reponses ;
- likes/dislikes sur posts et commentaires ;
- edition et suppression des posts par leur auteur ;
- profil utilisateur ;
- page parametres ;
- filtres par tag, populaires, posts crees et posts likes ;
- migrations SQLite lancees automatiquement au demarrage ;
- Dockerfile et Docker Compose.

## Stack technique

- Go ;
- SQLite ;
- HTML templates ;
- CSS ;
- `github.com/mattn/go-sqlite3` ;
- `golang.org/x/crypto/bcrypt` ;
- `github.com/google/uuid`.

## Structure du projet

- `cmd/server/` : point d'entree du serveur Go ;
- `internal/db/` : connexion SQLite et execution des migrations ;
- `internal/db/migrations/` : scripts SQL ;
- `internal/handlers/` : logique des routes HTTP ;
- `internal/templates/` : rendu des templates HTML ;
- `web/templates/` : pages HTML ;
- `web/static/` : fichiers CSS et images.

## Routes principales

- `GET /` : accueil ;
- `GET /login` et `POST /login` : connexion ;
- `GET /register` et `POST /register` : inscription ;
- `GET /logout` : deconnexion ;
- `GET /createPost` et `POST /createPost` : creation de post ;
- `GET /createTag` et `POST /createTag` : creation de tag ;
- `GET /post/{id}` : detail d'un post ;
- `GET /post/like?id={id}` : like d'un post ;
- `GET /post/dislike?id={id}` : dislike d'un post ;
- `POST /post/comment` : ajout de commentaire ;
- `GET /comment/like?id={id}` : like d'un commentaire ;
- `GET /comment/dislike?id={id}` : dislike d'un commentaire ;
- `GET /post/edit?id={id}` et `POST /post/edit?id={id}` : edition d'un post ;
- `GET /post/delete?id={id}` : suppression d'un post ;
- `GET /user/@{username}` : profil public ;
- `GET /settings` et `POST /settings` : parametres utilisateur.

## Lancer le projet en local

Sur Windows avec `go-sqlite3`, `cgo` et `gcc` doivent etre disponibles.

```powershell
$env:Path = "C:\Program Files\Go\bin;C:\msys64\ucrt64\bin;" + $env:Path
$env:CGO_ENABLED = "1"
$env:COOKIE_SECURE = "false"
go run .\cmd\server
```

Le serveur demarre sur :

```text
http://localhost:5500
```

## Tester le projet

```powershell
$env:Path = "C:\Program Files\Go\bin;C:\msys64\ucrt64\bin;" + $env:Path
$env:CGO_ENABLED = "1"
go test ./...
```

## Docker

Construire l'image :

```powershell
docker compose build
```

Lancer le conteneur :

```powershell
docker compose up
```

Le site est disponible ici :

```text
http://localhost:5500
```

La base SQLite du conteneur est stockee dans un volume Docker nomme `forum-data`.

Arreter le conteneur :

```powershell
docker compose down
```

Repartir avec une base Docker vide :

```powershell
docker compose down -v
```

## Variables d'environnement

- `DB_PATH` : chemin du fichier SQLite. Par defaut : `forum.db`.
- `COOKIE_SECURE` : mettre `true` en production HTTPS. En local HTTP, laisser `false`.
- `APP_ENV=production` : active aussi les cookies securises.

## Base de donnees

Les migrations SQL sont dans :

```text
internal/db/migrations/
```

Au demarrage, l'application ouvre la base SQLite, lance les migrations non appliquees, puis cree les tables complementaires necessaires aux commentaires et aux votes.

Les fichiers `.db` sont ignores par Git car ils sont generes localement et peuvent contenir des donnees privees.

## Workflow Git

Chaque membre travaille sur sa branche.

Avant de commencer :

```powershell
git switch main
git pull origin main
git switch -c feat/nom-de-la-tache
```

Avant de proposer une merge request :

```powershell
go test ./...
git status
```

Format de commit :

```text
[TAG] description courte en francais
```

Tags utilises :

- `[ADD]` : ajout de fichier ou fonctionnalite ;
- `[MODIF]` : modification d'un fichier existant ;
- `[FIX]` : correction de bug ;
- `[INIT]` : initialisation ;
- `[WIP]` : travail en cours ;
- `[MERGE]` : fusion de branche.

## Reste a faire

- Ajouter plus de tests automatises sur les handlers HTTP ;
- mieux valider les formulaires cote serveur ;
- ajouter une vraie expiration de session en base ;
- ameliorer les messages d'erreur utilisateur ;
- finaliser le responsive mobile ;
- verifier accessibilite et securite avant rendu final ;
- deployer une version de demonstration.

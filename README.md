# Forum Project

Projet de forum web réalisé en Go, SQLite, HTML, CSS et JavaScript.

## Objectif

Créer une application web de type Reddit permettant aux utilisateurs de discuter avec des posts, des commentaires, des catégories, des likes/dislikes et des filtres.

## Documentation

- Cahier des charges : `docs/cahier-des-charges.md`
- Proposition de fondation équipe : `docs/proposition-equipe.md`

## Structure actuelle

- `cmd/` : point d'entrée du serveur Go
- `internal/` : logique interne de l'application
- `web/` : fichiers HTML, CSS, JavaScript et assets
- `docs/` : documentation du projet

## Lancer le projet

```powershell
go run .\cmd\server
```

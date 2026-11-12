# JUL – Discographie

C'est mon projet pour le module « Développement WEB – consommation d'APIs ».
Le but : un petit site en Go qui affiche des titres de JUL, récupérés directement depuis l'API de Spotify.

Quelques règles du sujet que j'ai respectées : Go avec la bibliothèque standard uniquement, des templates `html/template`, aucun JavaScript, et les filtres écrits à la main.

## Lancer le projet

Il faut Go 1.22 minimum, parce que les routes du genre `/tracks/{id}` ne marchent pas avant cette version. Il faut aussi une application Spotify créée sur le site développeur (https://developer.spotify.com) pour avoir un Client ID et un Client Secret.

1. Récupérer le projet :

        git clone <lien-du-depot>
        cd jul-discographie

2. Copier le fichier d'exemple :

        cp .env.example .env

3. Ouvrir `.env` et mettre ses vrais identifiants à la place des valeurs bidon :

        SPOTIFY_CLIENT_ID=ton_client_id_ici
        SPOTIFY_CLIENT_SECRET=ton_secret_ici

   Pas de guillemets autour des valeurs. Ce fichier est dans le `.gitignore`, donc il ne part jamais sur Gitea.

4. Démarrer :

        go run .

5. Aller sur http://localhost:8080

## Ce que fait le site

- `/` : la liste des titres sous forme de cartes (pochette, nom, artistes, date), avec les filtres sur le côté.
- `/tracks/{id}` : la page d'un titre (pochette, album, date de sortie, durée, un bouton pour l'ouvrir sur Spotify).
- Toute autre adresse : une page 404 à mon design.
- Si Spotify ne répond pas et que je n'ai aucun titre à montrer : une page 503, le serveur ne plante pas.

## Les filtres

Ils passent dans l'adresse (méthode GET), donc on peut les mettre en favori ou les partager. Exemple :

    http://localhost:8080/?annee=2026&type=single&feat=oui

- `annee` : l'année de sortie (la liste des années est construite à partir des titres)
- `type` : `single` ou `album`
- `feat` : `oui` si le titre a plusieurs artistes, `non` sinon

On peut tous les combiner. Après validation, les menus gardent ce qu'on a choisi. Tout le tri est fait en Go, dans `filtres.go`, sans aucune bibliothèque.

## Le cache

Je limite les appels à Spotify pour éviter l'erreur 429 (trop de requêtes). Au premier lancement, le serveur va chercher 10 albums et les 4 premiers titres de chacun (donc 40 titres), puis les écrit dans `cache_titres.json`. Aux lancements suivants, il lit juste ce fichier.

Pour avoir des données fraîches : arrêter le serveur, supprimer `cache_titres.json`, relancer.

## Quand ça se passe mal

- Spotify met trop de temps à répondre : au bout de 10 secondes, j'abandonne l'appel.
- Le token a expiré (il dure 1h) : Spotify répond 401, le serveur en redemande un tout seul et réessaie une fois.
- Trop de requêtes (429) : je regarde le header `Retry-After`. Si c'est 10 secondes ou moins, j'attends et je réessaie. Si c'est plus long, j'arrête et j'affiche une erreur propre.
- Dans tous les cas, le site affiche une page d'erreur au lieu de planter.

J'ai testé le cas du token expiré en mettant un faux token en mémoire. Je n'ai pas pu tester le 429 en vrai, je ne peux pas le provoquer quand je veux.

## Postman

Dans le dépôt il y a `JUL-Spotify.postman_collection.json`. Dedans : récupérer un token, chercher JUL, sa discographie, le détail de « La bandite » et 3 cas d'erreur.

Le secret n'est pas dans le fichier. Pour la rejouer, il faut créer un environnement Postman avec `client_id` et `client_secret` (la variable `token` se remplit après la requête 1).

## Les fichiers

- `main.go` : le serveur, les routes et la page d'erreur
- `models.go` : la struct `Titre` et ses petites méthodes (durée, année, date en français, feat)
- `spotify.go` : le token et les appels à l'API
- `cache.go` : le cache local
- `filtres.go` : les filtres
- `env.go` : la lecture du `.env`
- `templates/` : les pages `.gohtml`
- `styles/` : le CSS
- `NOTES.md` : mes notes et mes galères, avec ce que j'ai testé

## Ce que j'améliorerais avec plus de temps

- Aller au-delà de 10 albums avec la pagination.
- Ajouter une recherche par nom de titre.
- Tester le 429 pour de vrai.
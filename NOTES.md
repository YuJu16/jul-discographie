# NOTES.md – Projet JUL (Spotify + Go)

Ce fichier garde une trace de ce que j'ai fait, de ce que j'ai compris et de mes galères.

## Partie I – Postman

Variables d'environnement Postman (`Spotify`) : `client_id`, `client_secret`, `token`, `artist_id`, `track_id`.
Le secret n'est jamais écrit dans la collection exportée.

| Étape | Requête | Ce que j'ai compris |
|---|---|---|
| 1 | `POST accounts.spotify.com/api/token` (Basic Auth, `grant_type=client_credentials`) | Je récupère un token valable 1h |
| 2 | `GET /v1/search?q=JUL&type=artist` | Je trouve l'ID de JUL : `3IW7ScrzXmPvZhB27hmfgy` |
| 3 | `GET /v1/artists/{id}/albums` | La discographie (albums et singles) |
| 4 | `GET /v1/search?q=La bandite&type=track` puis `GET /v1/tracks/{id}` | Le détail d'un titre |

Le token se met dans l'en-tête `Authorization: Bearer <token>`.

### Les 3 erreurs testées

- **E1 – mauvais secret** : code `À COMPLÉTER`, message `À COMPLÉTER`.
- **E2 – token invalide ou expiré** : code `À COMPLÉTER` (401), message `À COMPLÉTER`.
- **E3 – mauvais `grant_type`** : `unsupported_grant_type`.

### Galères Postman

- Je cliquais sur Mock servers au lieu d'Environments.
- Mon `client_secret` était identique au `client_id` : erreur `invalid_client` (400).
- Après avoir dupliqué une requête, l'authentification avait disparu : 401.
- `{{track_id}}` restait en rouge (variable non sauvegardée) : 403.
- Il manquait le `0` au début d'un ID.
- Le token expire au bout d'1h : j'ai dû relancer la requête 1.

## Partie II – Le serveur Go

### Ce qui est fait

- Serveur sur `localhost:8080`, bibliothèque standard uniquement, pas de JavaScript.
- Routes : `/` (liste), `/tracks/{id}` (détail), `/styles/` (CSS).
- Pages d'erreur personnalisées : 404 (route inconnue ou titre introuvable) et 503 (Spotify ne répond pas).
- Filtres écrits à la main en Go (`filtres.go`) : année, type (single/album), feat (oui/non). Les menus gardent le choix après l'envoi.
- Templates `html/template` en `.gohtml` : `index`, `detail`, `erreur`.
- Identifiants lus dans `.env` (`env.go`), jamais dans le code.

### Fichiers

| Fichier | Rôle |
|---|---|
| `main.go` | Serveur, routes, page d'erreur |
| `models.go` | Struct `Titre` et ses méthodes (durée, année, date FR, feat) |
| `spotify.go` | Token, appels à l'API, renouvellement et gestion du 429 |
| `cache.go` | Cache `cache_titres.json` |
| `filtres.go` | Lecture et application des filtres |
| `env.go` | Lecture du `.env` |

### Gestion des erreurs

- **Timeout** : le client HTTP coupe après 10 secondes.
- **401 (token expiré)** : on vide le token en mémoire, on en redemande un, on réessaie une fois.
- **429 (trop de requêtes)** : on lit `Retry-After`. Si c'est 10 secondes ou moins, on attend et on réessaie. Sinon on abandonne avec un message clair.
- **Spotify indisponible** : le serveur ne plante pas, la page affiche une erreur 503.

### Tests de robustesse

- **401** : j'ai mis un faux token en mémoire. Résultat : le serveur a reçu un 401, a redemandé un token tout seul et a récupéré les 40 titres.
- **503** : j'ai cassé volontairement mon secret dans `.env`. Résultat : page 503 affichée, sans crash.
- **429** : non testé en vrai, car je ne peux pas le provoquer à volonté. Le mécanisme est écrit et je sais l'expliquer.

### Choix techniques

- **Cache `cache_titres.json`** : je ne retape pas Spotify à chaque lancement, ce qui évite le 429. Pour rafraîchir, je supprime le fichier. Il est dans le `.gitignore`.
- **10 albums, 4 titres par album** : la recherche est limitée à 10 résultats depuis les changements Spotify de 2026, et ça limite le nombre de requêtes.
- **Pas de popularité ni de followers** : Spotify les a retirés.

### Galères Go

- **`spotify.go` vide** : j'avais écrit le code sans sauvegarder, Go disait « expected 'package', found 'EOF' ». Le point blanc sur l'onglet VS Code veut dire « pas sauvegardé ».
- **Erreur 400 au token** : mon `.env` avait des guillemets autour des valeurs. J'ai nettoyé guillemets et espaces dans `chargerEnv()` avec `strings.Trim`.
- **`:=` qui redéclare** : « no new variables on left side ». Si la variable existe déjà, il faut `=`.
- **Appel dans la boucle** : j'avais mis `recupererTitres` dans la boucle des albums, ça faisait trop de requêtes. Je l'ai sorti de la boucle.
- **`chargerEnv()` dans la mauvaise fonction** : je l'avais mis dans `fauxTitres()`, il fallait l'appeler dans `main()`.
- **Reste de l'ancienne fonction** dans `cache.go` après un copier-coller : ça ne compilait plus.
- **Collection Postman avec un espace dans le nom** : elle avait déjà été commit. J'ai vérifié avec `git show` qu'elle ne contenait aucun secret.

## Sécurité

- `.env` et `cache_titres.json` sont dans `.gitignore`.
- `.env.example` contient uniquement des valeurs bidon.
- Vérifié : aucun secret ni token dans le code, les notes, le README ou la collection Postman.

## Ce qui reste / pistes

- Pagination pour afficher plus que 10 albums.
- Tester le 429 en vrai si possible.
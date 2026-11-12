package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// On laisse 10 secondes max à Spotify pour répondre, sinon on coupe
var clientHTTP = &http.Client{Timeout: 10 * time.Second}

// On cible uniquement le token dans le JSON renvoyé par Spotify
type reponseToken struct {
	AccessToken string `json:"access_token"`
}

// demanderToken fait le taf de notre test "1 - Token" sur Postman
func demanderToken() (string, error) {
	// Le corps de la requête (plus simple que de s'embêter avec url.Values)
	corps := strings.NewReader("grant_type=client_credentials")

	// On prépare la requête POST vers Spotify
	demande, erreur := http.NewRequest("POST", "https://accounts.spotify.com/api/token", corps)
	if erreur != nil {
		return "", erreur
	}

	// L'équivalent de l'onglet "Authorization" (Basic Auth) sur Postman
	demande.SetBasicAuth(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	demande.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// On envoie la requête
	reponse, erreur := clientHTTP.Do(demande)
	if erreur != nil {
		return "", erreur
	}
	defer reponse.Body.Close()

	// Si ce n'est pas un code 200 (ex: 400, 401...), Spotify nous a recalé
	if reponse.StatusCode != http.StatusOK {
		return "", fmt.Errorf("erreur API Spotify : statut %d", reponse.StatusCode)
	}

	// On décode le JSON pour chopper la valeur de l'access_token
	var token reponseToken
	if erreur := json.NewDecoder(reponse.Body).Decode(&token); erreur != nil {
		return "", erreur
	}

	return token.AccessToken, nil
}

// L'ID de JUL qu'on a trouvé tout à l'heure sur Postman
const idJul = "3IW7ScrzXmPvZhB27hmfgy"

// Les images (pochettes). Spotify en renvoie plusieurs, on garde juste l'URL
type Image struct {
	URL string `json:"url"`
}

// Ce qu'on garde d'un album ou d'un single (on zappe le reste des infos inutiles)
type Album struct {
	ID         string  `json:"id"`
	Nom        string  `json:"name"`
	DateSortie string  `json:"release_date"`
	Type       string  `json:"album_type"`
	Images     []Image `json:"images"`
}

// Le format de la réponse de Spotify (notre liste d'albums est rangée dans "items")
type reponseAlbums struct {
	Items []Album `json:"items"`
}

// Fonction outil pour faire n'importe quel appel GET à Spotify sans réécrire le code
func appelSpotify(token string, adresse string, cible any) error {
	// Préparation de la requête GET
	demande, erreur := http.NewRequest("GET", adresse, nil)
	if erreur != nil {
		return erreur
	}

	// L'équivalent de l'onglet "Authorization" > "Bearer Token" sur Postman
	demande.Header.Set("Authorization", "Bearer "+token)

	// On lance l'appel
	reponse, erreur := clientHTTP.Do(demande)
	if erreur != nil {
		return erreur
	}
	defer reponse.Body.Close()

	// Si on a spammé l'API (erreur 429), on chope le temps d'attente imposé
	if reponse.StatusCode == http.StatusTooManyRequests {
		attente := reponse.Header.Get("Retry-After")
		return fmt.Errorf("trop de requêtes, Spotify nous bloque. Réessayer dans %s sec", attente)
	}

	// Si c'est pas un code 200, il y a une embrouille
	if reponse.StatusCode != http.StatusOK {
		return fmt.Errorf("erreur API Spotify : statut %d", reponse.StatusCode)
	}

	// On injecte le JSON directement dans notre variable "cible"
	return json.NewDecoder(reponse.Body).Decode(cible)
}

// Récupère les 10 derniers albums/singles de JUL
func recupererAlbums(token string) ([]Album, error) {
	url := "https://api.spotify.com/v1/artists/" + idJul + "/albums?include_groups=album,single&market=FR&limit=10"

	var donnees reponseAlbums

	// On fait l'appel et on range le JSON décodé dans "donnees"
	if erreur := appelSpotify(token, url, &donnees); erreur != nil {
		return nil, erreur
	}

	return donnees.Items, nil
}

// On garde juste le nom de l'artiste (Spotify renvoie plein de trucs inutiles ici)
type Artiste struct {
	Nom string `json:"name"`
}

// Un son (track) tel qu'il sort de l'API
type Piste struct {
	ID       string    `json:"id"`
	Nom      string    `json:"name"`
	DureeMs  int       `json:"duration_ms"`
	Artistes []Artiste `json:"artists"`
	Liens    struct {
		Spotify string `json:"spotify"` // Pratique pour le bouton "Écouter sur Spotify"
	} `json:"external_urls"`
}

// Le format de la réponse (la liste des sons est rangée dans "items")
type reponsePistes struct {
	Items []Piste `json:"items"`
}

// Récupère les 4 premiers sons de chaque album de notre liste
func recupererTitres(token string, albums []Album) ([]Titre, error) {
	var titres []Titre // Notre tableau final qui va stocker tous les sons

	// On boucle sur chaque album qu'on a récupéré avant
	for _, album := range albums {
		url := "https://api.spotify.com/v1/albums/" + album.ID + "/tracks?market=FR&limit=4"

		var donnees reponsePistes

		// On réutilise notre super fonction outil de tout à l'heure !
		if erreur := appelSpotify(token, url, &donnees); erreur != nil {
			return titres, erreur
		}

		// La pochette en meilleure qualité, c'est toujours la première du tableau
		pochette := ""
		if len(album.Images) > 0 {
			pochette = album.Images[0].URL
		}

		// On boucle sur les sons de cet album
		for _, piste := range donnees.Items {

			// En JS on ferait un .map(), en Go on boucle pour extraire juste les noms des artistes
			var nomsArtistes []string
			for _, artiste := range piste.Artistes {
				nomsArtistes = append(nomsArtistes, artiste.Nom)
			}

			// On "mixe" les infos de l'album et de la piste pour créer notre objet Titre final
			titres = append(titres, Titre{
				ID:          piste.ID,
				Nom:         piste.Nom,
				Artistes:    nomsArtistes,
				DateSortie:  album.DateSortie,
				DureeMs:     piste.DureeMs,
				Type:        album.Type,
				Pochette:    pochette,
				Album:       album.Nom,
				LienSpotify: piste.Liens.Spotify,
			})
		}
	}

	return titres, nil
}

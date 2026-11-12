package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Timeout de 10 secondes max pour les requêtes vers Spotify
var clientHTTP = &http.Client{Timeout: 10 * time.Second}

// Variable globale pour stocker notre token et éviter de spammer la route d'auth
var tokenActuel string

type reponseToken struct {
	AccessToken string `json:"access_token"`
}

// Génère un nouveau token via les identifiants du .env
func demanderToken() (string, error) {
	corps := strings.NewReader("grant_type=client_credentials")

	demande, erreur := http.NewRequest("POST", "https://accounts.spotify.com/api/token", corps)
	if erreur != nil {
		return "", erreur
	}

	demande.SetBasicAuth(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	demande.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	reponse, erreur := clientHTTP.Do(demande)
	if erreur != nil {
		return "", erreur
	}
	defer reponse.Body.Close()

	if reponse.StatusCode != http.StatusOK {
		return "", fmt.Errorf("erreur API Spotify : statut %d", reponse.StatusCode)
	}

	var token reponseToken
	if erreur := json.NewDecoder(reponse.Body).Decode(&token); erreur != nil {
		return "", erreur
	}

	return token.AccessToken, nil
}

// Renvoie le token en mémoire, ou en demande un neuf si c'est vide
func obtenirToken() (string, error) {
	if tokenActuel != "" {
		return tokenActuel, nil
	}

	token, erreur := demanderToken()
	if erreur != nil {
		return "", erreur
	}

	tokenActuel = token
	return token, nil
}

// ID Spotify de JUL
const idJul = "3IW7ScrzXmPvZhB27hmfgy"

type Image struct {
	URL string `json:"url"`
}

type Album struct {
	ID         string  `json:"id"`
	Nom        string  `json:"name"`
	DateSortie string  `json:"release_date"`
	Type       string  `json:"album_type"`
	Images     []Image `json:"images"`
}

type reponseAlbums struct {
	Items []Album `json:"items"`
}

// Fonction générique pour faire nos requêtes GET avec retry automatique (token expiré / rate limit)
func appelSpotify(adresse string, cible any) error {
	// On s'autorise 2 essais max par appel
	for essai := 1; essai <= 2; essai++ {
		token, erreur := obtenirToken()
		if erreur != nil {
			return erreur
		}

		demande, erreur := http.NewRequest("GET", adresse, nil)
		if erreur != nil {
			return erreur
		}

		demande.Header.Set("Authorization", "Bearer "+token)

		reponse, erreur := clientHTTP.Do(demande)
		if erreur != nil {
			return erreur
		}

		// Erreur 401 : notre token a expiré. On le reset pour que la boucle en redemande un neuf.
		if reponse.StatusCode == http.StatusUnauthorized && essai == 1 {
			reponse.Body.Close()
			fmt.Println("🔄 Token expiré, refresh en cours...")
			tokenActuel = ""
			continue
		}

		// Erreur 429 : On a tapé l'API trop vite. On regarde le Header Retry-After.
		if reponse.StatusCode == http.StatusTooManyRequests {
			reponse.Body.Close()
			attente, _ := strconv.Atoi(reponse.Header.Get("Retry-After"))

			// Si la punition est de 10s ou moins, on met le code en pause (Sleep) et on réessaie.
			if essai == 1 && attente > 0 && attente <= 10 {
				fmt.Println("⏳ Rate limit atteint, pause de", attente, "secondes...")
				time.Sleep(time.Duration(attente) * time.Second)
				continue
			}
			return fmt.Errorf("trop de requêtes (429), réessayer dans %d secondes", attente)
		}

		if reponse.StatusCode != http.StatusOK {
			reponse.Body.Close()
			return fmt.Errorf("erreur API Spotify : statut %d", reponse.StatusCode)
		}

		// Décodage du JSON directement dans le pointeur de structure passé en paramètre
		erreur = json.NewDecoder(reponse.Body).Decode(cible)
		reponse.Body.Close()
		return erreur
	}

	return fmt.Errorf("échec de l'appel Spotify après les retrys")
}

// Récupère les 10 derniers albums/singles
func recupererAlbums() ([]Album, error) {
	url := "https://api.spotify.com/v1/artists/" + idJul + "/albums?include_groups=album,single&market=FR&limit=10"

	var donnees reponseAlbums
	if erreur := appelSpotify(url, &donnees); erreur != nil {
		return nil, erreur
	}
	return donnees.Items, nil
}

type Artiste struct {
	Nom string `json:"name"`
}

type Piste struct {
	ID       string    `json:"id"`
	Nom      string    `json:"name"`
	DureeMs  int       `json:"duration_ms"`
	Artistes []Artiste `json:"artists"`
	Liens    struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
}

type reponsePistes struct {
	Items []Piste `json:"items"`
}

// Récupère les 4 premiers sons pour chaque album de la liste
func recupererTitres(albums []Album) ([]Titre, error) {
	var titres []Titre

	for _, album := range albums {
		url := "https://api.spotify.com/v1/albums/" + album.ID + "/tracks?market=FR&limit=4"

		var donnees reponsePistes
		if erreur := appelSpotify(url, &donnees); erreur != nil {
			return titres, erreur
		}

		pochette := ""
		if len(album.Images) > 0 {
			pochette = album.Images[0].URL // On chope la version HQ (toujours la première)
		}

		for _, piste := range donnees.Items {
			var nomsArtistes []string
			for _, artiste := range piste.Artistes {
				nomsArtistes = append(nomsArtistes, artiste.Nom)
			}

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

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

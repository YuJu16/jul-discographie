package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// Handler de l'accueil : récupère la requête (demande) et écrit la vue dans (reponse)
func accueil(reponse http.ResponseWriter, demande *http.Request) {
	// Chargement du fichier HTML
	page, erreur := template.ParseFiles("templates/index.gohtml")
	if erreur != nil { // Si erreur n'est pas vide (nil), fichier introuvable -> erreur 500 et on stoppe (return)
		http.Error(reponse, "Impossible de charger la page", 500)
		return
	}

	// Envoi du HTML au navigateur (nil = aucune donnée injectée pour l'instant, futur emplacement des données Spotify)
	page.Execute(reponse, nil)
}

func main() {
	// Sert les fichiers statiques (CSS) depuis le dossier styles/
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))
	//petit explication pour comprendre la ligne au dessus :
	// http.Handle("/styles/", ...) : définit une route pour servir les fichiers statiques
	// http.StripPrefix("/styles/", ...) : retire le préfixe "/styles/" de l'URL avant de servir le fichier
	// http.FileServer(http.Dir("styles")) : crée un serveur de fichiers pour le dossier "styles"

	// Route principale "/" reliée à la fonction accueil
	http.HandleFunc("/", accueil)
	fmt.Println("Serveur lancé sur http://localhost:8080")

	// Lancement du serveur sur le port 8080 (nil = routeur par défaut de Go)
	erreur := http.ListenAndServe(":8080", nil)
	if erreur != nil {
		fmt.Println("Erreur serveur :", erreur)
	}
}

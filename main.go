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

// fauxTitre renvoie des titres d'exemple en attendant l'api Spotify
func fauxTitre() []Titre {
	return []Titre{
		{ID: "1", Nom: "Titre 1", Artist: []string{"Artiste A"}, DateSortie: "2024-05-10", DureeMS: 222000, Type: "single"},
		{ID: "2", Nom: "Titre 2", Artist: []string{"Artiste A", "Artiste B"}, DateSortie: "2023-11-03", DureeMS: 198000, Type: "album"},
		{ID: "3", Nom: "Titre 3", Artist: []string{"Artiste A"}, DateSortie: "2022-02-18", DureeMS: 245000, Type: "album"},
	}
}

package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// fauxTitres renvoie des titres d'exemple, en attendant Spotify
func fauxTitres() []Titre {
	return []Titre{
		{ID: "1", Nom: "Titre 1", Artistes: []string{"Artiste A"}, DateSortie: "2024-05-10", DureeMs: 222000, Type: "single"},
		{ID: "2", Nom: "Titre 2", Artistes: []string{"Artiste A", "Artiste B"}, DateSortie: "2023-11-03", DureeMs: 198000, Type: "album"},
		{ID: "3", Nom: "Titre 3", Artistes: []string{"Artiste A"}, DateSortie: "2022-02-18", DureeMs: 245000, Type: "album"},
	}
}

// Handler de l'accueil : récupère la requête (demande) et écrit la vue dans (reponse)
func accueil(reponse http.ResponseWriter, demande *http.Request) {
	// Chargement du fichier HTML
	page, erreur := template.ParseFiles("templates/index.gohtml")
	if erreur != nil { // si erreur n'est pas vide (nil), fichier introuvable -> erreur 500 et on stoppe
		http.Error(reponse, "Impossible de charger la page", http.StatusInternalServerError)
		return
	}

	// La liste des titres (faux pour l'instant)
	titres := []Titre{} // test : liste vide

	// Envoi du HTML au navigateur, avec les titres injectés dans le template
	erreur = page.Execute(reponse, titres)
	if erreur != nil {
		http.Error(reponse, "Impossible d'afficher la page", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Route principale "/" reliée à la fonction accueil
	http.HandleFunc("/", accueil)

	// Les fichiers du dossier styles sont servis sur /styles/
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))

	fmt.Println("Serveur lancé sur http://localhost:8080")

	// Lancement du serveur sur le port 8080 (nil = routeur par défaut de Go)
	erreur := http.ListenAndServe(":8080", nil)
	if erreur != nil {
		fmt.Println("Erreur serveur :", erreur)
	}
}

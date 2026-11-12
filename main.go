package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// Les infos qu'on va envoyer à notre page d'erreur HTML pour qu'elle soit dynamique
type PageErreur struct {
	Code    int
	Titre   string
	Message string
}

// Fonction outil pour afficher notre belle page d'erreur au lieu d'une page blanche moche
func afficherErreur(reponse http.ResponseWriter, code int, titre string, message string) {
	page, erreur := template.ParseFiles("templates/erreur.gohtml")
	if erreur != nil {
		// Si notre fichier d'erreur est cassé ou introuvable, on renvoie un texte brut de secours
		http.Error(reponse, titre, code)
		return
	}

	// Important en Go : il faut envoyer le code (ex: 404) AVANT d'envoyer le HTML
	reponse.WriteHeader(code)
	page.Execute(reponse, PageErreur{Code: code, Titre: titre, Message: message})
}

// Handler de la page d'accueil
func accueil(reponse http.ResponseWriter, demande *http.Request) {
	// Piège classique en Go : la route "/" agit comme un joker et attrape tout.
	// Donc si l'URL n'est pas EXACTEMENT "/", ça veut dire que la page n'existe pas -> 404.
	if demande.URL.Path != "/" {
		afficherErreur(reponse, http.StatusNotFound, "Page introuvable", "Cette page n'existe pas, le J l'a volée.")
		return
	}

	// Chargement et affichage de l'accueil
	page, erreur := template.ParseFiles("templates/index.gohtml")
	if erreur != nil {
		http.Error(reponse, "Impossible de charger l'accueil", http.StatusInternalServerError)
		return
	}

	// On injecte notre variable globale "titresJul" dans le HTML
	if erreur := page.Execute(reponse, titresJul); erreur != nil {
		http.Error(reponse, "Erreur d'affichage", http.StatusInternalServerError)
	}
}

// Handler pour la page détail d'un titre précis (ex: /tracks/12345)
func detail(reponse http.ResponseWriter, demande *http.Request) {
	// On récupère l'ID depuis l'URL (l'équivalent de req.params.id en Express)
	id := demande.PathValue("id")

	page, erreur := template.ParseFiles("templates/detail.gohtml")
	if erreur != nil {
		http.Error(reponse, "Impossible de charger la page détail", http.StatusInternalServerError)
		return
	}

	// On fouille dans notre mémoire (titresJul) pour trouver le son qui correspond à l'ID
	for _, titre := range titresJul {
		if titre.ID == id {
			// Titre trouvé : on l'envoie au template et on arrête la fonction (return)
			page.Execute(reponse, titre)
			return
		}
	}

	// Si on a fini la boucle sans faire de return, c'est qu'on a rien trouvé
	afficherErreur(reponse, http.StatusNotFound, "Titre introuvable", "Ce son n'est plus dispo.")
}

func main() {
	// 1. Setup de l'environnement (pour cacher la clé Spotify)
	chargerEnv()

	// 2. Chargement des données au démarrage (depuis le cache ou l'API)
	if erreur := chargerTitres(); erreur != nil {
		fmt.Println("❌ Aïe, erreur de chargement des titres :", erreur)
	}

	// 3. Déclaration de nos routes
	http.HandleFunc("/", accueil)
	// Nouvelle syntaxe Go 1.22 super pratique pour les variables dans l'URL
	http.HandleFunc("/tracks/{id}", detail)

	// 4. On sert les fichiers statiques (notre CSS) pour qu'ils soient accessibles sur /styles/
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))

	fmt.Println("🚀 Serveur lancé sur http://localhost:8080")

	// 5. On allume le serveur
	if erreur := http.ListenAndServe(":8080", nil); erreur != nil {
		fmt.Println("❌ Erreur critique du serveur :", erreur)
	}
}

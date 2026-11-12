package main

import "net/http"

// Les filtres choisis par l'utilisateur (récupérés dans l'URL)
type Filtres struct {
	Annee string
	Type  string
	Feat  string
}

// L'objet complet qu'on va envoyer à notre template index.gohtml
type PageAccueil struct {
	Titres  []Titre  // La liste des sons (après avoir appliqué le filtre)
	Annees  []string // Pour remplir le menu déroulant des années sans coder en dur
	Filtres Filtres  // On renvoie les choix actuels pour que le menu reste sur la bonne valeur après le clic
}

// Chopons les paramètres dans l'URL (ex: ?year=2025&type=album&feat=oui)
func lireFiltres(demande *http.Request) Filtres {
	// demande.URL.Query().Get() c'est exactement le req.query de Express en JS
	return Filtres{
		Annee: demande.URL.Query().Get("year"),
		Type:  demande.URL.Query().Get("type"),
		Feat:  demande.URL.Query().Get("feat"),
	}
}

// Le moteur de recherche : on trie les sons selon ce que l'utilisateur a demandé
func filtrerTitres(titres []Titre, filtres Filtres) []Titre {
	var resultat []Titre

	for _, titre := range titres {
		// Si un filtre est vide (""), ça veut dire que l'utilisateur veut "Tout" voir.
		// Le mot-clé "continue" permet de zapper le son actuel et de passer direct au suivant.

		if filtres.Annee != "" && titre.Annee() != filtres.Annee {
			continue // Pas la bonne année ? On zappe.
		}
		if filtres.Type != "" && titre.Type != filtres.Type {
			continue // Pas le bon type (album/single) ? On zappe.
		}
		if filtres.Feat == "oui" && !titre.AFeat() {
			continue // On veut des feats mais y'en a pas ? On zappe.
		}
		if filtres.Feat == "non" && titre.AFeat() {
			continue // On ne veut PAS de feat mais y'en a un ? On zappe.
		}

		// Si le son a survécu à tous les 'continue', c'est qu'il matche ! On le garde.
		resultat = append(resultat, titre)
	}

	return resultat
}

// Extrait toutes les années uniques pour fabriquer notre menu déroulant
func listerAnnees(titres []Titre) []string {
	var annees []string

	for _, titre := range titres {
		annee := titre.Annee()

		// En JS on ferait juste un [...new Set(tableau)] pour virer les doublons.
		// En Go, on le fait à l'ancienne : on vérifie manuellement si on l'a déjà ajouté.
		dejaVue := false
		for _, existante := range annees {
			if existante == annee {
				dejaVue = true
				break // Petite opti : dès qu'on l'a trouvé, on arrête de chercher
			}
		}

		if !dejaVue {
			annees = append(annees, annee)
		}
	}

	return annees
}

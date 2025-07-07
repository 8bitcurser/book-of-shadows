package wizard

import (
	"book-of-shadows/models"
	"book-of-shadows/storage"
	"book-of-shadows/views"
	"log"
	"net/http"
	"strings"
)

func HandleBaseStep(w http.ResponseWriter, r *http.Request) {
	component := views.BaseStep(nil)
	err := component.Render(r.Context(), w)
	if err != nil {
		log.Println(err)
	}
}

func HandleAttrStep(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/wizard/attributes/")
	cm := storage.NewInvestigatorCookieConfig()
	investigator, err := cm.GetInvestigatorCookie(r, key)
	if err != nil {
		log.Println(err)
		http.Error(w, "Investigator not found", http.StatusNotFound)
		return
	}
	components := views.AttrStep(investigator)
	err = components.Render(r.Context(), w)
	if err != nil {
		log.Println(err)
	}
}

func HandleSkillForm(w http.ResponseWriter, r *http.Request) {
	investigator := &models.Investigator{}
	components := views.SkillStep(investigator)
	err := components.Render(r.Context(), w)
	if err != nil {
		log.Println(err)
	}
}

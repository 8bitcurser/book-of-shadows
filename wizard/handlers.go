package wizard

import (
	"book-of-shadows/storage"
	"book-of-shadows/views"
	"log"
	"net/http"
	"strings"
)

func HandleBaseStep(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/wizard/base/")
	component := views.BaseStep(nil)
	if key != "" && key != "new" {
		cm := storage.NewInvestigatorCookieConfig()
		investigator, _ := cm.GetInvestigatorCookie(r, key)
		component = views.BaseStep(investigator)
	}
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
	key := strings.TrimPrefix(r.URL.Path, "/wizard/skills/")
	cm := storage.NewInvestigatorCookieConfig()
	investigator, err := cm.GetInvestigatorCookie(r, key)
	if err != nil {
		log.Println(err)
		http.Error(w, "Investigator not found", http.StatusNotFound)
		return
	}
	components := views.SkillStep(investigator)
	err = components.Render(r.Context(), w)
	if err != nil {
		log.Println(err)
	}
}

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/eugenefoxx/savePictureAOIMertec/pkg/logging"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	//"google.golang.org/genproto/googleapis/cloud/metastore/logging/v1"
)

var (
	qr         = "2AQC622457019707"
	sourceFile = "C:/Users/Евгений/Code/github.com/eugenefoxx/savePictureAOIMertec/internal/folder/metrec/SmallTop.jpg"
	checkDIR   = "C:/Users/Евгений/Code/github.com/eugenefoxx/savePictureAOIMertec/internal/folder/saveFoto/"
	//checkDIR   = "/home/eugenearch/Code/github.com/eugenefoxx/savePictureAOIMertec/internal/folder/saveFoto/"
	//sourceFile = "/home/eugenearch/Code/github.com/eugenefoxx/savePictureAOIMertec/internal/folder/metrec/SmallTop.jpg"
	//destinationFile = checkDIR + qr //"/home/eugenearch/Code/github.com/eugenefoxx/savePictureAOIMertec/internal/folder/saveFoto/"
)

type Handler struct {
	*chi.Mux
	logger *logging.Logger
}

func init() {
	logging.Init()
	//	logger := logging.GetLogger()

	/*	err := godotenv.Load()
		if err != nil {
			logger.Fatal(err.Error)

			}*/
}

func main() {
	//logger := logging.GetLogger()

	r := chi.NewRouter()
	h := &Handler{
		Mux: chi.NewMux(),
	}
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/savefotoAOI", h.PagesavefotoAOI())
		r.Post("/savefotoAOI", h.SavefotoAOI())

	})
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "web"))
	FileServer(r, "/web", filesDir)
	//open.StartWith("http://localhost:3010/savefotoAOI", "chrome.exe")
	http.ListenAndServe(":3010", r)

}

func (h *Handler) PagesavefotoAOI() http.HandlerFunc {

	//tpl, err := template.New("").ParseFiles(viper.GetString("web.html"))
	tpl, err := template.New("").ParseFiles("C:/Users/Евгений/Code/github.com/eugenefoxx/savePictureAOIMertec/internal/web/index.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		err = tpl.ExecuteTemplate(w, "index.html", nil)

	}

}

type Data struct {
	Responce  string
	Directory string
	File      string
}

func (h *Handler) SavefotoAOI() http.HandlerFunc {
	type searchBy struct {
		QR string `json:"qr"`
	}
	//tpl, err := template.New("").ParseFiles(viper.GetString("web.html"))
	tpl, err := template.New("").ParseFiles("C:/Users/Евгений/Code/github.com/eugenefoxx/savePictureAOIMertec/internal/web/index.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {

		search := &searchBy{}
		search.QR = r.FormValue("qr")
		if search.QR == "" {
			fmt.Println("Пусто")
		}
		qr := search.QR

		// copy file
		input, err := ioutil.ReadFile(sourceFile)
		if err != nil {
			h.logger.Errorf(err.Error())
			return
		}

		checkingDIR := checkDIR + qr + "/" //filepath.Dir(checkDIR + qr + "/")
		fmt.Println(checkingDIR)
		if _, err := os.Stat(checkingDIR); os.IsNotExist(err) {
			os.Mkdir(checkingDIR, 0755)
		}

		destinationFile := checkingDIR
		fmt.Println(destinationFile)
		currentDate := time.Now()
		strcurrentTime := currentDate.Format("2006-01-02 15-04-05")

		//err = ioutil.WriteFile("C:/Users/Евгений/Code/github.com/eugenefoxx/savePictureAOIMertec/internal/folder/saveFoto/"+qr+"/"+qr+"-time-"+strcurrentTime+".jpg", input, 0) //0644
		err = ioutil.WriteFile(destinationFile+qr+"-time-"+strcurrentTime+".jpg", input, 0)
		if err != nil {
			h.logger.Errorln("Error creating", destinationFile)
			pathtosavefile := fmt.Sprintf("Error creating", destinationFile)

			params := []Data{
				{
					Responce: pathtosavefile,
				},
			}
			data := map[string]interface{}{
				"GetParams": params,
			}
			tpl.ExecuteTemplate(w, "index.html", data)
			h.logger.Errorf(err.Error())
			//return
		}
		pathtosavefile := destinationFile
		directory := qr
		file := qr + "-time-" + strcurrentTime + ".jpg"
		params := []Data{
			{
				Responce:  pathtosavefile,
				Directory: directory,
				File:      file,
			},
		}
		data := map[string]interface{}{
			"GetParams": params,
		}
		qr = ""
		err = tpl.ExecuteTemplate(w, "index.html", data)

	}
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, CacheControlWrapper(http.FileServer(root)))
		fs.ServeHTTP(w, r)
	})
}

func CacheControlWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=100") // "max-age=2592000" 30 days
		h.ServeHTTP(w, r)
	})
}

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

	"github.com/eugenefoxx/savePictureAOIMertec/internal/config"
	"github.com/eugenefoxx/savePictureAOIMertec/pkg/logging"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"github.com/skratchdot/open-golang/open"
	//"google.golang.org/genproto/googleapis/cloud/metastore/logging/v1"
)

type Handler struct {
	*chi.Mux
	//logger logging.Logger
}

func init() {
	logging.Init()
	logger := logging.GetLogger()

	//PortPath := os.Getenv("port")
	logger.Printf("Запускаем конфигурацию переменной окружения.")
	err := godotenv.Load()
	if err != nil {
		logger.Fatalf("No .env file found", err.Error)

	}
}

func main() {
	//logger := logging.GetLogger()
	conf := config.New()

	r := chi.NewRouter()

	h := &Handler{
		Mux: chi.NewMux(),
		//Mux: chi.NewMux(),
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
	filesDir := http.Dir(filepath.Join(workDir, conf.HTTP.WorkDirWEB))
	FileServer(r, "/web", filesDir)
	open.StartWith(conf.HTTP.StartURL, "chrome.exe")
	http.ListenAndServe(conf.HTTP.Port, r)

}

func (h *Handler) PagesavefotoAOI() http.HandlerFunc {
	conf := config.New()
	logger := logging.GetLogger()

	//tpl, err := template.New("").ParseFiles(viper.GetString("web.html"))
	tpl, err := template.New("").ParseFiles(conf.HTTP.IndexHTML)
	if err != nil {
		logger.Fatalf(err.Error())
		//panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		err = tpl.ExecuteTemplate(w, conf.HTTP.IndexExecute, nil)
		if err != nil {
			logger.Errorf("%s", err.Error())
		}

	}

}

type DataResponse struct {
	Response  string
	Directory string
	File      string
}

func (h *Handler) SavefotoAOI() http.HandlerFunc {

	conf := config.New()
	logger := logging.GetLogger()
	type searchBy struct {
		QR string `json:"qr"`
	}
	//tpl, err := template.New("").ParseFiles(viper.GetString("web.html"))
	tpl, err := template.New("").ParseFiles(conf.HTTP.IndexHTML)
	if err != nil {
		logger.Fatalf("%s", err.Error())

	}
	return func(w http.ResponseWriter, r *http.Request) {

		//conf.Files.SourceFile

		search := &searchBy{}
		search.QR = r.FormValue("qr")
		qr := search.QR

		// copy file from source Metrec
		//input, err := ioutil.ReadFile(sourceFile)
		input, err := ioutil.ReadFile(conf.Files.SourceFile)
		if err != nil {
			logger.Errorf(err.Error())
			return
		}

		checkingDIR := conf.Files.HomeSavingDir + qr + "/" //filepath.Dir(checkDIR + qr + "/")
		//fmt.Println(checkingDIR)
		if _, err := os.Stat(checkingDIR); os.IsNotExist(err) {
			os.Mkdir(checkingDIR, 0755)
		}

		destinationFile := checkingDIR
		fmt.Println(destinationFile)

		currentDate := time.Now()
		strcurrentTime := currentDate.Format("2006-01-02 15-04-05")

		//err = ioutil.WriteFile("C:/Users/Евгений/Code/github.com/eugenefoxx/savePictureAOIMertec/internal/folder/saveFoto/"+qr+"/"+qr+"-time-"+strcurrentTime+".jpg", input, 0) //0644
		err = ioutil.WriteFile(destinationFile+qr+"-time-"+strcurrentTime+".jpg", input, 0644)
		if err != nil {
			logger.Errorln("Error creating", destinationFile)
			/*pathtosavefile := fmt.Sprintf("Error creating", destinationFile)

			params := []DataResponse{
				{
					Response: pathtosavefile,
				},
			}
			data := map[string]interface{}{
				"GetParams": params,
			}
			tpl.ExecuteTemplate(w, "index.html", data)*/
			logger.Errorf(err.Error())

			//return
		}
		logger.Printf("Save file: %s", fmt.Sprintf(destinationFile+qr+"-time-"+strcurrentTime+".jpg"))

		pathtosavefile := destinationFile
		directory := qr
		file := qr + "-time-" + strcurrentTime + ".jpg"
		params := []DataResponse{
			{
				Response:  pathtosavefile,
				Directory: directory,
				File:      file,
			},
		}
		data := map[string]interface{}{
			"GetParams": params,
		}
		//qr = ""
		err = tpl.ExecuteTemplate(w, conf.HTTP.IndexExecute, data)
		if err != nil {
			logger.Errorf("%s", err.Error())
		}

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

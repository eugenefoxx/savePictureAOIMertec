package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/eugenefoxx/savePictureAOIMertec/pkg/logging"
	//"google.golang.org/genproto/googleapis/cloud/metastore/logging/v1"
)

const (
	qr         = "2AQC622457019706"
	checkDIR   = "/home/eugenearch/Code/github.com/eugenefoxx/savePictureAOIMertec/internal/folder/saveFoto/"
	sourceFile = "/home/eugenearch/Code/github.com/eugenefoxx/savePictureAOIMertec/internal/folder/metrec/SmallTop.jpg"
	//destinationFile = checkDIR + qr //"/home/eugenearch/Code/github.com/eugenefoxx/savePictureAOIMertec/internal/folder/saveFoto/"
)

func init() {
	logging.Init()
	//	logger := logging.GetLogger()

	/*	err := godotenv.Load()
		if err != nil {
			logger.Fatal(err.Error)
		}*/
}

func main() {
	logger := logging.GetLogger()

	// copy file
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}

	checkingDIR := checkDIR + qr + "/"
	fmt.Println(checkingDIR)
	if _, err := os.Stat(checkingDIR); os.IsNotExist(err) {
		os.Mkdir(checkingDIR, 0755)
	}

	destinationFile := checkingDIR
	fmt.Println(destinationFile)
	currentDate := time.Now()
	strcurrentTime := currentDate.Format("2006-01-02 15:04:05")
	err = ioutil.WriteFile(destinationFile+qr+"-time-"+strcurrentTime+".jpg", input, 0644)
	if err != nil {
		logger.Errorln("Error creating", destinationFile)
		logger.Errorf(err.Error())
		return
	}

}

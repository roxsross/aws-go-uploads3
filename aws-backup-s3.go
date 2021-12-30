/*
- Creado por: http://github.com/roxsross
- configure sus credenciales de AWS
- después de configurar, puede ejecutar el script
- si es necesario, cambie la constante 'Región'
- para ejecutar: vaya a ejecutar main.go name-of-bucket / absolute / path / to / folder / filename.zip
- puede usar cron para automatizar sus copias de seguridad
  con cron, puedes personalizar la hora exacta
  o intervalo de tiempo que desee entre una copia de seguridad
  y otro
*/

package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// La región debe cambiarse si es necesario
const Region = "us-east-1"

func zipWriter(path, filename string) {
	baseFolder := path

	outFile, err := os.Create(filename)
	checkErr(err)
	defer outFile.Close()

	// Crear un nuevo archivo zip.
	w := zip.NewWriter(outFile)

	// Agregar algunos archivos.
	addFiles(w, baseFolder, "")

	// Make sure to check the error on Close.
	err = w.Close()
	checkErr(err)
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	// abrir directorio
	files, err := ioutil.ReadDir(basePath)
	checkErr(err)

	for _, file := range files {
		log.Println("[INFO] " + basePath + file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			checkErr(err)

			f, err := w.Create(baseInZip + file.Name())
			checkErr(err)

			_, err = f.Write(dat)
			checkErr(err)
		} else if file.IsDir() {

			// Recurse
			newBase := basePath + file.Name() + "/"
			log.Println("[INFO] Recursing and Adding SubDir: " + file.Name())
			log.Println("[INFO] Recursing and Adding SubDir: " + newBase)

			addFiles(w, newBase, file.Name()+"/")
		}
	}
}

/*
  Funcion usada para subir zip en s3
*/
func uploadArchive(filename, bucket string) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(Region),
	})
	checkErr(err)

	newUpload := s3manager.NewUploader(sess)

	file, err := os.Open(filename)
	checkErr(err)

	result, err := newUpload.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filepath.Base(filename)),
		Body:   file,
	})
	checkErr(err)

	log.Println("[INFO] Subida successfully! Path del archivo:", result.Location)
}

func checkErr(err error) {
	if err != nil {
		log.Printf("[ERROR] %s", err)
		os.Exit(1)
	}
}

func main() {
	var (
		bucketNamePtr     = flag.String("bucket", "", "Bucket name")
		folderToUploadPtr = flag.String("path", "", "Folder to upload")
		zipNamePtr        = flag.String("zip", "", "Name of the compressed file")
		bucket            string
		path              string
		filename          string
	)

	flag.StringVar(bucketNamePtr, "b", "", "Bucket name")
	flag.StringVar(folderToUploadPtr, "p", "", "Folder to upload")
	flag.StringVar(zipNamePtr, "z", "", "Name of the compressed file")
	flag.Parse()

	bucket = *bucketNamePtr
	path = *folderToUploadPtr
	filename = *zipNamePtr

	if bucket == "" || path == "" || filename == "" {
		log.Println("[WARNING] Sintaxis Correcta: aws-backup-s3 -b name-of-bucket -p absolute/volumes/path/ -f filename.zip")
		os.Exit(1)
	}

	t := time.Now()

	fmt.Println("=== STARTING NEW BACKUP ====")
	log.Println("[INFO] Time: " + t.String())

	log.Println("[INFO] Comprimir archivos	...")
	zipWriter(path, filename)
	log.Println("[INFO] Archivos comprimidos!")

	log.Printf("[INFO] Uploading %s...\n", filename)
	uploadArchive(filename, bucket)

}

# Script de backup de S3 usando Golang

_Un script de Golang creado para realizar backup en AWS S3_

## Comenzando 🚀

### Como Usar el Script 📋

```
aws-backup --bucket bucket-name --path / folder / to / upload --zip name-of-the-zip.zip
```

_El script comprime la carpeta que se pasa como parámetro, después de esto, carga el archivo comprimido en AWS S3_

## Cómo automatizar el script para realizar backup todos los días ⚙️

_Una solución para automatizar la copia de seguridad de archivos es utilizar un Cron_

### Cree un script de Go en un ejecutable ⌨️
```
go build aws-backup-s3.go
```

### Ahora, instale el cron (si no lo tiene instalado) ⌨️
```
apt-get install cron
```
### Después de eso, acceda al archivo para configurar un horario.⌨️

```
crontab -e
```

¡Programe, guarde y diviértase RoxsRoss !

---
⌨️ con ❤️ por [RoxsRoss](https://github.com/roxsross) 😊
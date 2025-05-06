package ws

import (
	"github.com/arxtect/ArxBackend/golangp/common/utils"
	"github.com/arxtect/ArxBackend/golangp/common/xminio"
	"github.com/arxtect/ArxBackend/golangp/config"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/toheart/functrace"
)

func HandleGetStaticResource(filename string) func(http.ResponseWriter, *http.Request) {
	defer functrace.Trace([]interface {
	}{filename})()
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(415)
			return
		}
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			utils.WriteJSON(w, 400, err.Error())
			return
		}
		w.WriteHeader(200)
		io.WriteString(w, string(b))
	}
}
func HandleCreateRoom(fileId string) (string, string, error) {
	defer functrace.Trace([]interface {
	}{fileId})()

	existence, errByMinio := xminio.NewS3Manager(config.Env.MinioBucket, config.Env.MinioAccessKey, config.Env.MinioSecretKey, config.Env.MinioBucketUrl).CheckFileExistence(fileId)
	if !existence || errByMinio != nil {
		return "", "", errByMinio
	}

	room, err := roomService.NewRoom(fileId)
	if err != nil {
		return "", "", err
	}

	Invitation := utils.RoomIdCreate(6)

	return room.ID, Invitation, nil
}

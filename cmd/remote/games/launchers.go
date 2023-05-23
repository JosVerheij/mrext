package games

import (
	"encoding/json"
	"github.com/wizzomafizzo/mrext/pkg/config"
	"github.com/wizzomafizzo/mrext/pkg/games"
	"github.com/wizzomafizzo/mrext/pkg/mister"
	"github.com/wizzomafizzo/mrext/pkg/service"
	"github.com/wizzomafizzo/mrext/pkg/utils"
	"net/http"
	"path/filepath"
	"strings"
)

func LaunchGame(logger *service.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var args struct {
			Path string `json:"path"`
		}

		err := json.NewDecoder(r.Body).Decode(&args)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Error("launch game: decoding request: %s", err)
			return
		}

		syss := games.FolderToSystems(args.Path)
		if len(syss) == 0 {
			http.Error(w, "no system found for game", http.StatusBadRequest)
			logger.Error("launch game: no system found for game: %s (%s)", args.Path, syss[0].Id)
			return
		}

		err = mister.LaunchGame(syss[0], args.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("launch game: during launch: %s", err)
			return
		}
	}
}

func LaunchFile(logger *service.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var args struct {
			Path string `json:"path"`
		}

		err := json.NewDecoder(r.Body).Decode(&args)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Error("launch file: decoding request: %s", err)
			return
		}

		err = mister.LaunchGenericFile(args.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("launch file: during launch: %s", err)
			return
		}
	}
}

func LaunchMenu(w http.ResponseWriter, _ *http.Request) {
	err := mister.LaunchMenu()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type CreateLauncherRequest struct {
	GamePath string `json:"gamePath"`
	Folder   string `json:"folder"`
	Name     string `json:"name"`
}

type CreateLauncherResponse struct {
	Path string `json:"path"`
}

func CreateLauncher(logger *service.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var args CreateLauncherRequest

		err := json.NewDecoder(r.Body).Decode(&args)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Error("create launcher: decoding request: %s", err)
			return
		}

		//file, err := os.Stat(args.GamePath)
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	logger.Error("create launcher: path is not accessible: %s", err)
		//	return
		//}
		//
		//if file.IsDir() {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	logger.Error("create launcher: path is a directory")
		//	return
		//}

		systems := games.FolderToSystems(args.GamePath)
		if len(systems) == 0 {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("create launcher: unknown file type or folder")
			return
		}

		if !strings.HasPrefix(args.Folder, config.SdFolder) {
			args.Folder = filepath.Join(config.SdFolder, args.Folder)
		}

		args.Name = utils.StripBadFileChars(args.Name)

		system := systems[0]

		mglPath, err := mister.CreateLauncher(
			&system,
			args.GamePath,
			args.Folder,
			args.Name,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("create launcher: creation: %s", err)
			return
		} else {
			err = json.NewEncoder(w).Encode(CreateLauncherResponse{
				Path: mglPath,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				logger.Error("create launcher: encoding response: %s", err)
				return
			}
		}
	}
}
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	// "github.com/spf13/viper"
	"github.com/coyle/bridge/server/routes"
	"github.com/coyle/bridge/server/routes/buckets"
	"github.com/coyle/bridge/server/routes/contacts"
	"github.com/coyle/bridge/server/routes/files"
	"github.com/coyle/bridge/server/routes/frames"
	"github.com/coyle/bridge/server/routes/keys"
	"github.com/coyle/bridge/server/routes/reports"
	"github.com/coyle/bridge/server/routes/users"
	"github.com/coyle/bridge/storage/mongodb"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = level.NewFilter(logger, level.AllowAll())

	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	storageClient, err := mongodb.NewClient(os.Getenv("MONGO"))
	if err != nil {
		level.Error(logger).Log("Error connecting", err)
		// TODO(coyle): handle appropriately
		// return
	}

	handler := routes.Handler{
		Logger: logger,
		User:   users.NewServer(storageClient, logger),
	}

	http.ListenAndServe(":8080", start(&handler))

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signalChan

	level.Info(logger).Log("Server Stopping", "bridger-server", "sig", sig)
}

func start(handler *routes.Handler) *httprouter.Router {
	router := httprouter.New()
	// Bucket specific routes
	router.GET("/buckets", buckets.Get)
	router.GET("/buckets/:id", buckets.GetByID)
	router.GET("/bucket-ids/:name", buckets.GetByID)
	router.POST("/buckets", buckets.Create)
	router.DELETE("/buckets/:id", buckets.DestroyByID)
	router.PATCH("/buckets/:id", buckets.UpdateByID)
	router.POST("/buckets/:id/tokens", buckets.CreateToken)
	// File specific routes
	router.GET("/buckets/:id/files", files.List)
	router.GET("/buckets/:id/file-ids/:name", files.GetID)
	router.GET("/buckets/:id/files/:file", files.Get)
	router.DELETE("/buckets/:id/files/:file", files.Delete)
	router.GET("/buckets/:id/files/:file/info", files.GetInfo)
	router.POST("/buckets/:id/files", files.CreateEntryFromFrame)
	router.GET("/buckets/:id/files/:file/mirrors", files.ListMirrorsForFile)
	// Contact specific routes
	router.GET("/contacts", contacts.GetList)
	router.GET("/contacts/:nodeID", contacts.GetByNodeID)
	router.PATCH("/contacts/:nodeID", contacts.PatchByNodeID)
	router.POST("/contacts", contacts.Create)
	router.POST("/contacts/challenges", contacts.CreateChallenge)
	// Frames specific routes
	router.POST("/frames", frames.Create)
	router.PUT("/frames:frame", frames.AddShard)
	router.DELETE("/frames/:frame", frames.RemoveByID)
	router.GET("/frames", frames.Get)
	router.GET("/frames/:frame", frames.GetByID)
	// Public Key specific routes
	router.GET("/keys", keys.Get)
	router.POST("/keys", keys.Add)
	router.DELETE("/keys/:pubkey", keys.Remove)
	// Report specific routes
	router.POST("/reports/exchanges", reports.Create)
	// User specific routes
	router.POST("/users", handler.User.Create)
	router.POST("/activations", handler.User.Reactivate)
	router.GET("/activations/:token", handler.User.ConfirmActivation)
	router.DELETE("/users/:id", handler.User.Remove)
	router.GET("/deactivations/:token", handler.User.ConfirmDeactivation)
	router.PATCH("/users/:id", handler.User.CreatePasswordResetToken)
	router.POST("/resets/:token", handler.User.ConfirmPasswordReset)
	// DEBUG specific endpoints
	router.GET("/health", health)

	return router
}

func health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO(coyle): actually check something
	fmt.Fprintln(w, "OK")
}

package server

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"usermanagement/app/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

var (
	cfgFile        string
	configurations config.Config
)

var rootCmd = &cobra.Command{
	Short: "User Management Service",
	Long:  `REST APIs for User Management microservice`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := viper.Unmarshal(&configurations); err != nil {
			log.WithField("err", err).Error("unmarshal config")
		}
	},
}

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.usermanagement.yml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName(".usermanagement")
		viper.SetConfigType("yml")
	}
	viper.SetEnvPrefix("ms")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		log.WithField("err", err).Error("reading config")
	}
}

type Server struct {
	context            context.Context
	shutdownFn         context.CancelFunc
	childRoutines      *errgroup.Group
	shutdownReason     string
	shutdownInProgress bool
	app                *config.AppConfiguration
}

func NewServer() *Server {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCtx, shutdownFn := context.WithCancel(context.Background())
	childRoutines, childCtx := errgroup.WithContext(rootCtx)
	app := config.NewAppService(configurations)
	RegisterService(app)
	return &Server{
		context:       childCtx,
		shutdownFn:    shutdownFn,
		childRoutines: childRoutines,
		app:           app,
	}
}

func (server *Server) Run() (err error) {
	services := GetServices()
	for _, service := range services {
		if err := service.Instance.Init(); err != nil {
			log.WithField("err", err).Fatal("starting service")
			return err
		}
	}

	for _, svc := range services {
		service, ok := svc.Instance.(BackgroundService)
		if !ok {
			continue
		}

		descriptor := svc
		server.childRoutines.Go(func() error {
			if server.shutdownInProgress {
				return nil
			}

			err := service.Run(server.context)
			server.shutdownInProgress = true
			if err != nil {
				log.WithField("reason", err.Error()).Errorf("stopped  %s", descriptor.Name)
				return err
			}
			return nil
		})
	}

	defer func() {
		if waitErr := server.childRoutines.Wait(); waitErr != nil && reflect.TypeOf(waitErr) != reflect.TypeOf(context.Canceled) {
			log.WithField("err", waitErr).Error("service failed")
			if err == nil {
				err = waitErr
			}
		}
	}()
	return
}

func (server *Server) Shutdown(reason string) {

	server.shutdownReason = reason
	server.shutdownInProgress = true
	server.shutdownFn()

	if err := server.childRoutines.Wait(); err != nil && reflect.TypeOf(err) != reflect.TypeOf(context.Canceled) {
		log.WithField("err", err).Error("shutdown failed")
	}
}

func (server *Server) ExitCode(reason error) int {
	code := 1
	if reason == context.Canceled || server.shutdownReason != "" {
		log.WithField("reason", server.shutdownReason).Info("shutting down")
		return 0
	}
	return code
}

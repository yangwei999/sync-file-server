package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/ghodss/yaml"
	"github.com/opensourceways/community-robot-lib/kafka"
	"github.com/opensourceways/community-robot-lib/logrusutil"
	"github.com/opensourceways/community-robot-lib/mq"
	"github.com/opensourceways/community-robot-lib/secret"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/sync-file-server/grpc/server"
	"github.com/opensourceways/sync-file-server/models"
)

type options struct {
	port              string
	endpoint          string
	platform          string
	platformTokenPath string
	topic             string
	configFile        string
	concurrentSize    int
}

func (o *options) addFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.port, "port", "8888", "Port to listen on.")
	fs.StringVar(&o.endpoint, "endpoint", "", "The endpoint of repo file cache")
	fs.StringVar(&o.platform, "platform", "gitee", "The code platform which implements rpc service. Currently only gitee is supported")
	fs.StringVar(&o.platformTokenPath, "platform-token-path", "/etc/platform/oauth", "The path to the token file which is used to access code platform.")
	fs.StringVar(&o.topic, "topic", "", "The topic to which jobs need to be published ")
	fs.StringVar(&o.configFile, "config-file", "", "Path to the config file.")

	fs.IntVar(&o.concurrentSize, "concurrent-size", 2000, "The grpc server goroutine pool size.")
}

func (o *options) validate() error {
	v, err := url.Parse(o.endpoint)
	if err != nil {
		return err
	}
	o.endpoint = v.String()

	if o.concurrentSize <= 0 {
		return fmt.Errorf("concurrent size must be bigger than 0")
	}

	if o.configFile == "" {
		return fmt.Errorf("config-file must be set")
	}

	if o.topic == "" {
		return fmt.Errorf("topic must be set")
	}

	return nil
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options
	o.addFlags(fs)
	_ = fs.Parse(args)
	return o
}

const component = "sync-file-server"

func main() {
	logrusutil.ComponentInit("sync-file-server")

	o := gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if err := o.validate(); err != nil {
		logrus.WithError(err).Fatal("Invalid options")
	}

	logrus.WithField("platform", o.platform).Info("Starts sync file server.")

	secretAgent := new(secret.Agent)
	if err := secretAgent.Start([]string{o.platformTokenPath}); err != nil {
		logrus.WithError(err).Fatal("Error to start secret agent.")
	}
	defer secretAgent.Stop()

	getToken := secretAgent.GetTokenGenerator(o.platformTokenPath)

	backend, err := newBackend(o.endpoint, o.platform, getToken)

	if err != nil {
		logrus.WithError(err).Fatal("Error to generate backend")
	}

	if err = initBroker(o.configFile); err != nil {
		logrus.WithError(err).Fatal("error to init broker")
	}

	s, err := models.SyncFromMQ(o.topic, component)
	if err != nil {
		logrus.WithError(err).Fatal("error to subscribe")
	}
	defer s.Unsubscribe()

	if err := server.Start(":"+o.port, o.concurrentSize, backend); err != nil {
		logrus.WithError(err).Fatal("Error to start grpc server.")
	}
}

func initBroker(configFile string) error {
	cfg, err := loadConfig(configFile)
	if err != nil {
		return fmt.Errorf("load config:%s", configFile)
	}

	err = kafka.Init(
		mq.Addresses(cfg.MQConfig.Addresses...),
		mq.Log(logrus.WithField("module", "broker")),
	)

	if err != nil {
		return err
	}

	return kafka.Connect()
}

func loadConfig(path string) (*configuration, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	v := new(configuration)
	if err := yaml.Unmarshal(b, v); err != nil {
		return nil, err
	}

	return v, nil
}

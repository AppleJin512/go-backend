package config

import (
	"fmt"
	"moonbite/trending/internal/blockchain"
	"strings"

	"github.com/centrifugal/gocent/v3"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/sirupsen/logrus"
)

type Configuration struct {
	Debug     bool
	LogLevel  logrus.Level
	Secret    string
	Instances string
}

func (c *Configuration) SetLogLevel(l string) error {
	lvl, err := logrus.ParseLevel(l)
	if err != nil {
		return fmt.Errorf("set log level error: %v", err)
	}
	c.LogLevel = lvl
	return nil
}

var Config Configuration

var RpcClient *rpc.Client

func InitRpcClient(instances string) (*rpc.Client, error) {
	instancesSplit := strings.Split(instances, ":")
	if len(instancesSplit) == 0 {
		return nil, fmt.Errorf("invalid instances, must be set one")
	}
	nodeSetting, ok := blockchain.NodeSettings[instancesSplit[0]]
	if !ok {
		return nil, fmt.Errorf("unknown node setting %s", instancesSplit[0])
	}
	rpcClient := rpc.New(nodeSetting.NodeUrlRpc)
	return rpcClient, nil
}

var GoCent *gocent.Client

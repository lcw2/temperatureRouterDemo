package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/d2r2/go-dht"
	logger "github.com/d2r2/go-logger"
	"github.com/d2r2/go-shell"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
	"os"
	"strconv"
	"syscall"
	"time"
)

var lg = logger.NewPackageLogger("main", logger.DebugLevel)

//BaseMessage the base struct of event message
type BaseMessage struct {
	EventID   string `json:"event_id"`
	Timestamp int64  `json:"timestamp"`
}

type ValuesMessage struct {
	Temperature string `json:"temperature"`
	Humidity    string `json:"humidity"`
}

func main() {
	defer logger.FinalizeLogger()

	logger.ChangePackageLogLevel("dht", logger.InfoLevel)
	//create context with cancellation possibility
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	defer close(done)
	//build actual signal list to control
	signals := []os.Signal{os.Kill, os.Interrupt}
	if shell.IsLinuxMacOSFreeBSD() {
		signals = append(signals, syscall.SIGTERM)
	}

	// run goroutine waiting for OS terminate events, including keyboard Ctrl+C
	shell.CloseContextOnSignals(cancel, done, signals...)

	sensorType := dht.DHT11
	pin := 11
	totalRetried := 0
	totalMeasured := 0
	totalFailed := 0
	term := false

	cli := connectToMqtt()

	for {
		temperature, humidity, retried, err :=
			dht.ReadDHTxxWithContextAndRetry(ctx, sensorType, pin, false, 10)
		totalMeasured++
		totalRetried += retried
		if err != nil && ctx.Err() == nil {
			totalFailed++
			lg.Error(err)
			continue
		}

		//print temperature and humidity
		if ctx.Err() == nil {
			lg.Infof("Sensor = %v: Temperature = %v*C, Humidity= %v%% (retried %d times)", sensorType, temperature, humidity, retried)
		}

		//publish temperature
		publishToMqtt(cli, temperature, humidity)

		select {
		case <-ctx.Done():
			lg.Errorf("Termination pending: %s", ctx.Err())
			term = true
		case <-time.After(2000 * time.Millisecond):
		}

		if term {
			break
		}

	}
	lg.Info("exited")
}

func publishToMqtt(cli *client.Client, temperature float32, humidity float32) {
	deviceValueUpdate := "default/device/temperature/and/humidity"

	updateMessage := createActualUpdateMessage(strconv.Itoa(int(temperature))+"C", strconv.Itoa(int(humidity))+"%")
	UpdateBody, _ := json.Marshal(updateMessage)

	cli.Publish(&client.PublishOptions{
		QoS:       mqtt.QoS0,
		TopicName: []byte(deviceValueUpdate),
		Message:   UpdateBody,
	})

}

func createActualUpdateMessage(actualTeperature string, actualHumidity string) ValuesMessage {
	var deviceUpdateMessage ValuesMessage
	deviceUpdateMessage = ValuesMessage{
		Temperature: actualTeperature,
		Humidity:    actualHumidity,
	}
	return deviceUpdateMessage
}

func connectToMqtt() *client.Client {
	cli := client.New(&client.Options{
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})
	defer cli.Terminate()

	err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  "localhost:1883",
		ClientID: []byte("receive-client"),
	})

	if err != nil {
		panic(err)
	}
	return cli
}

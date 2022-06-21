package nats

import (
	"encoding/json"

	"io/ioutil"
	"log"
	"nats/internal/model"
	"os"

	stan "github.com/nats-io/stan.go"
)

const urlConent string = "http://localhost:4223"

func natsCon() stan.Conn {
	sc, err := stan.Connect("test-cluster", "client1", stan.NatsURL(urlConent))
	if err != nil {
		log.Fatalf(err.Error())
	}
	return sc
}

func NatsPublisher() error {
	s, err := readMsgNats()
	if err != nil {
		return err
	}
	sc := natsCon()

	err = sc.Publish("foo10", s)
	if err != nil {
		sc.Close()
		return err
	}
	sc.Close()
	return nil
}

var b []byte

func NatsSubscriber() (model.ModelNats, error) {
	var JsonModel model.ModelNats

	sc := natsCon()

	sub, err := sc.Subscribe("foo10", func(m *stan.Msg) {
		b = m.Data

	}, stan.DeliverAllAvailable())
	if err != nil {
		return model.ModelNats{}, err
	}

	sub.Unsubscribe()

	sc.Close()

	JsonModel, err = modelUnmarshal(b)
	if err != nil {
		return model.ModelNats{}, err
	}

	return JsonModel, nil
}

func readMsgNats() ([]byte, error) {
	fileName := "model.json"
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	readFile, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return readFile, nil
}

func modelUnmarshal(b []byte) (model.ModelNats, error) {
	var model model.ModelNats

	if err := json.Unmarshal(b, &model); err != nil {
		return model, err
	}

	return model, nil
}

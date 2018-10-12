package channel

import (
	"context"
	"fmt"

	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/engine/channels"
	"github.com/project-flogo/core/support/logger"
	"github.com/project-flogo/core/trigger"

)

var log = logger.GetLogger("trigger-channel")

var triggerMd = trigger.NewMetadata(&HandlerSettings{}, &Output{})

func init() {
	trigger.Register(&Trigger{}, &Factory{})
}

type Factory struct {
}

// Metadata implements trigger.Factory.Metadata
func (*Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// New implements trigger.Factory.New
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	return &Trigger{}, nil
}


// Trigger CHANNEL trigger struct
type Trigger struct {
	metadata *trigger.Metadata
	config   *trigger.Config
	handlers []*trigger.Handler
}


func (t *Trigger) Initialize(ctx trigger.InitContext) error {

	// validate handlers
	for _, handler := range ctx.GetHandlers() {

		s := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), s, true)
		if err != nil {
			return err
		}

		ch := channels.Get(s.Channel)
		if ch == nil {
			return fmt.Errorf("unknown engine channel '%s'", s.Channel)
		}

		l := &Listener{handler: handler}
		ch.RegisterCallback(l.OnMessage)
	}

	return nil
}

// Stop implements util.Managed.Start
func (t *Trigger) Start() error {
	return nil
}

// Stop implements util.Managed.Stop
func (t *Trigger) Stop() error {
	return nil
}

type Listener struct {
	handler trigger.Handler
}

func (l *Listener) OnMessage(msg interface{}) {
	triggerData := make(map[string]interface{})

	if vals, ok := msg.(map[string]interface{}); ok {
		triggerData[ovData] = vals
	} else {
		triggerData[ovData] = msg
	}

	//todo what should we do with the results?
	_, err := l.handler.Handle(context.TODO(), triggerData)

	if err != nil {
		log.Error(err)
	}
}
package db

import (
	"github.com/golang/glog"
)

// EventReceiver is a sentinel EventReceiver; use it if the caller doesn't supply one
type EventReceiver struct{}

// Event receives a simple notification when various events occur
func (n *EventReceiver) Event(eventName string) {

}

// EventKv receives a notification when various events occur along with
// optional key/value data
func (n *EventReceiver) EventKv(eventName string, kvs map[string]string) {
}

// EventErr receives a notification of an error if one occurs
func (n *EventReceiver) EventErr(eventName string, err error) error {
	return err
}

// EventErrKv receives a notification of an error if one occurs along with
// optional key/value data
func (n *EventReceiver) EventErrKv(eventName string, err error, kvs map[string]string) error {
	glog.Errorf("%+v", err)
	glog.Errorf("%s: %+v", eventName, kvs)
	return err
}

// Timing receives the time an event took to happen
func (n *EventReceiver) Timing(eventName string, nanoseconds int64) {

}

// TimingKv receives the time an event took to happen along with optional key/value data
func (n *EventReceiver) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {
	// TODO: Change logger level to debug
	glog.Infof("%s spend %.2fms: %+v", eventName, float32(nanoseconds)/1000000, kvs)
}
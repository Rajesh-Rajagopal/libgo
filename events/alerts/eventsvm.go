package alerts

import (
	log "github.com/Sirupsen/logrus"
	ldb "github.com/megamsys/libgo/db"
	constants "github.com/megamsys/libgo/utils"
	"time"
)

const (
	EVENTSBUCKET = "events_for_vms"
	EVENTVM_JSONCLAZ = "Megam::EventsVm"

)
type EventsVm struct {
	EventType  string   `json:"event_type" cql:"event_type"`
	AccountId  string   `json:"account_id" cql:"account_id"`
	AssemblyId string   `json:"assembly_id" cql:"assembly_id"`
	Data       []string `json:"data" cql:"data"`
	CreatedAt  string   `json:"created_at" cql:"created_at"`
	JsonClaz   string   `json:"json_claz" cql:"json_claz"`
}

func (s *Scylla) NotifyVm(eva EventAction, edata EventData) error {
	if !s.satisfied(eva) {
		return nil
	}
	s_data := parseMapToOutputFormat(edata)
	ops := ldb.Options{
		TableName:   EVENTSBUCKET,
		Pks:         []string{constants.EVENT_TYPE, constants.CREATED_AT},
		Ccms:        []string{constants.ASSEMBLY_ID, constants.ACCOUNT_ID},
		Hosts:       s.scylla_host,
		Keyspace:    s.scylla_keyspace,
		PksClauses:  map[string]interface{}{constants.EVENT_TYPE: edata.M[constants.EVENT_TYPE], constants.CREATED_AT: s_data.CreatedAt},
		CcmsClauses: map[string]interface{}{constants.ASSEMBLY_ID: edata.M[constants.ASSEMBLY_ID], constants.ACCOUNT_ID: edata.M[constants.ACCOUNT_ID]},
	}
	if err := ldb.Storedb(ops, s_data); err != nil {
		log.Debugf(err.Error())
		return err
	}

	return nil
}

func parseMapToOutputFormat(edata EventData) EventsVm {
	return EventsVm{
		EventType:  edata.M[constants.EVENT_TYPE],
		AccountId:  edata.M[constants.ACCOUNT_ID],
		AssemblyId: edata.M[constants.ASSEMBLY_ID],
		Data:       edata.D,
		CreatedAt:  time.Now().String(),
    JsonClaz: EVENTVM_JSONCLAZ,
	}
}

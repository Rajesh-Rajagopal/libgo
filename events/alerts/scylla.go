package alerts

import (
	constants "github.com/megamsys/libgo/utils"
	"strings"
)

type Scylla struct {
	scylla_host     []string
	scylla_keyspace string
	scylla_username string
	scylla_password string
}

func NewScylla(m map[string]string) Notifier {
	return &Scylla{
		scylla_host:  strings.Split(m[constants.SCYLLAHOST], ","),
		scylla_keyspace: m[constants.SCYLLAKEYSPACE],
		scylla_username: m[constants.SCYLLAUSERNAME],
		scylla_password: m[constants.SCYLLAPASSWORD],
	}
}

func (s *Scylla) satisfied(eva EventAction) bool {
	if eva == STATUS {
		return true
	}
	return false
}

func (s *Scylla) Notify(eva EventAction, edata EventData) error {
	value := edata.M[constants.EVENT_TYPE]
	et := strings.Split(value, ".")
	if et[0] == "compute" && et[1] == "instance"{
		return s.NotifyVm(eva, edata)
	} else if et[0] == "bill" {
		return s.NotifyBill(eva, edata)
	} else if et[0] == "storage" {
		return s.NotifyStorage(eva, edata)
	} else if et[0] == "obc" {
		return s.NotifyOBC(eva, edata)
	} else if et[0] == "compute" && et[1] == "container" {
		return s.NotifyContainer(eva, edata)
	} else {
		return nil
	}
	return nil
}

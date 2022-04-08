package monitor

import (
	"fmt"
	"time"

	g "github.com/gosnmp/gosnmp"
)

type C_snmp struct {
	snmp_conn *g.GoSNMP

	s_snmp__ip   string
	s_snmp__port string
	s_snmp__oid  string

	arrs_snmp__oid []string
}

func (t *C_snmp) Run__snmp() error {

	// SNMP 필요 OID 값 정의
	oids := []string{"1.3.6.1.2.1.1.4.0"}
	result, err := t.snmp_conn.Get(oids)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for i, variable := range result.Variables {
		fmt.Printf("%d: oid: %s ", i, variable.Name)

		switch variable.Type {
		case g.OctetString:
			bytes := variable.Value.([]byte)
			fmt.Printf("string: %s\n", string(bytes))
		default:
			fmt.Printf("number: %d\n", g.ToBigInt(variable.Value))
		}
	}
	return nil
}

// SNMP Connection 시작
func (t *C_snmp) Snmp__conn(_s_snmp__ip string, _s_snmp__port uint16) error {

	t.snmp_conn = &g.GoSNMP{
		Port:      _s_snmp__port,
		Target:    _s_snmp__ip,
		Community: "public",
		Version:   g.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		Retries:   3,
	}

	err := t.snmp_conn.Connect()
	if err != nil {
		return err
	}

	return nil
}

// SNMP Session 종료
func (t *C_snmp) Snmp__close() error {
	return t.snmp_conn.Conn.Close()

}

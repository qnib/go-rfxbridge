package rfx // import "github.com/qnib/go-rfxbridge/rfx"


import (
	"fmt"
	"github.com/barnybug/gorfxtrx"
	"reflect"
	"strings"
)

type Devices struct {
	usbDev 	string
	debug 	bool
	mapping map[string]string
	trans 	chan gorfxtrx.Packet
	data 	map[string]string
}

func NewDevices(usbDev string, debug bool, m map[string]string) Devices {
	d := Devices{
		usbDev: usbDev,
		debug: debug,
		mapping: m,
		data: make(map[string]string),
		trans: make(chan gorfxtrx.Packet),
	}
	return d
}

func (d *Devices) UpdateData() {
	for {
		select {
		case p := <- d.trans:
			t := reflect.TypeOf(p)
			switch t.String() {
			case "*gorfxtrx.LightingX10":
				continue
			case "*gorfxtrx.LightingHE":
				dev := p.(*gorfxtrx.LightingHE)
				key := fmt.Sprintf("%v.%v", dev.HouseCode, dev.UnitCode)
				key = d.EvalKey(key)
				d.data[key] = dev.Command()
				fmt.Printf("LightingHE: %v.%v %s\n", dev.HouseCode, dev.UnitCode, dev.Command())
			default:
				fmt.Printf("Received (%s): %v", t, p)
			}
		}
	}
}

func (d *Devices) String() string {
	res := []string{}
	for k,v := range d.data {
		res = append(res, fmt.Sprintf("%s:%s", k, v))
	}
	res = append(res, "")
	return strings.Join(res, "\n")
}

func (d *Devices) EvalKey(key string) string {
	if k, ok := d.mapping[key];ok {
		return k
	}
	return key
}

func (d *Devices) GetKey(key string) (res string,err error) {
	key = d.EvalKey(key)
	if res, ok := d.data[key];ok {
		return res, nil
	}
	keys := []string{}
	for k, _ := range d.data {
		keys = append(keys, k)
		}
	return res, fmt.Errorf("No key '%s' (%s)", key, strings.Join(keys,","))
}

func (d *Devices) WatchRFX() {
	dev, err := gorfxtrx.Open(d.usbDev, d.debug)
	if err != nil {
		panic("Error opening device")
	}

	for {
		packet, err := dev.Read()
		if err != nil {
			continue
		}
		d.trans <- packet
	}
	dev.Close()
}
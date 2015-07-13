package engine

import (
	"github.com/astaxie/beego"

	"log"
	"google.golang.org/grpc"
	bd "adexchange/engine/baidu/mobads_api"
//	"github.com/golang/protobuf/proto"

//	"github.com/franela/goreq"
//	"time"
//	"bytes"
//	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"time"
	"github.com/franela/goreq"
)

func invokeBD2(demand *Demand) {

	//	address := "http://220.181.163.105/api"
	address := "http://mobads.baidu.com:80/api"

	conn, err := grpc.Dial(address)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Set up a connection to the server.
//	c := bd.NewBDServiceClient(conn)
//
//	r, err := c.RequestAd(context.Background(), &bd.BidRequest{})
//
//	if err != nil {
//		log.Fatalf("could not get add from server: %v", err)
//	}
//	log.Printf("Greeting: %s", r.ErrorCode)

	beego.Debug("invoke BD..XXXXXXXXXXXXXXXXXXXXXXXXXXX")

//	beego.Debug(r)

}



func invokeBD(demand *Demand) {

	beego.Debug("target url: ", demand.URL)


	apiVersion := &bd.Version{
		Major: pUint32(4),
		Minor: pUint32(0),
	}
	appVersion := &bd.Version{
		Major: pUint32(4),
		Minor: pUint32(0),
	}


	var appId string = "c5702dff";			// required
	var appBundleId string = "com.xxxxx"	// required
	var appName string = "";
	userPermissionType := bd.App_UserPermission_ACCESS_FINE_LOCATION
	userPermissionStatus := bd.App_UserPermission_UNKNOWN
	appCategories := []uint32{uint32(1)}			// required
	app := &bd.App{
		Id: &appId,
		StaticInfo: &bd.App_StaticInfo{
			BundleId: &appBundleId,
			Name: &appName,
			Categories: appCategories,
		},
		Version: appVersion,
		UserPermission: []*bd.App_UserPermission{
			&bd.App_UserPermission{
				Type: &userPermissionType,		// required
				Status: &userPermissionStatus,					// required
			},
		},
	}

	devOsVersion := &bd.Version{
		Major: pUint32(8),
		Minor: pUint32(4),
	}
	devModel := "Iphone6"
	devVendor := "China MObile"
	devId := "123abc"
	devUdid := bd.Device_UdId{
		Idfa: &devId,
	}
	devType := bd.Device_PHONE
	devOs := bd.Device_IOS
	dev := &bd.Device{
		Type: &devType,					// required. Mobile, Tablet, TV
		Os: &devOs,						// required. android or IOS
		OsVersion: devOsVersion,		// required. OS version
		Vendor: &devVendor,				// required.
		Model: &devModel,				// required.
		Udid: &devUdid,					// required. ios: idfa, mac, android: imei, mac, tv: imei, mac, idfv
	}

	nt := &bd.Network{

	}

	adSpaceId := "L000015a"
	adWidth := 600
	adHeight := 400
	adSize := bd.Size{
		Width: pUint32(adWidth),		// required
		Height: pUint32(adHeight),		// required
	}
	adType := bd.AdSlot_StaticInfo_BANNER
	adStaticInfo := bd.AdSlot_StaticInfo{
		Type: &adType,
	}

	var requestId string = "12341234123413241234";
	req := bd.BidRequest{
		RequestId: &requestId,
//		IsDebug: &true,
		ApiVersion: apiVersion,
		App: app,
		Device: dev,
		Network: nt,
		Adslots: []*bd.AdSlot{
			&bd.AdSlot{
				Id: &adSpaceId,				// required.
				Size: &adSize,				// required
				StaticInfo: &adStaticInfo,	// required.
			},
		},
	}

	data, err := proto.Marshal(&req)

	resp, err := goreq.Request{
		Method: "POST",
		Uri:         demand.URL,
		Timeout: time.Duration(demand.Timeout) * time.Millisecond,
		Proxy: "http://localhost:8888",
		Body: data,
	}.Do()

	bidResp := &bd.BidResponse{}
	respStr, err := resp.Body.ToString()
	err = proto.Unmarshal([]byte(respStr), bidResp)

	if err != nil {
		beego.Debug("error: ", err)
	} else {
		beego.Debug("-------- get response from baidu: ", bidResp)
	}

}

func pUint32(v int) *uint32 {
	p := new(uint32)
	*p = uint32(v)
	return p
}

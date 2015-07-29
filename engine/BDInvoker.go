package engine

import (
	bd "adexchange/engine/baidu/mobads_api"
	"fmt"
	"github.com/astaxie/beego"
	//"google.golang.org/grpc"
	//"log"
	//	"github.com/golang/protobuf/proto"

	//	"github.com/franela/goreq"
	//	"time"
	//	"bytes"
	//	"github.com/golang/protobuf/proto"
	"adexchange/lib"
	m "adexchange/models"
	"github.com/franela/goreq"
	"github.com/golang/protobuf/proto"
	"strings"
	"time"
)

var (
	OS_MAP = map[int]bd.Device_Os{
		0: bd.Device_ANDROID,
		1: bd.Device_IOS,
		2: bd.Device_IOS,
		3: bd.Device_IOS,
	}
)

//func invokeBD2(demand *Demand) {

//	address := "http://220.181.163.105/api"

//	conn, err := grpc.Dial(address)
//	if err != nil {
//		log.Fatalf("failed to connect: %v", err)
//	}
//	defer conn.Close()

//	// Set up a connection to the server.
////	c := bd.NewBDServiceClient(conn)
////
////	r, err := c.RequestAd(context.Background(), &bd.BidRequest{})
////
////	if err != nil {
////		log.Fatalf("could not get add from server: %v", err)
////	}
////	log.Printf("Greeting: %s", r.ErrorCode)

//	beego.Debug("invoke BD..XXXXXXXXXXXXXXXXXXXXXXXXXXX")

////	beego.Debug(r)

//}

func invokeBD(demand *Demand) {

	// current baidu api version is 4.0
	// TODO move this to conf file
	apiVersion := &bd.Version{
		Major: pUint32(4),
		Minor: pUint32(0),
	}
	//	appVersion := &bd.Version{
	//		Major: pUint32(4),
	//		Minor: pUint32(0),
	//	}

	/*  App   (required)*/

	// TODO put appid in secret_key 字段???
	var appId string = demand.AdSecretKey // required

	// below are optional
	//	var appBundleId string = "com.xxxxx"	// required
	//	var appName string = "";
	//	appCategories := []uint32{uint32(1)}			// required
	//	userPermissionType := bd.App_UserPermission_ACCESS_FINE_LOCATION
	//	userPermissionStatus := bd.App_UserPermission_UNKNOWN
	app := &bd.App{
		Id: &appId,

		// optional.  because there's no categories, so don't provide static info
		//		StaticInfo: &bd.App_StaticInfo{
		//			BundleId: &appBundleId,		// required
		//			Name: &appName,
		//			Categories: appCategories,	// required
		//		},
		//		Version: appVersion,
		//		UserPermission: []*bd.App_UserPermission{
		//			&bd.App_UserPermission{
		//				Type: &userPermissionType,		// required
		//				Status: &userPermissionStatus,					// required
		//			},
		//		},
	}

	/*  Device  (required)*/
	stringArr := strings.Split(demand.AdRequest.Osv, ".")

	var strOsvMinor string
	var strOsvMajor string
	if len(stringArr) > 0 {
		strOsvMajor = stringArr[0]
	}
	if len(stringArr) > 1 {
		strOsvMinor = stringArr[1]
	}
	devOsVersion := &bd.Version{
		Major: pUint32(lib.ConvertStrToInt(strOsvMajor)),
		Minor: pUint32(lib.ConvertStrToInt(strOsvMinor)),
	}
	devModel := demand.AdRequest.Device // IPhone5s
	var devVendor string                // Apple
	if demand.AdRequest.Os == 1 {
		devVendor = "Apple"
	} else {
		devVendor = "Google"
	}

	devUdid := bd.Device_UdId{}

	if demand.AdRequest.Idfa != "" {
		devUdid.Idfa = &demand.AdRequest.Idfa
	}
	if demand.AdRequest.Imei != "" {
		devUdid.Imei = &demand.AdRequest.Imei
	}
	if demand.AdRequest.Wma != "" {
		devUdid.Mac = &demand.AdRequest.Wma
	}

	devType := bd.Device_PHONE
	devOs := OS_MAP[demand.AdRequest.Os]
	dev := &bd.Device{
		Type:      &devType,     // required. Mobile, Tablet, TV
		Os:        &devOs,       // required. android or IOS
		OsVersion: devOsVersion, // required. OS version
		Vendor:    &devVendor,   // required.
		Model:     &devModel,    // required.
		Udid:      &devUdid,     // required. ios: idfa, mac, android: imei, mac, tv: imei, mac, idfv
	}

	/*  Network  (required) */
	nt := &bd.Network{
		Ipv4: &demand.AdRequest.Ip,
	}

	/* Adslot  (required) */
	adSpaceId := demand.AdspaceKey
	adWidth := demand.AdRequest.Width
	adHeight := demand.AdRequest.Height
	adSize := bd.Size{
		Width:  pUint32(lib.ConvertStrToInt(adWidth)),  // required
		Height: pUint32(lib.ConvertStrToInt(adHeight)), // required
	}
	//	adType := bd.AdSlot_StaticInfo_BANNER
	//	adStaticInfo := bd.AdSlot_StaticInfo{
	//		Type: &adType,
	//	}

	var requestId string = demand.AdRequest.Bid
	req := bd.BidRequest{
		RequestId:  &requestId,
		ApiVersion: apiVersion,
		App:        app,
		Device:     dev,
		Network:    nt,
		Adslots: []*bd.AdSlot{
			&bd.AdSlot{
				Id:   &adSpaceId, // required.
				Size: &adSize,    // required
				//				StaticInfo: &adStaticInfo,
			},
		},
	}

	beego.Debug("baidu request: ", req.String())

	data, err := proto.Marshal(&req)

	if err != nil {
		generateErrorResp(lib.ERROR_BD_MARSHAL_REQ, "failed to marshal bd request", err, demand)
	} else {
		adResponse := new(m.AdResponse)
		adResponse.Bid = demand.AdRequest.Bid
		adResponse.SetDemandAdspaceKey(demand.AdspaceKey)
		adResponse.SetResponseTime(time.Now().Unix())

		resp, err := goreq.Request{
			Method:  "POST",
			Uri:     demand.URL,
			Timeout: time.Duration(demand.Timeout) * time.Millisecond,
			Body:    data,
		}.Do()

		if serr, ok := err.(*goreq.Error); ok {
			beego.Critical(err.Error())
			if serr.Timeout() {
				generateErrorResp(lib.ERROR_TIMEOUT_ERROR, "failed to send request to baidu", err, demand)

			} else {
				generateErrorResp(lib.ERROR_BD_SERVER, "failed to send request to baidu", err, demand)

			}

		} else {
			bidResp := &bd.BidResponse{}
			respStr, err := resp.Body.ToString()

			defer resp.Body.Close()

			if err != nil {
				generateErrorResp(lib.ERROR_BD_FAILED_RES, "failed to get bd response body", err, demand)
			} else {

				err = proto.Unmarshal([]byte(respStr), bidResp)

				if err != nil {
					generateErrorResp(lib.ERROR_BD_MARSHAL_RES, "failed to unmarshal response body", err, demand)
				} else {
					beego.Debug("baidu response: ", bidResp.String())
					mapBDResponse(bidResp, adResponse)
					beego.Debug("map to pmp response successfully.")
					demand.Result <- adResponse
				}
			}
		}
	}
}

func pUint32(v int) *uint32 {
	p := new(uint32)
	*p = uint32(v)
	return p
}

func mapBDResponse(bdResp *bd.BidResponse, adResponse *m.AdResponse) {
	beego.Debug("start mapping.....")
	//	adResponse.StatusCode, _ = strconv.Atoi(strconv.FormatUint(*bdResp.ErrorCode,10))
	//	beego.Debug("error code: ", *bdResp.ErrorCode)
	adResponse.SetResponseTime(time.Now().Unix())

	//	if adResponse.StatusCode == 0 {
	adUnit := new(m.AdUnit)
	adResponse.Adunit = adUnit
	if len(bdResp.GetAds()) > 0 {

		ad := bdResp.GetAds()[0]
		adMeta := ad.MaterialMeta
		//			adUnit.Cid = *ad.AdslotId
		adUnit.ClickUrl = *adMeta.ClickUrl
		//todo hardcode 3 for MH, only support picture ad
		adUnit.Cid = fmt.Sprint(ad.AdId)
		//adUnit.CreativeType = 3
		adUnit.ClickUrl = *adMeta.ClickUrl
		adUnit.CreativeUrls = []string{*adMeta.MediaUrl}

		//adUnit.ImpTrackingUrls = []string{*adMeta.MediaUrl}
		adUnit.ImpTrackingUrls = adMeta.WinNoticeUrl
		//			 baidu doens't need the tracking url
		adUnit.ClkTrackingUrls = nil
		adUnit.AdWidth = int(*adMeta.MediaWidth)
		adUnit.AdHeight = int(*adMeta.MediaHeight)
		adResponse.StatusCode = lib.STATUS_SUCCESS
	} else {
		// no ads returned from baidu
		adResponse.StatusCode = lib.ERROR_NOAD
	}
	//	}
}

func generateErrorResp(errorCode int, message string, err error, demand *Demand) {
	beego.Critical(err.Error())
	adResponse := generateErrorResponse(demand.AdRequest, demand.AdspaceKey, errorCode)
	demand.Result <- adResponse
}

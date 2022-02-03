// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package ngsetup

import (
	"fmt"
	"gnbsim/factory"
	"gnbsim/gnodeb"
	"gnbsim/gnodeb/context"
	profCtx "gnbsim/profile/context"
	"net"
)

func NgSetup_test(profile *profCtx.Profile) {
	// create amf

	gnb, err := factory.AppConfig.Configuration.GetGNodeB(profile.GnbName)
	if err != nil {
		profile.Log.Errorln("GetGNodeB returned:", err)
	}

	addrs, err := net.LookupHost("amf")
	if err != nil {
		fmt.Println("Failed to resolve amf")
		return
	}

	gnbamf := context.NewGnbAmf(addrs[0], context.NGAP_SCTP_PORT)

	err = gnb.CpTransport.ConnectToPeer(gnbamf)
	if err != nil {
		profile.Log.Errorln("ConnectToAmf returned:", err)
		return
	}

	successFulOutcome, err := gnodeb.PerformNgSetup(gnb, gnbamf)
	if err != nil {
		profile.Log.Errorln("PerformNGSetup returned:", err)
	} else if !successFulOutcome {
		profile.Log.Infoln("Result: FAIL, Expected SuccessfulOutcome, received UnsuccessfulOutcome")
		return
	}

	profile.Log.Infoln("Result: PASS")
}

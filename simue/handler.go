// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package simue

import (
	"fmt"
	"gnbsim/common"
	"gnbsim/simue/context"
	"time"
)

func HandleProfileStartEvent(ue *context.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	ue.Procedure = ue.ProfileCtx.GetFirstProcedure()
	ue.Log.Infoln("Updated procedure to", ue.Procedure)
	HandleProcedure(ue)
	return nil
}

func HandleRegRequestEvent(ue *context.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	SendToGnbUe(ue, intfcMsg)
	return nil
}

func HandleRegRejectEvent(ue *context.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	err = ue.ProfileCtx.CheckCurrentEvent(common.REG_REQUEST_EVENT,
		intfcMsg.GetEventType())
	if err != nil {
		ue.Log.Errorln("CheckCurrentEvent returned:", err)
		return err
	}

	return nil
}

func HandleAuthRequestEvent(ue *context.SimUe,
	intfMsg common.InterfaceMessage) (err error) {

	msg := intfMsg.(*common.UeMessage)
	// checking as per profile if Authentication Request Message is expected
	// from 5G Core against Registration Request message sent by RealUE
	err = ue.ProfileCtx.CheckCurrentEvent(common.REG_REQUEST_EVENT, msg.Event)
	if err != nil {
		ue.Log.Errorln("CheckCurrentEvent returned:", err)
		return err
	}
	nextEvent, err := ue.ProfileCtx.GetNextEvent(msg.Event)
	if err != nil {
		ue.Log.Errorln("GetNextEvent returned:", err)
		return err
	}
	ue.Log.Infoln("Next Event:", common.GetEvtString(nextEvent))
	msg.Event = nextEvent
	SendToRealUe(ue, msg)
	return nil
}

func HandleAuthResponseEvent(ue *context.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UuMessage)
	// Checking if RealUe has sent expected message as per profile against
	// Authentication Request message recevied from 5G Core
	err = ue.ProfileCtx.CheckCurrentEvent(common.AUTH_REQUEST_EVENT, msg.Event)
	if err != nil {
		ue.Log.Errorln("CheckCurrentEvent returned:", err)
		return err
	}

	msg.Event = common.UL_INFO_TRANSFER_EVENT
	SendToGnbUe(ue, msg)
	ue.Log.Traceln("Sending Authentication Response to the network")
	return nil
}

func HandleSecModCommandEvent(ue *context.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	// TODO: Should check if SecModCommandEvent event is expected

	msg := intfcMsg.(*common.UeMessage)
	nextEvent, err := ue.ProfileCtx.GetNextEvent(msg.Event)
	if err != nil {
		ue.Log.Errorln("GetNextEvent returned:", err)
		return err
	}
	msg.Event = nextEvent
	SendToRealUe(ue, msg)
	return nil
}

func HandleSecModCompleteEvent(ue *context.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	ue.Log.Traceln("Handling Security Mode Complete Event")

	msg := intfcMsg.(*common.UuMessage)
	err = ue.ProfileCtx.CheckCurrentEvent(common.SEC_MOD_COMMAND_EVENT,
		msg.Event)
	if err != nil {
		ue.Log.Errorln("CheckCurrentEvent returned:", err)
		return err
	}

	msg.Event = common.UL_INFO_TRANSFER_EVENT
	SendToGnbUe(ue, msg)
	ue.Log.Traceln("Sent Security Mode Complete to the network")
	return nil
}

func HandleRegAcceptEvent(ue *context.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UeMessage)
	// TODO: Should check if Registration Accept event is expected
	nextEvent, err := ue.ProfileCtx.GetNextEvent(msg.Event)
	if err != nil {
		ue.Log.Errorln("GetNextEvent returned:", err)
		return err
	}
	msg.Event = nextEvent
	SendToRealUe(ue, msg)
	return nil
}

func HandleRegCompleteEvent(ue *context.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UuMessage)
	err = ue.ProfileCtx.CheckCurrentEvent(common.REG_ACCEPT_EVENT, msg.Event)
	if err != nil {
		ue.Log.Errorln("CheckCurrentEvent returned:", err)
		return err
	}

	msg.Event = common.UL_INFO_TRANSFER_EVENT
	SendToGnbUe(ue, msg)
	ue.Log.Traceln("Sent Registration Complete to the network")

	ChangeProcedure(ue)
	return nil
}

func HandleDeregRequestEvent(ue *context.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UuMessage)
	msg.Event = common.UL_INFO_TRANSFER_EVENT
	SendToGnbUe(ue, msg)
	ue.Log.Traceln("Sent Deregistration Request to the network")

	return nil
}

func HandleDeregAcceptEvent(ue *context.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	return nil
}

func HandlePduSessEstRequestEvent(ue *context.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UuMessage)
	msg.Event = common.UL_INFO_TRANSFER_EVENT
	SendToGnbUe(ue, msg)
	return nil
}

func HandlePduSessEstAcceptEvent(ue *context.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UeMessage)
	err = ue.ProfileCtx.CheckCurrentEvent(common.PDU_SESS_EST_REQUEST_EVENT, msg.Event)
	if err != nil {
		ue.Log.Errorln("CheckCurrentEvent returned:", err)
		return err
	}
	nextEvent, err := ue.ProfileCtx.GetNextEvent(msg.Event)
	if err != nil {
		ue.Log.Errorln("GetNextEvent returned:", err)
		return err
	}
	ue.Log.Infoln("Next Event:", common.GetEvtString(nextEvent))
	msg.Event = nextEvent
	SendToRealUe(ue, msg)
	return nil
}

func HandlePduSessEstRejectEvent(ue *context.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	err = ue.ProfileCtx.CheckCurrentEvent(common.PDU_SESS_EST_REQUEST_EVENT,
		intfcMsg.GetEventType())
	if err != nil {
		ue.Log.Errorln("CheckCurrentEvent returned:", err)
		return err
	}

	return nil
}

func HandleDlInfoTransferEvent(ue *context.SimUe,
	msg common.InterfaceMessage) (err error) {

	SendToRealUe(ue, msg)
	return nil
}

func HandleDataBearerSetupRequestEvent(ue *context.SimUe,
	msg common.InterfaceMessage) (err error) {

	SendToRealUe(ue, msg)
	return nil
}

func HandleDataBearerSetupResponseEvent(ue *context.SimUe,
	msg common.InterfaceMessage) (err error) {

	SendToGnbUe(ue, msg)

	ChangeProcedure(ue)
	return nil
}

func HandleDataPktGenSuccessEvent(ue *context.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	ChangeProcedure(ue)
	return nil
}

func HandleDataPktGenFailureEvent(ue *context.SimUe,
	msg common.InterfaceMessage) (err error) {

	SendToProfile(ue, common.PROFILE_FAIL_EVENT, msg.GetErrorMsg())
	return nil
}

func HandleServiceRequestEvent(ue *context.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	err = ConnectToGnb(ue)
	if err != nil {
		return fmt.Errorf("failed to connect gnb:", err)
	}

	SendToGnbUe(ue, intfcMsg)

	ue.Log.Traceln("Sent Service Request Event to the network")
	return nil
}

func HandleServiceAcceptEvent(ue *context.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	err = ue.ProfileCtx.CheckCurrentEvent(common.SERVICE_REQUEST_EVENT,
		intfcMsg.GetEventType())
	if err != nil {
		ue.Log.Errorln("CheckCurrentEvent returned:", err)
		return err
	}

	return nil
}

func HandleConnectionReleaseRequestEvent(ue *context.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {
	msg := intfcMsg.(*common.UuMessage)

	if ue.Procedure == common.AN_RELEASE_PROCEDURE {
		err = ue.ProfileCtx.CheckCurrentEvent(common.TRIGGER_AN_RELEASE_EVENT,
			common.CONNECTION_RELEASE_REQUEST_EVENT)
		if err != nil {
			return err
		}
	}

	ue.WriteGnbUeChan = nil

	if msg.TriggeringEvent == common.DEREG_REQUEST_UE_ORIG_EVENT {
		msg := &common.UuMessage{}
		msg.Event = common.QUIT_EVENT
		ue.ReadChan <- msg
		// Once UE is deregistered, Sim UE is not expecting any further
		// procedures
		SendToProfile(ue, common.PROFILE_PASS_EVENT, nil)
		return nil
	}

	SendToRealUe(ue, msg)
	ChangeProcedure(ue)

	return nil
}

func HandleErrorEvent(ue *context.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	SendToProfile(ue, common.PROFILE_FAIL_EVENT, intfcMsg.GetErrorMsg())

	msg := &common.UuMessage{}
	msg.Event = common.QUIT_EVENT
	HandleQuitEvent(ue, msg)
	return nil
}

func HandleQuitEvent(ue *context.SimUe,
	msg common.InterfaceMessage) (err error) {
	if ue.WriteGnbUeChan != nil {
		SendToGnbUe(ue, msg)
	}
	SendToRealUe(ue, msg)
	ue.WriteRealUeChan = nil
	ue.WaitGrp.Wait()
	ue.Log.Infoln("Sim UE terminated")
	return nil
}

func ChangeProcedure(ue *context.SimUe) {
	nextProcedure := ue.ProfileCtx.GetNextProcedure(ue.Procedure)
	if nextProcedure != 0 {
		ue.Procedure = nextProcedure
		ue.Log.Infoln("Updated procedure to", common.GetProcString(nextProcedure))
		HandleProcedure(ue)
	} else {
		SendToProfile(ue, common.PROFILE_PASS_EVENT, nil)
		evt, err := ue.ProfileCtx.GetNextEvent(common.PROFILE_PASS_EVENT)
		if err != nil {
			ue.Log.Errorln("GetNextEvent failed:", err)
			return
		}
		if evt == common.QUIT_EVENT {
			msg := &common.DefaultMessage{}
			msg.Event = common.QUIT_EVENT
			ue.ReadChan <- msg
		}
	}
}

func HandleProcedure(ue *context.SimUe) {
	switch ue.Procedure {
	case common.REGISTRATION_PROCEDURE:
		ue.Log.Infoln("Initiating Registration Procedure")
		msg := &common.UeMessage{}
		msg.Event = common.REG_REQUEST_EVENT
		SendToRealUe(ue, msg)
	case common.PDU_SESSION_ESTABLISHMENT_PROCEDURE:
		ue.Log.Infoln("Initiating UE Requested PDU Session Establishment Procedure")
		msg := &common.UeMessage{}
		msg.Event = common.PDU_SESS_EST_REQUEST_EVENT
		SendToRealUe(ue, msg)
	case common.USER_DATA_PKT_GENERATION_PROCEDURE:
		ue.Log.Infoln("Initiating User Data Packet Generation Procedure")
		msg := &common.UeMessage{}
		msg.UserDataPktCount = ue.ProfileCtx.DataPktCount
		msg.Event = common.DATA_PKT_GEN_REQUEST_EVENT

		time.Sleep(500 * time.Millisecond)
		/* TODO: Solve timing issue. Currently UE may start sending user data
		 * before gnb has successfuly sent PDU Session Resource Setup Response
		 * or before 5g core has processed it
		 */
		ue.Log.Infoln("Please wait, initiating uplink user data in 3 seconds ...")
		time.Sleep(3 * time.Second)

		SendToRealUe(ue, msg)
	case common.UE_INITIATED_DEREGISTRATION_PROCEDURE:
		ue.Log.Infoln("Initiating UE Initiated Deregistration Procedure")
		msg := &common.UeMessage{}
		msg.Event = common.DEREG_REQUEST_UE_ORIG_EVENT
		SendToRealUe(ue, msg)
	case common.AN_RELEASE_PROCEDURE:
		ue.Log.Infoln("Initiating AN Release Procedure")
		msg := &common.UeMessage{}
		msg.Event = common.TRIGGER_AN_RELEASE_EVENT
		SendToGnbUe(ue, msg)
	case common.UE_TRIGGERED_SERVICE_REQUEST_PROCEDURE:
		ue.Log.Infoln("Initiating UE Triggered Service Request Procedure")
		msg := &common.UeMessage{}
		msg.Event = common.SERVICE_REQUEST_EVENT
		SendToRealUe(ue, msg)
	}
}

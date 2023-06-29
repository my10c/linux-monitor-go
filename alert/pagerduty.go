// Copyright (c) 2017 - 2017 badassops
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//	* Redistributions of source code must retain the above copyright
//	notice, this list of conditions and the following disclaimer.
//	* Redistributions in binary form must reproduce the above copyright
//	notice, this list of conditions and the following disclaimer in the
//	documentation and/or other materials provided with the distribution.
//	* Neither the name of the <organization> nor the
//	names of its contributors may be used to endorse or promote products
//	derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSEcw
// ARE DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
// DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Version		:	0.1
//
// Date			:	June 4, 2017
//
// History	:
// 	Date:			Author:		Info:
//	June 4, 2017	LIS			First Go release
//
// TODO:

package alerts

import (
	"fmt"
	"time"

	myGlobal	"github.com/my10c/linux-monitor-go/global"

	"github.com/PagerDuty/go-pagerduty"
)

// Function to post an alert in pagerduty
func alertPD(message string, tag string, hostname string) error {
	var validTime string
	// create the IncidentKey key value, based on time to that the event is being
	// todo so we use the global setting whatever we using hour or minute
	// we do not check invalid choice, it either hour or minute, with hour being default!
	currTime := time.Now().Local()
	if myGlobal.DefaultPDValidUnit == "minute" {
		validTime = fmt.Sprintf("%s", currTime.Format("2006-01-02-04"))
	} else {
		validTime = fmt.Sprintf("%s", currTime.Format("2006-01-02-15"))
	}
	incidentKey := fmt.Sprintf("%s-%s-%s-%s", validTime, tag, hostname, myGlobal.MyProgname)
	// create the description based on the give message
	descriptionEvent := fmt.Sprintf("%s-%s-%s", myGlobal.DefaultPD["pdevent"], myGlobal.MyProgname, hostname)
	// and now the detail
	detailsEvent := message
	// the payload/event to be sent
	eventPD := pagerduty.Event{
		ServiceKey:		myGlobal.DefaultPD["pdservicekey"],
		Type:			"trigger",
		IncidentKey:	incidentKey,
		Description:	descriptionEvent,
		Details:		detailsEvent,
	}
	_, err := pagerduty.CreateEvent(eventPD)
	return err
}

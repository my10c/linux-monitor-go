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
	"strings"

	myGlobal	"github.com/my10c/linux-monitor-go/global"

	"github.com/nlopes/slack"
)

// Function to post an alert in slack
func alertSlack(message string) error {
	slackAPI := slack.New(myGlobal.DefaultSlack["slackservicekey"])
	slackMsg := fmt.Sprintf(":imp: `- %s error: %s` :disappointed:\n", myGlobal.MyProgname, message)
	// remove all carriage return
	slackMsg = strings.TrimSuffix(strings.Replace(slackMsg, "\n", " - ", -1), " - ")
	// need to build a minimum config, the user profile
	slackUserProfile := slack.PostMessageParameters{
		Username:		myGlobal.DefaultSlack["slackuser"],
		AsUser:			false,
		Parse:			"",
		LinkNames:		0,
		Attachments:	nil,
		UnfurlLinks:	false,
		UnfurlMedia:	true,
		IconURL:		"",
		IconEmoji:		myGlobal.DefaultSlack["iconemoji"],
		Markdown:		true,
		EscapeText:		true,
	}
	_, _, err := slackAPI.PostMessage(myGlobal.DefaultSlack["slackchannel"], slackMsg, slackUserProfile)
	return err
}

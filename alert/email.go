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
	"net/mail"
	"net/smtp"

	myGlobal	"github.com/my10c/linux-monitor-go/global"
)

// Function to send an alert email
func alertEmail(message string, subject string) error {
	// if authEmail - empty then no authentication is required
	// for now only support PlainAuth
	authEmail := smtp.PlainAuth("",
			myGlobal.DefaultEmail["emailuser"],
			myGlobal.DefaultEmail["emailpass"],
			myGlobal.DefaultEmail["emailhost"],
	)
	// build the email component
	emailTo := mail.Address{myGlobal.DefaultEmail["emailtoname"], myGlobal.DefaultEmail["emailto"]}
	emailFrom := mail.Address{myGlobal.DefaultEmail["emailfromname"], myGlobal.DefaultEmail["emailfrom"]}
	emailHost := fmt.Sprintf("%s:%s", myGlobal.DefaultEmail["emailhost"], myGlobal.DefaultEmail["emailhostport"])
	fromAndBody := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n",
			emailTo.String(), subject, message)
	// send the email
	err := smtp.SendMail(emailHost, authEmail, emailFrom.String(), []string{emailTo.String()}, []byte(fromAndBody))
	return err
}

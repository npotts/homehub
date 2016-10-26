/*
Copyright (c) 2016 Nick Potts
Licensed to You under the GNU GPLv3
See the LICENSE file at github.com/npotts/homehub/LICENSE

This file is part of the HomeHub project
*/

package homehub

/*An Attendant performs the function of listening for data messages and forwarding them to a backend to store*/
type Attendant interface {
	Use(Backend) //where do we aim messages
	Stop()       //cease operations and exit
}

/*A Backend supports multiple attendants that gather
messages by various means*/
type Backend interface {
	Register(Datam) //registers a datam
	Store(Datam)    //Stores datam
	Stop()          //cease operations
}

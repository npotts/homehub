/*
 GNU GENERAL PUBLIC LICENSE
                       Version 3, 29 June 2007

 Copyright (C) 2007 Free Software Foundation, Inc. <http://fsf.org/>
 Everyone is permitted to copy and distribute verbatim copies
 of this license document, but changing it is not allowed.*/

package brainiac

import ()

/*function call to both register and store data*/
type regstore *func(table string, data map[string]interface{}) error

/*listener listen on something and call registerFxn when*/
type stoppable interface {
	stop()
}

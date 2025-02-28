package main

import (
	"fmt"

	"github.com/webbben/note-utils/pkg/cleanup"
)

var sampleNote = `
Had a call with John today and he was quite unhappy about the last shipment of hamster food we sent.
He said the box was just full of a bunch of cracked seed shells and there was a hole in side.

Did those damn hamsters get in again? We gotta hunt down those guys, they're ruining our business!
Who would've thought a rogue gang of street hamsters would be our company's biggest issue...

I told him that we are making every effort to correct the situation and we will cover the costs.

Next week, we need to do the following: make a call to the chamber of hamster affairs to file a complaint.
Then, we will box a new shipment of Squeaky Deluxe and send it to John. We should also contact that cat guy and see if he has any battle cats for rent.
We will need 24/hr surveillance of the warehouse and those battle cats are the best of the best. Let's also send out a notice to the other distributors
of hamster food products that there is a rogue gang of hamsters pillaging the area.
`

func main() {
	fmt.Println("input:")
	fmt.Println(sampleNote)
	out, err := cleanup.CleanNoteWithOpts(sampleNote, cleanup.CleanNoteOpts{})
	if err != nil {
		panic(err)
	}

	fmt.Println("\n\noutput:")
	fmt.Println(out)
}

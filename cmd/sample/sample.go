package main

import (
	"fmt"

	notecleanup "github.com/webbben/note-utils/pkg/note-cleanup"
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

var sampleNote1 = `
# Meeting Notes 2/6/2025

Met with the cybersecurity team again about the data breach issue again. For some reason they are really upset about it. I can't really understand why, I mean
this time it was only 1.2k accounts! Much better than the 20k from last time... sheesh.

Anyway, they made the following requests:

- next week have another meeting on tuesday to review vulnerabilities again
- reach out to backend team to review why our server is giving out admin passwords when you query it's IP without auth tokens
- implement a code review and testing process, to hopefully mitigate futre issues.

Man these guys are a pain in the ass!`

func main() {
	out, err := notecleanup.CleanNoteWithOpts(sampleNote, notecleanup.CleanNoteOpts{})
	if err != nil {
		panic(err)
	}

	fmt.Println(out)
}

package main

import (
	"encoding/json"
	"math/rand"
	"time"
)

var ticker *time.Ticker

type dataMsg1 struct {
	Key HashKey
}

type dataMsg2 struct {
	Key   HashKey
	Value string
}

type doMsgs struct {
	Do string
}

type joinRingMsg struct {
	Do      string
	Sponsor HashKey
}

type leaveRingMsg struct {
	Do   string
	Mode string
}

type doRespondToMsgs struct {
	Do        string
	RespondTO HashKey
}

type findRingSPMsg struct {
	Do        string
	RespondTO HashKey
	TargetID  HashKey
}

type putMsg struct {
	Do   doRespondToMsgs
	Data dataMsg1
}

type getRemMsgs struct {
	Do   doRespondToMsgs
	Data dataMsg2
}

func injectRequests() {
	ticker = time.NewTicker(1500 * time.Millisecond)
	c := 0
	go func() {
		for range ticker.C {
			//coordinateChan <- generateMessages(t)
			if c == 0 {
				coordinateChan <- generateRandomMessage()
				c++
			}
		}
	}()
}

func generateMessages(timeSeed time.Time) int {
	return randomGenerator(timeSeed, 0, 21)
}

//Generate messages
func generateRandomMessage() []byte {
	//generate random message
	//1. join-ring msg
	sponsorKey := nodeList[rand.Intn(len(nodeList))]

	msg1 := &joinRingMsg{
		Do:      "join-ring",
		Sponsor: sponsorKey,
	}
	message, _ := json.Marshal(msg1)

	/*choice := rand.Intn(12)
	switch choice {
	case 0:
		return joinRingMessage
	case 1:
		return leaveRingMessage
	case 2:
		return stbRingMessage
	case 3:
		return initRingMessage
	case 4:
		return fixRingMessage
	case 5:
		return ringNotifyMessage
	case 6:
		return getRingFinMessage
	case 7:
		return findSMessage
	case 8:
		return findPMessage
	case 9:
		return putMessage
	case 10:
		return getMessage
	case 11:
		return remMessage
	}
	*/
	return message
}

func randomGenerator(timeSeed time.Time, min int, max int) int {
	rand.Seed(timeSeed.UTC().UnixNano())
	return rand.Intn(max-min) + min
}

func triggerSuccesorMessage(sponsor HashKey, recipient HashKey) []byte {
	findSuccesorM := &findRingSPMsg{
		Do:        "find-ring-successor",
		RespondTO: sponsor,
		TargetID:  recipient,
	}
	fsMessage, _ := json.Marshal(findSuccesorM)
	return fsMessage
}

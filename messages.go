package main

import (
	"encoding/json"
	"math/rand"
	"time"
)

var ticker *time.Ticker

type doMsgs struct {
	Do string
}

type joinRingMsg struct {
	Do      string
	Sponsor HashKey
}

type leaveRingMsg struct {
	Do        string
	Mode      string
	Recipient HashKey
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
	Do        string
	RespondTO HashKey
	Data      dataMsg2
}

type getRemMsgs struct {
	Do        string
	RespondTO HashKey
	Data      dataMsg1
}

type dataMsg1 struct {
	Key HashKey
}

type dataMsg2 struct {
	Key   HashKey
	Value string
}

type updateBucketsMsg struct {
	Do         string
	BucketData map[HashKey]string
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
func generateRandomMessage() string {
	//generate random message

	//1. join-ring msg
	/*


		msg1 := &joinRingMsg{
			Do:      "join-ring",
			Sponsor: sponsorKey,
		}
		marshalledMessage, _ := json.Marshal(msg1)

		message := string(marshalledMessage)
	*/
	// msg1 := &joinRingMsg{
	// 	Do:      "join-ring",
	// 	Sponsor: sponsorKey,
	// }
	sponsorKey := nodeList[rand.Intn(len(nodeList))]
	msg1 := &leaveRingMsg{
		Do:        "leave-ring",
		Recipient: sponsorKey,
		Mode:      "orderly",
	}

	// marshalledMessage, _ := json.Marshal(msg1)

	// //bucket list messages
	// sponsor := nodeList[rand.Intn(len(nodeList))]
	// key := genKey(randString())
	/*
		//9. put msg
		datMsg := dataMsg2{
			Key:   key,
			Value: "val",
		}
		msg9 := &putMsg{
			Do:        "put",
			RespondTO: sponsor,
			Data:      datMsg,
		}
		marshalledMessage, _ := json.Marshal(msg9)
		message := string(marshalledMessage)
	*/
	//10. get msg
	// datMsg := dataMsg1{
	// 	Key: key,
	// }
	// msg10 := &getRemMsgs{
	// 	Do:        "get",
	// 	RespondTO: sponsor,
	// 	Data:      datMsg,
	// }
	marshalledMessage, _ := json.Marshal(msg1)
	message := string(marshalledMessage)

	/*
		//11. remove msg
		datMsg := dataMsg1{
			Key: key,
		}
		msg11 := &getRemMsgs{
			Do:        "remove",
			RespondTO: sponsor,
			Data:      datMsg,
		}
		marshalledMessage, _ := json.Marshal(msg11)
		message := string(marshalledMessage)
	*/
	/*
		choice := rand.Intn(12)
		switch choice {
		case 0:
			{
				sponsorKey := nodeList[rand.Intn(len(nodeList))]

				msg0 := &joinRingMsg{
					Do:      "join-ring",
					Sponsor: sponsorKey,
				}
				joinRingMessage, _ := json.Marshal(msg0)
				return joinRingMessage
			}
		case 1:
			{
				ch := rand.Intn(2)
				switch ch {
				case 0:
					{
						mode := "immediate"
						msg1 := &leaveRingMsg{
							Do:   "leave-ring",
							Mode: mode,
						}
						leaveRingMessage, _ := json.Marshal(msg1)
						return leaveRingMessage
					}
				case 1:
					{
						mode := "orderly"
						msg1 := &leaveRingMsg{
							Do:   "leave-ring",
							Mode: mode,
						}
						leaveRingMessage, _ := json.Marshal(msg1)
						return leaveRingMessage

					}
				}
			}
		case 2:
			{
				msg2 := &doMsgs{
					Do: "stabilize-ring",
				}
				stbRingMessage, _ := json.Marshal(msg2)
				return stbRingMessage
			}
		case 3:
			{
				msg3 := &doMsgs{
					Do: "init-ring-fingers",
				}
				initRingMessage, _ := json.Marshal(msg3)
				return initRingMessage
			}
		case 4:
			{
				msg4 := &doMsgs{
					Do: "fix-ring-fingers",
				}
				fixRingMessage, _ := json.Marshal(msg4)
				return fixRingMessage
			}
		case 5:
			{
				key := nodeList[rand.Intn(len(nodeList))]
				msg5 := &doRespondToMsgs{
					Do:        "ring-notify",
					RespondTO: key,
				}
				ringNotifyMessage, _ := json.Marshal(msg5)
				return ringNotifyMessage
			}
		case 6:
			{
				key := nodeList[rand.Intn(len(nodeList))]
				msg6 := &doRespondToMsgs{
					Do:        "get-ring-fingers",
					RespondTO: key,
				}
				getRingFinMessage, _ := json.Marshal(msg6)
				return getRingFinMessage
			}
		case 7:
			{
				sponsor := nodeList[rand.Intn(len(nodeList))]
				key := nodeList[rand.Intn(len(nodeList))]
				msg7 := &findRingSPMsg{
					Do:        "find-ring-successor",
					RespondTO: sponsor,
					TargetID:  key,
				}
				findSMessage, _ := json.Marshal(msg7)
				return findSMessage
			}
		case 8:
			{
				sponsor := nodeList[rand.Intn(len(nodeList))]
				key := nodeList[rand.Intn(len(nodeList))]
				msg8 := &findRingSPMsg{
					Do:        "find-ring-predecessor",
					RespondTO: sponsor,
					TargetID:  key,
				}
				findPMessage, _ := json.Marshal(msg8)
				return findPMessage
			}
		case 9:
			{
				sponsor := nodeList[rand.Intn(len(nodeList))]
				key := nodeList[rand.Intn(len(nodeList))]
				datMsg := dataMsg2{
					Key:   key,
					Value: "val",
				}
				msg9 := &putMsg{
					Do:        "put",
					RespondTO: sponsor,
					Data:      datMsg,
				}
				putMessage, _ := json.Marshal(msg9)
				return putMessage
			}
		case 10:
			{
				sponsor := nodeList[rand.Intn(len(nodeList))]
				key := nodeList[rand.Intn(len(nodeList))]
				datMsg := dataMsg1{
					Key: key,
				}
				msg10 := &getRemMsgs{
					Do:        "get",
					RespondTO: sponsor,
					Data:      datMsg,
				}
				getMessage, _ := json.Marshal(msg10)
				return getMessage
			}
		case 11:
			{
				sponsor := nodeList[rand.Intn(len(nodeList))]
				key := nodeList[rand.Intn(len(nodeList))]
				datMsg := dataMsg1{
					Key: key,
				}
				msg11 := &getRemMsgs{
					Do:        "remove",
					RespondTO: sponsor,
					Data:      datMsg,
				}
				remMessage, _ := json.Marshal(msg11)
				return remMessage
			}
		}
	*/
	return message
}

func randomGenerator(timeSeed time.Time, min int, max int) int {
	rand.Seed(timeSeed.UTC().UnixNano())
	return rand.Intn(max-min) + min
}

func triggerSuccesorMessage(sponsor HashKey, recipient HashKey) string {
	findSuccesorM := &findRingSPMsg{
		Do:        "find-ring-successor",
		RespondTO: sponsor,
		TargetID:  recipient,
	}
	fsMessage, _ := json.Marshal(findSuccesorM)
	return string(fsMessage)
}

func initRingFingMessage() string {
	msg3 := &doMsgs{
		Do: "init-ring-fingers",
	}
	initRingMessage, _ := json.Marshal(msg3)
	return string(initRingMessage)
}

func getRingFingMessage(key HashKey) string {

	msg6 := &doRespondToMsgs{
		Do:        "get-ring-fingers",
		RespondTO: key,
	}
	getRingFinMessage, _ := json.Marshal(msg6)
	return string(getRingFinMessage)

}

func updateBucketMessage(bucketData map[HashKey]string) string {

	msg := &updateBucketsMsg{
		Do:         "update-bucket",
		BucketData: bucketData,
	}
	updateBucketMessage, _ := json.Marshal(msg)
	return string(updateBucketMessage)

}

func triggerPredecessorMessage(sponsor HashKey, recipient HashKey) string {
	findPredecessorM := &findRingSPMsg{
		Do:        "find-ring-predecessor",
		RespondTO: sponsor,
		TargetID:  recipient,
	}
	fsMessage, _ := json.Marshal(findPredecessorM)
	return string(fsMessage)
}

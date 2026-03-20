
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
)

// # Message Queue
// Input  : Message string to enqueue (hardcoded in this implementation)
// Output : Error code on success/fail

/*
Since the queue and database have to be stored on separate data storage, this
implements the connection as well as the message enqueue. From the OS (set up 
on personal machine, native on azure), it gets the enviroment variables to
create a URL, credential and then to queueing a base 64 string to the message
queue for the later implemented stats function.
*/

var accountName = os.Getenv(ACCOUNT_NAME_VAR)
var accountKey = os.Getenv(ACCOUNT_KEY_VAR)

// ## Private/Helper Functiosn ------------------------------------------------

func getQueueURL(queue string) (string) {
	return fmt.Sprintf(
		"https://%s.queue.core.windows.net/%s",
		accountName,
		queue,
	)
}

func getCredential() (*(azqueue.SharedKeyCredential), error) {
	cred, err := azqueue.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	return cred, nil
}

func getQueueClient(
	credential *(azqueue.SharedKeyCredential),
) (*(azqueue.QueueClient), error) {
	queueClient, err := azqueue.NewQueueClientWithSharedKeyCredential(
		getQueueURL(QUEUE), credential, nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create queue client: %w", err)
	}

	return queueClient, nil
}

func encodeMsg(msg string) string {
	return base64.StdEncoding.EncodeToString([]byte(msg))
}

// ## Public/Principle Function ===============================================

func enqueueMessage(message string) error {
	ctx := context.Background()

	cred, err := getCredential()
	if err != nil { return err }

	queueClient, err := getQueueClient(cred)
	if err != nil { return err }

	msg := encodeMsg(message)
	_, err = queueClient.EnqueueMessage(ctx, msg, nil)
	if err != nil {
		return fmt.Errorf("failed to enqueue message: %w", err)
	}

	log.Printf("Enqueued message to queue '%s': %s", QUEUE, message)
	return nil
}
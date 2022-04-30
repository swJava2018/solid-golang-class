package consumers

import "context"

type Consumer interface {
	//create Consumer instance
	Create() error

	//read data
	Read(ctx context.Context, stream chan interface{}, errc chan error, shutdown chan bool) error

	//delete Consumer instance
	Delete() error
}

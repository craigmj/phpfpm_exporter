# phpfpm_exporter
Prometheus exporter for php-fpm status information.

# Installation

You will need to have Go http://golang.org.

Checkout the source:

    git clone https://github.com/craigmj/phpfpm_exporter

Build the code:

    cd phpfpm_exporter
    ./build.sh

Now you've got the binary in `bin/phpfpm_exporter`

# Running

Because this is a Go application, you only need the binary. So you can copy that to your server, which needs to be configured to report php-fpm status information.

Start the exporter with

	phpfpm_exporter

# Command Line Options

The exporter accepts 3 command line options

`status.url`	The URL from which to scrape the php-fpm status. Defaults to `http://localhost/status`

`listen.address`	The address on which to serve the php-fpm metrics to prometheus. Defaults to `http://127.0.0.1:9099/`

`scrape.interval`	The interval between scraps of the php-fpm status. Defaults to `5m`. You can set any string value that can be parsed by Go's Duration class: https://golang.org/pkg/time/#ParseDuration

# Metrics

The phpfpm_exporter exports the following metrics to prometheus:

## phpfpm_acceptedconnections_count
The number of connections accepted by the `pool`.
### Time Series labels
`pool` - the fpm pool

## phpfpm_listenqueue_size
The size of the listen queue for each pool.
### Time Series labels
`pool` - the fpm pool
`metric` - one of 
  `current`  is the `listen queue` value from php-fpm: the number of requests in the queue of pending connections
  `max` is the `max listen queue` value from php-fpm: the maximum number of requests in the queue of pending connections since FPM started
  `len` is the `listen queue len` value from php-fpm: the size of the socket queue of pending connections
  
## phpfpm_processes_count
The number of processes in each pool.
### Time Series labels
`pool` - the fpm pool
`state` - the state of the proceses. One of
   `idle` - the number of idle processes
   `active` - the number of active processes
   `max_active` - the maximum number of active processes in this pool since fpm started
   Note that we don't reflect the `total processes` as reported by php-fpm, since that is simply `idle+active` and can thus be calculated.

## phpfpm_maxchildren_count
The maximum number of child processes reached in the pool.
### Time Series labels
`pool` - the fpm pool
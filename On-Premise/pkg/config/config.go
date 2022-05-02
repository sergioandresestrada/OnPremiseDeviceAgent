package config

// NumberOfRetries refers to the number of times that a message will be tried
// to be delivered in case of failure
const NumberOfRetries = 5

// InitialTimeBetweenRetries refers to the number of seconds of waiting time
// before the first retry in case of failure while delivering a message
const InitialTimeBetweenRetries = 15

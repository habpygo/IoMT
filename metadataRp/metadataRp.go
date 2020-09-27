/* Copyright 2018 Harry Boer - DappDevelopment.com
Licensed under the MIT License, see LICENCE file for details.
*/

/* Package metadata contains the paths, seeds, addresses, and provider */
package metadataRp

var Seed = "BMYDBMQIKMPAVPY9CGQWSAQEVV9JAYDIZUOIWLQPGKLBIXEGY9S9BJBCQUSQODIYFUIEKPWAHUPNXQKAH"

//var Seed = "BKQMEXEUHFLFEZJMDXHIZIUUMWMQCZWMEZJIMIXFDZLACHTEXCNKTQONVADEGBTTAPLIWJFXUEIJWHGSB"

var Address = "HZIQ9FHDEQEQLMKFVXLMLEWXBUZQORXGWIWFMSOCWA9KMIXRJ9HTTUPTQMUDFWBCUVUZSACHECQQGPHNC"

//var Address = "ZFDABOULAUQAEDLXMHKRVAUKNCJZKZWQCWYKMDFKNWFKGQRT9KTAMFMKIYDTHKAQUKHXMZYJOCB9SDVD9"

var Provider = "https://nodes.devnet.thetangle.org:443"

// "https://nodes.spamnet.iota.org"
//var Provider = "https://nodes.testnet.iota.org:443"

var MWM int64 = 14 // In the real world set this to 14 or 15!

var SecurityLevel int = 2

var GraphSize int = 20

// Harry
//var DATAFILE = "tmp/dat/w1_slave"

// Siert1
var DATAFILE = "/sys/bus/w1/devices/28-00000ac8371f/w1_slave"

// Siert2
//var DATAFILE = "./w1_slave"

/*
Some alternative test nodes one can use

https://nodes.testnet.thetangle.org:443
https://nodes.testnet.iota.org:443
https://testnet.thetangle.org:443
https://node01.testnet.iotatoken.nl:16265
*/

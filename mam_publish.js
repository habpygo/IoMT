// ==========================================================================================
// IOTA Tutorial 19
// MAM Publish

// ==========================================================================================
// MAM object example:

// {
//     "state": {
//         "subscribed": [],
//         "channel": {
//             "side_key": null,
//             "mode": "public",
//             "next_root": "SJLONMWLYHTGAEEL9ORJYWFCIMAKHAEJBKFAAIDBQNFRMNYLXDYMZTZIYKJGPWNWHPEW9GXDCPFPHRC9T",
//             "security": "2",
//             "start": 1,
//             "count": 1,
//             "next_count": 1,
//             "index": 0
//         },
//         "seed": "EQR9LNZAXMWAJI9TZLSGYIEWIWCRHHANNORQAHYADIKLPDJHCOPIRNADXMTWQFJCC9XACTXMRWVCJQHNV"
//     },
//     "payload": "AHBANSN9MUOOZ9PMKJPUCABEZ9YHXDXLMVTDCGYYLCWXWSQDLRXFBOTZDLABBZPFYC9BOPAMVHXIGQNRAQYI9POY...
//     "root": "FOJFWFAECPMXYBEKMMOPBNPQLXATCBWJQNQGCAKWORZGDXREORLOHAAEJBDKJQJYVSDVNHIRODHIYQZXK",
//     "address": "FOJFWFAECPMXYBEKMMOPBNPQLXATCBWJQNQGCAKWORZGDXREORLOHAAEJBDKJQJYVSDVNHIRODHIYQZXK"
// }
//==========================================================================================
//Publish data to the Tangle:


<script type="text/javascript" src="../../scripts/iota/v0.4.7/iota.min.js"></script>
<script type="text/javascript" src="../../scripts/iota/mam.web.js"></script>

let securityLevel = 2;
let seed = 'EQR9LNZAXMWAJI9TZLSGYIEWIWCRHHANNORQAHYADIKLPDJHCOPIRNADXMTWQFJCC9XACTXMRWVCJQHNV';
let host 'https://field.carriota.com';
let port = '443';
let iota = new IOTA({'host': host, 'port': port});
let keyString = 'mysecret';
let key = iota.utils.toTrytes(keyString);
let timeInterval = 15;

let Mam = require('mam.web.js');
let mamState = Mam.init(iota,seed,securityLevel);

// channel mode default public
//mamState = Mam.changeMode(mamState, 'private');
//mamState = Mam.changeMode(mamState, 'restricted', key);

// Publish to Tangle
const publish = async function(packet) {
    // Create MAM Payload
    let trytes = iota.utils.toTrytes(JSON.stringify(packet));

    let message = Mam.create(mamState, trytes);

    // Example message:
    // {
    //     "state": {
    //         "subscribed": [],
    //         "channel": {
    //             "side_key": null,
    //             "mode": "public",
    //             "next_root": "DXLG...FUMX",
    //             "security": "3",
    //             "start": 1,
    //             "count": 1,
    //             "next_count": 1,
    //             "index": 0
    //         },
    //         "seed": "EQR9...QHNV"
    //     },
    //     "payload": "AHBA...CBA9",
    //     "root": "ZHWO...WWUC",
    //     "address": "ZHWO...WWUC"
    // }

    // ***********************************************************************************
    // 1. https://github.com/iotaledger/mam.client.js/blob/master/src/index.js
    //    See: const create = (state, message)
    // 2. https://raw.githubusercontent.com/iotaledger/mam.client.js/master/lib/mam.web.js
    //    See: function createMessage(SEED, MESSAGE, SIDE_KEY, CHANNEL)
    //
    //    a. Set up merkle tree:
    //    See IOTA Tutorial 19 (slide #A)
    //    var root_merkle = iota_merkle_create(SEED_trits, START, COUNT, SECURITY);
    //    var next_root_merkle = iota_merkle_create(SEED_trits, NEXT_START, NEXT_COUNT, SECURITY);
    //    var root_branch = iota_merkle_branch(root_merkle, INDEX);
    //    var root_siblings = iota_merkle_siblings(root_branch);
    //    var next_root_branch = iota_merkle_branch(next_root_merkle, INDEX);
    //    var root = iota_merkle_slice(root_merkle);
    //    var next_root = iota_merkle_slice(next_root_merkle);
    //
    //    b. Calling iota_merkle_create creates the subseed, digest and address:
    //    https://github.com/iotaledger/iota.rs/blob/master/bindings/src/merkle/simple.rs
    //    pub fn iota_merkle_create
    //    But the actual implementation is:
    //    https://github.com/iotaledger/iota.rs/blob/master/merkle/src/simple.rs
    //    pub fn create
    //    In the create function the Merkle tree is setup where the leaves are created:
    //    leaf(seed, index, security, c1, c2, c3)
    //    In the leaf function the following is done:
    //    - Calculate subseed:  iss::subseed(&seed, index, &mut subseed, c1)
    //    - Calculate digest:   iss::subseed_to_digest(&subseed, security, &mut hash, c1, c2, c3)
    //    - Calculate address:  iss::address(&mut hash, c1)
    //    The iss:: function is located at: https://github.com/iotaledger/iota.rs/blob/master/sign/src/iss.rs
    //
    //    c. The address is hashed to create the leaf (and not the root).
    //    See IOTA Tutorial 19 (slide #B)
    //    https://github.com/iotaledger/iota.rs/blob/master/merkle/src/simple.rs
    //    pub fn create
    //
    // 3. Create masked_payload:
    //    See IOTA Tutorial 20
    //    var masked_payload = iota_mam_create(SEED_trits, MESSAGE_trits, SIDE_KEY_trits, root, root_siblings, next_root, START, INDEX, SECURITY);
    //    https://github.com/iotaledger/MAM/blob/master/bindings/src/mam/simple.rs
    //    pub fn iota_mam_create
    //    But the actual implementation is:
    //    https://github.com/iotaledger/MAM/blob/master/mam/src/mam.rs
    //    pub fn create<C, CB, H>(seed: &[Trit],....);
    //    The resulting payload structure looks like:
    //
    //         [
    //             Encoded Index,
    //             Encoded Message Length,
    //             encrypted[
    //                 Message,
    //                 Nonce,
    //                 Signature,
    //                 Encoded Number of Siblings,
    //                 Siblings
    //             ]
    //         ]
    //
    //    where message contains the following data: {"payload": “ODGD..GAQD”, "next_root": “SJLO..RC9T”}
    //
    // 4. Continue with the create function in:
    //    See IOTA Tutorial 19
    //    https://github.com/iotaledger/mam.client.js/blob/master/src/index.js
    //    See: const create = (state, message)
    //    Calculate: - address based on the root.
    //                 If channel mode = public: address = root
    //                 If channel mode = private or restricted: address = HASH(root)
    //               - index and start
    //               - Create object:
    //                 {
    //                   state,
    //                   payload: masked_payload,
    //                   root: root,
    //                   address: address
    //                 }
    // ***********************************************************************************


    console.log(JSON.stringify(message,null,"\t"));

    // Save new mamState
    mamState = message.state;

    // Attach the payload to the Tangle
    await Mam.attach(message.payload, message.address);

    // ***********************************************************************************
    // 5. https://github.com/iotaledger/mam.client.js/blob/master/src/index.js
    //    See: const attach = async (trytes, root, depth = 6, mwm = 14)
    //    - Create var transfers = [{address: root, value: 0,  message: trytes}]
    //    - Create a rondom seed to make a transaction:
    //    - Call iota.api.sendTransfer(keyGen(81), depth = 6, mwm = 14, transfers)
    //      https://github.com/iotaledger/iota.lib.js/blob/master/lib/api/api.js
    //      The sendTransfer function requires a seed, but this seed is not used.
    //      Therefore a dummy valid seed is generated.
    //    - Call iota.api.prepareTransfers = function(seed, transfers, options, callback)
    //      https://github.com/iotaledger/iota.lib.js/blob/master/lib/api/api.js
    //      to create the signatureFragments
    //
    // 6. A transaction bundle is created.
    //    See IOTA Tutorial 19 (slide #C, #D, #E)
    // ***********************************************************************************

    return message.root;
}

// Create simulated sensor data
// For example json = {"data":40,"dateTime":"23/02/2018 10:54:34"}
const generateDummyJSON = function(){
    let randomNumber = Math.floor((Math.random()*89)+10);
    let dateTime = getDateAndTime();

    // For example: json = { "data": 54, "dateTime": "23/02/2018 10:30:07" }
    let json = {"data": randomNumber, "dateTime":dateTime};
    return json;
}

const executeDataPublishing = async function() {
    let json = generateDummyJSON();
    let root = await publish(json);

    console.log("dateTime: "+json.dateTime+", "+"data: "+json.data+", "+"root: "+root);
}

// Time interval set to every 15 seconds
intervalInstance = setInterval(executeDataPublishing, parseInt(timeInterval)*1000);
// ==========================================================================================
// IOTA Tutorial 19
// MAM Receive

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
// ==========================================================================================
// Receive data from the Tangle:

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

//mamState = Mam.changeMode(mamState, 'private');
//mamState = Mam.changeMode(mamState, 'restricted', key);

let root = 'FOJFWFAECPMXYBEKMMOPBNPQLXATCBWJQNQGCAKWORZGDXREORLOHAAEJBDKJQJYVSDVNHIRODHIYQZXK';
let channelMode = 'public';   // public, private, restricted

const executeDataRetrieval = async function(rootVal, keyVal) {

    // ***********************************************************************************
    // 1. https://github.com/iotaledger/mam.client.js/blob/master/src/index.js
    //    const fetch = async (address, mode, sidekey, callback, rounds = 81)
    //    Calculate address based on the root.
    //    If channel mode = public: address = root
    //    If channel mode = private or restricted: address = HASH(root)
    // 2. https://github.com/iotaledger/iota.lib.js/blob/master/lib/api/api.js
    //    Call the api.prototype.findTransactions = function(searchValues, callback) and
    //    pass the address as searchValue.
    //    This will return a list of transaction hashes.
    // 3. https://github.com/iotaledger/iota.lib.js/blob/master/lib/api/api.js
    //    Call api.prototype.getTransactionsObjects = function(hashes, callback)
    //    This will get the transaction objects from a list of transaction hashes.
    // 4. https://github.com/iotaledger/mam.client.js/blob/master/src/index.js
    //    Call const txHashesToMessages = async hashes
    //    This will get the signatureMessageFragments from the transaction objects.
    //    Create the original masked payload.
    // 5. https://github.com/iotaledger/mam.client.js/blob/master/src/index.js
    //    Call const decode = (payload, side_key, root)
    //    https://github.com/iotaledger/mam.client.js/blob/master/lib/mam.web.js
    //    Call function decodeMessage(PAYLOAD, SIDE_KEY, ROOT)
    //
    //    See IOTA Tutorial 20
    //    https://github.com/iotaledger/MAM/blob/master/mam/src/mam.rs
    //    Call pub fn parse<C>(payload: &mut [Trit], side_key: &[Trit], root: &[Trit],..)
    //    This function parses an encrypted 'payload', by first decrypting it with a
    //    side_key and root as initialization vector.
    //    Then checks that the signature is valid and with sibling hashes in the payload
    //    resolves to the merkle 'root'.
    //    Returns the unmasked message and the next root.
    // ***********************************************************************************


    let resp = await Mam.fetch(rootVal, channelMode, keyVal, function(data) {

        // For example: json = { "data": 54, "dateTime": "23/02/2018 10:30:07" }
        let json = JSON.parse(iota.utils.fromTrytes(data));
        console.log(JSON.stringify(json,null,"\t"));
    });

    // The resp looks like:
    // {
    // "nextRoot": "MIXGFXOGWRSJIEPOOXCOXIZGUKASMSTQTSGOSHNICLZHAHJPEHDCPIBHADLRLZWYBWFHDB9AHXPOXLBEY"
    // }

    // ***********************************************************************************
    // 6. The resp contains the next root.
    //    The function executeDataRetrieval is again executed with the next root and
    //    side key as input.
    //    This process keeps repeating itself.
    // ***********************************************************************************

    executeDataRetrieval(resp.nextRoot, keyVal);
}

executeDataRetrieval(root, key);
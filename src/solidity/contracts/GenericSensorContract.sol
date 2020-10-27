pragma solidity >=0.6.4 <=0.7.3;

import "./GeneralContract.sol";

contract GenericSensorContract {
    /* Structs */
    struct Record {
        // Day and second of
        uint16 day;
        uint24 second;
        uint8 valueStored;
    }

    /* Attributes of the contract in context of the organization */
    // General contract of the organization
    GeneralContract gc;
    // Address of person responsible to this sensor
    address responsible;

    /* State variables */
    // Value picked up by the sensor
    uint8 lastValueRead;
    // Rate of reading the sensor (ms)
    uint32 rate;
    // Number of items in the whole History
    uint historyLength;
    // History of records for anormal values
    mapping(uint => Record) History;
    // Number of record this month
    uint32 month;
    // Index of records, indexed per month (from 0 to current month)
    mapping(uint32 => uint) monthCount;

    /* Modifiers */
    modifier onlyContract{
        require(msg.sender == address(gc),"can only be called from the company contract");
        _;
    }

    /* Setters / getters */
    function getHistoryLength(address) public view returns (uint) {
        return historyLength;
    }

    function getRate() public view returns (uint32) {
        return rate;
    }

    function setRate(uint32 newRate) public onlyContract {
        // TODO: authentication in setters
        rate = newRate;
    }

    /* Functionality */
    // This function registers a abnormal value in the Record
    // It's passed the month so that that computation is made on the client
    //function registerValue(uint8 value, uint32 month) public onlyContract {
        //Record r;
        //
        //// Increment index
        //historyLength++;
        //monthCount[month]++;
    //} 
}

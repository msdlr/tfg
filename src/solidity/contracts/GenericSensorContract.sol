pragma solidity >=0.6.4 <=0.7.3;

import "./GeneralContract.sol";

abstract contract GenericSensorContract {
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
    // Number of items in the whole History
    uint historyLength;
    // History of records for anormal values
    mapping(uint => Record) History;
    // Number of record this month
    uint32 monthCount;
    // Index of records, indexed per month (from 0 to current month)
    mapping(uint32 => uint) monthlyRecord;

    /* Business logic variables */

    // Rate of reading the sensor (ms)
    uint32 rate;
    // Average values for the constant
    uint32 constant avgValue=0;
    uint32 constant maxOk=0; // If the value is above this, we have to register it and notify
    uint32 constant minOk=0; // If the value is under this, we have to register it and notify
    uint32 constant warningMin=0; // Under this value, we only register it
    uint32 constant warningMax=0; // Above this value, we only register it

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
    function registerValue(uint8 value) public onlyContract {
        // Create entry in the sensor history
        Record memory r;
        // Append data
        r.day = getToday();
        r.second = secondOfDay();
        r.valueStored = value;

        // Increment index
        historyLength++;
        monthlyRecord[monthCount]++;
    }

    // The client sends the value picked up by the sensor and this function
    // evaluates wether to store it and/or notify
    function storeNotify(uint8 value) public view returns (bool mustStore, bool mustNotify) {
        // Evaluate if we need to notify
        mustNotify = (value >= maxOk || value <= minOk);
        // If we have to notify we also have to write the record
        if (mustNotify) return(true, true);
        // Check if only a write in the record is needed
        else mustStore = (value >= warningMax || value <= warningMin);
    }


    /* AUXILIARY FUNCTIONS */
    //Day number since 1/1/2020 (UNIX time + 50 years)
    function getToday() private view returns(uint16 today){
        uint day = (block.timestamp / 1 days) - (50*365 days);
        return uint16(day);
    }
    function secondOfDay() private view returns(uint24 sec){
        sec = uint24(block.timestamp % getToday());
    }
}

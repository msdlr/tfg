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
    GeneralContract private gc;
    // Address of person responsible to this sensor
    address private responsible;

    /* State variables */

    // Value picked up by the sensor
    uint32 private lastValueRead;
    // Number of items in the whole History
    uint32 private historyLength;
    // History of records for anormal values
    mapping(uint => Record) private History;
    // Number of record this month
    uint32 private  monthCount;
    // Index of records, indexed per month (from 0 to current month)
    mapping(uint32 => uint) private monthlyRecord;

    /* Business logic variables */

    // Rate of reading the sensor (ms)
    uint32 private defaultRate;
    uint32 private rate;
    // Average values for the constant
    uint32 private avgValue;
    uint32 private maxOk; // If the value is above this, we have to register it and notify
    uint32 private minOk; // If the value is under this, we have to register it and notify
    uint32 private warningMin; // Under this value, we only register it
    uint32 private warningMax; // Above this value, we only register it

    /* Modifiers */
    modifier onlyContract{
        require(msg.sender == address(gc),"can only be called from the company contract");
        _;
    }

    /* Setters / getters */
    function getHistoryLength() public view returns (uint32) {
        return historyLength;
    }

    function getRate() public view returns (uint32) {
        return rate;
    }

    function getDefaultRate() public view returns (uint32) {
        return defaultRate;
    }

    function setRate(uint32 _newRate) public onlyContract {
        // TODO: authentication in setters
        rate = _newRate;
    }
    
    function getLastValueRead() public view returns (uint32) {
        return lastValueRead;
    }

    /* Functionality */
    // This function registers a abnormal value in the Record
    // It's passed the month so that that computation is made on the client
    function registerValue(uint8 _value) public onlyContract {
        // Create entry in the sensor history
        Record memory r;
        // Append data
        r.day = getToday();
        r.second = secondOfDay();
        r.valueStored = _value;

        // Increment index
        historyLength++;
        monthlyRecord[monthCount]++;
    }

    // The client sends the value picked up by the sensor and this function
    // evaluates wether to store it and/or notify
    function storeNotify(uint8 _value) public view returns (bool mustStore, bool mustNotify) {
        // Evaluate if we need to notify
        mustNotify = (_value >= maxOk || _value <= minOk);
        // If we have to notify we also have to write the record
        if (mustNotify) return(true, true);
        // Check if only a write in the record is needed
        else mustStore = (_value >= warningMax || _value <= warningMin);
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

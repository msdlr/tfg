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
    uint8 valueRead;
    // Rate of reading the sensor (ms)
    uint32 rate;
    // History of records for anormal values
    mapping(uint32 => Record) History;
}

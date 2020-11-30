pragma solidity >=0.6.4;

import "./GeneralContract.sol";

contract AuthContract {

    // This contract is created by the administrators
    // for every person that use the login system

    /* STRUCTS */
    struct OTP {
        bytes32 passHash;
        // Timestamp relative to today
        uint24 time; // seconds in a day: 2 ^ 16.39
        bool valid;
        uint16 ttl;
        // The OTP can expire the next day it's issued (p.e. 00:01)
        uint16 date; // 2^16 days is about 179  years
    }

    /* MODIFIERS */
 
    modifier onlyContract{
        require(msg.sender == address(gc),"can only be called from the company contract");
        _;
    }

    modifier validOTP{
        require(eOTP.valid == true, "This OTP is not valid");
        require((eOTP.date) * 1 days + secondOfDay() + (eOTP.ttl) > getToday() * 1 days + secondOfDay(), "timestamp expired");
        _;
    }

    /* STATUS VARIABLES */
    GeneralContract gc;
    address employee;
    OTP eOTP;

    constructor (GeneralContract _genContract, address _employee) public payable{
        // Set status data
        gc = _genContract;
        employee = _employee;
        // No need to generate an invalid token
        // default values already provide it
    }

    /* FUNCTIONS */
    function tryLogin(uint16 _pass) public validOTP returns(bool success) {
        // We just revert if the OTP is not valid
        success = false;
        require (keccak256(abi.encode(_pass)) == eOTP.passHash, "The password is not correct");
        eOTP.valid = false;
        return true;
    }

    // Returns the generated pass and generate the OTP struture
    function newOTP() public onlyContract returns (uint16 pass_){
        // Generate the OTP number
        uint16 p = uint16(uint256(keccak256(abi.encode(block.timestamp, msg.sender))) % 9998) +1;

        //Fill the OTP fields:
        // Timestamp: relative to today instead of 1970
        eOTP.time = uint24(block.timestamp % 1 days); // Timestamp relative to the day
        // TTL
        eOTP.ttl = uint16(5 minutes);
        // OTP day
        eOTP.date = getToday();
        // Used flag
        eOTP.valid = true;
        // Pass Hash
        eOTP.passHash = keccak256(abi.encode(p));
        
        // Return p to be retreived by the interface
        return p;
    }

    function terminate() external onlyContract{
        // We already require that msg.sender is the general contract
        selfdestruct(msg.sender);
    }

    //Day number since 1/1/2020 (UNIX time + 50 years)
    function getToday() private view returns(uint16 today){
        uint day = (block.timestamp / 1 days) - (50*365 days);
        return uint16(day);
    }
    function secondOfDay() private view returns(uint24 sec){
        sec = uint24(block.timestamp % getToday());
    }
}

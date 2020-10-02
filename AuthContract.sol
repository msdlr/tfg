pragma solidity >=0.4.22 <0.7.0;
import "./General.sol";

contract AuthContract {

    // This contract is created by the administrators
    // for every person that use the login system

    /* STRUCTS */
    struct OTP {
        bytes32 passHash;
        // Timestamp relative to today
        uint24 time; // seconds in a day: 2 ^ 16.39
        bool isUsed;
        uint16 ttl;
        // The OTP can expire the next day it's issued (p.e. 00:01)
        uint32 date; // 2^32 days is about 136 years
    }

    /* MODIFIERS */
 
    modifier onlyContract{
        require(msg.sender == gc,"can only be called from the company contract");
        _;
    }

    modifier validOTP{
        require(eOTP.isUsed == false);
        require(eOTP.date days + ttl seconds <= getToday() days + getSeconds() seconds);
    }

    /* STATUS VARIABLES */
    GeneralContract gc;
    address employee;
    OTP eOTP;

    constructor (GeneralContract _genContract, address _employee) public payable{
        // Set status data
        gc = _genContract;
        employee = _employee;
        // Generate an invalid token
        eOTP = OTP(0,0,true,0,0);
    }

    /* FUNCTIONS */
    function tryLogin(uint16 _pass) public validOTP {
        // We just revert if the OTP is not valid
        require(checkValid());
        require (keccak256(abi.encode(_pass)) == eOTP.passHash, "The password is not correct");
    }

    // Returns the generated pass and generate the OTP struture
    function newOTP() public onlyContract returns (uint16 pass_){
        // Generate the OTP number
        uint16 p = uint16(uint256(keccak256(abi.encode(now, msg.sender))) % 9999);

        //Fill the OTP fields:
        // Timestamp: relative to today instead of 1970
        eOTP.time = uint24(now % 1 days); // Timestamp relative to the day
        // TTL
        eOTP.ttl = uint16(5 minutes);
        // OTP day
        eOTP.date = getToday();
        // Used flag
        eOTP.isUsed = false;
        // Pass Hash
        eOTP.passHash = keccak256(abi.encode(p));
        
        // Return p to be retreived by the interface
        return p;
    }

    function terminate() external fromgContract{
        // We already require that msg.sender is the general contract
        selfdestruct(msg.sender);
    }

    //Day number since 1/1/2020 (UNIX time + 50 years)
    function getToday() private view returns(uint8 today){
        uint day = (now / 1 days) - (50*365 days);
        return uint8(day);
    }
    function secondOfDay() private view returns(uint24 sec){
        uint24 sec = now % getToday();
        return sec;
    }
}